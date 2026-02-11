package fetcher

import (
	"context"
	"time"
)

// MockFetcher simulates a slow data source for concurrency demonstration
type MockFetcher struct {
	delay time.Duration
}

func NewMockFetcher(delay time.Duration) *MockFetcher {
	return &MockFetcher{delay: delay}
}

func (f *MockFetcher) Name() string {
	return "mock"
}

func (f *MockFetcher) Fetch(ctx context.Context) ([]Job, error) {
	// Simulate network delay
	time.Sleep(f.delay)

	// Return mock job vacancies
	return []Job{
		{
			ID:       "mock:1",
			Position: "Senior Go Developer",
			Company:  "MockTech Inc",
			URL:      "https://mocktech.example.com/job/1",
			Tags:     []string{"go", "backend", "remote"},
			Location: "Remote",
			Date:     time.Now().Format("2006-01-02"),
		},
		{
			ID:       "mock:2",
			Position: "DevOps Engineer",
			Company:  "FakeOps Ltd",
			URL:      "https://fakeops.example.com/careers/2",
			Tags:     []string{"kubernetes", "aws", "docker"},
			Location: "Europe",
			Date:     time.Now().Format("2006-01-02"),
		},
	}, nil
}