package sbanken

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Account struct {
	ID          string  `json:"accountId"`
	Name        string  `json:"name"`
	Type        string  `json:"accountType"`
	Number      string  `json:"accountNumber"`
	Available   float32 `json:"available"`
	Balance     float32 `json:"balance"`
	CreditLimit float32 `json:"creditLimit"`
}

func (c *Client) ListAccounts(ctx context.Context) ([]Account, error) {
	url := fmt.Sprintf("%s/v1/Accounts", c.baseURL)

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
		Accounts []Account `json:"items"`
	}{}

	if err := json.Unmarshal(res, &data); err != nil {
		return data.Accounts, err
	}

	return data.Accounts, nil
}

func (c *Client) ReadAccount(ctx context.Context, accountID string) (Account, error) {
	url := fmt.Sprintf("%s/v1/Accounts/%s", c.baseURL, accountID)

	res, sc, err := c.request(ctx, &httpRequest{
		method: http.MethodGet,
		url:    url,
	})
	if err != nil {
		return Account{}, err
	}

	if sc != http.StatusOK {
		return Account{}, fmt.Errorf("unexpected status code: %d", sc)
	}

	data := struct {
		Account Account `json:"item"`
	}{}

	if err := json.Unmarshal(res, &data); err != nil {
		return data.Account, err
	}

	return data.Account, nil
}
