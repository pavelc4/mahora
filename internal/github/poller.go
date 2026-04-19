package github

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	tele "gopkg.in/telebot.v3"

	dbgen "github.com/pavelc4/mahora/internal/db/gen"
)

type Poller struct {
	queries  *dbgen.Queries
	bot      *tele.Bot
	interval time.Duration
	clients  map[int64]*Client
	mu       sync.Mutex
}

func NewPoller(queries *dbgen.Queries, bot *tele.Bot, interval time.Duration) *Poller {
	return &Poller{
		queries:  queries,
		bot:      bot,
		interval: interval,
		clients:  make(map[int64]*Client),
	}
}

func (p *Poller) getClient(telegramID int64, token string) *Client {
	p.mu.Lock()
	defer p.mu.Unlock()
	if c, ok := p.clients[telegramID]; ok {
		return c
	}
	c := NewClient(token)
	p.clients[telegramID] = c
	return c
}

func (p *Poller) Run(ctx context.Context) {
	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()

	slog.Info("github poller started", "interval", p.interval)

	for {
		select {
		case <-ctx.Done():
			slog.Info("github poller stopped")
			return
		case <-ticker.C:
			if err := p.poll(ctx); err != nil {
				slog.Error("github poller tick failed", "err", err)
			}
		}
	}
}

func (p *Poller) poll(ctx context.Context) error {
	users, err := p.queries.ListUsersWithToken(ctx)
	if err != nil {
		return fmt.Errorf("poll list users: %w", err)
	}

	for _, user := range users {
		if err = p.pollUser(ctx, user); err != nil {
			slog.Warn("poll user failed",
				"telegram_id", user.TelegramID,
				"err", err,
			)
		}
	}
	return nil
}

func (p *Poller) pollUser(ctx context.Context, user dbgen.User) error {
	client := NewClient(user.GithubToken.String)

	repos, err := p.queries.ListReposByUser(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("pollUser list repos: %w", err)
	}

	for _, repo := range repos {
		if err = p.pollRepo(ctx, user, client, repo); err != nil {
			slog.Warn("poll repo failed",
				"repo", repo.Owner+"/"+repo.Name,
				"err", err,
			)
		}
	}
	return nil
}

func (p *Poller) pollRepo(ctx context.Context, user dbgen.User, client *Client, repo dbgen.Repo) error {
	events, err := client.GetEvents(ctx, repo.Owner, repo.Name)
	if err != nil {
		if errors.Is(err, ErrUnauthorized) {
			_ = p.queries.ClearGitHubToken(ctx, user.TelegramID)
			p.bot.Send(&tele.User{ID: user.TelegramID},
				"⚠️ Your GitHub token expired. Please /login again.",
				&tele.SendOptions{ParseMode: tele.ModeHTML},
			)
		}
		return fmt.Errorf("pollRepo get events: %w", err)
	}

	repoFull := repo.Owner + "/" + repo.Name

	for _, event := range events {
		notif, err := ParseEvent(event)
		if err != nil || notif == nil {
			continue
		}

		count, err := p.queries.HasNotification(ctx, dbgen.HasNotificationParams{
			UserID:    user.ID,
			RepoFull:  repoFull,
			EventType: notif.Type,
			EventID:   notif.EventID,
		})
		if err != nil || count > 0 {
			continue
		}

		msg := fmt.Sprintf("📌 <b>%s</b>\n\n%s", repoFull, notif.Message)
		if _, err = p.bot.Send(&tele.User{ID: user.TelegramID}, msg,
			&tele.SendOptions{ParseMode: tele.ModeHTML},
		); err != nil {
			slog.Warn("send notification failed", "err", err)
			continue
		}

		_ = p.queries.InsertNotification(ctx, dbgen.InsertNotificationParams{
			UserID:    user.ID,
			RepoFull:  repoFull,
			EventType: notif.Type,
			EventID:   notif.EventID,
		})
	}
	return nil
}
