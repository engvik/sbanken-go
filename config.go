package sbanken

// Config represents Sbanken client config.
type Config struct {
	// ClientID is required.
	ClientID string
	// ClientSecret is required.
	ClientSecret string
	// CustomerID is required.
	CustomerID string
	// UserAgent is for optionally setting a custom user agent.
	UserAgent string
	skipAuth  bool
}

func (c *Config) validate() error {
	if c.ClientID == "" {
		return ErrMissingClientID
	}

	if c.ClientSecret == "" {
		return ErrMissingClientSecret
	}

	if c.CustomerID == "" {
		return ErrMissingCustomerID
	}

	return nil
}
