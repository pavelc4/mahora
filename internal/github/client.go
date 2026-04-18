package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const baseURL = "https://api.github.com"

type Client struct {
	http  *http.Client
	token string
}

func NewClient(token string) *Client {
	return &Client{
		token: token,
		http:  &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *Client) GetEvents(ctx context.Context, owner, repo string) ([]Event, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/events?per_page=30", baseURL, owner, repo)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("github.GetEvents build request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("github.GetEvents do: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusNotModified:
		return nil, nil
	case http.StatusUnauthorized:
		return nil, fmt.Errorf("github.GetEvents: %w", ErrUnauthorized)
	case http.StatusForbidden:
		return nil, fmt.Errorf("github.GetEvents: %w", ErrRateLimited)
	case http.StatusNotFound:
		return nil, fmt.Errorf("github.GetEvents: %w", ErrNotFound)
	default:
		return nil, fmt.Errorf("github.GetEvents unexpected status: %d", resp.StatusCode)
	}

	var events []Event
	if err = json.NewDecoder(resp.Body).Decode(&events); err != nil {
		return nil, fmt.Errorf("github.GetEvents decode: %w", err)
	}
	return events, nil
}
