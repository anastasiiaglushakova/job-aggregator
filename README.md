# Job Aggregator

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

Concurrent job aggregator written in Go that fetches vacancies from public APIs (RemoteOK) and sends Telegram notifications about new opportunities. Built with idiomatic Go practices: explicit error handling, goroutines for concurrency, and clean internal architecture.

## ✨ Features

- **Concurrent fetching** — Goroutines with `sync.WaitGroup` for parallel data collection
- **Deduplication** — Unique URL constraint in PostgreSQL prevents duplicates
- **Telegram Bot API** — HTML-formatted notifications about new vacancies
- **Ethical scraping** — Proper User-Agent header, rate limiting respected
- **Idiomatic architecture** — Clean separation: `cmd/`, `internal/db`, `internal/fetcher`, `internal/notifier`
- **Production-ready** — Context timeouts, explicit error wrapping, structured logging
- **Tested** — 9 unit tests covering core logic with `httptest` mocks

## 🚀 Quick Start

### Prerequisites
- Go 1.22+
- PostgreSQL 15+
- Telegram Bot token (from [@BotFather](https://t.me/BotFather))

### Setup

```bash
# Clone repository
git clone https://github.com/anastasiiaglushakova/job-aggregator.git
cd job-aggregator

# Create PostgreSQL database
sudo -u postgres psql -c "CREATE DATABASE job_aggregator;"

# Apply migration
sudo -u postgres psql -d job_aggregator -f internal/db/migrations/001_init.sql

# Configure environment
cp .env.example .env
nano .env  # Fill in TELEGRAM_BOT_TOKEN and TELEGRAM_CHAT_ID

# Install dependencies
go mod tidy

# Run aggregator
go run cmd/aggregator/main.go
```

### Example Output

```
$ go run cmd/aggregator/main.go
2026/02/12 10:30:00 Starting concurrent aggregation with 2 fetchers...
2026/02/12 10:30:00 → Fetching from remoteok...
2026/02/12 10:30:00 → Fetching from mock...
2026/02/12 10:30:03 ✓ remoteok fetched 95 jobs in 2.8s
2026/02/12 10:30:05 ✓ mock fetched 2 jobs in 2.0s

Aggregation completed in 2.8s
Total new jobs: 97

[OK] Added new jobs: 97
  [NEW] Senior Full Stack Engineer — Acme Corp (Remote)
  [NEW] DevOps Specialist — TechFlow Inc (Europe)
  [NEW] Senior Go Developer — MockTech Inc (Remote)
  [NEW] Backend Engineer — StartupXYZ (USA)
  [NEW] Site Reliability Engineer — CloudCo (Remote)
... and 92 more

Sending notification to Telegram...
[OK] Telegram notification sent successfully!
```

### Running Tests

```bash
go test ./... -v
```

✅ 9 unit tests passing (fetchers, notifier, HTML escaping)  
✅ No external dependencies required (all APIs mocked with `httptest`)  
✅ 100% isolated tests (no real Telegram/RemoteOK calls)

## 📂 Project Structure

```
job-aggregator/
├── cmd/aggregator/            # Application entry point
├── internal/
│   ├── db/                    # PostgreSQL repository + migrations
│   ├── fetcher/               # Job fetchers (RemoteOK, mock)
│   └── notifier/              # Telegram notifications
├── .github/workflows/         # GitHub Actions (scheduled runs)
├── .env.example               # Environment variables template
├── .gitignore                 # Git ignore rules
├── go.mod                     # Go module definition
├── go.sum                     # Dependency checksums
├── LICENSE                    # MIT License
└── README.md                  # This file
```

## ⚖️ Ethical Guidelines

This project strictly follows ethical data collection practices:

✅ Only public APIs with explicit permission (RemoteOK)  
✅ Proper User-Agent header with repository link  
✅ Rate limiting respected (1s delay between requests)  
✅ No credentials or secrets committed (`.env` in `.gitignore`)  
❌ No scraping of commercial platforms without permission

## 🤖 Telegram Bot Setup

1. Talk to [@BotFather](https://t.me/BotFather) → `/newbot` → get token
2. Message your bot once
3. Get your `chat_id` via [@userinfobot](https://t.me/userinfobot)
4. Fill `.env`:
   ```env
   TELEGRAM_BOT_TOKEN="123456789:AAH_ABC..."
   TELEGRAM_CHAT_ID="987654321"
   ```

## 📜 License

MIT — see [LICENSE](LICENSE) for details.