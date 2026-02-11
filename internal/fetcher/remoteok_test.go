package fetcher

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRemoteOKFetcher_Fetch(t *testing.T) {
	// Mock RemoteOK API response with metadata + 2 jobs
	mockResponse := []map[string]interface{}{
		{"id": "", "position": "Remote OK"}, // Metadata row (will be skipped)
		{
			"id":          "1",
			"position":    "Senior Go Developer",
			"company":     "TechCorp",
			"url":         "https://remoteok.com/jobs/1",
			"tags":        []string{"go", "backend"},
			"location":    "Remote",
			"date":        "2026-02-12",
		},
		{
			"id":          "2",
			"position":    "Frontend Engineer",
			"company":     "WebInc",
			"url":         "https://remoteok.com/jobs/2",
			"tags":        []string{"react", "typescript"},
			"location":    "USA",
			"date":        "2026-02-11",
		},
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	fetcher := &RemoteOKFetcher{
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
		userAgent: "TestAgent/1.0",
		baseURL:   mockServer.URL, // Override for tests
	}

	ctx := context.Background()
	jobs, err := fetcher.Fetch(ctx)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(jobs) != 2 {
		t.Fatalf("expected 2 jobs (after skipping metadata), got %d", len(jobs))
	}

	// Verify first job
	if jobs[0].ID != "1" {
		t.Errorf("expected ID '1', got %q", jobs[0].ID)
	}
	if jobs[0].Position != "Senior Go Developer" {
		t.Errorf("expected Position 'Senior Go Developer', got %q", jobs[0].Position)
	}
	if jobs[0].Company != "TechCorp" {
		t.Errorf("expected Company 'TechCorp', got %q", jobs[0].Company)
	}

	// Verify second job
	if jobs[1].ID != "2" {
		t.Errorf("expected ID '2', got %q", jobs[1].ID)
	}
	if jobs[1].Company != "WebInc" {
		t.Errorf("expected Company 'WebInc', got %q", jobs[1].Company)
	}
}

func TestRemoteOKFetcher_Name(t *testing.T) {
	fetcher := NewRemoteOKFetcher("TestAgent/1.0")
	if fetcher.Name() != "remoteok" {
		t.Errorf("expected name 'remoteok', got %q", fetcher.Name())
	}
}

func TestRemoteOKFetcher_Fetch_HTTPError(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer mockServer.Close()

	fetcher := &RemoteOKFetcher{
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
		userAgent: "TestAgent/1.0",
		baseURL:   mockServer.URL,
	}

	ctx := context.Background()
	_, err := fetcher.Fetch(ctx)
	if err == nil {
		t.Fatal("expected error for 500 response, got nil")
	}
	if err.Error() != "remoteok returned status 500" {
		t.Errorf("expected error 'remoteok returned status 500', got %q", err.Error())
	}
}