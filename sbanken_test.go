package sbanken

import (
	"log"
	"testing"
)

func TestNewClient(t *testing.T) {
	c := NewClient(&Config{})
	log.Println(c)
}
