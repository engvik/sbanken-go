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
	OriginalAmount      float32 `json:"originalAmount"`
	UpdatedAmount       float32 `json:"updatedAmount"`
	MinimumAmount       float32 `json:"minimumAmount"`
	CreditAccountNumber int     `json:"creditAccountNumber"`
}

type EfakturaListQuery struct {
	StartDate time.Time
	EndDate   time.Time
	Status    string
	Index     string
	Length    string
}

func (q *EfakturaListQuery) QueryString(u string) (string, error) {
	parsedURL, err := url.Parse(u)
	if err != nil {
		return u, err
	}

	query := parsedURL.Query()

	if !q.StartDate.IsZero() {
		query.Add("startDate", q.StartDate.String())
	}

	if !q.EndDate.IsZero() {
		query.Add("endDate", q.EndDate.String())
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

type EfakturaPayQuery struct {
	ID                   string `json:"eFakturaId"`
	AccountID            string `json:"accountId"`
	PayOnlyMinimumAmount bool   `json:"PayOnlyMinimumAmount"`
}

func (c *Client) ListEfakturas(ctx context.Context, q *EfakturaListQuery) ([]Efaktura, error) {
	url := fmt.Sprintf("%s/v1/Efakturas", c.baseURL)

	return c.listEfakturas(ctx, url, q, "ListEfakturas")
}

func (c *Client) PayEfaktura(ctx context.Context, q *EfakturaPayQuery) error {
	if q == nil {
		return ErrMissingEfakturaQuery
	}

	payload, err := json.Marshal(q)
	if err != nil {
		return fmt.Errorf("Marshal: %w", err)
	}

	url := fmt.Sprintf("%s/v1/Efakturas", c.baseURL)

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

func (c *Client) ListNewEfakturas(ctx context.Context, q *EfakturaListQuery) ([]Efaktura, error) {
	url := fmt.Sprintf("%s/v1/Efakturas/new", c.baseURL)

	if !q.StartDate.IsZero() {
		return nil, ErrNotValidOptionStartDate
	}

	if !q.EndDate.IsZero() {
		return nil, ErrNotValidOptionEndDate
	}

	if q.Status != "" {
		return nil, ErrNotValidOptionStatus
	}

	return c.listEfakturas(ctx, url, q, "ListNewEfakturas")
}

func (c *Client) ReadEfaktura(ctx context.Context, efakturaID string) (Efaktura, error) {
	if efakturaID == "" {
		return Efaktura{}, ErrMissingEfakturaID
	}

	url := fmt.Sprintf("%s/v1/Efakturas/%s", c.baseURL, efakturaID)

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
