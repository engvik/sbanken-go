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

// Transaction represents a transaction.
// Sbanken API documentation: https://api.sbanken.no/exec.bank/swagger/index.html?urls.primaryName=Transactions%20v1
type Transaction struct {
	CardDetails                 CardDetails        `json:"cardDetails"`
	TransactionDetails          TransactionDetails `json:"transactionDetails"`
	AccountingDate              string             `json:"accountingDate"`
	InterestDate                string             `json:"interestDate"`
	OtherAccountNumber          string             `json:"otherAccountNumber"`
	Text                        string             `json:"text"`
	TransactionType             string             `json:"transactionType"`
	TransactionTypeText         string             `json:"transactionTypeText"`
	ReservationType             string             `json:"reservationType"`
	Source                      string             `json:"source"`
	Amount                      float32            `json:"amount"`
	TransactionTypeCode         int                `json:"transactionTypeCode"`
	IsReservation               bool               `json:"isReservation"`
	CardDetailsSpecified        bool               `json:"cardDetailsSpecified"`
	OtherAccountNumberSpecified bool               `json:"otherAccountNumberSpecified"`
	TransactionDetailSpecified  bool               `json:"transactionDetailSpecified"`
}

// CardDetails contains card details about the card used.
type CardDetails struct {
	CardNumber                  string  `json:"cardNumber"`
	MerchantCategoryCode        string  `json:"merchantCategoryCode"`
	MerchantCategoryDescription string  `json:"merchantCategoryDescription"`
	MerchantCity                string  `json:"merchantCity"`
	MerchantName                string  `json:"merchantName"`
	OriginalCurrencyCode        string  `json:"originalCurrencyCode"`
	PurchaseDate                string  `json:"purchaseDate"`
	TransactionID               string  `json:"transactionId"`
	CurrencyAmount              float32 `json:"currencyAmount"`
	CurrencyRate                float32 `json:"currencyRate"`
}

// TransactionDetails contains transaction details.
type TransactionDetails struct {
	ID                     string `json:"transactionId"`
	FormattedAccountNumber string `json:"formattedAccountNumber"`
	CID                    string `json:"cid"`
	AmountDescription      string `json:"amountDescription"`
	ReceiverName           string `json:"receiverName"`
	PayerName              string `json:"payerName"`
	RegistrationDate       string `json:"registrationDate"`
	NumericReference       int    `json:"numericReference"`
}

// TransactionListQuery represents query parameters for querying transactions.
type TransactionListQuery struct {
	StartDate time.Time
	EndDate   time.Time
	Index     string
	Length    string
}

// QueryString translates the query into a query string.
func (q *TransactionListQuery) QueryString(u string) (string, error) {
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

	if q.Index != "" {
		query.Add("index", q.Index)
	}

	if q.Length != "" {
		query.Add("length", q.Length)
	}

	return query.Encode(), nil
}

// ListTransactions returns the latest transactions of the given account.
func (c *Client) ListTransactions(ctx context.Context, accountID string, q *TransactionListQuery) ([]Transaction, error) {
	url := fmt.Sprintf("%s/v1/Transactions/%s", c.baseURL, accountID)

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
		Transactions []Transaction `json:"items"`
		transport.HTTPResponse
	}{}

	if err := json.Unmarshal(res, &data); err != nil {
		return data.Transactions, fmt.Errorf("Unmarshal: %w", err)
	}

	if data.IsError || sc != http.StatusOK {
		return data.Transactions, &Error{
			"ListTransactions",
			data.ErrorType,
			data.ErrorMessage,
			data.ErrorCode,
			sc,
		}
	}

	return data.Transactions, nil
}
