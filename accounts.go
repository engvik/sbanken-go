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
		return nil, fmt.Errorf("request: %w", err)
	}

	data := struct {
		Accounts []Account `json:"items"`
		httpResponse
	}{}

	if err := json.Unmarshal(res, &data); err != nil {
		return data.Accounts, fmt.Errorf("Unmarshal: %w", err)
	}

	if data.IsError || sc != http.StatusOK {
		return nil, &Error{
			"ListAccounts",
			data.ErrorType,
			data.ErrorMessage,
			data.ErrorCode,
			sc,
		}
	}

	return data.Accounts, nil
}

func (c *Client) ReadAccount(ctx context.Context, accountID string) (Account, error) {
	if accountID == "" {
		return Account{}, ErrMissingAccountID
	}

	url := fmt.Sprintf("%s/v1/Accounts/%s", c.baseURL, accountID)

	res, sc, err := c.request(ctx, &httpRequest{
		method: http.MethodGet,
		url:    url,
	})
	if err != nil {
		return Account{}, fmt.Errorf("request: %w", err)
	}

	data := struct {
		Account Account `json:"item"`
		httpResponse
	}{}

	if err := json.Unmarshal(res, &data); err != nil {
		return data.Account, fmt.Errorf("Unmarshal: %w", err)
	}

	if data.IsError || sc != http.StatusOK {
		return data.Account, &Error{
			"ReadAccount",
			data.ErrorType,
			data.ErrorMessage,
			data.ErrorCode,
			sc,
		}
	}

	return data.Account, nil
}
