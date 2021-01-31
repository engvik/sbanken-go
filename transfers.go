package sbanken

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/engvik/sbanken-go/internal/transport"
)

// TransferQuery represents the query for transferring between accounts.
type TransferQuery struct {
	FromAccountID string  `json:"fromAccountId"`
	ToAccountID   string  `json:"toAccountId"`
	Message       string  `json:"message"`
	Amount        float32 `json:"amount"`
}

// Transfer  executes a transfer between two accounts.
// Sbanken API documentation: https://api.sbanken.no/exec.bank/swagger/index.html?urls.primaryName=Transfers%20v1
func (c *Client) Transfer(ctx context.Context, q *TransferQuery) error {
	if q == nil {
		return ErrMissingTransferQuery
	}

	payload, err := json.Marshal(q)
	if err != nil {
		return fmt.Errorf("Marshal: %w", err)
	}

	url := fmt.Sprintf("%s/v1/Transfers", c.bankBaseURL)

	res, sc, err := c.transport.Request(ctx, &transport.HTTPRequest{
		Method:      http.MethodPost,
		URL:         url,
		PostPayload: payload,
	})
	if err != nil {
		return fmt.Errorf("request: %w", err)
	}

	var data transport.HTTPResponse
	if err := json.Unmarshal(res, &data); err != nil {
		return fmt.Errorf("Unmarshal: %w", err)
	}

	if data.IsError || sc != http.StatusOK {
		return &Error{
			"Transfer",
			data.ErrorType,
			data.ErrorMessage,
			data.ErrorCode,
			sc,
		}
	}

	return nil
}
