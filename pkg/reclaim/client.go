package reclaim

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	apiUrl = "https://api.app.reclaim.ai"
)

type Client struct {
	h      http.Client
	apiKey string
}

func New() *Client {
	return &Client{h: http.Client{
		Timeout: time.Second * 10,
	},
		apiKey: os.Getenv("RECLAIM_API_KEY")}
}

func (c *Client) CreateTask(title string, minChunkSize int, maxChunkSize int, timeChunksRequired int) error {
	requestBody := fmt.Sprintf(`{
		"title": "%s",
		"status": "NEW",
		"minChunkSize": %d,
		"maxChunkSize": %d,
		"timeChunksRequired": %d,
		"eventCategory": "WORK"
}`, title, minChunkSize, maxChunkSize, timeChunksRequired)
	request, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/tasks", apiUrl), strings.NewReader(requestBody))

	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	response, err := c.h.Do(request)
	if err != nil {
		return err
	}

	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code %d when creating a task:\n%s", response.StatusCode, requestBody)
	}

	return nil
}
