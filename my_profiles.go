package sbanken

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/engvik/sbanken-go/internal/transport"
)

// Type ContactInformation represents my profile's contact information.
// Same as sbanken.Customers
// Sbanken API documentation: https://api.sbanken.no/exec.customers/swagger/index.html?urls.primaryName=MyProfiles%20v1
type ContactInformation struct {
	CustomerID    string        `json:"customerID"`
	FirstName     string        `json:"firstName"`
	LastName      string        `json:"lastName"`
	EmailAddress  string        `json:"emailAddress"`
	DateOfBirth   string        `json:"dateOfBirth"`
	PostalAddress Address       `json:"postalAddress"` // Address is defined in customers.go
	StreetAddress Address       `json:"streetAddress"`
	PhoneNumbers  []PhoneNumber `json:"phoneNumbers"` // PhoneNumber is defined in customers.go
}

func (c *Client) GetContactInformation(ctx context.Context) (ContactInformation, error) {
	url := fmt.Sprintf("%s/v1/MyProfiles/contactinformation", c.customersBaseURL)

	res, sc, err := c.transport.Request(ctx, &transport.HTTPRequest{
		Method: http.MethodGet,
		URL:    url,
		QueryParams: map[string]string{
			"customerId": c.customerID,
		},
	})
	if err != nil {
		return ContactInformation{}, fmt.Errorf("request: %w", err)
	}
	log.Println(string(res), sc)

	data := struct {
		ContactInformation `json:"item"`
		transport.HTTPResponse
	}{}

	if err := json.Unmarshal(res, &data); err != nil {
		return data.ContactInformation, fmt.Errorf("Unmarshal: %w", err)
	}

	if data.IsError || sc != http.StatusOK {
		return ContactInformation{}, &Error{
			"Customers",
			data.ErrorType,
			data.ErrorMessage,
			data.ErrorCode,
			sc,
		}
	}

	return data.ContactInformation, nil
}
