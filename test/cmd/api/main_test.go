package main_test

import (
	"backend-go/internal/api"
	"net/http"
	"testing"
	"time"
)

func TestMainServerIntegration(t *testing.T) {
	srv := api.NewApiServer()
	if srv == nil {
		t.Fatal("NewApiServer() returned nil")
	}

	serverAddr := "localhost:3001" // Use different port to avoid conflicts
	errChan := make(chan error, 1)
	go func() {
		errChan <- srv.ListenAndServe(serverAddr)
	}()

	// Give server time to start
	time.Sleep(100 * time.Millisecond)

	resp, err := http.Get("http://" + serverAddr + "/mailing_list")
	if err != nil {
		t.Fatalf("Failed to make request to server: %v", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			t.Logf("Failed to close response body: %v", closeErr)
		}
	}()

	// Should return method not allowed (GET on POST endpoint) but server is running
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Expected status %d for GET on POST endpoint, got %d", http.StatusMethodNotAllowed, resp.StatusCode)
	}
}
