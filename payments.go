package sbanken

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/engvik/sbanken-go/internal/transport"
)

// Payment represents a payment.
// Sbanken API documentation: https://api.sbanken.no/exec.bank/swagger/index.html?urls.primaryName=Payments%20v1
type Payment struct {
	AllowedNewStatusTypes  []string `json:"allowedNewStatusTypes"`
	ID                     string   `json:"paymentId"`
	RecipientAccountNumber string   `json:"recipientAccountNumber"`
	DueDate                string   `json:"dueDate"`
	KID                    string   `json:"kid"`
	Text                   string   `json:"text"`
	Status                 string   `json:"status"`
	StatusDetails          string   `json:"statusDetails"`
	ProductType            string   `json:"productType"`
	PaymentType            string   `json:"paymentType"`
	BeneficiaryName        string   `json:"beneficiaryName"`
	Amount                 float32  `json:"amount"`
	PaymentNumber          int      `json:"paymentNumber"`
	IsActive               bool     `json:"isActive"`
}

// PaymentListQuery represents query parameters for querying payments.
type PaymentListQuery struct {
	Index  string
	Length string
}

// QueryString translates the query into a query string.
func (q *PaymentListQuery) QueryString(u string) (string, error) {
	parsedURL, err := url.Parse(u)
	if err != nil {
		return u, err
	}

	query := parsedURL.Query()

	if q.Index != "" {
		query.Add("index", q.Index)
	}

	if q.Length != "" {
		query.Add("length", q.Length)
	}

	return query.Encode(), nil
}

// ListPayments list the payments. The accountID are required.
func (c *Client) ListPayments(ctx context.Context, accountID string, q *PaymentListQuery) ([]Payment, error) {
	if accountID == "" {
		return nil, ErrMissingAccountID
	}

	url := fmt.Sprintf("%s/v1/Payments/%s", c.bankBaseURL, accountID)

	if q != nil {
		qs, err := q.QueryString(url)
		if err != nil {
			return nil, fmt.Errorf("QueryString: %w", err)
		}

		url = fmt.Sprintf("%s?%s", url, qs)
	}

	res, sc, err := c.transport.Request(ctx, &transport.HTTPRequest{
		Method: http.MethodGet,
		URL:    url,
	})
	if err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}

	data := struct {
		Payments []Payment `json:"items"`
		transport.HTTPResponse
	}{}

	if err := json.Unmarshal(res, &data); err != nil {
		return data.Payments, fmt.Errorf("Unmarshal: %w", err)
	}

	if data.IsError || sc != http.StatusOK {
		return data.Payments, &Error{
			"ListPayments",
			data.ErrorType,
			data.ErrorMessage,
			data.ErrorCode,
			sc,
		}
	}

	return data.Payments, nil
}

// ReadPayment reads a payment. The accountID and paymentID are required.
func (c *Client) ReadPayment(ctx context.Context, accountID string, paymentID string) (Payment, error) {
	if accountID == "" {
		return Payment{}, ErrMissingAccountID
	}

	if paymentID == "" {
		return Payment{}, ErrMissingPaymentID
	}

	url := fmt.Sprintf("%s/v1/Payments/%s/%s", c.bankBaseURL, accountID, paymentID)

	res, sc, err := c.transport.Request(ctx, &transport.HTTPRequest{
		Method: http.MethodGet,
		URL:    url,
	})
	if err != nil {
		return Payment{}, fmt.Errorf("request: %w", err)
	}

	data := struct {
		Payment Payment `json:"item"`
		transport.HTTPResponse
	}{}

	if err := json.Unmarshal(res, &data); err != nil {
		return data.Payment, fmt.Errorf("Unmarshal: %w", err)
	}

	if data.IsError || sc != http.StatusOK {
		return data.Payment, &Error{
			"ReadPayment",
			data.ErrorType,
			data.ErrorMessage,
			data.ErrorCode,
			sc,
		}
	}

	return data.Payment, nil
}
