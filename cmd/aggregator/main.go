package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"

	"job-aggregator/internal/db"
	"job-aggregator/internal/fetcher"
	"job-aggregator/internal/notifier"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		log.Printf("Warning: .env not found or failed to load: %v", err)
	}

	// Connect to PostgreSQL database
	database, err := db.Connect("")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	repo := db.NewRepository(database)

	// Initialize fetchers with proper User-Agent
	userAgent := "JobAggregator/1.0 (+https://github.com/anastasiiaglushakova/job-aggregator)"
	fetchers := []fetcher.Fetcher{
		fetcher.NewRemoteOKFetcher(userAgent),
		fetcher.NewMockFetcher(2 * time.Second),
	}

	// Initialize Telegram bot if credentials are provided
	var tgBot *notifier.Telegram
	tgToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	tgChatID := os.Getenv("TELEGRAM_CHAT_ID")

	if tgToken != "" && tgChatID != "" {
		tgBot = notifier.NewTelegram(tgToken, tgChatID)
		log.Printf("Telegram bot initialized for chat %s", tgChatID)
	} else {
		log.Println("Telegram bot disabled (no token or chat_id)")
	}

	// Run concurrent aggregation with 90s timeout
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	newJobs, err := aggregateConcurrently(ctx, fetchers, repo)
	if err != nil {
		log.Fatalf("Aggregation failed: %v", err)
	}

	// Print results to console
	fmt.Printf("\n[OK] Added new jobs: %d\n", len(newJobs))
	for i, job := range newJobs {
		if i >= 5 {
			fmt.Printf("... and %d more\n", len(newJobs)-5)
			break
		}
		fmt.Printf("  [NEW] %s — %s (%s)\n", job.Title, job.Company, job.Location)
	}

	// Send Telegram notification if bot is configured and new jobs exist
	if tgBot != nil && len(newJobs) > 0 {
		log.Println("Sending notification to Telegram...")

		infos := make([]notifier.JobInfo, len(newJobs))
		for i, job := range newJobs {
			infos[i] = notifier.JobInfo{
				Title:    job.Title,
				Company:  job.Company,
				Location: job.Location,
				URL:      job.URL,
			}
		}

		if err := tgBot.SendNewJobs(ctx, infos); err != nil {
			log.Printf("Failed to send Telegram notification: %v", err)
		} else {
			log.Println("[OK] Telegram notification sent successfully!")
		}
	}
}

// aggregateConcurrently fetches jobs from multiple sources in parallel using goroutines
func aggregateConcurrently(
	ctx context.Context,
	fetchers []fetcher.Fetcher,
	repo *db.Repository,
) ([]db.Job, error) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var allNewJobs []db.Job
	var errs []error
	var errMu sync.Mutex

	start := time.Now()
	log.Printf("Starting concurrent aggregation with %d fetchers...", len(fetchers))

	for _, f := range fetchers {
		wg.Add(1)
		go func(fetcher fetcher.Fetcher) {
			defer wg.Done()

			sourceName := fetcher.Name()
			log.Printf("→ Fetching from %s...", sourceName)

			fetchStart := time.Now()
			jobs, err := fetcher.Fetch(ctx)
			fetchDuration := time.Since(fetchStart)

			if err != nil {
				errMu.Lock()
				errs = append(errs, fmt.Errorf("%s: %w", sourceName, err))
				errMu.Unlock()
				log.Printf("✗ %s failed in %v: %v", sourceName, fetchDuration, err)
				return
			}

			log.Printf("✓ %s fetched %d jobs in %v", sourceName, len(jobs), fetchDuration)

			// Convert to DB format
			var dbJobs []db.Job
			for _, rj := range jobs {
				publishedAt, _ := time.Parse("2006-01-02", rj.Date)
				if publishedAt.IsZero() {
					publishedAt = time.Now()
				}

				dbJobs = append(dbJobs, db.Job{
					ID:          fmt.Sprintf("%s:%s", sourceName, rj.ID),
					Source:      sourceName,
					Title:       rj.Position,
					Company:     rj.Company,
					URL:         rj.URL,
					Description: rj.Description,
					Tags:        rj.Tags,
					Location:    rj.Location,
					PublishedAt: publishedAt,
					FetchedAt:   time.Now(),
				})
			}

			// Save to DB with deduplication
			newJobs, err := repo.InsertJobs(ctx, dbJobs)
			if err != nil {
				errMu.Lock()
				errs = append(errs, fmt.Errorf("%s insert: %w", sourceName, err))
				errMu.Unlock()
				log.Printf("✗ %s insert failed: %v", sourceName, err)
				return
			}

			log.Printf("✓ %s inserted %d new jobs", sourceName, len(newJobs))

			// Safely append to shared slice
			mu.Lock()
			allNewJobs = append(allNewJobs, newJobs...)
			mu.Unlock()
		}(f)
	}

	wg.Wait()
	duration := time.Since(start)

	log.Printf("\nAggregation completed in %v", duration)
	log.Printf("Total new jobs: %d", len(allNewJobs))

	if len(errs) > 0 {
		for _, e := range errs {
			log.Printf("  Error: %v", e)
		}
		return allNewJobs, fmt.Errorf("aggregation had %d error(s)", len(errs))
	}

	return allNewJobs, nil
}