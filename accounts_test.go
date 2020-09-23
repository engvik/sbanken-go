package sbanken

import (
	"context"
	"encoding/json"
	"net/http"
	"reflect"
	"testing"

	"github.com/engvik/sbanken-go/internal/transport"
)

var testAccount = Account{
	ID:          "test-account",
	Name:        "My account",
	Type:        "Account",
	Number:      "123456789",
	Available:   123.45,
	Balance:     123.45,
	CreditLimit: 0.0,
}

func testListAccountsEndpointResponse() ([]byte, int, error) {
	d := struct {
		Accounts []Account `json:"items"`
	}{
		[]Account{testAccount},
	}

	b, err := json.Marshal(d)
	if err != nil {
		return nil, 0, err
	}

	return b, http.StatusOK, nil
}

func testReadAccountEndpointResponse(shouldFail bool) ([]byte, int, error) {
	d := struct {
		Account Account `json:"item"`
		transport.HTTPResponse
	}{
		testAccount,
		transport.HTTPResponse{},
	}

	b, err := json.Marshal(d)
	if err != nil {
		return nil, 0, err
	}

	if shouldFail {
	}

	return b, http.StatusOK, nil
}

func TestListAccounts(t *testing.T) {
	ctx := context.Background()
	c, err := newTestClient(ctx, t)
	if err != nil {
		t.Fatalf("error setting up test: %v", err)
	}

	tests := []struct {
		name   string
		exp    []Account
		expErr error
	}{
		{
			name:   "should list accounts",
			exp:    []Account{testAccount},
			expErr: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			a, err := c.ListAccounts(ctx)
			if err != nil {
				if err != tc.expErr {
					t.Errorf("unexpected error: got %s, exp %s", err, tc.expErr)
				}

				return
			}

			if !reflect.DeepEqual(a, tc.exp) {
				t.Errorf("unexpected accounts: got %v, exp %v", a, tc.exp)
			}
		})
	}
}

func TestReadAccount(t *testing.T) {
	ctx := context.Background()
	c, err := newTestClient(ctx, t)
	if err != nil {
		t.Fatalf("error setting up test: %v", err)
	}

	tests := []struct {
		name      string
		accountID string
		exp       Account
		expErr    error
	}{
		{
			name:   "should fail when no accountID is passed",
			expErr: ErrMissingAccountID,
		},
		{
			name:      "should return error when error occurs",
			accountID: "test-account",
			exp:       testAccount,
			expErr:    nil,
		},

		{
			name:      "should return account",
			accountID: "test-account",
			exp:       testAccount,
			expErr:    nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			a, err := c.ReadAccount(ctx, tc.accountID)
			if err != nil {
				if err != tc.expErr {
					t.Errorf("unexpected error: got %s, exp %s", err, tc.expErr)
				}

				return
			}

			if !reflect.DeepEqual(a, tc.exp) {
				t.Errorf("unexpected account: got %v, exp %v", a, tc.exp)
			}
		})
	}
}
