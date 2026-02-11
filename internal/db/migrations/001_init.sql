-- Jobs table stores aggregated job vacancies from public APIs
-- Deduplication is enforced via UNIQUE constraint on url column
CREATE TABLE IF NOT EXISTS jobs (
    id TEXT PRIMARY KEY,                    -- Composite ID: "source:remote_id"
    source TEXT NOT NULL,                   -- Source name: "remoteok", "mock", etc.
    title TEXT NOT NULL,                    -- Job title
    company TEXT NOT NULL,                  -- Company name
    url TEXT UNIQUE NOT NULL,               -- Unique job URL (deduplication key)
    description TEXT,                       -- Full job description
    tags JSONB DEFAULT '[]',                -- Skills/tags as JSON array
    location TEXT,                          -- Job location (city/country/remote)
    published_at TIMESTAMP NOT NULL,        -- Original publication timestamp
    fetched_at TIMESTAMP NOT NULL DEFAULT NOW()  -- When job was fetched by aggregator
);

-- Indexes for query performance
CREATE INDEX IF NOT EXISTS idx_jobs_url ON jobs(url);
CREATE INDEX IF NOT EXISTS idx_jobs_published ON jobs(published_at DESC);