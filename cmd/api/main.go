package main

import (
	"backend-go/internal/api"
	"backend-go/internal/config"
	"backend-go/internal/repositories"
	"log"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize SQLite repository
	repo, err := repositories.NewSqliteMailingListRepository(cfg.DatabasePath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer repo.Close()

	// Create and start server
	srv := api.NewApiServer(repo)
	if err := srv.ListenAndServe("localhost:3000"); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
