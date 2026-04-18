package bot

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	tele "gopkg.in/telebot.v3"

	dbgen "github.com/pavelc4/mahora/internal/db/gen"
)

func (b *Bot) handleStart(c tele.Context) error {
	return c.Send(
		"👋 Hey! I'm <b>Mahora</b> — your GitHub notification bot.\n\nType /login to get started.",
		htmlOpt,
	)
}

func (b *Bot) handleHelp(c tele.Context) error {
	return c.Send(
		"📖 <b>Commands:</b>\n\n"+
			"/login — Connect your GitHub account\n"+
			"/logout — Disconnect your GitHub account\n"+
			"/watch <code>owner/repo</code> — Watch a repo\n"+
			"/unwatch <code>owner/repo</code> — Unwatch a repo\n"+
			"/list — List your watched repos",
		htmlOpt,
	)
}

func (b *Bot) handleLogin(c tele.Context) error {
	stateToken := uuid.NewString()
	loginURL := fmt.Sprintf("%s/auth/github?state=%s", b.cfg.WorkerURL, stateToken)

	if err := c.Send(fmt.Sprintf(
		"🔐 Click the link below to login with GitHub:\n\n<a href=\"%s\">Authorize Mahora</a>\n\n<i>Link expires in 5 minutes.</i>",
		loginURL,
	), htmlOpt); err != nil {
		return fmt.Errorf("handleLogin send: %w", err)
	}

	go b.pollAndSaveToken(c.Sender().ID, stateToken)
	return nil
}

func (b *Bot) pollAndSaveToken(telegramID int64, stateToken string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	tok, err := b.poller.Poll(ctx, stateToken)
	if err != nil {
		slog.Warn("pollAndSaveToken timeout", "telegram_id", telegramID, "err", err)
		b.tele.Send(
			&tele.User{ID: telegramID},
			"⏰ Login timed out. Type /login to try again.",
			htmlOpt,
		)
		return
	}

	_, err = b.queries.UpsertUser(ctx, dbgen.UpsertUserParams{
		TelegramID:  telegramID,
		GithubLogin: sql.NullString{String: tok.GitHubLogin, Valid: true},
		GithubToken: sql.NullString{String: tok.AccessToken, Valid: true},
	})
	if err != nil {
		slog.Error("pollAndSaveToken upsert failed", "telegram_id", telegramID, "err", err)
		b.tele.Send(
			&tele.User{ID: telegramID},
			"❌ Something went wrong saving your token. Try /login again.",
			htmlOpt,
		)
		return
	}

	b.tele.Send(
		&tele.User{ID: telegramID},
		fmt.Sprintf("✅ Logged in as <b>%s</b>!", tok.GitHubLogin),
		htmlOpt,
	)
}

func (b *Bot) handleLogout(c tele.Context) error {
	if err := b.queries.ClearGitHubToken(context.Background(), c.Sender().ID); err != nil {
		return fmt.Errorf("handleLogout: %w", err)
	}
	return c.Send("👋 Logged out. Your GitHub token has been removed.", htmlOpt)
}

func (b *Bot) handleWatch(c tele.Context) error {
	return c.Send("🚧 <i>Coming soon.</i>", htmlOpt)
}

func (b *Bot) handleUnwatch(c tele.Context) error {
	return c.Send("🚧 <i>Coming soon.</i>", htmlOpt)
}

func (b *Bot) handleList(c tele.Context) error {
	return c.Send("🚧 <i>Coming soon.</i>", htmlOpt)
}
