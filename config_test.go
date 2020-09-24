package sbanken

import "testing"

func TestValidate(t *testing.T) {
	tests := []struct {
		name string
		cfg  *Config
		exp  error
	}{
		{
			name: "should not validate when ClientID is missing",
			cfg:  &Config{},
			exp:  ErrMissingClientID,
		},
		{
			name: "should not validate when ClientSecret is missing",
			cfg:  &Config{ClientID: "client-id"},
			exp:  ErrMissingClientSecret,
		},
		{
			name: "should not validate when CustomerID is missing",
			cfg:  &Config{ClientID: "client-id", ClientSecret: "client-secret"},
			exp:  ErrMissingCustomerID,
		},
		{
			name: "should validate when all required parameters are set",
			cfg:  &Config{ClientID: "client-id", ClientSecret: "client-secret", CustomerID: "customer-id"},
			exp:  nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.cfg.validate(); err != tc.exp {
				t.Errorf("unexpected result: got %s, exp %s", err, tc.exp)
			}
		})
	}
}
