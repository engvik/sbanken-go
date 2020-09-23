package sbanken

import (
	"errors"
	"fmt"
)

var (
	ErrMissingAccountID        = errors.New("accountID must be set")
	ErrMissingPaymentID        = errors.New("paymentID must be set")
	ErrMissingEfakturaID       = errors.New("efakturaID must be set")
	ErrMissingTransferQuery    = errors.New("TransferQuery must be set")
	ErrMissingEfakturaQuery    = errors.New("EfakturaPayQuery must be set")
	ErrMissingPostPayload      = errors.New("postPayload must be set")
	ErrMissingClientID         = errors.New("ClientID must be set")
	ErrMissingClientSecret     = errors.New("ClientSecret must be set")
	ErrMissingCustomerID       = errors.New("CustomerID must be set")
	ErrNotValidOptionStartDate = errors.New("StartDate is not valid option for this method")
	ErrNotValidOptionEndDate   = errors.New("EndDate is not valid option for this method")
	ErrNotValidOptionStatus    = errors.New("Status is not valid option for this method")
)

type Error struct {
	ErrorString string
	Type        string
	Message     string
	Code        int
	StatusCode  int
}

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
