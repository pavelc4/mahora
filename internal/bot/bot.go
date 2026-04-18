package bot

import (
	"fmt"
	"time"

	tele "gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"

	"github.com/pavelc4/mahora/internal/auth"
	"github.com/pavelc4/mahora/internal/config"
	dbgen "github.com/pavelc4/mahora/internal/db/gen"
)

var htmlOpt = &tele.SendOptions{ParseMode: tele.ModeHTML}

type Bot struct {
	tele    *tele.Bot
	queries *dbgen.Queries
	poller  auth.Poller
	cfg     *config.Config
}

func New(cfg *config.Config, queries *dbgen.Queries, poller auth.Poller) (*Bot, error) {
	pref := tele.Settings{
		Token:  cfg.TelegramToken,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		return nil, fmt.Errorf("bot.New: %w", err)
	}

	bot := &Bot{
		tele:    b,
		queries: queries,
		poller:  poller,
		cfg:     cfg,
	}

	bot.registerRoutes()
	return bot, nil
}

func (b *Bot) Start() {
	b.tele.Start()
}

func (b *Bot) Stop() {
	b.tele.Stop()
}

func (b *Bot) registerRoutes() {
	b.tele.Use(middleware.Recover())
	b.tele.Use(b.logMiddleware())

	b.tele.Handle("/start", b.handleStart)
	b.tele.Handle("/help", b.handleHelp)
	b.tele.Handle("/login", b.handleLogin)

	protected := map[string]tele.HandlerFunc{
		"/logout":  b.handleLogout,
		"/watch":   b.handleWatch,
		"/unwatch": b.handleUnwatch,
		"/list":    b.handleList,
	}
	for cmd, handler := range protected {
		b.tele.Handle(cmd, b.requireAuth(handler))
	}
}
