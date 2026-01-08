package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/ilhammramadhan/gabble/internal/config"
	"github.com/ilhammramadhan/gabble/internal/database"
	"github.com/ilhammramadhan/gabble/internal/handlers"
	"github.com/ilhammramadhan/gabble/internal/middleware"
	"github.com/ilhammramadhan/gabble/internal/websocket"
)

func main() {
	// Debug: print all env vars that start with common prefixes
	log.Println("=== Environment Variables Debug ===")
	for _, env := range os.Environ() {
		// Only log non-sensitive variable names
		if len(env) > 0 {
			parts := strings.SplitN(env, "=", 2)
			if len(parts) == 2 {
				log.Printf("ENV: %s (len=%d)", parts[0], len(parts[1]))
			}
		}
	}
	log.Println("=== End Environment Variables ===")

	cfg := config.Load()

	// Debug: log if DATABASE_URL is set
	if cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}
	log.Printf("Connecting to database (URL length: %d)", len(cfg.DatabaseURL))

	db, err := database.New(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Migrate(context.Background()); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	hub := websocket.NewHub(db)
	go hub.Run()

	authHandler := handlers.NewAuthHandler(db, cfg)
	roomHandler := handlers.NewRoomHandler(db)
	wsHandler := handlers.NewWebSocketHandler(hub, db, cfg)

	r := chi.NewRouter()

	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Timeout(60 * time.Second))

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{cfg.FrontendURL, "http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	r.Get("/auth/github", authHandler.GithubLogin)
	r.Get("/auth/github/callback", authHandler.GithubCallback)

	r.Route("/api", func(r chi.Router) {
		r.Get("/rooms", roomHandler.GetRooms)

		r.Group(func(r chi.Router) {
			r.Use(middleware.AuthMiddleware(db, cfg.JWTSecret))

			r.Get("/auth/me", authHandler.GetCurrentUser)

			r.Post("/rooms", roomHandler.CreateRoom)
			r.Get("/rooms/{id}", roomHandler.GetRoom)
			r.Delete("/rooms/{id}", roomHandler.DeleteRoom)
			r.Get("/rooms/{id}/messages", roomHandler.GetMessages)
		})
	})

	r.Get("/ws", wsHandler.HandleWebSocket)

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Server starting on port %s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
