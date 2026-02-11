package notifier

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Telegram sends notifications to Telegram chat
type Telegram struct {
	botToken   string
	chatID     string
	client     *http.Client
	apiBaseURL string
}

func NewTelegram(botToken, chatID string) *Telegram {
	return &Telegram{
		botToken:   botToken,
		chatID:     chatID,
		client:     &http.Client{Timeout: 10 * time.Second},
		apiBaseURL: "https://api.telegram.org",
	}
}

// SendNewJobs sends notification about new job vacancies
func (t *Telegram) SendNewJobs(ctx context.Context, jobs []JobInfo) error {
	if len(jobs) == 0 {
		return nil
	}

	var messages []string
	for _, job := range jobs {
		msg := fmt.Sprintf(
			"🆕 <b>%s</b>\n🏢 %s\n📍 %s\n🔗 <a href=\"%s\">Apply</a>",
			escapeHTML(job.Title),
			escapeHTML(job.Company),
			escapeHTML(job.Location),
			job.URL,
		)
		messages = append(messages, msg)
	}

	if len(messages) > 5 {
		messages = messages[:5]
	}

	text := strings.Join(messages, "\n\n")
	if len(jobs) > 5 {
		text += fmt.Sprintf("\n\n... and %d more new jobs", len(jobs)-5)
	}

	return t.sendMessage(ctx, text)
}

// JobInfo contains minimal data required for Telegram notification
type JobInfo struct {
	Title    string
	Company  string
	Location string
	URL      string
}

func (t *Telegram) sendMessage(ctx context.Context, text string) error {
	apiURL := fmt.Sprintf("%s/bot%s/sendMessage", t.apiBaseURL, url.PathEscape(t.botToken))

	data := url.Values{}
	data.Set("chat_id", t.chatID)
	data.Set("text", text)
	data.Set("parse_mode", "HTML")
	data.Set("disable_web_page_preview", "true")

	req, err := http.NewRequestWithContext(ctx, "POST", apiURL, strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := t.client.Do(req)
	if err != nil {
		return fmt.Errorf("send failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram api returned %d", resp.StatusCode)
	}

	return nil
}

// escapeHTML escapes special characters for safe HTML parsing in Telegram
func escapeHTML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	return s
}