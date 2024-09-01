package controllers

import (
	"context"
	"github.com/GitCMDR/microblogreposter-bot/internal/gateways"
	
	"gopkg.in/telebot.v3"
)

type Controller struct {
	MastodonGateway *gateways.MastodonGateway
}

func NewController(mastodonGateway *gateways.MastodonGateway) *Controller {
	return &Controller{
		MastodonGateway: mastodonGateway,
	}
}

func (c *Controller) ProcessMessage (tCtx telebot.Context) error {
	// set up some dummy context, need to find a way to integrate this onto the main request lifecycle, but this suffices for now
	ctx := context.Background()
	
	// post the message text to Mastodon
	status, err := c.MastodonGateway.PostStatus(ctx, tCtx.Text())
	if err != nil {
		return tCtx.Send("Failed to post status to Mastodon: " + err.Error())
	}

	return tCtx.Send("Posted to Mastodon: " + status.URL)
}

func (c *Controller) StartCommand (tCtx telebot.Context) error {
	return tCtx.Send("Ready to post.")
}

func (c *Controller) HelpCommand (tCtx telebot.Context) error {
	return tCtx.Send("You are on your own, kiddo *smirks*")
}