package sbanken

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
)

type httpRequest struct {
	method      string
	url         string
	postPayload []byte
}

type httpResponse struct {
	TraceID        string `json:"traceId"`
	ErrorType      string `json:"errorType"`
	ErrorMessage   string `json:"errorMessage"`
	AvailableItems int    `json:"availableItems"`
	ErrorCode      int    `json:"errorCode"`
	IsError        bool   `json:"isError"`
}

func (c *Client) request(ctx context.Context, r *httpRequest) ([]byte, int, error) {
	token, err := c.getToken(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("getToken: %w", err)
	}

	var req *http.Request

	switch r.method {
	case http.MethodGet:
		req, err = http.NewRequest(r.method, r.url, nil)
	case http.MethodPost:
		if r.postPayload == nil {
			return nil, 0, ErrMissingPostPayload
		}
		req, err = http.NewRequest(r.method, r.url, bytes.NewBuffer(r.postPayload))
	default:
		return nil, 0, fmt.Errorf("HTTP request method not supported: %s", r.method)
	}

	if err != nil {
		return nil, 0, fmt.Errorf("NewRequest: %w", err)
	}

	req = req.WithContext(ctx)

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("customerId", c.config.CustomerID)

	res, err := c.HTTP.Do(req)
	if err != nil {
		return nil, res.StatusCode, fmt.Errorf("Do: %w", err)
	}

	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, res.StatusCode, fmt.Errorf("ReadAll: %w", err)
	}

	return data, res.StatusCode, nil
}
