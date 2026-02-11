package fetcher

import (
	"context"
	"testing"
	"time"
)

func TestMockFetcher_Fetch(t *testing.T) {
	fetcher := NewMockFetcher(100 * time.Millisecond)
	ctx := context.Background()

	jobs, err := fetcher.Fetch(ctx)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(jobs) != 2 {
		t.Fatalf("expected 2 jobs, got %d", len(jobs))
	}

	// Verify first job
	if jobs[0].ID != "mock:1" {
		t.Errorf("expected ID 'mock:1', got %q", jobs[0].ID)
	}
	if jobs[0].Position != "Senior Go Developer" {
		t.Errorf("expected Position 'Senior Go Developer', got %q", jobs[0].Position)
	}
	if jobs[0].Company != "MockTech Inc" {
		t.Errorf("expected Company 'MockTech Inc', got %q", jobs[0].Company)
	}
	if jobs[0].URL != "https://mocktech.example.com/job/1" {
		t.Errorf("expected URL 'https://mocktech.example.com/job/1', got %q", jobs[0].URL)
	}

	// Verify second job
	if jobs[1].ID != "mock:2" {
		t.Errorf("expected ID 'mock:2', got %q", jobs[1].ID)
	}
	if jobs[1].Position != "DevOps Engineer" {
		t.Errorf("expected Position 'DevOps Engineer', got %q", jobs[1].Position)
	}

	// Verify date format
	_, err = time.Parse("2006-01-02", jobs[0].Date)
	if err != nil {
		t.Errorf("expected valid date format '2006-01-02', got %q", jobs[0].Date)
	}
}

func TestMockFetcher_Name(t *testing.T) {
	fetcher := NewMockFetcher(0)
	if fetcher.Name() != "mock" {
		t.Errorf("expected name 'mock', got %q", fetcher.Name())
	}
}