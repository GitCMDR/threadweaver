package bot

import (
	"time"

	"gopkg.in/telebot.v3"
	"github.com/GitCMDR/microblogreposter-bot/internal/config"
	"github.com/GitCMDR/microblogreposter-bot/internal/handlers"
	"github.com/GitCMDR/microblogreposter-bot/internal/controllers"
)

type Bot struct {
	*telebot.Bot
}

func NewBot(cfg *config.Config) (*Bot, error) {
	// Initialize the bot
	b, err := telebot.NewBot(telebot.Settings{
		Token:  cfg.TelegramToken,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		return nil, err
	}

	// Set up controller
	controller := controllers.NewController()

	// Set up handler
	handler := handlers.NewHandler(controller)

	// Handle messages and commands
	b.Handle(telebot.OnText, handler.HandleMessage)
	b.Handle("/start", handler.HandleStart)
	b.Handle("/help", handler.HandleHelp)

	return &Bot{b}, nil
}

func (b *Bot) Start() {
	b.Bot.Start()
}
