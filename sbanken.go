package sbanken

import (
	"context"
	"fmt"
	"net/http"

	"github.com/engvik/sbanken-go/internal/transport"
)

type transportClient interface {
	Request(context.Context, *transport.HTTPRequest) ([]byte, int, error)
}

type Client struct {
	baseURL   string
	transport transportClient
}

func NewClient(ctx context.Context, cfg *Config, httpClient *http.Client) (*Client, error) {
	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}

	tCfg := &transport.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		CustomerID:   cfg.CustomerID,
	}
	t, err := transport.New(ctx, tCfg, httpClient)
	if err != nil {
		return nil, fmt.Errorf("NewClient: %w", err)
	}

	c := &Client{
		baseURL:   "https://api.sbanken.no/exec.bank/api",
		transport: t,
	}

	return c, nil
}
