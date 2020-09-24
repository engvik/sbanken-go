package sbanken

import "testing"

func getTestError(str string) *Error {
	return &Error{
		Code:        100,
		StatusCode:  500,
		ErrorString: str,
		Message:     "an error occurred",
		Type:        "Error",
	}
}

func TestError(t *testing.T) {
	tests := []struct {
		name string
		err  *Error
		exp  string
	}{
		{
			"should format error correctly",
			&Error{
				ErrorString: "TestError",
				Type:        "Test",
				StatusCode:  500,
				Message:     "an error occurred",
			},
			"TestError error: Test (StatusCode: 500 / ErrorCode: 0): an error occurred",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if str := tc.err.Error(); str != tc.exp {
				t.Errorf("unexpected result: got %s, exp %s", str, tc.exp)
			}
		})
	}
}
