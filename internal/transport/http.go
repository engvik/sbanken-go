package transport

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

// HTTPRequest represents a http request.
type HTTPRequest struct {
	Method      string
	URL         string
	PostPayload []byte
}

// HTTPResponse represents a http response.
type HTTPResponse struct {
	TraceID        string `json:"traceId"`
	ErrorType      string `json:"errorType"`
	ErrorMessage   string `json:"errorMessage"`
	AvailableItems int    `json:"availableItems"`
	ErrorCode      int    `json:"errorCode"`
	IsError        bool   `json:"isError"`
}

// Request performs the HTTP request.
func (c *Client) Request(ctx context.Context, r *HTTPRequest) ([]byte, int, error) {
	token, err := c.getToken(ctx)
	if err != nil {
		return nil, 0, err
	}

	var req *http.Request

	switch r.Method {
	case http.MethodGet:
		req, err = http.NewRequest(r.Method, r.URL, nil)
	case http.MethodPost:
		if r.PostPayload == nil {
			return nil, 0, errors.New("Post payload missing from POST")
		}
		req, err = http.NewRequest(r.Method, r.URL, bytes.NewBuffer(r.PostPayload))
		req.Header.Set("Content-Type", "application/json")
	default:
		return nil, 0, fmt.Errorf("Invalid HTTP request method: %s", r.Method)
	}

	if err != nil {
		return nil, 0, err
	}

	req = req.WithContext(ctx)

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.userAgent)

	res, err := c.http.Do(req)
	if err != nil {
		return nil, res.StatusCode, err
	}

	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, res.StatusCode, err
	}

	return data, res.StatusCode, nil
}
