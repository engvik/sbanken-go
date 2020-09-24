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

var testEfaktura = Efaktura{
	ID:                  "test-efaktura",
	IssuerID:            "issuer",
	Reference:           "ref",
	DocumentType:        "doctype",
	Status:              "NEW",
	KID:                 "000098765432123456789",
	OriginalDueDate:     time.Now().String(),
	NotificationDate:    time.Now().String(),
	IssuerName:          "Hello",
	OriginalAmount:      133.33,
	MinimumAmount:       100.00,
	CreditAccountNumber: 998877665544332211,
}

func testListPayEfakturasResponses(behavior string) ([]byte, int, error) {
	if behavior == "pay" {
		pay := transport.HTTPResponse{
			IsError: false,
		}

		b, err := json.Marshal(pay)
		if err != nil {
			return nil, 0, err
		}

		return b, http.StatusOK, nil
	}

	d := struct {
		Efakturas []Efaktura `json:"items"`
		transport.HTTPResponse
	}{
		Efakturas: []Efaktura{testEfaktura},
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

func testReadEfakturaResponse(behavior string) ([]byte, int, error) {
	d := struct {
		Efaktura Efaktura `json:"item"`
		transport.HTTPResponse
	}{
		Efaktura: testEfaktura,
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

func TestEfakturaQueryString(t *testing.T) {
	timestamp := time.Now()
	qsTimestamp := url.QueryEscape(timestamp.String())

	tests := []struct {
		name string
		q    *EfakturaListQuery
		exp  string
	}{
		{
			name: "should create correct query string",
			q:    &EfakturaListQuery{},
			exp:  "",
		},
		{
			name: "should create correct query string",
			q:    &EfakturaListQuery{StartDate: timestamp},
			exp:  fmt.Sprintf("startDate=%s", qsTimestamp),
		},
		{
			name: "should create correct query string",
			q: &EfakturaListQuery{
				StartDate: timestamp,
				EndDate:   timestamp,
			},
			exp: fmt.Sprintf("endDate=%s&startDate=%s", qsTimestamp, qsTimestamp),
		},
		{
			name: "should create correct query string",
			q:    &EfakturaListQuery{Status: "test"},
			exp:  "status=test",
		},
		{
			name: "should create correct query string",
			q:    &EfakturaListQuery{Index: "1", Length: "5"},
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

func TestListEfakturas(t *testing.T) {
	ctx := context.Background()
	c, err := newTestClient(ctx, t)
	if err != nil {
		t.Fatalf("error setting up test: %v", err)
	}

	tests := []struct {
		name     string
		q        *EfakturaListQuery
		behavior string
		exp      []Efaktura
		expErr   error
	}{
		{
			name:     "should return error when error occurs",
			q:        nil,
			behavior: "fail",
			exp:      nil,
			expErr:   getTestError("ListEfakturas"),
		},
		{
			name:   "should list efakturas without query",
			q:      nil,
			exp:    []Efaktura{testEfaktura},
			expErr: nil,
		},
		{
			name:   "should list efakturas with query",
			q:      &EfakturaListQuery{Index: "1"},
			exp:    []Efaktura{testEfaktura},
			expErr: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx = context.WithValue(ctx, testBehavior("test-behavior"), tc.behavior)

			a, err := c.ListEfakturas(ctx, tc.q)
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

func TestPayEfaktura(t *testing.T) {
	ctx := context.Background()
	c, err := newTestClient(ctx, t)
	if err != nil {
		t.Fatalf("error setting up test: %v", err)
	}

	tests := []struct {
		name     string
		q        *EfakturaPayQuery
		behavior string
		exp      error
	}{
		{
			name:     "should fail if no efaktura pay query",
			q:        nil,
			behavior: "pay",
			exp:      ErrMissingEfakturaPayQuery,
		},
		{
			name:     "should return error when error occurs",
			q:        &EfakturaPayQuery{ID: "efaktura-id", AccountID: "account-id"},
			behavior: "fail",
			exp:      getTestError("PayEfaktura"),
		},
		{
			name:     "should pay efaktura",
			q:        &EfakturaPayQuery{ID: "efaktura-id", AccountID: "account-id"},
			behavior: "pay",
			exp:      nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx = context.WithValue(ctx, testBehavior("test-behavior"), tc.behavior)

			err := c.PayEfaktura(ctx, tc.q)
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

func TestListNewEfakturas(t *testing.T) {
	ctx := context.Background()
	c, err := newTestClient(ctx, t)
	if err != nil {
		t.Fatalf("error setting up test: %v", err)
	}

	tests := []struct {
		name     string
		q        *EfakturaListQuery
		behavior string
		exp      []Efaktura
		expErr   error
	}{
		{
			name:     "should return error when error occurs",
			q:        nil,
			behavior: "fail",
			exp:      nil,
			expErr:   getTestError("ListNewEfakturas"),
		},
		{
			name:   "should fail with invalid query StartDate",
			q:      &EfakturaListQuery{StartDate: time.Now()},
			exp:    []Efaktura{testEfaktura},
			expErr: ErrNotValidOptionStartDate,
		},
		{
			name:   "should fail with invalid query EndDate",
			q:      &EfakturaListQuery{EndDate: time.Now()},
			exp:    []Efaktura{testEfaktura},
			expErr: ErrNotValidOptionEndDate,
		},
		{
			name:   "should fail with invalid query Status",
			q:      &EfakturaListQuery{Status: "NEW"},
			exp:    []Efaktura{testEfaktura},
			expErr: ErrNotValidOptionStatus,
		},
		{
			name:   "should list efakturas without query",
			q:      nil,
			exp:    []Efaktura{testEfaktura},
			expErr: nil,
		},
		{
			name:   "should list efaktura with valid query",
			q:      &EfakturaListQuery{Index: "1"},
			exp:    []Efaktura{testEfaktura},
			expErr: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx = context.WithValue(ctx, testBehavior("test-behavior"), tc.behavior)

			a, err := c.ListNewEfakturas(ctx, tc.q)
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

func TestReadEfaktura(t *testing.T) {
	ctx := context.Background()
	c, err := newTestClient(ctx, t)
	if err != nil {
		t.Fatalf("error setting up test: %v", err)
	}

	tests := []struct {
		name       string
		efakturaID string
		behavior   string
		exp        Efaktura
		expErr     error
	}{
		{
			name:   "should fail when no efakturaID is passed",
			expErr: ErrMissingEfakturaID,
		},
		{
			name:       "should return error when error occurs",
			efakturaID: "test-efaktura",
			behavior:   "fail",
			exp:        testEfaktura,
			expErr:     getTestError("ReadEfaktura"),
		},
		{
			name:       "should return efaktura",
			efakturaID: "test-efaktura",
			exp:        testEfaktura,
			expErr:     nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx = context.WithValue(ctx, testBehavior("test-behavior"), tc.behavior)

			a, err := c.ReadEfaktura(ctx, tc.efakturaID)
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
