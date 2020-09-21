package sbanken

import (
	"encoding/json"
	"fmt"
	"net/http"
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

func (c *Client) ListEfakturas() ([]Efaktura, error) {
	url := fmt.Sprintf("%s/v1/Efakturas", c.baseURL)

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
