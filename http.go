package sbanken

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
)

type httpRequest struct {
	method string
	url    string
}

func (c *Client) request(ctx context.Context, r *httpRequest) ([]byte, int, error) {
	token, err := c.getToken(ctx)
	if err != nil {
		return nil, 0, err
	}

	req, err := http.NewRequest(r.method, r.url, nil)
	if err != nil {
		return nil, 0, err
	}

	req = req.WithContext(ctx)

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("customerId", c.config.CustomerID)

	res, err := c.HTTP.Do(req)
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
