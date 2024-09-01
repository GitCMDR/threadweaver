package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

// Config is a pretty simple struct to hold config values
type Config struct {
	TelegramToken string
	MastodonToken string
	BlueskyToken string
	ThreadsToken string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load() // no need to specify path here, it will read stuff from the root dir anyway
	if err != nil {
		log.Fatalf("error loading .env file: %v", err) // let's kill the app if this fails
		return nil, err
	}

	return &Config{
		TelegramToken: os.Getenv("TELEGRAM_TOKEN"),
		MastodonToken: os.Getenv("MASTODON_TOKEN"),
		BlueskyToken: os.Getenv("BLUESKY_TOKEN"),
		ThreadsToken: os.Getenv("THREADS_TOKEN"),
	}, nil
}
