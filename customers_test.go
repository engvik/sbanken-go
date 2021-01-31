package sbanken

import (
	"context"
	"encoding/json"
	"net/http"
	"reflect"
	"testing"

	"github.com/engvik/sbanken-go/internal/transport"
)

var testCustomer = Customer{
	CustomerID:   "test-customer-id",
	FirstName:    "Testy",
	LastName:     "Tester",
	EmailAddress: "testy@tester.com",
	DateOfBirth:  "2021-01-31T10:05:54.590Z",
	PostalAddress: Address{
		AddressLine1: "Tester street 1",
	},
	StreetAddress: Address{
		AddressLine1: "Tester street 1",
	},
	PhoneNumbers: []PhoneNumber{
		{
			"1",
			"1337133713371337",
		},
	},
}

func testCustomersEndpointResponse(behavior string) ([]byte, int, error) {
	d := struct {
		Customer Customer `json:"item"`
		transport.HTTPResponse
	}{
		Customer: testCustomer,
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

func TestGetCustomer(t *testing.T) {
	ctx := context.Background()
	c, err := newTestClient(ctx, t)
	if err != nil {
		t.Fatalf("error setting up test: %v", err)
	}

	tests := []struct {
		name     string
		behavior string
		exp      Customer
		expErr   error
	}{
		{
			name:     "should return error when error occurs",
			behavior: "fail",
			exp:      testCustomer,
			expErr:   getTestError("Customers"),
		},
		{
			name:   "should return customer",
			exp:    testCustomer,
			expErr: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx = context.WithValue(ctx, testBehavior("test-behavior"), tc.behavior)
			cust, err := c.GetCustomer(ctx)
			if err != nil {
				errStr := err.Error()
				expErrStr := tc.expErr.Error()
				if errStr != expErrStr {
					t.Errorf("unexpected error: got %s, exp %s", errStr, expErrStr)
				}

				return
			}

			if !reflect.DeepEqual(cust, tc.exp) {
				t.Errorf("unexpected customer: got %v, exp %v", cust, tc.exp)
			}
		})
	}
}
