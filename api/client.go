package api

import (
	"net/http"
	"time"
)

func NewClient(url, apiKey string) *DatabaseHTTPClient {
	return &DatabaseHTTPClient{
		URL:    url,
		APIKey: apiKey,
		Client: &http.Client{Timeout: 5 * time.Second},
	}
}

func (c *DatabaseHTTPClient) SetHeaders(headers map[string]string) {
	c.Headers = headers
}
