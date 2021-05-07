package sbanken

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/engvik/sbanken-go/internal/transport"
)

// Customer represents a customer.
// Sbanken API documentation: https://api.sbanken.no/exec.customers/swagger/index.html?urls.primaryName=Customers%20v1
type Customer struct {
	CustomerID    string        `json:"customerID"`
	FirstName     string        `json:"firstName"`
	LastName      string        `json:"lastName"`
	EmailAddress  string        `json:"emailAddress"`
	DateOfBirth   string        `json:"dateOfBirth"`
	PostalAddress Address       `json:"postalAddress"`
	StreetAddress Address       `json:"streetAddress"`
	PhoneNumbers  []PhoneNumber `json:"phoneNumbers"`
}

// PhoneNumber represents a customer phone number.
type PhoneNumber struct {
	CountryCode string `json:"countryCode"`
	Number      string `json:"number"`
}

// GetCustomer lists customer information.
func (c *Client) GetCustomer(ctx context.Context) (Customer, error) {
	url := fmt.Sprintf("%s/v1/Customers", c.bankBaseURL)

	res, sc, err := c.transport.Request(ctx, &transport.HTTPRequest{
		Method: http.MethodGet,
		URL:    url,
	})
	if err != nil {
		return Customer{}, fmt.Errorf("request: %w", err)
	}

	data := struct {
		Customer `json:"item"`
		transport.HTTPResponse
	}{}

	if err := json.Unmarshal(res, &data); err != nil {
		return data.Customer, fmt.Errorf("Unmarshal: %w", err)
	}

	if data.IsError || sc != http.StatusOK {
		return Customer{}, &Error{
			"Customers",
			data.ErrorType,
			data.ErrorMessage,
			data.ErrorCode,
			sc,
		}
	}

	return data.Customer, nil
}
