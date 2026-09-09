package handlers

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/ijon6k/song-moodboard/internal/auth"
	"github.com/ijon6k/song-moodboard/internal/models"
)

type ExtractRequest struct {
	URL          string `json:"url"`
	MoodTag      string `json:"mood_tag"`
	CustomTitle  string `json:"custom_title"`
	CustomArtist string `json:"custom_artist"`
}

type ExtractJob struct {
	ID        string        `json:"id"`
	URL       string        `json:"url"`
	Status    string        `json:"status"` // "processing", "completed", "failed"
	Progress  int           `json:"progress"`
	Stage     string        `json:"stage"`
	Track     *models.Track `json:"track,omitempty"`
	Error     string        `json:"error,omitempty"`
	CreatedAt time.Time     `json:"created_at"`
}

var (
	jobsMu sync.RWMutex
	jobs   = make(map[string]*ExtractJob)
)

func getJob(id string) (*ExtractJob, bool) {
	jobsMu.RLock()
	defer jobsMu.RUnlock()
	j, ok := jobs[id]
	return j, ok
}

func setJob(j *ExtractJob) {
	jobsMu.Lock()
	defer jobsMu.Unlock()
	jobs[j.ID] = j
}

func updateJob(id string, status string, progress int, stage string, track *models.Track, errMsg string) {
	jobsMu.Lock()
	defer jobsMu.Unlock()
	if j, ok := jobs[id]; ok {
		j.Status = status
		j.Progress = progress
		j.Stage = stage
		if track != nil {
			j.Track = track
		}
		if errMsg != "" {
			j.Error = errMsg
		}
	}
}

// POST /api/tracks/extract
func (h *Handler) ExtractStream(w http.ResponseWriter, r *http.Request) {
	var req ExtractRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	req.URL = strings.TrimSpace(req.URL)
	if req.URL == "" {
		writeError(w, http.StatusBadRequest, "stream url is required")
		return
	}

	if req.MoodTag == "" {
		req.MoodTag = "Twilight"
	}

	// Resolve user ID: if authenticated, use logged in user; else fallback to default user
	var userID uuid.UUID
	claims, ok := auth.GetUserFromContext(r.Context())
	if ok && claims != nil {
		userID = claims.UserID
	} else {
		// Fallback query first user from DB
		var defaultID uuid.UUID
		err := h.DB.QueryRowContext(r.Context(), "SELECT id FROM users ORDER BY created_at ASC LIMIT 1;").Scan(&defaultID)
		if err == nil {
			userID = defaultID
		} else {
			userID = uuid.MustParse("17faed3f-08c8-430d-92b9-f856267289a6")
		}
	}

	jobID := uuid.New().String()
	job := &ExtractJob{
		ID:        jobID,
		URL:       req.URL,
		Status:    "processing",
		Progress:  5,
		Stage:     "Connecting to stream source...",
		CreatedAt: time.Now(),
	}
	setJob(job)

	// Run extraction asynchronously in background
	go h.processStreamExtraction(jobID, req, userID)

	writeJSON(w, http.StatusAccepted, map[string]interface{}{
		"job_id":  jobID,
		"message": "extraction job started",
		"status":  "processing",
	})
}

// GET /api/tracks/extract/{id}
func (h *Handler) GetExtractProgress(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "job id required")
		return
	}

	job, ok := getJob(id)
	if !ok {
		writeError(w, http.StatusNotFound, "job not found")
		return
	}

	writeJSON(w, http.StatusOK, job)
}

type ytdlpMetadata struct {
	Title     string  `json:"title"`
	Uploader  string  `json:"uploader"`
	Artist    string  `json:"artist"`
	Channel   string  `json:"channel"`
	Duration  float64 `json:"duration"`
	Thumbnail string  `json:"thumbnail"`
}

func (h *Handler) processStreamExtraction(jobID string, req ExtractRequest, userID uuid.UUID) {
	ctx, cancel := context.WithTimeout(context.Background(), 7*time.Minute)
	defer cancel()

	defer func() {
		// Clean up any remaining temp files for this job
		files, _ := filepath.Glob(fmt.Sprintf("/tmp/downloads/%s.*", jobID))
		for _, f := range files {
			_ = os.Remove(f)
		}
	}()

	// -------------------------------------------------------------
	// Stage 1: Inspect Metadata (5% - 20%)
	// -------------------------------------------------------------
	updateJob(jobID, "processing", 10, "Inspecting audio metadata...", nil, "")

	metaCmd := exec.CommandContext(ctx, "yt-dlp",
		"--dump-single-json",
		"--no-playlist",
		"--no-warnings",
		req.URL,
	)
	metaOut, err := metaCmd.Output()
	var meta ytdlpMetadata
	if err == nil && len(metaOut) > 0 {
		_ = json.Unmarshal(metaOut, &meta)
	}

	title := strings.TrimSpace(req.CustomTitle)
	if title == "" {
		title = strings.TrimSpace(meta.Title)
	}
	if title == "" {
		title = "Extracted Audio Track"
	}

	artist := strings.TrimSpace(req.CustomArtist)
	if artist == "" {
		if meta.Artist != "" {
			artist = strings.TrimSpace(meta.Artist)
		} else if meta.Uploader != "" {
			artist = strings.TrimSpace(meta.Uploader)
		} else if meta.Channel != "" {
			artist = strings.TrimSpace(meta.Channel)
		}
	}
	if artist == "" {
		artist = "Silent Stream"
	}

	duration := int(meta.Duration)
	if duration <= 0 {
		duration = 180
	}

	updateJob(jobID, "processing", 22, fmt.Sprintf("Extracting: %s - %s", title, artist), nil, "")

	// -------------------------------------------------------------
	// Stage 2: Download & Extract Audio to MP3 via yt-dlp & FFmpeg
	// -------------------------------------------------------------
	outputPathPattern := fmt.Sprintf("/tmp/downloads/%s.%%(ext)s", jobID)
	targetMp3Path := fmt.Sprintf("/tmp/downloads/%s.mp3", jobID)

	dlCmd := exec.CommandContext(ctx, "yt-dlp",
		"-x",
		"--audio-format", "mp3",
		"--audio-quality", "0",
		"--newline",
		"--no-playlist",
		"-o", outputPathPattern,
		req.URL,
	)

	stdoutPipe, err := dlCmd.StdoutPipe()
	if err != nil {
		updateJob(jobID, "failed", 0, "Failed to initialize audio extraction", nil, err.Error())
		return
	}
	dlCmd.Stderr = dlCmd.Stdout // merge stderr for detailed diagnostics

	if err := dlCmd.Start(); err != nil {
		updateJob(jobID, "failed", 0, "Failed to start extractor process", nil, err.Error())
		return
	}

	// Parse progress lines like "[download]  45.2% of ..."
	pctRegex := regexp.MustCompile(`\[download\]\s+([0-9.]+)%`)
	scanner := bufio.NewScanner(stdoutPipe)

	for scanner.Scan() {
		line := scanner.Text()
		matches := pctRegex.FindStringSubmatch(line)
		if len(matches) > 1 {
			if pct, err := strconv.ParseFloat(matches[1], 64); err == nil {
				// Map 0-100% download to 22-75% overall progress
				mapped := 22 + int(pct*0.53)
				if mapped > 75 {
					mapped = 75
				}
				updateJob(jobID, "processing", mapped, fmt.Sprintf("Extracting audio stream (%d%%)...", int(pct)), nil, "")
			}
		} else if strings.Contains(line, "[ExtractAudio]") {
			updateJob(jobID, "processing", 76, "Encoding high-fidelity MP3 with FFmpeg...", nil, "")
		}
	}

	if err := dlCmd.Wait(); err != nil {
		log.Printf("Extraction error for URL %s: %v", req.URL, err)
		updateJob(jobID, "failed", 0, "Extraction failed", nil, "Unable to extract audio from stream. URL may be restricted or unsupported.")
		return
	}

	// -------------------------------------------------------------
	// Stage 3: Store MP3 in SeaweedFS Object Storage (78% - 88%)
	// -------------------------------------------------------------
	updateJob(jobID, "processing", 80, "Saving audio to SeaweedFS storage...", nil, "")

	// Locate resulting mp3 file
	var finalAudioFile string
	if _, err := os.Stat(targetMp3Path); err == nil {
		finalAudioFile = targetMp3Path
	} else {
		// Fallback glob
		candidates, _ := filepath.Glob(fmt.Sprintf("/tmp/downloads/%s.*", jobID))
		for _, c := range candidates {
			if strings.HasSuffix(strings.ToLower(c), ".mp3") || strings.HasSuffix(strings.ToLower(c), ".m4a") || strings.HasSuffix(strings.ToLower(c), ".opus") {
				finalAudioFile = c
				break
			}
		}
	}

	if finalAudioFile == "" {
		updateJob(jobID, "failed", 0, "Converted audio file not found", nil, "No audio file produced by extractor.")
		return
	}

	f, err := os.Open(finalAudioFile)
	if err != nil {
		updateJob(jobID, "failed", 0, "Failed to read extracted audio", nil, err.Error())
		return
	}
	defer f.Close()

	stat, _ := f.Stat()
	audioObjName := fmt.Sprintf("audio/%s.mp3", uuid.New().String())
	audioPath, err := h.Storage.Upload(ctx, audioObjName, f, stat.Size(), "audio/mpeg")
	if err != nil {
		updateJob(jobID, "failed", 0, "Failed to upload audio to storage", nil, err.Error())
		return
	}

	// -------------------------------------------------------------
	// Stage 4: Process Artwork & Thumbnails (88% - 95%)
	// -------------------------------------------------------------
	updateJob(jobID, "processing", 89, "Processing album artwork & thumbnails...", nil, "")

	var origImgPath string
	var thumbImgPath string

	if meta.Thumbnail != "" {
		reqThumb, err := http.NewRequestWithContext(ctx, http.MethodGet, meta.Thumbnail, nil)
		if err == nil {
			reqThumb.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")
			client := &http.Client{Timeout: 15 * time.Second}
			resp, errResp := client.Do(reqThumb)
			if errResp == nil && resp.StatusCode == http.StatusOK {
				imgBytes, errRead := io.ReadAll(resp.Body)
				_ = resp.Body.Close()

				if errRead == nil && len(imgBytes) > 0 {
					imgUUID := uuid.New().String()
					origName := fmt.Sprintf("images/orig_%s.jpg", imgUUID)
					origImgPath, _ = h.Storage.Upload(ctx, origName, bytes.NewReader(imgBytes), int64(len(imgBytes)), "image/jpeg")

					// Generate 500px responsive thumbnail
					thumbBytes := generateThumbnail(imgBytes, 500)
					if len(thumbBytes) > 0 {
						thumbName := fmt.Sprintf("images/thumb_%s.jpg", imgUUID)
						thumbImgPath, _ = h.Storage.Upload(ctx, thumbName, bytes.NewReader(thumbBytes), int64(len(thumbBytes)), "image/jpeg")
					} else {
						thumbImgPath = origImgPath
					}
				}
			}
		}
	}

	// If no artwork was fetched, fallback to null or default
	if origImgPath == "" {
		origImgPath = "/images/hero.jpg"
		thumbImgPath = "/images/hero.jpg"
	}

	// -------------------------------------------------------------
	// Stage 5: Register Track & Moodboard in Database (95% - 100%)
	// -------------------------------------------------------------
	updateJob(jobID, "processing", 96, "Registering track into your library...", nil, "")

	var track models.Track
	insertTrackQuery := `INSERT INTO tracks (user_id, title, artist, audio_path, duration, mood_tag) 
	                     VALUES ($1, $2, $3, $4, $5, $6) 
	                     RETURNING id, user_id, title, artist, audio_path, duration, mood_tag, created_at;`
	err = h.DB.QueryRowContext(ctx, insertTrackQuery, userID, title, artist, audioPath, duration, req.MoodTag).
		Scan(&track.ID, &track.UserID, &track.Title, &track.Artist, &track.AudioPath, &track.Duration, &track.MoodTag, &track.CreatedAt)
	if err != nil {
		updateJob(jobID, "failed", 0, "Failed to save track record to database", nil, err.Error())
		return
	}

	// Insert artwork into moodboard items
	moodboardQuery := `INSERT INTO moodboard_items (track_id, item_type, content, media_path, thumb_path, sort_order) 
	                   VALUES ($1, 'image', $2, $3, $4, 0);`
	_, _ = h.DB.ExecContext(ctx, moodboardQuery, track.ID, title+" Artwork", origImgPath, thumbImgPath)

	// Invalidate Redis tracks cache
	if h.Cache != nil {
		_ = h.Cache.Del(ctx, "tracks:all")
	}

	// -------------------------------------------------------------
	// Final Stage: Success Complete (100%)
	// -------------------------------------------------------------
	updateJob(jobID, "completed", 100, "Track successfully extracted and added to your library!", &track, "")
}
