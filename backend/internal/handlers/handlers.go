package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/ijon6k/song-moodboard/internal/auth"
	"github.com/ijon6k/song-moodboard/internal/cache"
	"github.com/ijon6k/song-moodboard/internal/database"
	"github.com/ijon6k/song-moodboard/internal/models"
	"github.com/ijon6k/song-moodboard/internal/storage"
)

type Handler struct {
	DB      *database.DB
	Cache   *cache.Cache
	Storage *storage.Storage
	Auth    *auth.AuthManager
}

func New(db *database.DB, c *cache.Cache, st *storage.Storage, a *auth.AuthManager) *Handler {
	return &Handler{
		DB:      db,
		Cache:   c,
		Storage: st,
		Auth:    a,
	}
}

// Helpers
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

// Health Check
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"status":    "healthy",
		"service":   "song-moodboard-api",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

// Auth Handlers
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if strings.TrimSpace(req.Username) == "" || strings.TrimSpace(req.Email) == "" || len(req.Password) < 6 {
		writeError(w, http.StatusBadRequest, "username, email, and password (min 6 chars) are required")
		return
	}

	hash, err := h.Auth.HashPassword(req.Password)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to hash password")
		return
	}

	var user models.User
	query := `INSERT INTO users (username, email, password_hash) 
	          VALUES ($1, $2, $3) 
	          RETURNING id, username, email, created_at;`
	err = h.DB.QueryRowContext(r.Context(), query, req.Username, req.Email, hash).
		Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt)
	if err != nil {
		writeError(w, http.StatusConflict, "username or email already exists")
		return
	}

	token, err := h.Auth.GenerateToken(user.ID, user.Username, user.Email)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to generate token")
		return
	}

	writeJSON(w, http.StatusCreated, models.AuthResponse{
		Token: token,
		User:  user,
	})
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	var user models.User
	query := `SELECT id, username, email, password_hash, created_at FROM users WHERE email = $1;`
	err := h.DB.QueryRowContext(r.Context(), query, req.Email).
		Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "invalid email or password")
		return
	}

	if !h.Auth.CheckPassword(user.PasswordHash, req.Password) {
		writeError(w, http.StatusUnauthorized, "invalid email or password")
		return
	}

	token, err := h.Auth.GenerateToken(user.ID, user.Username, user.Email)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to generate token")
		return
	}

	writeJSON(w, http.StatusOK, models.AuthResponse{
		Token: token,
		User:  user,
	})
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.GetUserFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var user models.User
	query := `SELECT id, username, email, created_at FROM users WHERE id = $1;`
	err := h.DB.QueryRowContext(r.Context(), query, claims.UserID).
		Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt)
	if err != nil {
		writeError(w, http.StatusNotFound, "user not found")
		return
	}

	writeJSON(w, http.StatusOK, user)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	parts := strings.Split(authHeader, " ")
	if len(parts) == 2 && h.Cache != nil {
		_ = h.Cache.BlacklistToken(r.Context(), parts[1], 72*time.Hour)
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "logged out successfully"})
}

// Track Handlers
func (h *Handler) ListTracks(w http.ResponseWriter, r *http.Request) {
	// Try Redis Cache First
	cacheKey := "tracks:all"
	if h.Cache != nil {
		cached, err := h.Cache.Get(r.Context(), cacheKey)
		if err == nil && cached != "" {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("X-Cache", "HIT")
			_, _ = w.Write([]byte(cached))
			return
		}
	}

	query := `SELECT id, user_id, title, artist, audio_path, duration, mood_tag, created_at 
	          FROM tracks ORDER BY created_at DESC;`
	rows, err := h.DB.QueryContext(r.Context(), query)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to query tracks")
		return
	}
	defer rows.Close()

	tracks := make([]models.Track, 0)
	for rows.Next() {
		var t models.Track
		if err := rows.Scan(&t.ID, &t.UserID, &t.Title, &t.Artist, &t.AudioPath, &t.Duration, &t.MoodTag, &t.CreatedAt); err != nil {
			continue
		}
		tracks = append(tracks, t)
	}

	jsonData, _ := json.Marshal(tracks)
	if h.Cache != nil {
		_ = h.Cache.Set(r.Context(), cacheKey, string(jsonData), 30*time.Second)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Cache", "MISS")
	_, _ = w.Write(jsonData)
}

func (h *Handler) GetTrack(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	trackID, err := uuid.Parse(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid track id")
		return
	}

	cacheKey := fmt.Sprintf("track:%s", trackID)
	if h.Cache != nil {
		cached, err := h.Cache.Get(r.Context(), cacheKey)
		if err == nil && cached != "" {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("X-Cache", "HIT")
			_, _ = w.Write([]byte(cached))
			return
		}
	}

	var t models.Track
	trackQuery := `SELECT id, user_id, title, artist, audio_path, duration, mood_tag, created_at 
	               FROM tracks WHERE id = $1;`
	err = h.DB.QueryRowContext(r.Context(), trackQuery, trackID).
		Scan(&t.ID, &t.UserID, &t.Title, &t.Artist, &t.AudioPath, &t.Duration, &t.MoodTag, &t.CreatedAt)
	if err != nil {
		writeError(w, http.StatusNotFound, "track not found")
		return
	}

	// Query Moodboard Items
	itemsQuery := `SELECT id, track_id, item_type, content, media_path, COALESCE(thumb_path, ''), sort_order, created_at 
	               FROM moodboard_items WHERE track_id = $1 ORDER BY sort_order ASC, created_at ASC;`
	rows, err := h.DB.QueryContext(r.Context(), itemsQuery, trackID)
	if err == nil {
		defer rows.Close()
		t.MoodboardItems = make([]models.MoodboardItem, 0)
		for rows.Next() {
			var it models.MoodboardItem
			if err := rows.Scan(&it.ID, &it.TrackID, &it.ItemType, &it.Content, &it.MediaPath, &it.ThumbPath, &it.SortOrder, &it.CreatedAt); err == nil {
				if it.ThumbPath == "" && it.ItemType == "image" {
					it.ThumbPath = it.MediaPath
				}
				t.MoodboardItems = append(t.MoodboardItems, it)
			}
		}
	}

	jsonData, _ := json.Marshal(t)
	if h.Cache != nil {
		_ = h.Cache.Set(r.Context(), cacheKey, string(jsonData), 60*time.Second)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Cache", "MISS")
	_, _ = w.Write(jsonData)
}

func (h *Handler) CreateTrack(w http.ResponseWriter, r *http.Request) {
	var userID uuid.UUID
	claims, ok := auth.GetUserFromContext(r.Context())
	if ok && claims != nil {
		userID = claims.UserID
	} else {
		var defaultID uuid.UUID
		err := h.DB.QueryRowContext(r.Context(), "SELECT id FROM users ORDER BY created_at ASC LIMIT 1;").Scan(&defaultID)
		if err == nil {
			userID = defaultID
		} else {
			userID = uuid.MustParse("17faed3f-08c8-430d-92b9-f856267289a6")
		}
	}

	// Limit request size to 32MB
	_ = r.ParseMultipartForm(32 << 20)

	title := r.FormValue("title")
	artist := r.FormValue("artist")
	moodTag := r.FormValue("mood_tag")
	if moodTag == "" {
		moodTag = "twilight"
	}

	audioURL := strings.TrimSpace(r.FormValue("audio_url"))
	if audioURL == "" {
		audioURL = strings.TrimSpace(r.FormValue("audio_path"))
	}

	duration := 180
	if dStr := r.FormValue("duration"); dStr != "" {
		if d, err := strconv.Atoi(dStr); err == nil && d > 0 {
			duration = d
		}
	}

	artworkURL := strings.TrimSpace(r.FormValue("artwork"))

	if strings.TrimSpace(title) == "" || strings.TrimSpace(artist) == "" {
		writeError(w, http.StatusBadRequest, "title and artist are required")
		return
	}

	var audioPath string
	file, header, err := r.FormFile("audio")
	if err == nil && file != nil {
		defer file.Close()
		ext := filepath.Ext(header.Filename)
		objectName := fmt.Sprintf("audio/%s%s", uuid.New().String(), ext)
		audioPath, err = h.Storage.Upload(r.Context(), objectName, file, header.Size, header.Header.Get("Content-Type"))
		if err != nil {
			writeError(w, http.StatusInternalServerError, "failed to upload audio to SeaweedFS: "+err.Error())
			return
		}
	} else if audioURL != "" {
		// VALIDATION: Reject raw YouTube or stream links in direct upload
		lowerURL := strings.ToLower(audioURL)
		if strings.Contains(lowerURL, "youtube.com") || strings.Contains(lowerURL, "youtu.be") {
			writeError(w, http.StatusBadRequest, "Tautan YouTube harus diekstrak lewat menu 'Extract from link' agar diunduh & dikonversi menjadi audio MP3.")
			return
		}
		if strings.Contains(lowerURL, "soundcloud.com") {
			writeError(w, http.StatusBadRequest, "Tautan SoundCloud harus diekstrak lewat menu 'Extract from link'.")
			return
		}

		// Validate direct audio link: check extension or perform HEAD check
		isAudioExt := strings.HasSuffix(lowerURL, ".mp3") || strings.HasSuffix(lowerURL, ".m4a") ||
			strings.HasSuffix(lowerURL, ".aac") || strings.HasSuffix(lowerURL, ".wav") ||
			strings.HasSuffix(lowerURL, ".flac") || strings.HasSuffix(lowerURL, ".ogg") ||
			strings.HasPrefix(lowerURL, "/api/storage/") || strings.HasPrefix(lowerURL, "/demo-audio/")

		if !isAudioExt && strings.HasPrefix(lowerURL, "http") {
			client := &http.Client{Timeout: 5 * time.Second}
			headResp, headErr := client.Head(audioURL)
			if headErr != nil || (!strings.HasPrefix(headResp.Header.Get("Content-Type"), "audio/") && !strings.HasPrefix(headResp.Header.Get("Content-Type"), "video/")) {
				writeError(w, http.StatusBadRequest, "Tautan bukan file audio yang valid (.mp3, .wav, .aac, dll). Pastikan URL langsung menuju file audio.")
				return
			}
		}
		audioPath = audioURL
	} else {
		// Fallback sample quiet track
		audioPath = "/demo-audio/quiet-night.mp3"
	}

	var track models.Track
	insertQuery := `INSERT INTO tracks (user_id, title, artist, audio_path, duration, mood_tag) 
	                VALUES ($1, $2, $3, $4, $5, $6) 
	                RETURNING id, user_id, title, artist, audio_path, duration, mood_tag, created_at;`
	err = h.DB.QueryRowContext(r.Context(), insertQuery, userID, title, artist, audioPath, duration, moodTag).
		Scan(&track.ID, &track.UserID, &track.Title, &track.Artist, &track.AudioPath, &track.Duration, &track.MoodTag, &track.CreatedAt)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to save track to database")
		return
	}

	if artworkURL != "" {
		moodItemQuery := `INSERT INTO moodboard_items (track_id, item_type, content, media_path, thumb_path, sort_order) 
		                  VALUES ($1, 'image', $2, $3, $4, 0);`
		_, _ = h.DB.ExecContext(r.Context(), moodItemQuery, track.ID, title+" Artwork", artworkURL, artworkURL)
	}

	// Invalidate tracks cache in Redis
	if h.Cache != nil {
		_ = h.Cache.Del(r.Context(), "tracks:all")
	}

	writeJSON(w, http.StatusCreated, track)
}

// DELETE /api/tracks/{id}
func (h *Handler) DeleteTrack(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	trackID, err := uuid.Parse(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid track id")
		return
	}

	// Delete associated moodboard items first
	_, _ = h.DB.ExecContext(r.Context(), "DELETE FROM moodboard_items WHERE track_id = $1;", trackID)

	// Delete track
	res, err := h.DB.ExecContext(r.Context(), "DELETE FROM tracks WHERE id = $1;", trackID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to delete track: "+err.Error())
		return
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		writeError(w, http.StatusNotFound, "track not found")
		return
	}

	// Invalidate Redis caches
	if h.Cache != nil {
		_ = h.Cache.Del(r.Context(), "tracks:all")
		_ = h.Cache.Del(r.Context(), fmt.Sprintf("track:%s", trackID))
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"message": "track deleted successfully",
		"id":      trackID.String(),
	})
}

func (h *Handler) AddMoodboardItem(w http.ResponseWriter, r *http.Request) {
	_, ok := auth.GetUserFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	idStr := chi.URLParam(r, "id")
	trackID, err := uuid.Parse(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid track id")
		return
	}

	_ = r.ParseMultipartForm(16 << 20)

	itemType := r.FormValue("item_type")
	content := r.FormValue("content")
	var mediaPath string
	var thumbPath string

	if itemType == "image" {
		file, header, err := r.FormFile("image")
		if err == nil && file != nil {
			defer file.Close()
			ext := filepath.Ext(header.Filename)
			if ext == "" {
				ext = ".jpg"
			}
			fileBytes, err := io.ReadAll(file)
			if err == nil && len(fileBytes) > 0 {
				fileUUID := uuid.New().String()
				origObjectName := fmt.Sprintf("images/orig_%s%s", fileUUID, ext)
				mediaPath, err = h.Storage.Upload(r.Context(), origObjectName, bytes.NewReader(fileBytes), int64(len(fileBytes)), header.Header.Get("Content-Type"))
				if err != nil {
					writeError(w, http.StatusInternalServerError, "failed to upload moodboard image: "+err.Error())
					return
				}

				// Downsample lightweight thumbnail (~500px)
				thumbBytes := generateThumbnail(fileBytes, 500)
				if len(thumbBytes) > 0 {
					thumbObjectName := fmt.Sprintf("images/thumb_%s.jpg", fileUUID)
					thumbPath, _ = h.Storage.Upload(r.Context(), thumbObjectName, bytes.NewReader(thumbBytes), int64(len(thumbBytes)), "image/jpeg")
				}
				if thumbPath == "" {
					thumbPath = mediaPath
				}
			}
		}
	}

	var item models.MoodboardItem
	query := `INSERT INTO moodboard_items (track_id, item_type, content, media_path, thumb_path) 
	          VALUES ($1, $2, $3, $4, $5) 
	          RETURNING id, track_id, item_type, content, media_path, thumb_path, sort_order, created_at;`
	err = h.DB.QueryRowContext(r.Context(), query, trackID, itemType, content, mediaPath, thumbPath).
		Scan(&item.ID, &item.TrackID, &item.ItemType, &item.Content, &item.MediaPath, &item.ThumbPath, &item.SortOrder, &item.CreatedAt)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to save moodboard item")
		return
	}

	// Invalidate track detail cache
	if h.Cache != nil {
		_ = h.Cache.Del(r.Context(), fmt.Sprintf("track:%s", trackID))
	}

	writeJSON(w, http.StatusCreated, item)
}

// SeaweedFS Object Proxy (Streams media from SeaweedFS directly to client via Nginx 1122)
func (h *Handler) StreamStorageObject(w http.ResponseWriter, r *http.Request) {
	objectName := chi.URLParam(r, "*")
	if objectName == "" {
		http.NotFound(w, r)
		return
	}

	obj, contentType, err := h.Storage.GetObject(r.Context(), objectName)
	if err != nil {
		http.Error(w, "object not found in storage", http.StatusNotFound)
		return
	}
	defer obj.Close()

	if contentType != "" {
		w.Header().Set("Content-Type", contentType)
	}
	w.Header().Set("Cache-Control", "public, max-age=86400")
	_, _ = io.Copy(w, obj)
}

func generateThumbnail(data []byte, maxDim int) []byte {
	src, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil
	}
	b := src.Bounds()
	w := b.Dx()
	h := b.Dy()
	if w <= 0 || h <= 0 {
		return nil
	}

	targetW, targetH := w, h
	if w > maxDim || h > maxDim {
		if w > h {
			targetW = maxDim
			targetH = int(float64(h) * float64(maxDim) / float64(w))
		} else {
			targetH = maxDim
			targetW = int(float64(w) * float64(maxDim) / float64(h))
		}
	}
	if targetW < 1 {
		targetW = 1
	}
	if targetH < 1 {
		targetH = 1
	}

	dst := image.NewRGBA(image.Rect(0, 0, targetW, targetH))
	for y := 0; y < targetH; y++ {
		srcY := y * h / targetH
		for x := 0; x < targetW; x++ {
			srcX := x * w / targetW
			dst.Set(x, y, src.At(b.Min.X+srcX, b.Min.Y+srcY))
		}
	}

	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, dst, &jpeg.Options{Quality: 78}); err != nil {
		return nil
	}
	return buf.Bytes()
}
