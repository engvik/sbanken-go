package sbanken

import "testing"

func TestAddressString(t *testing.T) {
	tests := []struct {
		name string
		addr Address
		exp  string
	}{
		{
			name: "should handle empty address struct",
			addr: Address{},
			exp:  "",
		},
		{
			name: "should handle partly filled address struct",
			addr: Address{
				AddressLine1: "Testerstreet 1",
			},
			exp: "Testerstreet 1",
		},
		{
			name: "should handle partly filled address struct",
			addr: Address{
				AddressLine2: "Testerstreet 2",
			},
			exp: "Testerstreet 2",
		},
		{
			name: "should handle partly filled address struct",
			addr: Address{
				AddressLine1: "Testerstreet 1",
				Country:      "Norway",
				ZipCode:      "1337",
				City:         "Sandvika",
			},
			exp: "Testerstreet 1, 1337 Sandvika, Norway",
		},
		{
			name: "should handle partly filled address struct",
			addr: Address{
				AddressLine1: "Testerstreet 1",
				Country:      "Norway",
				City:         "Sandvika",
			},
			exp: "Testerstreet 1, Sandvika, Norway",
		},

		{
			name: "should handle completly filled address struct",
			addr: Address{
				AddressLine1: "c/o Test Testesen",
				AddressLine2: "Testerstreet 1",
				AddressLine3: "PO 13371337",
				Country:      "Norway",
				ZipCode:      "1337",
				City:         "Sandvika",
			},
			exp: "c/o Test Testesen, Testerstreet 1, PO 13371337, 1337 Sandvika, Norway",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			addr := tc.addr.String()
			if addr != tc.exp {
				t.Errorf("unexpected address string: got %s, exp %s", addr, tc.exp)
			}
		})
	}
}
