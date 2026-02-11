package fetcher

import "context"

// Job represents a job vacancy from external API
type Job struct {
	ID          string   `json:"id"`
	Position    string   `json:"position"`
	Company     string   `json:"company"`
	URL         string   `json:"url"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
	Location    string   `json:"location"`
	Date        string   `json:"date"` // Format: "2006-01-02"
}

// Fetcher interface for job data sources
type Fetcher interface {
	Name() string
	Fetch(ctx context.Context) ([]Job, error)
}