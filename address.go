package sbanken

import (
	"fmt"
	"strings"
)

// Address represents a customer address.
type Address struct {
	AddressLine1 string `json:"addressLine1"`
	AddressLine2 string `json:"addressLine2"`
	AddressLine3 string `json:"addressLine3"`
	Country      string `json:"country"`
	ZipCode      string `json:"zipCode"`
	City         string `json:"city"`
}

// String returns an address represented as a string.
func (a Address) String() string {
	var addrLines []string

	if a.AddressLine1 != "" {
		addrLines = append(addrLines, a.AddressLine1)
	}

	if a.AddressLine2 != "" {
		addrLines = append(addrLines, a.AddressLine2)
	}

	if a.AddressLine3 != "" {
		addrLines = append(addrLines, a.AddressLine3)
	}

	if a.ZipCode != "" && a.City != "" {
		addrLines = append(addrLines, fmt.Sprintf("%s %s", a.ZipCode, a.City))
	} else {
		if a.ZipCode != "" {
			addrLines = append(addrLines, a.ZipCode)
		}

		if a.City != "" {
			addrLines = append(addrLines, a.City)
		}
	}

	if a.Country != "" {
		addrLines = append(addrLines, a.Country)
	}

	return fmt.Sprintf("%s", strings.Join(addrLines, ", "))
}
