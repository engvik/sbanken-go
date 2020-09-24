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

func testListAccountsEndpointResponse(behavior string) ([]byte, int, error) {
	d := struct {
		Accounts []Account `json:"items"`
		transport.HTTPResponse
	}{
		Accounts: []Account{testAccount},
	}

	if behavior == "fail" {
		d.IsError = testHTTPResponseError.IsError
		d.ErrorCode = testHTTPResponseError.ErrorCode
		d.ErrorMessage = testHTTPResponseError.ErrorMessage
		d.ErrorType = testHTTPResponseError.ErrorType
	}

	b, err := json.Marshal(d)
	if err != nil {
		return nil, 0, err
	}

	if behavior == "fail" {
		return b, http.StatusInternalServerError, nil
	}

	return b, http.StatusOK, nil
}

func testReadAccountEndpointResponse(behavior string) ([]byte, int, error) {
	d := struct {
		Account Account `json:"item"`
		transport.HTTPResponse
	}{
		Account: testAccount,
	}

	if behavior == "fail" {
		d.IsError = testHTTPResponseError.IsError
		d.ErrorCode = testHTTPResponseError.ErrorCode
		d.ErrorMessage = testHTTPResponseError.ErrorMessage
		d.ErrorType = testHTTPResponseError.ErrorType
	}

	b, err := json.Marshal(d)
	if err != nil {
		return nil, 0, err
	}

	if behavior == "fail" {
		return b, http.StatusInternalServerError, nil
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
		name     string
		behavior string
		exp      []Account
		expErr   error
	}{
		{
			name:     "should return error when error occurs",
			behavior: "fail",
			exp:      nil,
			expErr:   getTestError("ListAccounts"),
		},
		{
			name:   "should list accounts",
			exp:    []Account{testAccount},
			expErr: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx = context.WithValue(ctx, testBehavior("test-behavior"), tc.behavior)

			a, err := c.ListAccounts(ctx)
			if err != nil {
				errStr := err.Error()
				expErrStr := tc.expErr.Error()
				if errStr != expErrStr {
					t.Errorf("unexpected error: got %s, exp %s", errStr, expErrStr)
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
		behavior  string
		exp       Account
		expErr    error
	}{
		{
			name:   "should fail when no accountID is passed",
			expErr: ErrMissingAccountID,
		},
		{
			name:      "should return error when error occurs",
			behavior:  "fail",
			accountID: "test-account",
			exp:       testAccount,
			expErr:    getTestError("ReadAccount"),
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
			ctx = context.WithValue(ctx, testBehavior("test-behavior"), tc.behavior)

			a, err := c.ReadAccount(ctx, tc.accountID)
			if err != nil {
				errStr := err.Error()
				expErrStr := tc.expErr.Error()
				if errStr != expErrStr {
					t.Errorf("unexpected error: got %s, exp %s", errStr, expErrStr)
				}

				return
			}

			if !reflect.DeepEqual(a, tc.exp) {
				t.Errorf("unexpected account: got %v, exp %v", a, tc.exp)
			}
		})
	}
}
