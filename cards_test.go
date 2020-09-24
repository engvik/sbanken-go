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

var testCard = Card{
	ID:            "test-card",
	Number:        "123456789",
	ExpiryDate:    time.Now().String(),
	Status:        "status",
	Type:          "type",
	ProductCode:   "code",
	VersionNumber: 5,
	AccountNumber: 987654321,
}

func testListCardsEndpointResponse(behavior string) ([]byte, int, error) {
	d := struct {
		Cards []Card `json:"items"`
		transport.HTTPResponse
	}{
		Cards: []Card{testCard},
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

func TestListCards(t *testing.T) {
	ctx := context.Background()
	c, err := newTestClient(ctx, t)
	if err != nil {
		t.Fatalf("error setting up test: %v", err)
	}

	tests := []struct {
		name     string
		behavior string
		exp      []Card
		expErr   error
	}{
		{
			name:     "should return error when error occurs",
			behavior: "fail",
			exp:      nil,
			expErr:   getTestError("ListCards"),
		},
		{
			name:   "should list cards",
			exp:    []Card{testCard},
			expErr: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx = context.WithValue(ctx, testBehavior("test-behavior"), tc.behavior)

			a, err := c.ListCards(ctx)
			if err != nil {
				errStr := err.Error()
				expErrStr := tc.expErr.Error()
				if errStr != expErrStr {
					t.Errorf("unexpected error: got %s, exp %s", errStr, expErrStr)
				}

				return
			}

			if !reflect.DeepEqual(a, tc.exp) {
				t.Errorf("unexpected cards: got %v, exp %v", a, tc.exp)
			}
		})
	}
}
