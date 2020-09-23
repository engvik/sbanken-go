package sbanken

import "fmt"

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
