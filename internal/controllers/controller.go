package controllers

import (
	"context"
	"fmt"
	"strings"

	"github.com/GitCMDR/microblogreposter-bot/internal/gateways"
	"gopkg.in/telebot.v3"
)

type Controller struct {
	MastodonGateway *gateways.MastodonGateway
	BlueskyGateway  *gateways.BlueskyGateway
}

func NewController(mastodonGateway *gateways.MastodonGateway, blueskyGateway *gateways.BlueskyGateway) *Controller {
	return &Controller{
		MastodonGateway: mastodonGateway,
		BlueskyGateway:  blueskyGateway,
	}
}

func (c *Controller) ProcessMessage (tCtx telebot.Context) error {
	// set up some dummy context, need to find a way to integrate this onto the main request lifecycle, but this suffices for now
	ctx := context.Background()
	
	// post the message text to Mastodon
	mastodonStatus, err := c.MastodonGateway.PostStatus(ctx, tCtx.Text())
	if err != nil {
		return tCtx.Send("Failed to post status to Mastodon: " + err.Error())
	}

	blueskyStatus, err := c.BlueskyGateway.PostStatus(tCtx.Text())
	if err != nil {
		return tCtx.Send("Failed to post status to Bluesky: " + err.Error())
	}
	blueSkyWebURL := c.convertBlueskyURIToWebURL(blueskyStatus.Uri)

	return tCtx.Send(fmt.Sprintf("Posted to Mastodon: %s\nPosted to Bluesky: %s", mastodonStatus.URL, blueSkyWebURL))
}

func (c *Controller) StartCommand (tCtx telebot.Context) error {
	return tCtx.Send("Ready to post.")
}

func (c *Controller) HelpCommand (tCtx telebot.Context) error {
	return tCtx.Send("You are on your own, kiddo *smirks*")
}

// Helper function to convert AT Protocol URI to a web URL
func (c *Controller) convertBlueskyURIToWebURL(uri string) string {
	// example: at://did:plc:kuhyi7lvatum5vomob6fvw5l/app.bsky.feed.post/3l33dphr2sj2g

	parts := strings.Split(uri, "/")
	if len(parts) < 2 {
		return uri // Fallback to original URI if something goes wrong
	}

	// assume i'll stick to this handle forever, TODO: get the value from somewhere else.
	handle := "lbnvds.bsky.social"

	// The post ID is the last part of the URI
	postID := parts[len(parts)-1]

	return fmt.Sprintf("https://bsky.app/profile/%s/post/%s", handle, postID)
}