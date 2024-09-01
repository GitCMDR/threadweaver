package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

// Config is a pretty simple struct to hold config values
type Config struct {
	TelegramToken string
	MastodonServerURL string
	MastodonClientID string
	MastodonClientSecret string
	MastodonAccessToken string
	BlueskyAPIURL      string
	BlueskyUsername    string
	BlueskyPassword    string
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
		MastodonServerURL: os.Getenv("MASTODON_SERVER_URL"),
		MastodonClientID: os.Getenv("MASTODON_CLIENT_ID"),
		MastodonClientSecret: os.Getenv("MASTODON_CLIENT_SECRET"),
		MastodonAccessToken: os.Getenv("MASTODON_ACCESS_TOKEN"),
		BlueskyAPIURL:      os.Getenv("BLUESKY_API_URL"),
		BlueskyUsername:    os.Getenv("BLUESKY_USERNAME"),
		BlueskyPassword:    os.Getenv("BLUESKY_PASSWORD"),
		ThreadsToken: os.Getenv("THREADS_TOKEN"),
	}, nil
}
