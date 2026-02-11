package notifier

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestTelegram_SendNewJobs(t *testing.T) {
	var receivedText string

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		values, _ := url.ParseQuery(string(body))
		receivedText = values.Get("text")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"ok": true})
	}))
	defer mockServer.Close()

	telegram := &Telegram{
		botToken:   "test_token",
		chatID:     "123456789",
		client:     &http.Client{Timeout: 5 * time.Second},
		apiBaseURL: mockServer.URL,
	}

	jobs := []JobInfo{
		{Title: "Go Developer", Company: "TechCorp", Location: "Remote", URL: "https://example.com/job/1"},
		{Title: "Frontend Engineer", Company: "WebInc", Location: "USA", URL: "https://example.com/job/2"},
	}

	err := telegram.SendNewJobs(context.Background(), jobs)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !strings.Contains(receivedText, "Go Developer") {
		t.Errorf("expected 'Go Developer' in message, got %q", receivedText)
	}
	if !strings.Contains(receivedText, "TechCorp") {
		t.Errorf("expected 'TechCorp' in message, got %q", receivedText)
	}
	if !strings.Contains(receivedText, "https://example.com/job/1") {
		t.Errorf("expected URL in message, got %q", receivedText)
	}
}

func TestTelegram_SendNewJobs_Empty(t *testing.T) {
	telegram := NewTelegram("test_token", "123456789")
	err := telegram.SendNewJobs(context.Background(), []JobInfo{})
	if err != nil {
		t.Errorf("expected no error for empty jobs, got %v", err)
	}
}

func TestTelegram_SendNewJobs_Limit(t *testing.T) {
	var receivedText string

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		values, _ := url.ParseQuery(string(body))
		receivedText = values.Get("text")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"ok": true})
	}))
	defer mockServer.Close()

	telegram := &Telegram{
		botToken:   "test_token",
		chatID:     "123456789",
		client:     &http.Client{Timeout: 5 * time.Second},
		apiBaseURL: mockServer.URL,
	}

	jobs := make([]JobInfo, 10)
	for i := range jobs {
		jobs[i] = JobInfo{
			Title:    "Job " + string(rune('A'+i)),
			Company:  "Company",
			Location: "Remote",
			URL:      "https://example.com/job/" + string(rune('A'+i)),
		}
	}

	err := telegram.SendNewJobs(context.Background(), jobs)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !strings.Contains(receivedText, "Job A") {
		t.Errorf("expected first job in message")
	}
	if strings.Contains(receivedText, "Job F") {
		t.Errorf("expected max 5 jobs, but found Job F")
	}
	if !strings.Contains(receivedText, "... and 5 more new jobs") {
		t.Errorf("expected 'and 5 more' suffix, got %q", receivedText)
	}
}

func TestTelegram_escapeHTML(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Normal text", "Normal text"},
		{"A & B", "A &amp; B"},
		{"<script>", "&lt;script&gt;"},
		{"A < B > C", "A &lt; B &gt; C"},
	}

	for _, tt := range tests {
		result := escapeHTML(tt.input)
		if result != tt.expected {
			t.Errorf("escapeHTML(%q) = %q, expected %q", tt.input, result, tt.expected)
		}
	}
}