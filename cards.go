package sbanken

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Card struct {
	ID            string `json:"cardId"`
	Number        string `json:"cardNumber"`
	ExpiryDate    string `json:"expiryDate"`
	Status        string `json:"status"`
	Type          string `json:"cardType"`
	ProductCode   string `json:"productCode"`
	VersionNumber int    `json:"cardVersionNumber"`
	AccountNumber int    `json:"accountNumber"`
}

func (c *Client) ListCards(ctx context.Context) ([]Card, error) {
	url := fmt.Sprintf("%s/v1/Cards", c.baseURL)

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
		Cards []Card `json:"items"`
	}{}

	json.Unmarshal(res, &data)

	return data.Cards, nil
}
