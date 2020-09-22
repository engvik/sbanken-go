package sbanken

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

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

type PaymentListQuery struct {
	Index  string
	Length string
}

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

func (c *Client) ListPayments(ctx context.Context, accountID string, q *PaymentListQuery) ([]Payment, error) {
	url := fmt.Sprintf("%s/v1/Payments/%s", c.baseURL, accountID)

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
		Payments []Payment `json:"items"`
	}{}

	if err := json.Unmarshal(res, &data); err != nil {
		return data.Payments, err
	}

	return data.Payments, nil
}
