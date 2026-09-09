package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type DB struct {
	*sql.DB
}

func Connect(databaseURL string) (*DB, error) {
	var db *sql.DB
	var err error

	// Retry loop for database connection (waiting for Postgres to be ready in Docker)
	for i := 0; i < 15; i++ {
		db, err = sql.Open("postgres", databaseURL)
		if err == nil {
			err = db.Ping()
			if err == nil {
				break
			}
		}
		log.Printf("Waiting for Postgres to be ready... (%d/15): %v", i+1, err)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)

	dbInstance := &DB{db}
	if err := dbInstance.AutoMigrate(); err != nil {
		return nil, fmt.Errorf("auto migration failed: %w", err)
	}

	return dbInstance, nil
}

func (db *DB) AutoMigrate() error {
	queries := []string{
		`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`,
		`CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			username VARCHAR(50) UNIQUE NOT NULL,
			email VARCHAR(255) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		);`,
		`CREATE TABLE IF NOT EXISTS tracks (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			user_id UUID REFERENCES users(id) ON DELETE CASCADE,
			title VARCHAR(150) NOT NULL,
			artist VARCHAR(100) NOT NULL,
			audio_path VARCHAR(500) NOT NULL,
			duration INT DEFAULT 0,
			mood_tag VARCHAR(50) DEFAULT 'twilight',
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		);`,
		`CREATE TABLE IF NOT EXISTS moodboard_items (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			track_id UUID REFERENCES tracks(id) ON DELETE CASCADE,
			item_type VARCHAR(20) NOT NULL, -- 'image', 'quote', 'note', 'palette'
			content TEXT,
			media_path VARCHAR(500),
			thumb_path VARCHAR(500),
			sort_order INT DEFAULT 0,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		);`,
		`ALTER TABLE moodboard_items ADD COLUMN IF NOT EXISTS thumb_path VARCHAR(500);`,
	}

	for _, q := range queries {
		if _, err := db.Exec(q); err != nil {
			return fmt.Errorf("migration query failed (%s): %w", q, err)
		}
	}

	log.Println("Database auto-migration completed successfully.")
	db.seedInitialData()
	return nil
}

func (db *DB) seedInitialData() {
	var userCount int
	err := db.QueryRow("SELECT COUNT(*) FROM users;").Scan(&userCount)
	if err != nil || userCount > 0 {
		return
	}

	log.Println("Seeding initial quiet dusk sample data...")
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("silentpassword123"), bcrypt.DefaultCost)

	var userID string
	err = db.QueryRow(
		`INSERT INTO users (username, email, password_hash) 
		 VALUES ($1, $2, $3) RETURNING id;`,
		"aimer_listener", "calm@song-moodboard.io", string(hashedPassword),
	).Scan(&userID)
	if err != nil {
		log.Printf("Failed to seed demo user: %v", err)
		return
	}

	// Seed Sample Track 1
	var trackID string
	err = db.QueryRow(
		`INSERT INTO tracks (user_id, title, artist, audio_path, duration, mood_tag)
		 VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;`,
		userID, "Rokutousei no Yoru (Six-Magnitude Star)", "Aimer / Ambient Piano",
		"/demo-audio/quiet-night.mp3", 240, "twilight",
	).Scan(&trackID)
	if err != nil {
		log.Printf("Failed to seed demo track: %v", err)
		return
	}

	// Seed Sample Moodboard Items for Track 1
	items := []struct {
		itemType  string
		content   string
		mediaPath string
		sortOrder int
	}{
		{
			itemType:  "image",
			content:   "Dusk mist over silent city lights",
			mediaPath: "/images/hero-dusk.jpg",
			sortOrder: 1,
		},
		{
			itemType:  "quote",
			content:   "Even the faint, dim stars continue to shine silently through the winter fog.",
			mediaPath: "",
			sortOrder: 2,
		},
		{
			itemType:  "palette",
			content:   "#EBF3FA,#8FB3D5,#5C82A6,#1E293B,#F8FAFC",
			mediaPath: "",
			sortOrder: 3,
		},
		{
			itemType:  "image",
			content:   "Rain drops on blue dusk window",
			mediaPath: "/images/mood-rain.jpg",
			sortOrder: 4,
		},
		{
			itemType:  "note",
			content:   "Listening at 5:30 AM before the city stirs. The air is cold, but the piano chords feel warm.",
			mediaPath: "",
			sortOrder: 5,
		},
		{
			itemType:  "image",
			content:   "Calm sea meeting mist horizon",
			mediaPath: "/images/mood-sea.jpg",
			sortOrder: 6,
		},
	}

	for _, it := range items {
		_, _ = db.Exec(
			`INSERT INTO moodboard_items (track_id, item_type, content, media_path, sort_order)
			 VALUES ($1, $2, $3, $4, $5);`,
			trackID, it.itemType, it.content, it.mediaPath, it.sortOrder,
		)
	}

	log.Println("Initial demo track & moodboard seeded successfully.")
}
