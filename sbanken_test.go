package sbanken

import (
	"context"
	"testing"

	"github.com/engvik/sbanken-go/internal/transport"
)

var (
	testListAccountsEndpoint    = "https://api.sbanken.no/exec.bank/api/v1/Accounts"
	testReadAccountEndpoint     = "https://api.sbanken.no/exec.bank/api/v1/Accounts/test-account"
	testReadAccountFailEndpoint = "https://api.sbanken.no/exec.bank/api/v1/Accounts/fail-account"
)

type testTransportClient struct{}

func (c testTransportClient) Authorize(ctx context.Context) error {
	return nil
}

func (c testTransportClient) Request(ctx context.Context, r *transport.HTTPRequest) ([]byte, int, error) {
	switch r.URL {
	case testListAccountsEndpoint:
		return testListAccountsEndpointResponse()
	case testReadAccountEndpoint:
		return testReadAccountEndpointResponse(false)
	case testReadAccountFailEndpoint:
		return testReadAccountEndpointResponse(true)
	}

	return nil, 0, nil
}

var testHTTPResponseError = transport.HTTPResponse{
	IsError:      true,
	ErrorCode:    100,
	ErrorMessage: "an error occured",
	ErrorType:    "Error",
}

func getTestError(str string) *Error {
	return &Error{
		Code:        100,
		StatusCode:  500,
		ErrorString: str,
		Message:     "an error occured",
		Type:        "Error",
	}
}

func newTestClient(ctx context.Context, t *testing.T) (*Client, error) {
	t.Helper()

	cfg := &Config{
		ClientID:     "some-client-id",
		ClientSecret: "some-client-secret",
		CustomerID:   "some-customer-id",
		skipAuth:     true,
	}

	c, err := NewClient(ctx, cfg, nil)
	if err != nil {
		return nil, err
	}

	c.transport = testTransportClient{}

	return c, err

}

func TestNewClient(t *testing.T) {
	ctx := context.Background()
	c, err := newTestClient(ctx, t)
	if err != nil {
		t.Fatalf("error setting up test: %v", err)
	}

	t.Run("should have baseURL set", func(t *testing.T) {
		exp := "https://api.sbanken.no/exec.bank/api"
		if c.baseURL != exp {
			t.Errorf("unexpected baseURL: got %s, exp %s", c.baseURL, exp)
		}
	})

	t.Run("should have transport set", func(t *testing.T) {
		if c.transport == nil {
			t.Errorf("expected transport to be set")
		}
	})
}
