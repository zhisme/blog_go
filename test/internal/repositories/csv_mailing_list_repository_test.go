package repositories_test

import (
	"backend-go/internal/dto"
	"backend-go/internal/repositories"
	"encoding/csv"
	"os"
	"testing"
	"time"
)

func TestNewCsvMailingListRepository(t *testing.T) {
	repo := repositories.NewCsvMailingListRepository("test.csv")

	if repo == nil {
		t.Fatal("NewCsvMailingListRepository() returned nil")
	}

	// Note: Cannot test unexported fields in black-box testing
}

func TestSave(t *testing.T) {
	testFile := "test_save.csv"
	defer func() { _ = os.Remove(testFile) }()

	repo := repositories.NewCsvMailingListRepository(testFile)

	t.Run("Save creates file with headers if it doesn't exist", func(t *testing.T) {
		_ = os.Remove(testFile) // Ensure clean state

		ml := &dto.MailingList{
			Username:  "testuser",
			Email:     "test@example.com",
			CreatedAt: time.Now(),
		}

		err := repo.Save(ml)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		// Verify file exists
		if _, statErr := os.Stat(testFile); os.IsNotExist(statErr) {
			t.Error("Expected file to be created, but it doesn't exist")
		}

		// Verify file contents
		file, err := os.Open(testFile)
		if err != nil {
			t.Fatalf("Failed to open test file: %v", err)
		}
		defer func() {
			if closeErr := file.Close(); closeErr != nil {
				t.Errorf("Failed to close file: %v", closeErr)
			}
		}()

		reader := csv.NewReader(file)
		records, err := reader.ReadAll()
		if err != nil {
			t.Fatalf("Failed to read CSV: %v", err)
		}

		if len(records) != 2 { // Header + 1 record
			t.Fatalf("Expected 2 records (header + data), got %d", len(records))
		}

		// Check headers
		expectedHeaders := []string{"Username", "Email", "CreatedAt"}
		for i, header := range records[0] {
			if header != expectedHeaders[i] {
				t.Errorf("Expected header %s, got %s", expectedHeaders[i], header)
			}
		}

		// Check data
		if records[1][0] != ml.Username {
			t.Errorf("Expected username %s, got %s", ml.Username, records[1][0])
		}
		if records[1][1] != ml.Email {
			t.Errorf("Expected email %s, got %s", ml.Email, records[1][1])
		}
	})

	t.Run("Save appends to existing file without duplicating headers", func(t *testing.T) {
		_ = os.Remove(testFile) // Clean state

		ml1 := &dto.MailingList{
			Username:  "user1",
			Email:     "user1@example.com",
			CreatedAt: time.Now(),
		}

		ml2 := &dto.MailingList{
			Username:  "user2",
			Email:     "user2@example.com",
			CreatedAt: time.Now(),
		}

		// Save first entry
		if err := repo.Save(ml1); err != nil {
			t.Fatalf("Failed to save first entry: %v", err)
		}

		// Save second entry
		if err := repo.Save(ml2); err != nil {
			t.Fatalf("Failed to save second entry: %v", err)
		}

		// Verify file contents
		file, err := os.Open(testFile)
		if err != nil {
			t.Fatalf("Failed to open test file: %v", err)
		}
		defer func() {
			if closeErr := file.Close(); closeErr != nil {
				t.Errorf("Failed to close file: %v", closeErr)
			}
		}()

		reader := csv.NewReader(file)
		records, err := reader.ReadAll()
		if err != nil {
			t.Fatalf("Failed to read CSV: %v", err)
		}

		if len(records) != 3 { // Header + 2 records
			t.Fatalf("Expected 3 records (header + 2 data), got %d", len(records))
		}

		// Verify both entries are present
		if records[1][1] != ml1.Email {
			t.Errorf("Expected first email %s, got %s", ml1.Email, records[1][1])
		}
		if records[2][1] != ml2.Email {
			t.Errorf("Expected second email %s, got %s", ml2.Email, records[2][1])
		}
	})

	t.Run("Save sets CreatedAt if not provided", func(t *testing.T) {
		_ = os.Remove(testFile)

		ml := &dto.MailingList{
			Username: "testuser",
			Email:    "test@example.com",
			// CreatedAt is zero value
		}

		before := time.Now()
		err := repo.Save(ml)
		after := time.Now()

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		// Read back the file
		file, err := os.Open(testFile)
		if err != nil {
			t.Fatalf("Failed to open test file: %v", err)
		}
		defer func() {
			if closeErr := file.Close(); closeErr != nil {
				t.Errorf("Failed to close file: %v", closeErr)
			}
		}()

		reader := csv.NewReader(file)
		records, err := reader.ReadAll()
		if err != nil {
			t.Fatalf("Failed to read CSV: %v", err)
		}

		// Parse the timestamp
		timestamp, err := time.Parse(time.RFC3339, records[1][2])
		if err != nil {
			t.Fatalf("Failed to parse timestamp: %v", err)
		}

		// Allow for some precision loss due to RFC3339 format (no nanoseconds)
		if timestamp.Before(before.Add(-time.Second)) || timestamp.After(after.Add(time.Second)) {
			t.Errorf("Timestamp %v is not between %v and %v", timestamp, before, after)
		}
	})

	t.Run("Save does not duplicate emails", func(t *testing.T) {
		_ = os.Remove(testFile)

		ml := &dto.MailingList{
			Username:  "user1",
			Email:     "duplicate@example.com",
			CreatedAt: time.Now(),
		}

		// Save first time
		err := repo.Save(ml)
		if err != nil {
			t.Fatalf("Expected no error on first save, got %v", err)
		}

		// Try to save duplicate
		ml2 := &dto.MailingList{
			Username:  "user2",
			Email:     "duplicate@example.com",
			CreatedAt: time.Now(),
		}

		err = repo.Save(ml2)
		if err != nil {
			t.Fatalf("Expected no error on duplicate save (should be silently handled), got %v", err)
		}

		// Verify file has only one entry
		file, err := os.Open(testFile)
		if err != nil {
			t.Fatalf("Failed to open test file: %v", err)
		}
		defer func() {
			if closeErr := file.Close(); closeErr != nil {
				t.Errorf("Failed to close file: %v", closeErr)
			}
		}()

		reader := csv.NewReader(file)
		records, err := reader.ReadAll()
		if err != nil {
			t.Fatalf("Failed to read CSV: %v", err)
		}

		if len(records) != 2 { // Header + 1 record (duplicate not saved)
			t.Errorf("Expected 2 records (header + 1 data), got %d", len(records))
		}
	})

	t.Run("Save preserves CreatedAt if provided", func(t *testing.T) {
		_ = os.Remove(testFile)

		specificTime := time.Date(2023, 5, 15, 10, 30, 0, 0, time.UTC)
		ml := &dto.MailingList{
			Username:  "testuser",
			Email:     "test@example.com",
			CreatedAt: specificTime,
		}

		err := repo.Save(ml)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		// Read back the file
		file, err := os.Open(testFile)
		if err != nil {
			t.Fatalf("Failed to open test file: %v", err)
		}
		defer func() {
			if closeErr := file.Close(); closeErr != nil {
				t.Errorf("Failed to close file: %v", closeErr)
			}
		}()

		reader := csv.NewReader(file)
		records, err := reader.ReadAll()
		if err != nil {
			t.Fatalf("Failed to read CSV: %v", err)
		}

		// Parse the timestamp
		timestamp, err := time.Parse(time.RFC3339, records[1][2])
		if err != nil {
			t.Fatalf("Failed to parse timestamp: %v", err)
		}

		if !timestamp.Equal(specificTime) {
			t.Errorf("Expected timestamp %v, got %v", specificTime, timestamp)
		}
	})
}

// Note: TestEmailExists removed - it tested unexported emailExists() function
// Duplicate detection is now tested through the public Save() API in TestSave
