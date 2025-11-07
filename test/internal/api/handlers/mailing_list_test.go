package handlers_test

import (
	"backend-go/internal/api/handlers"
	"backend-go/internal/dto"
	"backend-go/internal/repositories"
	"strings"
	"testing"
	"time"
)

func TestHandleCreate(t *testing.T) {
	// Create in-memory SQLite database for testing
	repo, err := repositories.NewSqliteMailingListRepository(":memory:")
	if err != nil {
		t.Fatalf("Failed to create test repository: %v", err)
	}
	defer func() {
		if closeErr := repo.Close(); closeErr != nil {
			t.Errorf("Failed to close repository: %v", closeErr)
		}
	}()

	t.Run("Valid mailing list entry is created successfully", func(t *testing.T) {
		input := dto.MailingList{
			Username: "testuser",
			Email:    "test@example.com",
		}

		result, err := handlers.HandleCreate(input, repo)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if result.Email != input.Email {
			t.Errorf("Expected email %s, got %s", input.Email, result.Email)
		}
		if result.Username != input.Username {
			t.Errorf("Expected username %s, got %s", input.Username, result.Username)
		}
		if result.CreatedAt.IsZero() {
			t.Error("Expected CreatedAt to be set, but it was zero")
		}
		if time.Since(result.CreatedAt) > time.Second {
			t.Error("Expected CreatedAt to be recent")
		}
	})

	t.Run("Invalid email returns validation error", func(t *testing.T) {
		input := dto.MailingList{
			Username: "testuser",
			Email:    "notanemail",
		}

		_, err := handlers.HandleCreate(input, repo)
		if err == nil {
			t.Fatal("Expected validation error, got nil")
		}

		if !strings.Contains(err.Error(), "invalid email format") {
			t.Errorf("Expected 'invalid email format' error, got %v", err)
		}
	})

	t.Run("Missing email returns validation error", func(t *testing.T) {
		input := dto.MailingList{
			Username: "testuser",
			Email:    "",
		}

		_, err := handlers.HandleCreate(input, repo)
		if err == nil {
			t.Fatal("Expected validation error, got nil")
		}

		if !strings.Contains(err.Error(), "email is required") {
			t.Errorf("Expected 'email is required' error, got %v", err)
		}
	})

	t.Run("Missing username returns validation error", func(t *testing.T) {
		input := dto.MailingList{
			Username: "",
			Email:    "test@example.com",
		}

		_, err := handlers.HandleCreate(input, repo)
		if err == nil {
			t.Fatal("Expected validation error, got nil")
		}

		if !strings.Contains(err.Error(), "username is required") {
			t.Errorf("Expected 'username is required' error, got %v", err)
		}
	})

	t.Run("Duplicate emails are handled gracefully", func(t *testing.T) {
		input := dto.MailingList{
			Username: "user1",
			Email:    "duplicate@example.com",
		}

		_, err := handlers.HandleCreate(input, repo)
		if err != nil {
			t.Fatalf("Expected no error on first create, got %v", err)
		}

		// Try to create duplicate
		input2 := dto.MailingList{
			Username: "user2",
			Email:    "duplicate@example.com",
		}

		result, err := handlers.HandleCreate(input2, repo)
		if err != nil {
			t.Fatalf("Expected no error on duplicate (should be handled gracefully), got %v", err)
		}

		// Result should still be returned even if it's a duplicate
		if result.Email != input2.Email {
			t.Errorf("Expected email %s, got %s", input2.Email, result.Email)
		}
	})

	t.Run("CreatedAt is set automatically if not provided", func(t *testing.T) {
		input := dto.MailingList{
			Username: "timetest",
			Email:    "timetest@example.com",
		}

		before := time.Now()
		result, err := handlers.HandleCreate(input, repo)
		after := time.Now()

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if result.CreatedAt.Before(before) || result.CreatedAt.After(after) {
			t.Errorf("CreatedAt %v is not between %v and %v", result.CreatedAt, before, after)
		}
	})

	t.Run("Various valid email formats are accepted", func(t *testing.T) {
		validEmails := []string{
			"simple@example.com",
			"user.name@example.com",
			"user+tag@example.co.uk",
			"123@example.com",
			"user_name@example-domain.com",
		}

		for i, email := range validEmails {
			t.Run(email, func(t *testing.T) {
				input := dto.MailingList{
					Username: "testuser" + string(rune(i)),
					Email:    email,
				}

				result, err := handlers.HandleCreate(input, repo)
				if err != nil {
					t.Errorf("Expected valid email %s to be accepted, got error: %v", email, err)
				}

				if result.Email != email {
					t.Errorf("Expected email %s, got %s", email, result.Email)
				}
			})
		}
	})

	t.Run("Various invalid email formats are rejected", func(t *testing.T) {
		invalidEmails := []string{
			"notanemail",
			"@example.com",
			"user@",
			"user @example.com",
			"user@example",
			"",
		}

		for _, email := range invalidEmails {
			t.Run(email, func(t *testing.T) {
				input := dto.MailingList{
					Username: "testuser",
					Email:    email,
				}

				_, err := handlers.HandleCreate(input, repo)
				if err == nil {
					t.Errorf("Expected invalid email %s to be rejected, but it was accepted", email)
				}
			})
		}
	})
}
