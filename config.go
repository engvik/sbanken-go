package sbanken

import "log"

// Config represents Sbanken client config.
type Config struct {
	// ClientID is required.
	ClientID string
	// ClientSecret is required.
	ClientSecret string
	// UserAgent is for optionally setting a custom user agent.
	CustomerID string
	UserAgent  string
	skipAuth   bool
}

func (c *Config) validate() error {
	if c.ClientID == "" {
		return ErrMissingClientID
	}

	if c.ClientSecret == "" {
		return ErrMissingClientSecret
	}

	if c.CustomerID != "" {
		log.Println("Customer ID is deprecated.")
	}

	return nil
}
