package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib" // Postgres Driver
	"github.com/joho/godotenv"         // Env loader

	// Your internal packages
	"github.com/Sheedy-T/huddle-backend/internal/core/services"
	"github.com/Sheedy-T/huddle-backend/internal/handlers"
	"github.com/Sheedy-T/huddle-backend/internal/repositories"
)

func main() {
	// 1. Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found")
	}

	logger := log.New(os.Stdout, "[HUDDLE-API] ", log.LstdFlags)

	// 2. Connect to Supabase
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		logger.Fatal("DATABASE_URL is not set")
	}

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		logger.Fatal("Cannot connect to DB: ", err)
	}
	logger.Println("✅ Connected to Supabase Postgres")

	// 3. Dependency Injection (Repo -> Service -> Handler)
	repo := repositories.NewPostgresRepository(db)
	service := services.NewHuddleService(repo)
	handler := handlers.NewHuddleHandler(service)

	// 4. Setup Routes
	mux := http.NewServeMux()
	
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})
	
	// Use POST method prefix (Go 1.22+)
	mux.HandleFunc("POST /huddle/start", handler.StartHuddle)
	mux.HandleFunc("POST /huddle/log", handler.LogActivity)
    mux.HandleFunc("/huddle/", handler.GetHuddleSummary)
	// 5. Start Server
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		logger.Println("🚀 Server starting on port :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal(err)
		}
	}()

	// 6. Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
	logger.Println("Server stopped")
};