package gateways

import (
	"time"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type BlueskyGateway struct {
	APIURL string
	Auth   *BlueskyAuthResponse
}

type BlueskyAuthResponse struct {
	AccessJWT string `json:"accessJwt"`
	Handle    string `json:"handle"`
	Did       string `json:"did"`
}

type BlueskyPostResponse struct {
	Uri  string `json:"uri"`
	Cid  string `json:"cid"`
	Did  string `json:"did"`
}

func NewBlueskyGateway(apiURL, username, password string) (*BlueskyGateway, error) {
	gateway := &BlueskyGateway{
		APIURL: apiURL,
	}

	// Authenticate with Bluesky
	if err := gateway.Authenticate(username, password); err != nil {
		return nil, err
	}

	return gateway, nil
}

func (g *BlueskyGateway) Authenticate(username, password string) error {
	url := fmt.Sprintf("%s/xrpc/com.atproto.server.createSession", g.APIURL)
	reqBody, _ := json.Marshal(map[string]string{
		"identifier": username,
		"password":   password,
	})

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to authenticate with Bluesky: %s", resp.Status)
	}

	var authResp BlueskyAuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return err
	}

	g.Auth = &authResp
	return nil
}

func (g *BlueskyGateway) PostStatus(status string) (*BlueskyPostResponse, error) {
	url := fmt.Sprintf("%s/xrpc/com.atproto.repo.createRecord", g.APIURL)
	payload := map[string]interface{}{
		"collection": "app.bsky.feed.post",
		"repo":       g.Auth.Did,
		"record": map[string]interface{}{
			"$type":    "app.bsky.feed.post",
			"text":     status,
			"createdAt": time.Now().Format(time.RFC3339),
		},
	}

	reqBody, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+g.Auth.AccessJWT)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to post status to Bluesky: %s", resp.Status)
	}

	var postResp BlueskyPostResponse
	if err := json.NewDecoder(resp.Body).Decode(&postResp); err != nil {
		return nil, err
	}

	return &postResp, nil
}
