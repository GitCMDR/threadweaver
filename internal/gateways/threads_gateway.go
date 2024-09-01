package gateways

import (
	"bytes"
	"io"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type ThreadsGateway struct {
	AccessToken string
	UserID      string
}

func NewThreadsGateway() *ThreadsGateway {
	return &ThreadsGateway{
		AccessToken: os.Getenv("THREADS_ACCESS_TOKEN"),
		UserID:      os.Getenv("THREADS_USER_ID"),
	}
}

// CreateMediaContainer creates a media container with text (and optionally, media)
func (g *ThreadsGateway) CreateMediaContainer(message string, mediaURL string, mediaType string) (string, error) {
	url := fmt.Sprintf("https://graph.threads.net/v1.0/%s/threads", g.UserID)

	payload := map[string]interface{}{
		"text": message,
	}

	// Add media type and corresponding URL if provided
	if mediaType != "" && mediaURL != "" {
		switch mediaType {
		case "IMAGE":
			payload["media_type"] = "IMAGE"
			payload["image_url"] = mediaURL
		case "VIDEO":
			payload["media_type"] = "VIDEO"
			payload["video_url"] = mediaURL
		default:
			return "", fmt.Errorf("invalid media type provided")
		}
	} else {
		// Default to text-only if no media is provided
		payload["media_type"] = "TEXT"
	}

	reqBody, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+g.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var bodyBytes []byte
		if resp.Body != nil {
			bodyBytes, _ = io.ReadAll(resp.Body)
		}
		return "", fmt.Errorf("failed to create media container on Threads: %s\n%s", resp.Status, string(bodyBytes))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	containerID, ok := result["id"].(string)
	if !ok {
		return "", fmt.Errorf("failed to retrieve container ID")
	}

	return containerID, nil
}

func (g *ThreadsGateway) PublishContainer(containerID string) (string, error) {
	url := fmt.Sprintf("https://graph.threads.net/v1.0/%s/threads_publish", g.UserID)

	payload := map[string]string{
		"creation_id": containerID,
	}

	reqBody, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+g.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var bodyBytes []byte
		if resp.Body != nil {
			bodyBytes, _ = io.ReadAll(resp.Body)
		}
		return "", fmt.Errorf("failed to publish container on Threads: %s\n%s", resp.Status, string(bodyBytes))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	threadID, ok := result["id"].(string)
	if !ok {
		return "", fmt.Errorf("failed to retrieve thread ID")
	}

	// Retrieve the username for the URL
	username, err := g.getUsername()
	if err != nil {
		return "", fmt.Errorf("failed to retrieve username: %v", err)
	}

	// Construct the correct URL
	threadURL := fmt.Sprintf("https://www.threads.net/@%s/post/%s", username, threadID)

	return threadURL, nil
}

// Function to get the username associated with the user ID
func (g *ThreadsGateway) getUsername() (string, error) {
	url := fmt.Sprintf("https://graph.threads.net/v1.0/%s?fields=username", g.UserID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+g.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to retrieve username: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	username, ok := result["username"].(string)
	if !ok {
		return "", fmt.Errorf("failed to retrieve username")
	}

	return username, nil
}