package transport

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

type auth struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
	ExpiresIn   int    `json:"expires_in"`
	expires     time.Time
}

func (c *Client) authorize(ctx context.Context) error {
	authURL := "https://auth.sbanken.no/identityserver/connect/token"
	payload := []byte("grant_type=client_credentials")

	req, err := http.NewRequest(http.MethodPost, authURL, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	req = req.WithContext(ctx)

	req.SetBasicAuth(url.QueryEscape(c.clientID), url.QueryEscape(c.clientSecret))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")
	req.Header.Set("Accept", "application/json")

	res, err := c.http.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		// TODO: Handle error from body
		return fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var a auth
	if err := json.Unmarshal(data, &a); err != nil {
		return err
	}

	exp := time.Now().Add(time.Second * time.Duration(a.ExpiresIn))
	a.expires = exp

	c.auth = &a

	return nil
}

func (c *Client) getToken(ctx context.Context) (string, error) {
	if time.Now().After(c.auth.expires) {
		err := c.authorize(ctx)
		if err != nil {
			return "", fmt.Errorf("error renewing token: %w", err)
		}
	}

	return c.auth.AccessToken, nil
}
