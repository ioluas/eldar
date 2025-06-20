package api

import "net/http"

type DatabaseHTTPClient struct {
	URL         string
	APIKey      string
	Client      *http.Client
	Headers     map[string]string
	AccessToken string
	UserID      string
	UserEmail   string
}
