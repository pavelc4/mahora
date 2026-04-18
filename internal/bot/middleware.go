package bot

import (
	"context"
	"log/slog"

	tele "gopkg.in/telebot.v3"
)

func (b *Bot) logMiddleware() tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			slog.Info("telegram update",
				"from", c.Sender().Username,
				"text", c.Text(),
			)
			return next(c)
		}
	}
}

func (b *Bot) requireAuth(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		user, err := b.queries.GetUserByTelegramID(context.Background(), c.Sender().ID)
		if err != nil || !user.GithubToken.Valid {
			return c.Send(" You're not logged in. Use /login first.", htmlOpt)
		}
		return next(c)
	}
}
