package tests

import (
	"net/http"
	"testing"
)

func TestHomePageAvailable(t *testing.T) {
	resp, err := http.Get("https://srefinal.onrender.com/")
	if err != nil {
		t.Fatalf("Failed to GET /: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %d", resp.StatusCode)
	}
}
