package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}

type Track struct {
	ID             uuid.UUID       `json:"id"`
	UserID         uuid.UUID       `json:"user_id"`
	Title          string          `json:"title"`
	Artist         string          `json:"artist"`
	AudioPath      string          `json:"audio_path"`
	Duration       int             `json:"duration"` // in seconds
	MoodTag        string          `json:"mood_tag"` // e.g. "twilight", "dusk", "fog"
	WaveformData   []float64       `json:"waveform_data,omitempty"`
	CreatedAt      time.Time       `json:"created_at"`
	MoodboardItems []MoodboardItem `json:"moodboard_items,omitempty"`
}

type MoodboardItem struct {
	ID        uuid.UUID `json:"id"`
	TrackID   uuid.UUID `json:"track_id"`
	ItemType  string    `json:"item_type"` // "image", "quote", "note", "palette"
	Content   string    `json:"content"`   // text quote / hex palette csv / memo
	MediaPath string    `json:"media_path"`// SeaweedFS url if image
	ThumbPath string    `json:"thumb_path,omitempty"` // SeaweedFS url for compressed thumbnail
	SortOrder int       `json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
}

// DTOs
type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type AddMoodboardItemRequest struct {
	ItemType string `json:"item_type"`
	Content  string `json:"content"`
}
