package sbanken

import (
	"errors"
	"fmt"
)

var (
	// ErrMissingAccountID are returned when accountID is not set.
	ErrMissingAccountID = errors.New("accountID must be set")
	// ErrMissingPaymentID are returned when paymentID is not set.
	ErrMissingPaymentID = errors.New("paymentID must be set")
	// ErrMissingEfakturaID are returned when efakturaID is not set.
	ErrMissingEfakturaID = errors.New("efakturaID must be set")
	// ErrMissingTransferQuery are returned when TransferQuery is not set.
	ErrMissingTransferQuery = errors.New("TransferQuery must be set")
	// ErrMissingEfakturaPayQuery are returned when EfakturaPayQuery is not set.
	ErrMissingEfakturaPayQuery = errors.New("EfakturaPayQuery must be set")
	// ErrMissingPostPayload are returned when PostPayload is not set.
	ErrMissingPostPayload = errors.New("PostPayload must be set")
	// ErrMissingClientID are returned when ClientID is not set.
	ErrMissingClientID = errors.New("ClientID must be set")
	// ErrMissingClientSecret are returned when ClientSecret is not set.
	ErrMissingClientSecret = errors.New("ClientSecret must be set")
	// ErrNotValidOptionStartDate are returned when StartDate is not allowed.
	ErrNotValidOptionStartDate = errors.New("StartDate is not valid option for this method")
	// ErrNotValidOptionEndDate are returned when EndDate is not allowed.
	ErrNotValidOptionEndDate = errors.New("EndDate is not valid option for this method")
	// ErrNotValidOptionStatus are returned when Status is not allowed.
	ErrNotValidOptionStatus = errors.New("Status is not valid option for this method")
)

// Error represents a standard error.
type Error struct {
	ErrorString string
	Type        string
	Message     string
	Code        int
	StatusCode  int
}

// Error returns the string representation of the error.
func (e *Error) Error() string {
	return fmt.Sprintf(
		"%s error: %s (StatusCode: %d / ErrorCode: %d): %s",
		e.ErrorString,
		e.Type,
		e.StatusCode,
		e.Code,
		e.Message,
	)
}
