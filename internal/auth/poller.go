package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

var (
	ErrTokenNotReady = errors.New("token not ready")
	ErrPollTimeout   = errors.New("poll timeout")
)

type TokenResponse struct {
	GitHubLogin string `json:"github_login"`
	AccessToken string `json:"access_token"`
}

type Poller interface {
	Poll(ctx context.Context, stateToken string, telegramID int64) (*TokenResponse, error)
}

type workerPoller struct {
	client     *http.Client
	workerURL  string
	secret     string
	interval   time.Duration
	telegramID int64
}

type Option func(*workerPoller)

func WithInterval(d time.Duration) Option {
	return func(p *workerPoller) { p.interval = d }
}

func WithHTTPClient(c *http.Client) Option {
	return func(p *workerPoller) { p.client = c }
}

func NewPoller(workerURL, secret string, opts ...Option) Poller {
	p := &workerPoller{
		client:    &http.Client{Timeout: 10 * time.Second},
		workerURL: workerURL,
		secret:    secret,
		interval:  3 * time.Second,
	}
	for _, o := range opts {
		o(p)
	}
	return p
}

func (p *workerPoller) Poll(ctx context.Context, stateToken string, telegramID int64) (*TokenResponse, error) {
	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("auth.Poll: %w", ErrPollTimeout)
		case <-ticker.C:
			tok, err := p.fetchToken(ctx, stateToken, telegramID)
			if err == nil {
				return tok, nil
			}
			if ctx.Err() != nil {
				return nil, fmt.Errorf("auth.Poll: %w", ErrPollTimeout)
			}
			if errors.Is(err, ErrTokenNotReady) {
				continue
			}
			return nil, fmt.Errorf("auth.Poll: %w", err)
		}
	}
}

func (p *workerPoller) fetchToken(ctx context.Context, stateToken string, telegramID int64) (*TokenResponse, error) {
	url := fmt.Sprintf("%s/auth/token?state=%s&tid=%d", p.workerURL, stateToken, telegramID)

	slog.Debug("fetchToken", "url", url) // tambah ini
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("fetchToken build request: %w", err)
	}
	req.Header.Set("X-Worker-Secret", p.secret)

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetchToken do request: %w", err)
	}
	defer resp.Body.Close()

	handlers := map[int]func() (*TokenResponse, error){
		http.StatusAccepted: func() (*TokenResponse, error) {
			return nil, ErrTokenNotReady // 202 = explicitly pending
		},
		http.StatusNotFound: func() (*TokenResponse, error) {
			return nil, ErrTokenNotReady // 404 = belum ada di KV, keep polling
		},
		http.StatusOK: func() (*TokenResponse, error) {
			return decodeTokenResponse(resp)
		},
	}

	if h, ok := handlers[resp.StatusCode]; ok {
		return h()
	}
	return nil, fmt.Errorf("fetchToken unexpected status: %d", resp.StatusCode)
}

func decodeTokenResponse(resp *http.Response) (*TokenResponse, error) {
	var tok TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tok); err != nil {
		return nil, fmt.Errorf("decodeTokenResponse: %w", err)
	}
	if tok.AccessToken == "" {
		return nil, fmt.Errorf("decodeTokenResponse: empty access_token")
	}
	return &tok, nil
}
