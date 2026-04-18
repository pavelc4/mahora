package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	TelegramToken      string `env:"BOT_TOKEN"            env-required:"true"`
	GitHubClientID     string `env:"GITHUB_CLIENT_ID"`
	GitHubClientSecret string `env:"GITHUB_CLIENT_SECRET"`
	WorkerURL          string `env:"WORKER_URL"`
	WorkerSecret       string `env:"WORKER_SECRET"`
	DBPath             string `env:"DB_PATH"              env-default:"./mahora.db"`
	Env                string `env:"ENV"                  env-default:"development"`
	LogLevel           string `env:"LOG_LEVEL"            env-default:"info"`
}

func Load() (*Config, error) {
	var cfg Config
	if err := cleanenv.ReadConfig(".env", &cfg); err != nil {
		return nil, fmt.Errorf("config.Load: %w", err)
	}
	return &cfg, nil
}
