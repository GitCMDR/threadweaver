package bot

import (
	"time"

	"gopkg.in/telebot.v3"
	"github.com/GitCMDR/threadweaver/internal/config"
	"github.com/GitCMDR/threadweaver/internal/handlers"
	"github.com/GitCMDR/threadweaver/internal/controllers"
	"github.com/GitCMDR/threadweaver/internal/gateways"
)

type Bot struct {
	*telebot.Bot
}

func NewBot(cfg *config.Config) (*Bot, error) {
	// initialise the bot
	b, err := telebot.NewBot(telebot.Settings{
		Token:  cfg.TelegramToken,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		return nil, err
	}

	// initialise gateways
	mastodonGateway := gateways.NewMastodonGateway(cfg.MastodonServerURL, cfg.MastodonClientID, cfg.MastodonClientSecret, cfg.MastodonAccessToken)
	blueskyGateway, err := gateways.NewBlueskyGateway(cfg.BlueskyAPIURL, cfg.BlueskyUsername, cfg.BlueskyPassword)
	if err != nil {
		return nil, err
	}

	threadsGateway := gateways.NewThreadsGateway() // okay this is iffy; lets rectify in future update (we should pass the env values here)

	// set up controller
	controller := controllers.NewController(mastodonGateway, blueskyGateway, threadsGateway)

	// set up handler
	handler := handlers.NewHandler(controller)

	// handler messages and commands
	b.Handle(telebot.OnText, handler.HandleMessage)
	b.Handle("/start", handler.HandleStart)
	b.Handle("/help", handler.HandleHelp)

	return &Bot{b}, nil
}

func (b *Bot) Start() {
	b.Bot.Start()
}
