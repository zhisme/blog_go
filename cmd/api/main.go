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
	defer func() {
		if closeErr := repo.Close(); closeErr != nil {
			log.Printf("Error closing database: %v", closeErr)
		}
	}()

	// Create and start server
	srv := api.NewApiServer(repo)
	err = srv.ListenAndServe(cfg.ServerAddr)
	if err != nil {
		log.Printf("Server error: %v", err)
	}
}
