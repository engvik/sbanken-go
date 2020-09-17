package sbanken

import (
	"log"
	"testing"
)

func TestNewClient(t *testing.T) {
	c, err := NewClient(&Config{}, nil)
	if err != nil {
		t.Fatalf("error setting up test: %v", err)
	}
	log.Println(c)
}
