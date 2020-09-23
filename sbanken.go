package sbanken

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	HTTP    *http.Client
	config  *Config
	auth    *auth
	baseURL string
}

type auth struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
	ExpiresIn   int    `json:"expires_in"`
	expires     time.Time
}

func NewClient(ctx context.Context, cfg *Config, httpClient *http.Client) (*Client, error) {
	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}

	c := &Client{}
	c.setHTTPClient(httpClient)
	c.config = cfg
	c.baseURL = "https://api.sbanken.no/exec.bank/api"

	if err := c.authorize(ctx, cfg); err != nil {
		return nil, fmt.Errorf("authorize: %w", err)
	}

	return c, nil
}

func (c *Client) setHTTPClient(httpClient *http.Client) {
	if httpClient == nil {
		c.HTTP = http.DefaultClient
		return
	}

	c.HTTP = httpClient
}
