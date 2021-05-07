package sbanken

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/engvik/sbanken-go/internal/transport"
)

// Card represents a card.
// Sbanken API documentation: https://api.sbanken.no/exec.bank/swagger/index.html?urls.primaryName=Cards%20v1
type Card struct {
	ID            string `json:"cardId"`
	Number        string `json:"cardNumber"`
	ExpiryDate    string `json:"expiryDate"`
	Status        string `json:"status"`
	Type          string `json:"cardType"`
	ProductCode   string `json:"productCode"`
	AccountNumber string `json:"accountNumber"`
	AccountOwner  string `json:"accountOwner"`
	CustomerID    string `json:"customerId"`
	VersionNumber int    `json:"cardVersionNumber"`
}

// ListCards lists the cards.
func (c *Client) ListCards(ctx context.Context) ([]Card, error) {
	url := fmt.Sprintf("%s/v1/Cards", c.bankBaseURL)

	res, sc, err := c.transport.Request(ctx, &transport.HTTPRequest{
		Method: http.MethodGet,
		URL:    url,
	})
	if err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}

	data := struct {
		Cards []Card `json:"items"`
		transport.HTTPResponse
	}{}

	if err := json.Unmarshal(res, &data); err != nil {
		return data.Cards, fmt.Errorf("Unmarshal: %w", err)
	}

	if data.IsError || sc != http.StatusOK {
		return data.Cards, &Error{
			"ListCards",
			data.ErrorType,
			data.ErrorMessage,
			data.ErrorCode,
			sc,
		}
	}

	return data.Cards, nil
}
