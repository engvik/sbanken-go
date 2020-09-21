package sbanken

type ConfigError struct {
	ErrorString string
}

func (ce *ConfigError) Error() string {
	return ce.ErrorString
}

var (
	ErrNoClientID     = &ConfigError{"No client ID"}
	ErrNoClientSecret = &ConfigError{"No client secret"}
	ErrNoCustomerID   = &ConfigError{"No customer ID"}
)

type Config struct {
	ClientID     string
	ClientSecret string
	CustomerID   string
}

func (c *Config) validate() error {
	if c.ClientID == "" {
		return ErrNoClientID
	}

	if c.ClientSecret == "" {
		return ErrNoClientSecret
	}

	if c.CustomerID == "" {
		return ErrNoCustomerID
	}

	return nil
}
