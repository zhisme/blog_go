package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run cmd/dbinspect/main.go <database-path>")
		os.Exit(1)
	}

	dbPath := os.Args[1]

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			log.Printf("Error closing database: %v", closeErr)
		}
	}()

	// Query all records
	rows, err := db.Query("SELECT id, username, email, created_at FROM mailing_list ORDER BY created_at DESC")
	if err != nil {
		log.Fatalf("Failed to query database: %v", err)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Printf("Error closing rows: %v", closeErr)
		}
	}()

	fmt.Println("\n=== Mailing List Records ===")
	fmt.Printf("%-5s %-20s %-30s %-25s\n", "ID", "Username", "Email", "Created At")
	fmt.Println("--------------------------------------------------------------------------------")

	count := 0
	for rows.Next() {
		var id int
		var username, email, createdAt string

		err := rows.Scan(&id, &username, &email, &createdAt)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}

		fmt.Printf("%-5d %-20s %-30s %-25s\n", id, username, email, createdAt)
		count++
	}

	if err = rows.Err(); err != nil {
		log.Fatalf("Error iterating rows: %v", err)
	}

	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Printf("Total records: %d\n\n", count)
}
