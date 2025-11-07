package main

import (
	"backend-go/internal/dto"
	"backend-go/internal/repositories"
	"encoding/csv"
	"flag"
	"log"
	"os"
	"time"
)

func main() {
	csvPath := flag.String("csv", "mailing_list.csv", "Path to CSV file")
	dbPath := flag.String("db", "blog.db", "Path to SQLite database")
	flag.Parse()

	log.Printf("Starting migration from %s to %s", *csvPath, *dbPath)

	// Check if CSV file exists
	if _, err := os.Stat(*csvPath); os.IsNotExist(err) {
		log.Fatalf("CSV file does not exist: %s", *csvPath)
	}

	// Open CSV file
	csvFile, err := os.Open(*csvPath)
	if err != nil {
		log.Fatalf("Failed to open CSV file: %v", err)
	}
	defer func() {
		if closeErr := csvFile.Close(); closeErr != nil {
			log.Printf("Error closing CSV file: %v", closeErr)
		}
	}()

	// Read CSV data
	reader := csv.NewReader(csvFile)
	records, err := reader.ReadAll()
	if err != nil {
		log.Printf("Failed to read CSV file: %v", err)
		return
	}

	if len(records) == 0 {
		log.Println("CSV file is empty, nothing to migrate")
		return
	}

	// Initialize SQLite repository
	repo, err := repositories.NewSqliteMailingListRepository(*dbPath)
	if err != nil {
		log.Printf("Failed to initialize SQLite repository: %v", err)
		return
	}
	defer func() {
		if closeErr := repo.Close(); closeErr != nil {
			log.Printf("Error closing database: %v", closeErr)
		}
	}()

	// Track statistics
	imported := 0
	skipped := 0
	errors := 0

	// Skip header row and process data
	for i, record := range records {
		if i == 0 {
			// Validate header
			if len(record) < 3 {
				log.Printf("Invalid CSV header: expected at least 3 columns, got %d", len(record))
				return
			}
			log.Printf("CSV Header: %v", record)
			continue
		}

		if len(record) < 3 {
			log.Printf("Skipping invalid record at line %d: insufficient columns", i+1)
			skipped++
			continue
		}

		username := record[0]
		email := record[1]
		createdAtStr := record[2]

		// Parse timestamp
		createdAt, parseErr := time.Parse(time.RFC3339, createdAtStr)
		if parseErr != nil {
			log.Printf("Warning: Failed to parse timestamp '%s' at line %d, using current time: %v", createdAtStr, i+1, parseErr)
			createdAt = time.Now()
		}

		// Import into SQLite
		mailingListEntry := &dto.MailingList{
			Username:  username,
			Email:     email,
			CreatedAt: createdAt,
		}

		// Use the repository to save (handles duplicates gracefully)
		err = repo.Save(mailingListEntry)

		if err != nil {
			log.Printf("Error importing record at line %d (%s): %v", i+1, email, err)
			errors++
		} else {
			imported++
		}
	}

	log.Printf("\n=== Migration Complete ===")
	log.Printf("Total records processed: %d", len(records)-1) // -1 for header
	log.Printf("Successfully imported: %d", imported)
	log.Printf("Skipped (duplicates): %d", skipped)
	log.Printf("Errors: %d", errors)

	if errors == 0 {
		log.Printf("\nMigration completed successfully!")
		log.Printf("Database saved to: %s", *dbPath)
		log.Printf("\nYou can now update your application to use the SQLite database.")
	} else {
		log.Printf("\nMigration completed with errors. Please review the log above.")
	}
}
