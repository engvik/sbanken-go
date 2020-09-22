package sbanken

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
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

	if q != nil {
		qs, err := q.QueryString(url)
		if err != nil {
			return nil, err
		}

		url = fmt.Sprintf("%s?%s", url, qs)
	}

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
		Efakturas []Efaktura `json:"items"`
	}{}

	json.Unmarshal(res, &data)

	return data.Efakturas, nil
}

func (c *Client) PayEfaktura(ctx context.Context, q *EfakturaPayQuery) error {
	if q == nil {
		return errors.New("No EfakturaPayQuery passed")
	}

	payload, err := json.Marshal(q)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/v1/Efakturas", c.baseURL)

	res, sc, err := c.request(ctx, &httpRequest{
		method:      http.MethodPost,
		url:         url,
		postPayload: payload,
	})
	log.Println(string(res))
	if err != nil {
		return err
	}

	if sc != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", sc)
	}

	return nil
}
