package main

import (
	"backend-go/internal/api"
	"net/http"
	"testing"
	"time"
)

func TestMainServerIntegration(t *testing.T) {
	// Create server
	srv := api.NewApiServer()
	if srv == nil {
		t.Fatal("NewApiServer() returned nil")
	}

	// Start server in goroutine
	serverAddr := "localhost:3001" // Use different port to avoid conflicts
	errChan := make(chan error, 1)
	go func() {
		errChan <- srv.ListenAndServe(serverAddr)
	}()

	// Give server time to start
	time.Sleep(100 * time.Millisecond)

	// Test that server is responding
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

	// Note: In a real integration test, you'd properly shut down the server
	// For this test, the server will be cleaned up when the test process exits
}
