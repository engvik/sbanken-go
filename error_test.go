package sbanken

import "testing"

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
				Message:     "an error occured",
			},
			"TestError error: Test (StatusCode: 500 / ErrorCode: 0): an error occured",
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
