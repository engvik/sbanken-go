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

	data := struct {
		Cards []Card `json:"items"`
		httpResponse
	}{}

	if err := json.Unmarshal(res, &data); err != nil {
		return data.Cards, err
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
