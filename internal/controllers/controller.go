package controllers

import (
	"context"
	"fmt"
	"strings"
	"os"

	"github.com/GitCMDR/microblogreposter-bot/internal/gateways"
	"gopkg.in/telebot.v3"
)

type Controller struct {
	MastodonGateway *gateways.MastodonGateway
	BlueskyGateway  *gateways.BlueskyGateway
	ThreadsGateway  *gateways.ThreadsGateway
}

func NewController(mastodonGateway *gateways.MastodonGateway, blueskyGateway *gateways.BlueskyGateway, threadsGateway *gateways.ThreadsGateway) *Controller {
	return &Controller{
		MastodonGateway: mastodonGateway,
		BlueskyGateway:  blueskyGateway,
		ThreadsGateway:  threadsGateway,
	}
}

func (c *Controller) ProcessMessageWithText (tCtx telebot.Context, txt string) error {
	// set up some dummy context, need to find a way to integrate this onto the main request lifecycle, but this suffices for now
	ctx := context.Background()
	
	// post the message text to Mastodon
	mastodonStatus, err := c.MastodonGateway.PostStatus(ctx, txt)
	if err != nil {
		return tCtx.Send("Failed to post status to Mastodon: " + err.Error())
	}

	blueskyStatus, err := c.BlueskyGateway.PostStatus(txt)
	if err != nil {
		return tCtx.Send("Failed to post status to Bluesky: " + err.Error())
	}
	blueSkyWebURL := c.convertBlueskyURIToWebURL(blueskyStatus.Uri)

	containerID, err := c.ThreadsGateway.CreateMediaContainer(txt, "", "")
	if err != nil {
		return tCtx.Send("Failed to create media container on Threads: " + err.Error())
	}

	threadsURL, err := c.ThreadsGateway.PublishContainer(containerID)
	if err != nil {
		return tCtx.Send("Failed to publish on Threads: " + err.Error())
	}

	return tCtx.Send(fmt.Sprintf("Posted to Mastodon: %s\nPosted to Bluesky: %s\nPosted to Threads: %s", mastodonStatus.URL, blueSkyWebURL, threadsURL))
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

	// assume i'll stick to this handle forever, TODO: get the value from somewhere else?
	handle := os.Getenv("BLUESKY_HANDLE")

	// The post ID is the last part of the URI
	postID := parts[len(parts)-1]

	return fmt.Sprintf("https://bsky.app/profile/%s/post/%s", handle, postID)
}