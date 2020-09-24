package sbanken

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/engvik/sbanken-go/internal/transport"
)

var testTransaction = Transaction{
	TransactionDetails: TransactionDetails{
		ID:                     "abc1",
		FormattedAccountNumber: "4987147291",
		CID:                    "fjsk39",
		AmountDescription:      "amount-desc",
		ReceiverName:           "name nameson",
		PayerName:              "pay payson",
		RegistrationDate:       time.Now().String(),
		NumericReference:       15,
	},
	AccountingDate:              time.Now().String(),
	InterestDate:                time.Now().String(),
	OtherAccountNumber:          "123141423",
	Text:                        "transaction",
	TransactionType:             "transaction",
	TransactionTypeText:         "asdf",
	ReservationType:             "reservation",
	Source:                      "source",
	Amount:                      999.99,
	IsReservation:               true,
	OtherAccountNumberSpecified: true,
	TransactionDetailSpecified:  true,
}

func testListTransactionsResponse(behavior string) ([]byte, int, error) {
	d := struct {
		Transactions []Transaction `json:"items"`
		transport.HTTPResponse
	}{
		Transactions: []Transaction{testTransaction},
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

func TestTransactionQueryString(t *testing.T) {
	timestamp := time.Now()
	qsTimestamp := url.QueryEscape(timestamp.String())

	tests := []struct {
		name string
		q    *TransactionListQuery
		exp  string
	}{
		{
			name: "should create correct query string",
			q:    &TransactionListQuery{},
			exp:  "",
		},
		{
			name: "should create correct query string",
			q:    &TransactionListQuery{StartDate: timestamp},
			exp:  fmt.Sprintf("startDate=%s", qsTimestamp),
		},
		{
			name: "should create correct query string",
			q: &TransactionListQuery{
				StartDate: timestamp,
				EndDate:   timestamp,
			},
			exp: fmt.Sprintf("endDate=%s&startDate=%s", qsTimestamp, qsTimestamp),
		},
		{
			name: "should create correct query string",
			q:    &TransactionListQuery{Index: "1", Length: "5"},
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

func TestListTransactions(t *testing.T) {
	ctx := context.Background()
	c, err := newTestClient(ctx, t)
	if err != nil {
		t.Fatalf("error setting up test: %v", err)
	}

	tests := []struct {
		name      string
		accountID string
		q         *TransactionListQuery
		behavior  string
		exp       []Transaction
		expErr    error
	}{
		{
			name:   "should fail when no accountID is passed",
			expErr: ErrMissingAccountID,
		},
		{
			name:      "should return error when error occurs",
			accountID: "test-account",
			q:         nil,
			behavior:  "fail",
			exp:       nil,
			expErr:    getTestError("ListTransactions"),
		},
		{
			name:      "should list transactions without query",
			accountID: "test-account",
			q:         nil,
			exp:       []Transaction{testTransaction},
			expErr:    nil,
		},
		{
			name:      "should list transactions with query",
			accountID: "test-account",
			q:         &TransactionListQuery{Index: "1"},
			exp:       []Transaction{testTransaction},
			expErr:    nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx = context.WithValue(ctx, testBehavior("test-behavior"), tc.behavior)

			a, err := c.ListTransactions(ctx, tc.accountID, tc.q)
			if err != nil {
				errStr := err.Error()
				expErrStr := tc.expErr.Error()
				if errStr != expErrStr {
					t.Errorf("unexpected error: got %s, exp %s", errStr, expErrStr)
				}

				return
			}

			if !reflect.DeepEqual(a, tc.exp) {
				t.Errorf("unexpected transaction: got %v, exp %v", a, tc.exp)
			}
		})
	}
}
