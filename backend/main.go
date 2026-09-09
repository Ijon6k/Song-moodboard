package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/ijon6k/song-moodboard/internal/auth"
	"github.com/ijon6k/song-moodboard/internal/cache"
	"github.com/ijon6k/song-moodboard/internal/config"
	"github.com/ijon6k/song-moodboard/internal/database"
	"github.com/ijon6k/song-moodboard/internal/handlers"
	"github.com/ijon6k/song-moodboard/internal/storage"
)

func main() {
	cfg := config.LoadConfig()

	log.Println("Starting Song Moodboard Backend API...")

	// 1. Connect PostgreSQL
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Fatal: Database initialization error: %v", err)
	}
	defer db.Close()

	// 2. Connect Redis
	redisCache, err := cache.Connect(cfg.RedisAddr, cfg.RedisPassword)
	if err != nil {
		log.Printf("Warning: Redis connection error: %v. Continuing without cache.", err)
	}

	// 3. Connect SeaweedFS Object Storage
	seaweedStorage, err := storage.Connect(cfg.SeaweedEndpoint, cfg.SeaweedAccessKey, cfg.SeaweedSecretKey, cfg.SeaweedBucket)
	if err != nil {
		log.Printf("Warning: SeaweedFS connection note: %v", err)
	}

	// 4. Initialize Auth & Handlers
	authMgr := auth.NewAuthManager(cfg.JWTSecret, redisCache)
	h := handlers.New(db, redisCache, seaweedStorage, authMgr)

	// 5. Router Setup
	r := chi.NewRouter()

	// Global Middlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// CORS Setup
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link", "X-Cache"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Routes
	r.Route("/api", func(api chi.Router) {
		// Public endpoints
		api.Get("/health", h.HealthCheck)
		api.Get("/storage/*", h.StreamStorageObject)
		api.Head("/storage/*", h.StreamStorageObject)

		// Auth (Public)
		api.Post("/auth/register", h.Register)
		api.Post("/auth/login", h.Login)

		// Protected endpoints (Requires JWT Bearer Token)
		api.Group(func(protected chi.Router) {
			protected.Use(authMgr.Middleware)

			// User
			protected.Get("/auth/me", h.Me)
			protected.Post("/auth/logout", h.Logout)

			// Personal Tracks
			protected.Get("/tracks", h.ListTracks)
			protected.Get("/tracks/{id}", h.GetTrack)
			protected.Post("/tracks", h.CreateTrack)
			protected.Delete("/tracks/{id}", h.DeleteTrack)

			// Audio Stream Extraction
			protected.Post("/tracks/extract", h.ExtractStream)
			protected.Get("/tracks/extract/{id}", h.GetExtractProgress)

			// Track Moodboard Items
			protected.Post("/tracks/{id}/items", h.AddMoodboardItem)
		})
	})

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Server run context for graceful shutdown
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig
		log.Println("Received termination signal, shutting down gracefully...")
		shutdownCtx, cancel := context.WithTimeout(serverCtx, 30*time.Second)
		defer cancel()

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("Graceful shutdown timed out.. forcing exit.")
			}
		}()

		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Fatalf("Server shutdown error: %v", err)
		}
		serverStopCtx()
	}()

	log.Printf("Song Moodboard Go API running on port %s", cfg.Port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("HTTP server error: %v", err)
	}

	<-serverCtx.Done()
	log.Println("Song Moodboard Go API terminated cleanly.")
}
