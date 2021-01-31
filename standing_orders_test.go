package sbanken

import (
	"context"
	"encoding/json"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/engvik/sbanken-go/internal/transport"
)

var testStandingOrder = StandingOrder{
	BeneficiaryName:        "name nameson",
	CID:                    "987654321",
	CreditAccountNumber:    "11111111111",
	DebitAccountNumber:     "2222222222",
	Frequency:              "monthly",
	LastPaymentDate:        time.Now().String(),
	NextDueDate:            time.Now().String(),
	StandingOrderEndDate:   time.Now().String(),
	StandingOrderStartDate: time.Now().String(),
	StandingOrderType:      "type",
	Amount:                 1337.00,
	StandingOrderID:        19,
}

func testListStandingOrdersEndpointResponse(behavior string) ([]byte, int, error) {
	d := struct {
		StandingOrders []StandingOrder `json:"items"`
		transport.HTTPResponse
	}{
		StandingOrders: []StandingOrder{testStandingOrder},
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

func TestListStandingOrders(t *testing.T) {
	ctx := context.Background()
	c, err := newTestClient(ctx, t)
	if err != nil {
		t.Fatalf("error setting up test: %v", err)
	}

	tests := []struct {
		name      string
		accountID string
		behavior  string
		exp       []StandingOrder
		expErr    error
	}{
		{
			name:   "should fail when no accountID is passed",
			expErr: ErrMissingAccountID,
		},
		{
			name:      "should return error when error occurs",
			accountID: "test-account",
			behavior:  "fail",
			exp:       nil,
			expErr:    getTestError("ListStandingOrders"),
		},
		{
			name:      "should list standing orders",
			accountID: "test-account",
			exp:       []StandingOrder{testStandingOrder},
			expErr:    nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx = context.WithValue(ctx, testBehavior("test-behavior"), tc.behavior)

			a, err := c.ListStandingOrders(ctx, tc.accountID)
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
