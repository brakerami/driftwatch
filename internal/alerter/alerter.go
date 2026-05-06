// Package alerter provides webhook-based alerting for drift findings.
// When drift is detected, alerts can be dispatched to an external HTTP endpoint.
package alerter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/user/driftwatch/internal/drift"
)

const defaultTimeout = 10 * time.Second

// Payload is the JSON body sent to the webhook endpoint.
type Payload struct {
	Timestamp string          `json:"timestamp"`
	Container string          `json:"container"`
	Findings  []drift.Finding `json:"findings"`
}

// Alerter dispatches drift findings to a configured webhook URL.
type Alerter struct {
	webhookURL string
	client     *http.Client
}

// New creates an Alerter that posts to the given webhook URL.
// Returns an error if webhookURL is empty.
func New(webhookURL string) (*Alerter, error) {
	if webhookURL == "" {
		return nil, fmt.Errorf("alerter: webhook URL must not be empty")
	}
	return &Alerter{
		webhookURL: webhookURL,
		client:     &http.Client{Timeout: defaultTimeout},
	}, nil
}

// Send posts the findings for the named container to the webhook endpoint.
// A non-2xx response is treated as an error.
func (a *Alerter) Send(ctx context.Context, container string, findings []drift.Finding) error {
	if len(findings) == 0 {
		return nil
	}

	payload := Payload{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Container: container,
		Findings:  findings,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("alerter: marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, a.webhookURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("alerter: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return fmt.Errorf("alerter: http post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("alerter: unexpected status %d from webhook", resp.StatusCode)
	}
	return nil
}
