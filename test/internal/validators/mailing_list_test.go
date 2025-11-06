package validators_test

import (
	"backend-go/internal/dto"
	"backend-go/internal/validators"
	"strings"
	"testing"
)

func TestNewMailingListValidator(t *testing.T) {
	validator := validators.NewMailingListValidator()

	if validator == nil {
		t.Fatal("NewMailingListValidator() returned nil")
	}
}

func TestValidate(t *testing.T) {
	validator := validators.NewMailingListValidator()

	t.Run("Valid mailing list passes validation", func(t *testing.T) {
		ml := &dto.MailingList{
			Username: "testuser",
			Email:    "test@example.com",
		}

		err := validator.Validate(ml)
		if err != nil {
			t.Errorf("Expected no error for valid mailing list, got %v", err)
		}
	})

	t.Run("Invalid email fails validation", func(t *testing.T) {
		ml := &dto.MailingList{
			Username: "testuser",
			Email:    "notanemail",
		}

		err := validator.Validate(ml)
		if err == nil {
			t.Error("Expected error for invalid email, got nil")
		}
	})

	t.Run("Empty username fails validation", func(t *testing.T) {
		ml := &dto.MailingList{
			Username: "",
			Email:    "test@example.com",
		}

		err := validator.Validate(ml)
		if err == nil {
			t.Error("Expected error for empty username, got nil")
		}
	})

	t.Run("Empty email fails validation", func(t *testing.T) {
		ml := &dto.MailingList{
			Username: "testuser",
			Email:    "",
		}

		err := validator.Validate(ml)
		if err == nil {
			t.Error("Expected error for empty email, got nil")
		}
	})
}

// Note: TestValidateEmail and TestValidateUsername removed - they tested unexported methods
// Email and username validation is now tested through the public Validate() API in TestValidate

func TestValidatorEdgeCases(t *testing.T) {
	validator := validators.NewMailingListValidator()

	t.Run("Both empty username and email should fail on email first", func(t *testing.T) {
		ml := &dto.MailingList{
			Username: "",
			Email:    "",
		}

		err := validator.Validate(ml)
		if err == nil {
			t.Error("Expected error for both empty, got nil")
		}
		// Should fail on email validation first
		if !strings.Contains(err.Error(), "email is required") {
			t.Errorf("Expected 'email is required' error, got %v", err)
		}
	})
}
