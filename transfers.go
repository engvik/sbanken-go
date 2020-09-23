package sbanken

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type TransferQuery struct {
	FromAccountID string  `json:"fromAccountId"`
	ToAccountID   string  `json:"toAccoundId"`
	Message       string  `json:"message"`
	Amount        float32 `json:"amount"`
}

func (c *Client) Transfer(ctx context.Context, q *TransferQuery) error {
	if q == nil {
		return ErrMissingTransferQuery
	}

	payload, err := json.Marshal(q)
	if err != nil {
		return fmt.Errorf("Marshal: %w", err)
	}

	url := fmt.Sprintf("%s/v1/Transfers", c.baseURL)

	res, sc, err := c.request(ctx, &httpRequest{
		method:      http.MethodPost,
		url:         url,
		postPayload: payload,
	})
	if err != nil {
		return fmt.Errorf("request: %w", err)
	}

	var data httpResponse
	if err := json.Unmarshal(res, &data); err != nil {
		return fmt.Errorf("Unmarshal: %w", err)
	}

	if data.IsError || sc != http.StatusOK {
		return &Error{
			"ListTransactions",
			data.ErrorType,
			data.ErrorMessage,
			data.ErrorCode,
			sc,
		}
	}

	return nil
}
