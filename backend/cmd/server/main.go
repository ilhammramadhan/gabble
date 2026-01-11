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
	cfg := config.Load()

	if cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

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
		AllowOriginFunc: func(r *http.Request, origin string) bool {
			// Allow localhost for development
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			// Allow all Vercel preview URLs
			if strings.HasSuffix(origin, ".vercel.app") {
				return true
			}
			// Allow configured frontend URL
			if origin == cfg.FrontendURL {
				return true
			}
			return false
		},
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
