package bot

import (
	"time"

	"gopkg.in/telebot.v3"
	"github.com/GitCMDR/microblogreposter-bot/internal/config"
	"github.com/GitCMDR/microblogreposter-bot/internal/handlers"
	"github.com/GitCMDR/microblogreposter-bot/internal/controllers"
	"github.com/GitCMDR/microblogreposter-bot/internal/gateways"
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

	// set up controller
	controller := controllers.NewController(mastodonGateway, blueskyGateway)

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
