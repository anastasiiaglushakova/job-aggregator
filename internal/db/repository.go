package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

// Job represents a job vacancy
type Job struct {
	ID          string    `json:"id"`
	Source      string    `json:"source"`
	Title       string    `json:"title"`
	Company     string    `json:"company"`
	URL         string    `json:"url"`
	Description string    `json:"description"`
	Tags        []string  `json:"tags"`
	Location    string    `json:"location"`
	PublishedAt time.Time `json:"published_at"`
	FetchedAt   time.Time `json:"fetched_at"`
}

// Repository encapsulates database access
type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// InsertJobs inserts jobs, returning only new ones (deduplication by URL)
func (r *Repository) InsertJobs(ctx context.Context, jobs []Job) ([]Job, error) {
	var newJobs []Job

	for _, job := range jobs {
		// Check existence by unique URL
		var exists bool
		err := r.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM jobs WHERE url = $1)", job.URL).
			Scan(&exists)
		if err != nil {
			return nil, fmt.Errorf("check existence for %s: %w", job.URL, err)
		}
		if exists {
			continue
		}

		// Serialize tags to JSON
		tagsJSON, err := json.Marshal(job.Tags)
		if err != nil {
			return nil, fmt.Errorf("marshal tags: %w", err)
		}

		// Insert new job
		_, err = r.db.ExecContext(ctx,
			`INSERT INTO jobs (
				id, source, title, company, url, description, tags, location, published_at, fetched_at
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
			ON CONFLICT (url) DO NOTHING`,
			job.ID, job.Source, job.Title, job.Company, job.URL,
			job.Description, tagsJSON, job.Location, job.PublishedAt, job.FetchedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("insert job %s: %w", job.ID, err)
		}

		newJobs = append(newJobs, job)
	}

	return newJobs, nil
}