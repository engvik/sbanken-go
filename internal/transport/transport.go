package transport

import (
	"context"
	"net/http"
)

// Client represents the transport client.
type Client struct {
	clientID     string
	clientSecret string
	userAgent    string
	http         *http.Client
	auth         *auth
}

// Config represents the transport config.
type Config struct {
	ClientID     string
	ClientSecret string
	UserAgent    string
}

// New returns a transport client.
func New(ctx context.Context, cfg *Config, httpClient *http.Client) *Client {
	c := &Client{
		clientID:     cfg.ClientID,
		clientSecret: cfg.ClientSecret,
		userAgent:    cfg.UserAgent,
	}

	c.setHTTPClient(httpClient)

	return c
}

func (c *Client) setHTTPClient(httpClient *http.Client) {
	if httpClient == nil {
		c.http = http.DefaultClient
		return
	}

	c.http = httpClient
}
