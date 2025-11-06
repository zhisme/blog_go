package validators

import (
	"backend-go/internal/dto"
	"strings"
	"testing"
)

func TestNewMailingListValidator(t *testing.T) {
	validator := NewMailingListValidator()

	if validator == nil {
		t.Fatal("NewMailingListValidator() returned nil")
	}
}

func TestValidate(t *testing.T) {
	validator := NewMailingListValidator()

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

func TestValidateEmail(t *testing.T) {
	validator := NewMailingListValidator()

	tests := []struct {
		name        string
		email       string
		errorMsg    string
		shouldError bool
	}{
		{
			name:        "Valid simple email",
			email:       "user@example.com",
			shouldError: false,
		},
		{
			name:        "Valid email with subdomain",
			email:       "user@mail.example.com",
			shouldError: false,
		},
		{
			name:        "Valid email with plus sign",
			email:       "user+tag@example.com",
			shouldError: false,
		},
		{
			name:        "Valid email with dots",
			email:       "first.last@example.com",
			shouldError: false,
		},
		{
			name:        "Valid email with numbers",
			email:       "user123@example.com",
			shouldError: false,
		},
		{
			name:        "Valid email with hyphens in domain",
			email:       "user@my-domain.com",
			shouldError: false,
		},
		{
			name:        "Valid email with multiple TLD parts",
			email:       "user@example.co.uk",
			shouldError: false,
		},
		{
			name:        "Empty email",
			email:       "",
			shouldError: true,
			errorMsg:    "email is required",
		},
		{
			name:        "Missing @ symbol",
			email:       "userexample.com",
			shouldError: true,
			errorMsg:    "invalid email format",
		},
		{
			name:        "Missing domain",
			email:       "user@",
			shouldError: true,
			errorMsg:    "invalid email format",
		},
		{
			name:        "Missing local part",
			email:       "@example.com",
			shouldError: true,
			errorMsg:    "invalid email format",
		},
		{
			name:        "Missing TLD",
			email:       "user@example",
			shouldError: true,
			errorMsg:    "invalid email format",
		},
		{
			name:        "Space in email",
			email:       "user @example.com",
			shouldError: true,
			errorMsg:    "invalid email format",
		},
		{
			name:        "Multiple @ symbols",
			email:       "user@@example.com",
			shouldError: true,
			errorMsg:    "invalid email format",
		},
		{
			name:        "Invalid characters",
			email:       "user!#$%@example.com",
			shouldError: true,
			errorMsg:    "invalid email format",
		},
		{
			name:        "TLD too short",
			email:       "user@example.c",
			shouldError: true,
			errorMsg:    "invalid email format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.validateEmail(tt.email)

			if tt.shouldError {
				if err == nil {
					t.Errorf("Expected error for email '%s', got nil", tt.email)
				} else if tt.errorMsg != "" && err.Error() != tt.errorMsg {
					t.Errorf("Expected error message '%s', got '%s'", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error for email '%s', got %v", tt.email, err)
				}
			}
		})
	}
}

func TestValidateUsername(t *testing.T) {
	validator := NewMailingListValidator()

	tests := []struct {
		name        string
		username    string
		errorMsg    string
		shouldError bool
	}{
		{
			name:        "Valid username",
			username:    "testuser",
			shouldError: false,
		},
		{
			name:        "Valid username with numbers",
			username:    "user123",
			shouldError: false,
		},
		{
			name:        "Valid username with special characters",
			username:    "user_name-123",
			shouldError: false,
		},
		{
			name:        "Empty username",
			username:    "",
			shouldError: true,
			errorMsg:    "username is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.validateUsername(tt.username)

			if tt.shouldError {
				if err == nil {
					t.Errorf("Expected error for username '%s', got nil", tt.username)
				} else if tt.errorMsg != "" && err.Error() != tt.errorMsg {
					t.Errorf("Expected error message '%s', got '%s'", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error for username '%s', got %v", tt.username, err)
				}
			}
		})
	}
}

func TestValidatorEdgeCases(t *testing.T) {
	validator := NewMailingListValidator()

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
