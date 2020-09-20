package sbanken

type ConfigError struct {
	ErrorString string
}

func (ce *ConfigError) Error() string {
	return ce.ErrorString
}

var (
	ErrNoClientID    = &ConfigError{"No client ID"}
	ErNoClientSecret = &ConfigError{"No client secret"}
)

type Config struct {
	ClientID     string
	ClientSecret string
}

func (c *Config) validate() error {
	if c.ClientID == "" {
		return ErrNoClientID
	}

	if c.ClientSecret == "" {
		return ErNoClientSecret
	}

	return nil
}
