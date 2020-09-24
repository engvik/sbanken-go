package sbanken

import (
	"context"
	"testing"

	"github.com/engvik/sbanken-go/internal/transport"
)

var (
	testListAccountsEndpoint          = "https://api.sbanken.no/exec.bank/api/v1/Accounts"
	testReadAccountEndpoint           = "https://api.sbanken.no/exec.bank/api/v1/Accounts/test-account"
	testListCardsEndpoint             = "https://api.sbanken.no/exec.bank/api/v1/Cards"
	testListEfakturasEndpoint         = "https://api.sbanken.no/exec.bank/api/v1/Efakturas"
	testListEfakturasQueryEndpoint    = "https://api.sbanken.no/exec.bank/api/v1/Efakturas?index=1"
	testPayEfakturaEndpoint           = "https://api.sbanken.no/exec.bank/api/v1/Efakturas"
	testListNewEfakturasEndpoint      = "https://api.sbanken.no/exec.bank/api/v1/Efakturas/new"
	testListNewEfakturasQueryEndpoint = "https://api.sbanken.no/exec.bank/api/v1/Efakturas/new?index=1"
)

type testBehavior string

type testTransportClient struct{}

func (c testTransportClient) Authorize(ctx context.Context) error {
	return nil
}

func (c testTransportClient) Request(ctx context.Context, r *transport.HTTPRequest) ([]byte, int, error) {
	switch r.URL {
	case testListAccountsEndpoint:
		return testListAccountsEndpointResponse(getTestBehavior(ctx))
	case testReadAccountEndpoint:
		return testReadAccountEndpointResponse(getTestBehavior(ctx))
	case testListCardsEndpoint:
		return testListCardsEndpointResponse(getTestBehavior(ctx))
	case testListEfakturasEndpoint:
		fallthrough
	case testListEfakturasQueryEndpoint:
		fallthrough
	case testPayEfakturaEndpoint:
		fallthrough
	case testListNewEfakturasEndpoint:
		fallthrough
	case testListNewEfakturasQueryEndpoint:
		return testEfakturasResponses(getTestBehavior(ctx))
	default:
		return nil, 0, nil
	}
}

func getTestBehavior(ctx context.Context) string {
	if v := ctx.Value(testBehavior("test-behavior")); v != nil {
		return v.(string)
	}

	return ""
}

var testHTTPResponseError = transport.HTTPResponse{
	IsError:      true,
	ErrorCode:    100,
	ErrorMessage: "an error occured",
	ErrorType:    "Error",
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
