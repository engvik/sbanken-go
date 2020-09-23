package sbanken

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

func (c *Client) authorize(ctx context.Context, cfg *Config) error {
	authURL := "https://auth.sbanken.no/identityserver/connect/token"
	payload := []byte("grant_type=client_credentials")

	req, err := http.NewRequest(http.MethodPost, authURL, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("NewRequest: %w", err)
	}

	req = req.WithContext(ctx)

	req.SetBasicAuth(url.QueryEscape(cfg.ClientID), url.QueryEscape(cfg.ClientSecret))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")
	req.Header.Set("Accept", "application/json")

	res, err := c.HTTP.Do(req)
	if err != nil {
		return fmt.Errorf("Do: %w", err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("ReadAll: %w", err)
	}

	var a auth
	if err := json.Unmarshal(data, &a); err != nil {
		return fmt.Errorf("Unmarshal: %w", err)
	}

	exp := time.Now().Add(time.Second * time.Duration(a.ExpiresIn))
	a.expires = exp

	c.auth = &a

	return nil
}

func (c *Client) getToken(ctx context.Context) (string, error) {
	if time.Now().After(c.auth.expires) {
		err := c.authorize(ctx, c.config)
		if err != nil {
			return "", fmt.Errorf("Error renewing token: %w", err)
		}
	}

	return c.auth.AccessToken, nil
}
