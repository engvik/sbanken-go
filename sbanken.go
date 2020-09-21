package sbanken

import (
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

func NewClient(cfg *Config, httpClient *http.Client) (*Client, error) {
	if err := cfg.validate(); err != nil {
		return nil, err
	}

	c := &Client{}
	c.setHTTPClient(httpClient)
	c.config = cfg
	c.baseURL = "https://api.sbanken.no/exec.bank/api"

	if err := c.authorize(cfg); err != nil {
		return nil, err
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
