package fetcher

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// RemoteOKFetcher fetches jobs from RemoteOK API
type RemoteOKFetcher struct {
	client    *http.Client
	userAgent string
	baseURL   string
}

func NewRemoteOKFetcher(userAgent string) *RemoteOKFetcher {
	return &RemoteOKFetcher{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		userAgent: userAgent,
		baseURL:   "https://remoteok.com/api",
	}
}

func (f *RemoteOKFetcher) Name() string {
	return "remoteok"
}

func (f *RemoteOKFetcher) Fetch(ctx context.Context) ([]Job, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", f.baseURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("User-Agent", f.userAgent)
	req.Header.Set("Accept", "application/json")

	resp, err := f.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("remoteok returned status %d", resp.StatusCode)
	}

	var jobs []Job
	if err := json.NewDecoder(resp.Body).Decode(&jobs); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	// Skip first element — RemoteOK metadata (not a job)
	if len(jobs) > 0 && jobs[0].ID == "" {
		jobs = jobs[1:]
	}

	// Respect rate limits: RemoteOK requests 1s delay between requests
	time.Sleep(time.Second)

	return jobs, nil
}