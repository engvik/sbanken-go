package sbanken

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/engvik/sbanken-go/internal/transport"
)

// Account represents an account.
// Sbanken API documentation: https://api.sbanken.no/exec.bank/swagger/index.html?urls.primaryName=Accounts%20v1
type Account struct {
	ID              string  `json:"accountId"`
	Name            string  `json:"name"`
	Type            string  `json:"accountType"`
	Number          string  `json:"accountNumber"`
	OwnerCustomerID string  `json:"ownerCustomerId"`
	Available       float32 `json:"available"`
	Balance         float32 `json:"balance"`
	CreditLimit     float32 `json:"creditLimit"`
}

// ListAccounts lists the accounts.
func (c *Client) ListAccounts(ctx context.Context) ([]Account, error) {
	url := fmt.Sprintf("%s/v1/Accounts", c.bankBaseURL)

	res, sc, err := c.transport.Request(ctx, &transport.HTTPRequest{
		Method: http.MethodGet,
		URL:    url,
	})
	if err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}

	data := struct {
		Accounts []Account `json:"items"`
		transport.HTTPResponse
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

// ReadAccount reads an account. The accountID are required.
func (c *Client) ReadAccount(ctx context.Context, accountID string) (Account, error) {
	if accountID == "" {
		return Account{}, ErrMissingAccountID
	}

	url := fmt.Sprintf("%s/v1/Accounts/%s", c.bankBaseURL, accountID)

	res, sc, err := c.transport.Request(ctx, &transport.HTTPRequest{
		Method: http.MethodGet,
		URL:    url,
	})
	if err != nil {
		return Account{}, fmt.Errorf("request: %w", err)
	}

	data := struct {
		Account Account `json:"item"`
		transport.HTTPResponse
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
