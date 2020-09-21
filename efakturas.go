package sbanken

import (
	"encoding/json"
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

type EfakturaQuery struct {
	StartDate time.Time
	EndDate   time.Time
	Status    string
	Index     string
	Length    string
}

func (q *EfakturaQuery) QueryString(u string) (string, error) {
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

func (c *Client) ListEfakturas(q *EfakturaQuery) ([]Efaktura, error) {
	url := fmt.Sprintf("%s/v1/Efakturas", c.baseURL)
	log.Println(url)
	if q != nil {
		qs, err := q.QueryString(url)
		if err != nil {
			return nil, err
		}

		url = fmt.Sprintf("%s?%s", url, qs)
		log.Println(url)
	}

	res, sc, err := c.request(url)
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
