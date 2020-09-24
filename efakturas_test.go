package sbanken

import (
	"fmt"
	"net/url"
	"testing"
	"time"
)

func TestQueryString(t *testing.T) {
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
			u, err := tc.q.QueryString("http://example.com")
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
