package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"backend-go/internal/api"
	"backend-go/internal/repositories"
)

func TestNewApiServer(t *testing.T) {
	repo, err := repositories.NewSqliteMailingListRepository(":memory:")
	if err != nil {
		t.Fatalf("Failed to create test repository: %v", err)
	}
	defer repo.Close()

	srv := api.NewApiServer(repo)

	if srv == nil {
		t.Fatal("NewApiServer() returned nil")
	}
}

func TestServerRouting(t *testing.T) {
	repo, err := repositories.NewSqliteMailingListRepository(":memory:")
	if err != nil {
		t.Fatalf("Failed to create test repository: %v", err)
	}
	defer repo.Close()

	srv := api.NewApiServer(repo)

	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
	}{
		{
			name:           "POST to /mailing_list should be handled",
			method:         http.MethodPost,
			path:           "/mailing_list",
			expectedStatus: http.StatusBadRequest, // Empty body
		},
		{
			name:           "GET to /mailing_list should return method not allowed",
			method:         http.MethodGet,
			path:           "/mailing_list",
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "POST to unknown path should return 404",
			method:         http.MethodPost,
			path:           "/unknown",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()

			srv.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestServerCORS(t *testing.T) {
	repo, err := repositories.NewSqliteMailingListRepository(":memory:")
	if err != nil {
		t.Fatalf("Failed to create test repository: %v", err)
	}
	defer repo.Close()

	srv := api.NewApiServer(repo)

	tests := []struct {
		name        string
		origin      string
		method      string
		shouldAllow bool
	}{
		{
			name:        "Should allow localhost:1313",
			origin:      "http://localhost:1313",
			method:      http.MethodPost,
			shouldAllow: true,
		},
		{
			name:        "Should allow zhisme.com",
			origin:      "https://zhisme.com/",
			method:      http.MethodPost,
			shouldAllow: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := bytes.NewBufferString(`{"email":"test@example.com","username":"testuser"}`)
			req := httptest.NewRequest(tt.method, "/mailing_list", body)
			req.Header.Set("Origin", tt.origin)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			srv.ServeHTTP(w, req)

			allowOrigin := w.Header().Get("Access-Control-Allow-Origin")
			if tt.shouldAllow && allowOrigin == "" {
				t.Errorf("Expected CORS headers for origin %s, but got none", tt.origin)
			}
		})
	}
}

func TestServerMailingListEndpoint(t *testing.T) {
	repo, err := repositories.NewSqliteMailingListRepository(":memory:")
	if err != nil {
		t.Fatalf("Failed to create test repository: %v", err)
	}
	defer repo.Close()

	srv := api.NewApiServer(repo)

	t.Run("Valid request creates mailing list entry", func(t *testing.T) {
		payload := map[string]string{
			"email":    "test@example.com",
			"username": "testuser",
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/mailing_list", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		srv.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("Expected status %d, got %d. Body: %s", http.StatusCreated, w.Code, w.Body.String())
		}

		var response map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if response["email"] != payload["email"] {
			t.Errorf("Expected email %s, got %v", payload["email"], response["email"])
		}
		if response["username"] != payload["username"] {
			t.Errorf("Expected username %s, got %v", payload["username"], response["username"])
		}
	})

	t.Run("Invalid JSON returns bad request", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/mailing_list", bytes.NewBufferString("invalid json"))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		srv.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("Empty body returns bad request", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/mailing_list", bytes.NewBufferString(""))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		srv.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}

		var response map[string]map[string]string
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to parse error response: %v", err)
		}

		if response["error"]["message"] != "request body is empty" {
			t.Errorf("Expected 'request body is empty' message, got %s", response["error"]["message"])
		}
	})
}
