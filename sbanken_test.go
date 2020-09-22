package sbanken

import (
	"context"
	"log"
	"os"
	"testing"
)

func TestNewClient(t *testing.T) {
	cfg := &Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		CustomerID:   os.Getenv("CUSTOMER_ID"),
	}

	ctx := context.Background()

	c, err := NewClient(ctx, cfg, nil)
	if err != nil {
		t.Fatalf("error setting up test: %v", err)
	}

	a, err := c.ListAccounts(ctx)
	if err != nil {
		t.Fatal(err)
	}

	log.Printf("%+v", a)

	e, err := c.ListStandingOrders(ctx, a[0].ID)
	if err != nil {
		t.Fatal(err)
	}

	log.Printf("%+v", e)
}
