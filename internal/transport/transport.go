package transport

import (
	"context"
	"net/http"
)

type Client struct {
	clientID     string
	clientSecret string
	customerID   string
	http         *http.Client
	auth         *auth
}

type Config struct {
	ClientID     string
	ClientSecret string
	CustomerID   string
}

func New(ctx context.Context, cfg *Config, httpClient *http.Client) *Client {
	c := &Client{
		clientID:     cfg.ClientID,
		clientSecret: cfg.ClientSecret,
		customerID:   cfg.CustomerID,
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
