package sbanken

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func (c *Client) request(url string) ([]byte, int, error) {
	token, err := c.getToken()
	if err != nil {
		return nil, 0, err
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, 0, err
	}

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
