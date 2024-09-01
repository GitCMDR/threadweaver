package gateways

import (
	"context"
	"github.com/mattn/go-mastodon"
)

type MastodonGateway struct {
	Client *mastodon.Client
}

func NewMastodonGateway(serverURL, clientID, clientSecret, accessToken string) *MastodonGateway {
	client := mastodon.NewClient(&mastodon.Config{
		Server:       serverURL,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		AccessToken:  accessToken,
	})
	return &MastodonGateway{Client: client}
}

func (g *MastodonGateway) PostStatus(ctx context.Context, status string) (*mastodon.Status, error) {
	return g.Client.PostStatus(ctx, &mastodon.Toot{
		Status: status,
	})
}
