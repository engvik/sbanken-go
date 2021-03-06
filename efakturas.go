package sbanken

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/engvik/sbanken-go/internal/transport"
)

// Efaktura represents an efaktura.
// Sbanken API documentation: https://publicapi.sbanken.no/openapi/apibeta/index.html#/Efaktura
type Efaktura struct {
	ID                  string  `json:"eFakturaId"`
	IssuerID            string  `json:"issuerId"`
	Reference           string  `json:"eFakturaReference"`
	DocumentType        string  `json:"documentType"`
	Status              string  `json:"status"`
	KID                 string  `json:"kid"`
	OriginalDueDate     string  `json:"originalDueDate"`
	UpdatedDueDate      string  `json:"updatedDueDate"`
	NotificationDate    string  `json:"notificationDate"`
	IssuerName          string  `json:"issuerName"`
	CreditAccountNumber string  `json:"creditAccountNumber"`
	OriginalAmount      float32 `json:"originalAmount"`
	UpdatedAmount       float32 `json:"updatedAmount"`
	MinimumAmount       float32 `json:"minimumAmount"`
}

// EfakturaListQuery represents query parameters for querying efakturas.
type EfakturaListQuery struct {
	StartDate time.Time
	EndDate   time.Time
	Status    string
	Index     string
	Length    string
}

// QueryString translates the query into a query string.
func (q *EfakturaListQuery) QueryString(u string) (string, error) {
	parsedURL, err := url.Parse(u)
	if err != nil {
		return u, err
	}

	query := parsedURL.Query()

	if !q.StartDate.IsZero() {
		query.Add("startDate", q.StartDate.Format(time.RFC3339))
	}

	if !q.EndDate.IsZero() {
		query.Add("endDate", q.EndDate.Format(time.RFC3339))
	}

	if q.Status != "" {
		query.Add("status", q.Status)
	}

	if q.Index != "" {
		query.Add("index", q.Index)
	}

	if q.Length != "" {
		query.Add("length", q.Length)
	}

	return query.Encode(), nil
}

// EfakturaPayQuery represents a payment query.
type EfakturaPayQuery struct {
	ID                   string `json:"eFakturaId"`
	AccountID            string `json:"accountId"`
	PayOnlyMinimumAmount bool   `json:"PayOnlyMinimumAmount"`
}

// ListEfakturas lists efakturas.
func (c *Client) ListEfakturas(ctx context.Context, q *EfakturaListQuery) ([]Efaktura, error) {
	url := fmt.Sprintf("%s/v1/Efakturas", c.bankBaseURL)

	return c.listEfakturas(ctx, url, q, "ListEfakturas")
}

// PayEfaktura pays an efaktura. The EfakturaPayQuery are required.
func (c *Client) PayEfaktura(ctx context.Context, q *EfakturaPayQuery) error {
	if q == nil {
		return ErrMissingEfakturaPayQuery
	}

	payload, err := json.Marshal(q)
	if err != nil {
		return fmt.Errorf("Marshal: %w", err)
	}

	url := fmt.Sprintf("%s/v1/Efakturas", c.bankBaseURL)

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
			"PayEfaktura",
			data.ErrorType,
			data.ErrorMessage,
			data.ErrorCode,
			sc,
		}
	}

	return nil
}

// ListNewEfakturas lists efakturas that have not yet been processed by the customer.
func (c *Client) ListNewEfakturas(ctx context.Context, q *EfakturaListQuery) ([]Efaktura, error) {
	url := fmt.Sprintf("%s/v1/Efakturas/new", c.bankBaseURL)

	if q != nil {
		if !q.StartDate.IsZero() {
			return nil, ErrNotValidOptionStartDate
		}

		if !q.EndDate.IsZero() {
			return nil, ErrNotValidOptionEndDate
		}

		if q.Status != "" {
			return nil, ErrNotValidOptionStatus
		}
	}

	return c.listEfakturas(ctx, url, q, "ListNewEfakturas")
}

// ReadEfaktura reads an efaktura. The efakturaID are required.
func (c *Client) ReadEfaktura(ctx context.Context, efakturaID string) (Efaktura, error) {
	if efakturaID == "" {
		return Efaktura{}, ErrMissingEfakturaID
	}

	url := fmt.Sprintf("%s/v1/Efakturas/%s", c.bankBaseURL, efakturaID)

	res, sc, err := c.transport.Request(ctx, &transport.HTTPRequest{
		Method: http.MethodGet,
		URL:    url,
	})
	if err != nil {
		return Efaktura{}, fmt.Errorf("request: %w", err)
	}

	data := struct {
		Efaktura Efaktura `json:"item"`
		transport.HTTPResponse
	}{}

	if err := json.Unmarshal(res, &data); err != nil {
		return data.Efaktura, fmt.Errorf("Unmarshal: %w", err)
	}

	if data.IsError || sc != http.StatusOK {
		return data.Efaktura, &Error{
			"ReadEfaktura",
			data.ErrorType,
			data.ErrorMessage,
			data.ErrorCode,
			sc,
		}
	}

	return data.Efaktura, nil
}

func (c *Client) listEfakturas(ctx context.Context, url string, q *EfakturaListQuery, caller string) ([]Efaktura, error) {
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
		Efakturas []Efaktura `json:"items"`
		transport.HTTPResponse
	}{}

	if err := json.Unmarshal(res, &data); err != nil {
		return data.Efakturas, fmt.Errorf("Unmarshal: %w", err)
	}

	if data.IsError || sc != http.StatusOK {
		return data.Efakturas, &Error{
			caller,
			data.ErrorType,
			data.ErrorMessage,
			data.ErrorCode,
			sc,
		}
	}

	return data.Efakturas, nil
}
