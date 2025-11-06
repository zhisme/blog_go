package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestCreateMailingList(t *testing.T) {
	// Use a test CSV file
	testCSVFile := "test_create_mailing_list.csv"
	defer func() { _ = os.Remove(testCSVFile) }()

	srv := NewApiServer()

	t.Run("Valid request returns 201 Created", func(t *testing.T) {
		payload := map[string]string{
			"email":    "valid@example.com",
			"username": "validuser",
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/mailing_list", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		srv.createMailingList(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("Expected status %d, got %d. Body: %s", http.StatusCreated, w.Code, w.Body.String())
		}

		contentType := w.Header().Get("Content-Type")
		if contentType != "application/json" {
			t.Errorf("Expected Content-Type 'application/json', got '%s'", contentType)
		}

		var response map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if response["email"] != payload["email"] {
			t.Errorf("Expected email %s, got %v", payload["email"], response["email"])
		}
	})

	t.Run("Empty body returns 400 Bad Request with specific message", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/mailing_list", bytes.NewBufferString(""))
		w := httptest.NewRecorder()

		srv.createMailingList(w, req)

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

	t.Run("Invalid JSON returns 400 Bad Request", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/mailing_list", bytes.NewBufferString("{invalid json}"))
		w := httptest.NewRecorder()

		srv.createMailingList(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}

		var response map[string]map[string]string
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to parse error response: %v", err)
		}

		if !strings.Contains(response["error"]["message"], "invalid JSON") {
			t.Errorf("Expected 'invalid JSON' in message, got %s", response["error"]["message"])
		}
	})

	t.Run("Missing email returns 400 Bad Request", func(t *testing.T) {
		payload := map[string]string{
			"username": "testuser",
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/mailing_list", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		srv.createMailingList(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}

		var response map[string]map[string]string
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to parse error response: %v", err)
		}

		if response["error"]["message"] != "email is required" {
			t.Errorf("Expected 'email is required' message, got %s", response["error"]["message"])
		}
	})

	t.Run("Invalid email format returns 400 Bad Request", func(t *testing.T) {
		payload := map[string]string{
			"email":    "notanemail",
			"username": "testuser",
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/mailing_list", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		srv.createMailingList(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}

		var response map[string]map[string]string
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to parse error response: %v", err)
		}

		if response["error"]["message"] != "invalid email format" {
			t.Errorf("Expected 'invalid email format' message, got %s", response["error"]["message"])
		}
	})

	t.Run("Missing username returns 400 Bad Request", func(t *testing.T) {
		payload := map[string]string{
			"email": "test@example.com",
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/mailing_list", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		srv.createMailingList(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}

		var response map[string]map[string]string
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to parse error response: %v", err)
		}

		if response["error"]["message"] != "username is required" {
			t.Errorf("Expected 'username is required' message, got %s", response["error"]["message"])
		}
	})

	t.Run("EOF on empty body is detected", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/mailing_list", http.NoBody)
		w := httptest.NewRecorder()

		srv.createMailingList(w, req)

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

	t.Run("Content-Type header is set for all responses", func(t *testing.T) {
		tests := []struct {
			body io.Reader
			name string
		}{
			{bytes.NewBufferString(`{"email":"test@example.com","username":"user"}`), "valid request"},
			{bytes.NewBufferString("invalid"), "invalid json"},
			{bytes.NewBufferString(""), "empty body"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				req := httptest.NewRequest(http.MethodPost, "/mailing_list", tt.body)
				w := httptest.NewRecorder()

				srv.createMailingList(w, req)

				contentType := w.Header().Get("Content-Type")
				if contentType != "application/json" {
					t.Errorf("Expected Content-Type 'application/json', got '%s'", contentType)
				}
			})
		}
	})
}
