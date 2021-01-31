// Package sbanken provides an easy way to work with the Sbanken API.
//
// Sbanken API information: https://sbanken.no/bruke/utviklerportalen/
// Sbanken API documentation: https://api.sbanken.no/exec.bank/swagger/index.html
package sbanken

import (
	"context"
	"fmt"
	"net/http"

	"github.com/engvik/sbanken-go/internal/transport"
)

type transportClient interface {
	Authorize(context.Context) error
	Request(context.Context, *transport.HTTPRequest) ([]byte, int, error)
}

// Client represents an Sbanken client.
type Client struct {
	bankBaseURL      string
	customersBaseURL string
	transport        transportClient
}

// NewClient returns a new Sbanken client. If httpClient is nil, http.DefaultClient will be used.
func NewClient(ctx context.Context, cfg *Config, httpClient *http.Client) (*Client, error) {
	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}

	tCfg := &transport.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		CustomerID:   cfg.CustomerID,
	}

	c := &Client{
		bankBaseURL:      "https://api.sbanken.no/exec.bank/api",
		customersBaseURL: "https://api.sbanken.no/exec.customers/api",
		transport:        transport.New(ctx, tCfg, httpClient),
	}

	if !cfg.skipAuth {
		if err := c.transport.Authorize(ctx); err != nil {
			return c, fmt.Errorf("Authorize: %w", err)
		}
	}

	return c, nil
}
