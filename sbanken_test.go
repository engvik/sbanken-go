package sbanken

import (
	"context"
	"testing"

	"github.com/engvik/sbanken-go/internal/transport"
)

const baseURL = "https://publicapi.sbanken.no/apibeta/api/"

var (
	testListAccountsEndpoint                  = baseURL + "v1/Accounts"
	testReadAccountEndpoint                   = baseURL + "v1/Accounts/test-account"
	testListCardsEndpoint                     = baseURL + "v1/Cards"
	testListEfakturasEndpoint                 = baseURL + "v1/Efakturas"
	testListEfakturasQueryEndpoint            = baseURL + "v1/Efakturas?index=1"
	testPayEfakturaEndpoint                   = baseURL + "v1/Efakturas"
	testListNewEfakturasEndpoint              = baseURL + "v1/Efakturas/new"
	testListNewEfakturasQueryEndpoint         = baseURL + "v1/Efakturas/new?index=1"
	testReadEfakturaEndpoint                  = baseURL + "v1/Efakturas/test-efaktura"
	testListPaymentsEndpoint                  = baseURL + "v1/Payments/test-account"
	testListPaymentsQueryEndpoint             = baseURL + "v1/Payments/test-account?index=1"
	testReadPaymentsEndpoint                  = baseURL + "v1/Payments/test-account/test-payment"
	testListStandingOrdersEndpoint            = baseURL + "v1/StandingOrders/test-account"
	testListTransactionsEndpoint              = baseURL + "v1/Transactions/test-account"
	testListTransactionsQueryEndpoint         = baseURL + "v1/Transactions/test-account?index=1"
	testListArchivedTransactionsEndpoint      = baseURL + "v1/Transactions/archive/test-account"
	testListArchivedTransactionsQueryEndpoint = baseURL + "v1/Transactions/archive/test-account?index=1"
	testTransferEndpoint                      = baseURL + "v1/Transfers"
	testCustomersEndpoint                     = baseURL + "v1/Customers"
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
		return testListPayEfakturasEndpointResponses(getTestBehavior(ctx))
	case testReadEfakturaEndpoint:
		return testReadEfakturaEndpointResponse(getTestBehavior(ctx))
	case testListPaymentsEndpoint:
		fallthrough
	case testListPaymentsQueryEndpoint:
		return testListPaymentEndpointResponses(getTestBehavior(ctx))
	case testReadPaymentsEndpoint:
		return testReadPaymentEndpointResponse(getTestBehavior(ctx))
	case testListStandingOrdersEndpoint:
		return testListStandingOrdersEndpointResponse(getTestBehavior(ctx))
	case testListTransactionsEndpoint:
		fallthrough
	case testListTransactionsQueryEndpoint:
		fallthrough
	case testListArchivedTransactionsEndpoint:
		fallthrough
	case testListArchivedTransactionsQueryEndpoint:
		return testListTransactionsEndpointResponse(getTestBehavior(ctx))
	case testTransferEndpoint:
		return testTransferEndpointResponse(getTestBehavior(ctx))
	case testCustomersEndpoint:
		return testCustomersEndpointResponse(getTestBehavior(ctx))
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
	ErrorMessage: "an error occurred",
	ErrorType:    "Error",
}

func newTestClient(ctx context.Context, t *testing.T) (*Client, error) {
	t.Helper()

	cfg := &Config{
		ClientID:     "some-client-id",
		ClientSecret: "some-client-secret",
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

	t.Run("should have bankBaseURL set", func(t *testing.T) {
		exp := "https://publicapi.sbanken.no/apibeta/api"
		if c.bankBaseURL != exp {
			t.Errorf("unexpected baseURL: got %s, exp %s", c.bankBaseURL, exp)
		}
	})

	t.Run("should have transport set", func(t *testing.T) {
		if c.transport == nil {
			t.Errorf("expected transport to be set")
		}
	})
}
