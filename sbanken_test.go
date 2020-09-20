package sbanken

import (
	"log"
	"os"
	"testing"
)

func TestNewClient(t *testing.T) {
	cfg := &Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
	}

	c, err := NewClient(cfg, nil)
	if err != nil {
		t.Fatalf("error setting up test: %v", err)
	}

	log.Println(c.getToken())
}
