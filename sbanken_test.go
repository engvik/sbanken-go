package sbanken

import (
	"context"
	"testing"
)

func TestNewClient(t *testing.T) {
	cfg := &Config{
		ClientID:     "some-client-id",
		ClientSecret: "some-client-secret",
		CustomerID:   "some-customer-id",
		skipAuth:     true,
	}

	ctx := context.Background()

	c, err := NewClient(ctx, cfg, nil)
	if err != nil {
		t.Fatalf("error setting up test: %v", err)
	}

	t.Run("should have baseURL set", func(t *testing.T) {
		exp := "https://api.sbanken.no/exec.bank/api"
		if c.baseURL != exp {
			t.Errorf("unexpected baseURL: got %s, exp %s", c.baseURL, exp)
		}
	})

	t.Run("should have transport set", func(t *testing.T) {
		if c.transport == nil {
			t.Errorf("expected transport to be set")
		}
	})
}
