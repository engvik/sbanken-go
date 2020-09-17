package sbanken

import (
	"net/http"
)

type Client struct {
	HTTP *http.Client
}

func NewClient(cfg *Config, httpClient *http.Client) (*Client, error) {
	if err := cfg.validate(); err != nil {
		return nil, err
	}

	c := &Client{}
	c.setHTTPClient(httpClient)

	c.authorize()

	return c, nil
}

func (c *Client) setHTTPClient(httpClient *http.Client) {
	if httpClient == nil {
		c.HTTP = http.DefaultClient
		return
	}

	c.HTTP = httpClient
}

func (c *Client) authorize() {
}
