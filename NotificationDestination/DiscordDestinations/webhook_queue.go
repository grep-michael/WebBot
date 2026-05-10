package discorddestinations

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

var webhookQueue = make(chan webhookJob, 500)

type webhookJob struct {
	ctx      context.Context
	url      string
	payload  []byte
	result   chan error
	botToken string
}

type rateLimitResponse struct {
	RetryAfter float64 `json:"retry_after"`
	Global     bool    `json:"global"`
	Message    string  `json:"message"`
}

func drainWebhookQueue() {
	client := &http.Client{Timeout: 10 * time.Second}
	for job := range webhookQueue {
		for {
			req, err := http.NewRequestWithContext(job.ctx, http.MethodPost, job.url, bytes.NewReader(job.payload))
			if err != nil {
				job.result <- fmt.Errorf("discord: build request: %w", err)
				break
			}
			req.Header.Set("Content-Type", "application/json")
			if job.botToken != "" {
				req.Header.Set("Authorization", fmt.Sprintf("Bot %s", job.botToken))
			}
			resp, err := client.Do(req)
			if err != nil {
				job.result <- fmt.Errorf("discord: send: %w", err)
				break
			}
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()

			if resp.StatusCode == http.StatusTooManyRequests {
				wait := parseRateLimit(resp, body)
				//log.Printf("discord: rate limited, retrying after %s", wait)
				select {
				case <-job.ctx.Done():
					job.result <- fmt.Errorf("discord: context cancelled while rate limited: %w", job.ctx.Err())
					goto nextjob
				case <-time.After(wait):
					// retry the same job
					continue
				}
			}

			if resp.StatusCode != http.StatusNoContent {
				job.result <- fmt.Errorf("discord: unexpected status %d: %s", resp.StatusCode, string(body))
				break
			}
			job.result <- nil
			break
		}
	nextjob:
	}
}

func parseRateLimit(resp *http.Response, body []byte) time.Duration {
	if xHeaderLimit := resp.Header.Get("X-RateLimit-Reset-After"); xHeaderLimit != "" {
		if secs, err := strconv.ParseFloat(xHeaderLimit, 64); err == nil && secs > 0 {
			return time.Duration(secs*1000) * time.Millisecond
		}
	}
	if headerLimit := resp.Header.Get("Retry-After"); headerLimit != "" {
		if secs, err := strconv.ParseFloat(headerLimit, 64); err == nil && secs > 0 {
			return time.Duration(secs*1000) * time.Millisecond
		}
	}

	var rl rateLimitResponse
	if err := json.Unmarshal(body, &rl); err == nil && rl.RetryAfter > 0 {
		return time.Duration(rl.RetryAfter*1000) * time.Millisecond
	}

	// Last resort
	return 1 * time.Second
}

func init() {
	go drainWebhookQueue()
}
