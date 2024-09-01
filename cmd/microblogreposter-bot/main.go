package main

import (
	"log"
    
	"github.com/GitCMDR/microblogreposter-bot/internal/bot"
	"github.com/GitCMDR/microblogreposter-bot/internal/config"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Initialize bot
	b, err := bot.NewBot(cfg)
	if err != nil {
		log.Fatalf("Error initializing bot: %v", err)
	}

	// Start the bot
	b.Start()
}