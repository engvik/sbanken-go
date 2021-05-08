package sbanken

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/engvik/sbanken-go/internal/transport"
)

// StandingOrder represents a standing order (repeated transfers and payments).
// Sbanken API documentation: https://publicapi.sbanken.no/openapi/apibeta/index.html#/StandingOrders
type StandingOrder struct {
	FreeTerms              []string `json:"freeTerms"`
	BeneficiaryName        string   `json:"beneficiaryName"`
	CID                    string   `json:"cId"`
	CreditAccountNumber    string   `json:"creditAccountNumber"`
	DebitAccountNumber     string   `json:"debitAccountNumber"`
	Frequency              string   `json:"frequency"`
	LastPaymentDate        string   `json:"lastPaymentDate"`
	NextDueDate            string   `json:"nextDueDate"`
	StandingOrderEndDate   string   `json:"standingOrderEndDate"`
	StandingOrderStartDate string   `json:"standingOrderStartDate"`
	StandingOrderType      string   `json:"standingOrderType"`
	Amount                 float32  `json:"amount"`
	StandingOrderID        int      `json:"standingOrderId"`
}

// ListStandingOrders lists the standing orders for repeated transfers and payments. The accoundID are required.
func (c *Client) ListStandingOrders(ctx context.Context, accountID string) ([]StandingOrder, error) {
	if accountID == "" {
		return nil, ErrMissingAccountID
	}

	url := fmt.Sprintf("%s/v1/StandingOrders/%s", c.bankBaseURL, accountID)

	res, sc, err := c.transport.Request(ctx, &transport.HTTPRequest{
		Method: http.MethodGet,
		URL:    url,
	})
	if err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}

	data := struct {
		StandingOrders []StandingOrder `json:"items"`
		transport.HTTPResponse
	}{}

	if err := json.Unmarshal(res, &data); err != nil {
		return data.StandingOrders, fmt.Errorf("Unmarshal: %w", err)
	}

	if data.IsError || sc != http.StatusOK {
		return data.StandingOrders, &Error{
			"ListStandingOrders",
			data.ErrorType,
			data.ErrorMessage,
			data.ErrorCode,
			sc,
		}
	}

	return data.StandingOrders, nil
}
