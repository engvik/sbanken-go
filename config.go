package sbanken

// Config represents Sbanken client config.
type Config struct {
	ClientID     string
	ClientSecret string
	CustomerID   string
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
