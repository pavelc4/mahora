package github

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

var (
	ErrUnauthorized = errors.New("github: unauthorized")
	ErrRateLimited  = errors.New("github: rate limited")
	ErrNotFound     = errors.New("github: repo not found")
)

type Event struct {
	ID      string          `json:"id"`
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type Notification struct {
	EventID string
	Type    string
	Message string
}

func ParseEvent(e Event) (*Notification, error) {
	parsers := map[string]func(Event) (*Notification, error){
		"PullRequestEvent": parsePullRequest,
		"ReleaseEvent":     parseRelease,
		"PushEvent":        parsePush,
	}

	parser, ok := parsers[e.Type]
	if !ok {
		return nil, nil
	}
	return parser(e)
}

type pullRequestPayload struct {
	Action      string `json:"action"`
	PullRequest struct {
		Title   string `json:"title"`
		HTMLURL string `json:"html_url"`
		Merged  bool   `json:"merged"`
		User    struct {
			Login string `json:"login"`
		} `json:"user"`
	} `json:"pull_request"`
}

type releasePayload struct {
	Action  string `json:"action"`
	Release struct {
		TagName string `json:"tag_name"`
		HTMLURL string `json:"html_url"`
		Name    string `json:"name"`
	} `json:"release"`
}

type pushPayload struct {
	Ref     string `json:"ref"`
	Commits []struct {
		Message string `json:"message"`
	} `json:"commits"`
}

func parsePush(e Event) (*Notification, error) {
	var p pushPayload
	if err := json.Unmarshal(e.Payload, &p); err != nil {
		return nil, fmt.Errorf("parsePush: %w", err)
	}
	if !strings.HasPrefix(p.Ref, "refs/heads/") {
		return nil, nil
	}
	if len(p.Commits) == 0 {
		return nil, nil
	}
	branch := strings.TrimPrefix(p.Ref, "refs/heads/")
	return &Notification{
		EventID: e.ID,
		Type:    "PushEvent",
		Message: fmt.Sprintf("📦 %d new commit(s) to <code>%s</code>", len(p.Commits), branch),
	}, nil
}

func parseRelease(e Event) (*Notification, error) {
	var p releasePayload
	if err := json.Unmarshal(e.Payload, &p); err != nil {
		return nil, fmt.Errorf("parseRelease: %w", err)
	}
	if p.Action != "published" {
		return nil, nil
	}
	return &Notification{
		EventID: e.ID,
		Type:    "ReleaseEvent",
		Message: fmt.Sprintf(
			"🚀 New release <b>%s</b>\n<a href=\"%s\">%s</a>",
			p.Release.TagName,
			p.Release.HTMLURL,
			p.Release.Name,
		),
	}, nil
}

func parsePullRequest(e Event) (*Notification, error) {
	var p pullRequestPayload
	if err := json.Unmarshal(e.Payload, &p); err != nil {
		return nil, fmt.Errorf("parsePullRequest: %w", err)
	}

	switch p.Action {
	case "opened":
		emoji := "🔔"
		if strings.HasPrefix(p.PullRequest.User.Login, "dependabot") ||
			strings.HasPrefix(p.PullRequest.User.Login, "renovate") {
			emoji = "📦"
		}
		return &Notification{
			EventID: e.ID,
			Type:    "PullRequestEvent",
			Message: fmt.Sprintf(
				"%s PR opened by <b>%s</b>\n<a href=\"%s\">%s</a>",
				emoji,
				p.PullRequest.User.Login,
				p.PullRequest.HTMLURL,
				p.PullRequest.Title,
			),
		}, nil

	case "closed":
		if !p.PullRequest.Merged {
			return nil, nil
		}
		return &Notification{
			EventID: e.ID,
			Type:    "PullRequestEvent",
			Message: fmt.Sprintf(
				"🔀 PR merged by <b>%s</b>\n<a href=\"%s\">%s</a>",
				p.PullRequest.User.Login,
				p.PullRequest.HTMLURL,
				p.PullRequest.Title,
			),
		}, nil

	default:
		return nil, nil
	}
}
