package sbanken

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/engvik/sbanken-go/internal/transport"
)

func testTransferResponse(behavior string) ([]byte, int, error) {
	d := transport.HTTPResponse{
		IsError: false,
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

func TestTransfer(t *testing.T) {
	ctx := context.Background()
	c, err := newTestClient(ctx, t)
	if err != nil {
		t.Fatalf("error setting up test: %v", err)
	}

	tests := []struct {
		name     string
		q        *TransferQuery
		behavior string
		exp      error
	}{
		{
			name:     "should fail if no transfer query",
			q:        nil,
			behavior: "pay",
			exp:      ErrMissingTransferQuery,
		},
		{
			name: "should return error when error occurs",
			q: &TransferQuery{
				FromAccountID: "133713371337",
				ToAccountID:   "leetleetleet",
				Message:       "transfer",
				Amount:        1337.13,
			},
			behavior: "fail",
			exp:      getTestError("Transfer"),
		},
		{
			name: "should transfer",
			q: &TransferQuery{
				FromAccountID: "133713371337",
				ToAccountID:   "leetleetleet",
				Message:       "transfer",
				Amount:        1337.13,
			},
			exp: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx = context.WithValue(ctx, testBehavior("test-behavior"), tc.behavior)

			err := c.Transfer(ctx, tc.q)
			if err != nil {
				errStr := err.Error()
				expStr := tc.exp.Error()
				if errStr != expStr {
					t.Errorf("unexpected error: got %s, exp %s", errStr, expStr)
				}

				return
			}
		})
	}
}
