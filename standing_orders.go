package sbanken

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

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

func (c *Client) ListStandingOrders(ctx context.Context, accountID string) ([]StandingOrder, error) {
	url := fmt.Sprintf("%s/v1/StandingOrders/%s", c.baseURL, accountID)

	res, sc, err := c.request(ctx, &httpRequest{
		method: http.MethodGet,
		url:    url,
	})
	if err != nil {
		return nil, err
	}

	if sc != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", sc)
	}

	data := struct {
		StandingOrders []StandingOrder `json:"items"`
	}{}

	if err := json.Unmarshal(res, &data); err != nil {
		return data.StandingOrders, err
	}

	return data.StandingOrders, nil
}
