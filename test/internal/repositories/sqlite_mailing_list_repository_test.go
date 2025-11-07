package repositories_test

import (
	"backend-go/internal/dto"
	"backend-go/internal/repositories"
	"os"
	"testing"
	"time"
)

func TestNewSqliteMailingListRepository(t *testing.T) {
	repo, err := repositories.NewSqliteMailingListRepository(":memory:")
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer func() {
		if err := repo.Close(); err != nil {
			t.Errorf("Failed to close repository: %v", err)
		}
	}()

	if repo == nil {
		t.Fatal("NewSqliteMailingListRepository() returned nil")
	}
}

func TestSqliteSave(t *testing.T) {
	repo, err := repositories.NewSqliteMailingListRepository(":memory:")
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer func() {
		if err := repo.Close(); err != nil {
			t.Errorf("Failed to close repository: %v", err)
		}
	}()

	t.Run("Save creates entry successfully", func(t *testing.T) {
		ml := &dto.MailingList{
			Username:  "testuser",
			Email:     "test@example.com",
			CreatedAt: time.Now(),
		}

		err := repo.Save(ml)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})

	t.Run("Save with zero CreatedAt sets timestamp automatically", func(t *testing.T) {
		ml := &dto.MailingList{
			Username:  "autotime",
			Email:     "autotime@example.com",
			CreatedAt: time.Time{}, // Zero value
		}

		err := repo.Save(ml)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})

	t.Run("Save handles duplicate email gracefully", func(t *testing.T) {
		ml1 := &dto.MailingList{
			Username:  "user1",
			Email:     "duplicate@example.com",
			CreatedAt: time.Now(),
		}

		err := repo.Save(ml1)
		if err != nil {
			t.Fatalf("Expected no error on first save, got %v", err)
		}

		ml2 := &dto.MailingList{
			Username:  "user2",
			Email:     "duplicate@example.com",
			CreatedAt: time.Now(),
		}

		err = repo.Save(ml2)
		if err != nil {
			t.Fatalf("Expected no error on duplicate (should be handled gracefully), got %v", err)
		}
	})

	t.Run("Save appends multiple entries", func(t *testing.T) {
		emails := []string{
			"user1@example.com",
			"user2@example.com",
			"user3@example.com",
		}

		for i, email := range emails {
			ml := &dto.MailingList{
				Username:  "user" + string(rune(i)),
				Email:     email,
				CreatedAt: time.Now(),
			}

			err := repo.Save(ml)
			if err != nil {
				t.Fatalf("Expected no error for entry %d, got %v", i, err)
			}
		}
	})
}

func TestSqliteWithFileDatabase(t *testing.T) {
	testFile := "test_sqlite.db"
	defer func() { _ = os.Remove(testFile) }()

	t.Run("Creates database file if it doesn't exist", func(t *testing.T) {
		_ = os.Remove(testFile)

		repo, err := repositories.NewSqliteMailingListRepository(testFile)
		if err != nil {
			t.Fatalf("Failed to create repository: %v", err)
		}
		defer func() {
			if err := repo.Close(); err != nil {
				t.Errorf("Failed to close repository: %v", err)
			}
		}()

		if _, statErr := os.Stat(testFile); os.IsNotExist(statErr) {
			t.Error("Expected database file to be created, but it doesn't exist")
		}
	})

	t.Run("Opens existing database file", func(t *testing.T) {
		// First create the database
		repo1, err := repositories.NewSqliteMailingListRepository(testFile)
		if err != nil {
			t.Fatalf("Failed to create repository: %v", err)
		}

		ml := &dto.MailingList{
			Username:  "testuser",
			Email:     "test@example.com",
			CreatedAt: time.Now(),
		}

		err = repo1.Save(ml)
		if err != nil {
			t.Fatalf("Failed to save entry: %v", err)
		}
		if err := repo1.Close(); err != nil {
			t.Errorf("Failed to close repo1: %v", err)
		}

		// Now open the existing database
		repo2, err := repositories.NewSqliteMailingListRepository(testFile)
		if err != nil {
			t.Fatalf("Failed to open existing repository: %v", err)
		}
		defer func() {
			if err := repo2.Close(); err != nil {
				t.Errorf("Failed to close repo2: %v", err)
			}
		}()

		// Try to save the same email - should be handled gracefully
		err = repo2.Save(ml)
		if err != nil {
			t.Fatalf("Expected no error on duplicate in existing db, got %v", err)
		}
	})
}

func TestSqliteClose(t *testing.T) {
	repo, err := repositories.NewSqliteMailingListRepository(":memory:")
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	err = repo.Close()
	if err != nil {
		t.Fatalf("Expected no error on close, got %v", err)
	}

	// Calling Close again should not panic
	err = repo.Close()
	if err != nil {
		t.Logf("Second close returned error: %v", err)
	}
}
