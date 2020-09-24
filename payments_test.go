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

var testPayment = Payment{
	ID:                     "test-payment",
	RecipientAccountNumber: "987654321",
	DueDate:                time.Now().String(),
	KID:                    "00000123456799",
	Text:                   "Hello, yes, this is Payment!",
	Status:                 "status",
	AllowedNewStatusTypes:  []string{"new-status"},
	StatusDetails:          "details",
	ProductType:            "product",
	PaymentType:            "payment",
	BeneficiaryName:        "name nameson",
	Amount:                 1337.00,
	PaymentNumber:          4,
	IsActive:               true,
}

func testListPaymentResponses(behavior string) ([]byte, int, error) {
	d := struct {
		Payments []Payment `json:"items"`
		transport.HTTPResponse
	}{
		Payments: []Payment{testPayment},
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

func testReadPaymentResponse(behavior string) ([]byte, int, error) {
	d := struct {
		Payment Payment `json:"item"`
		transport.HTTPResponse
	}{
		Payment: testPayment,
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

func TestPaymentQueryString(t *testing.T) {
	tests := []struct {
		name string
		q    *PaymentListQuery
		exp  string
	}{
		{
			name: "should create correct query string",
			q:    &PaymentListQuery{},
			exp:  "",
		},
		{
			name: "should create correct query string",
			q:    &PaymentListQuery{Index: "1", Length: "5"},
			exp:  "index=1&length=5",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u, err := tc.q.QueryString("https://example.com")
			if err != nil {
				t.Errorf("unexpected error: got %s", err)

				return
			}

			if u != tc.exp {
				t.Errorf("unexpected query string: got %s, exp %s", u, tc.exp)
			}
		})
	}
}

func TestListPayments(t *testing.T) {
	ctx := context.Background()
	c, err := newTestClient(ctx, t)
	if err != nil {
		t.Fatalf("error setting up test: %v", err)
	}

	tests := []struct {
		name      string
		accountID string
		q         *PaymentListQuery
		behavior  string
		exp       []Payment
		expErr    error
	}{
		{
			name:   "should fail when no accountID is passed",
			exp:    nil,
			expErr: ErrMissingAccountID,
		},
		{
			name:      "should return error when error occurs",
			accountID: "test-account",
			q:         nil,
			behavior:  "fail",
			exp:       nil,
			expErr:    getTestError("ListPayments"),
		},
		{
			name:      "should list payments without query",
			accountID: "test-account",
			q:         nil,
			exp:       []Payment{testPayment},
			expErr:    nil,
		},
		{
			name:      "should list payments with query",
			accountID: "test-account",
			q:         &PaymentListQuery{Index: "1"},
			exp:       []Payment{testPayment},
			expErr:    nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx = context.WithValue(ctx, testBehavior("test-behavior"), tc.behavior)

			a, err := c.ListPayments(ctx, tc.accountID, tc.q)
			if err != nil {
				errStr := err.Error()
				expErrStr := tc.expErr.Error()
				if errStr != expErrStr {
					t.Errorf("unexpected error: got %s, exp %s", errStr, expErrStr)
				}

				return
			}

			if !reflect.DeepEqual(a, tc.exp) {
				t.Errorf("unexpected efaktura: got %v, exp %v", a, tc.exp)
			}
		})
	}
}

func TestReadPayment(t *testing.T) {
	ctx := context.Background()
	c, err := newTestClient(ctx, t)
	if err != nil {
		t.Fatalf("error setting up test: %v", err)
	}

	tests := []struct {
		name      string
		accountID string
		paymentID string
		behavior  string
		exp       Payment
		expErr    error
	}{
		{
			name:   "should fail when no accountID is passed",
			expErr: ErrMissingAccountID,
		},
		{
			name:      "should fail when no paymentID is passed",
			accountID: "test-account",
			expErr:    ErrMissingPaymentID,
		},
		{
			name:      "should return error when error occurs",
			accountID: "test-account",
			paymentID: "test-payment",
			behavior:  "fail",
			exp:       testPayment,
			expErr:    getTestError("ReadPayment"),
		},
		{
			name:      "should return efaktura",
			accountID: "test-account",
			paymentID: "test-payment",
			exp:       testPayment,
			expErr:    nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx = context.WithValue(ctx, testBehavior("test-behavior"), tc.behavior)

			a, err := c.ReadPayment(ctx, tc.accountID, tc.paymentID)
			if err != nil {
				errStr := err.Error()
				expErrStr := tc.expErr.Error()
				if errStr != expErrStr {
					t.Errorf("unexpected error: got %s, exp %s", errStr, expErrStr)
				}

				return
			}

			if !reflect.DeepEqual(a, tc.exp) {
				t.Errorf("unexpected payment: got %v, exp %v", a, tc.exp)
			}
		})
	}
}
