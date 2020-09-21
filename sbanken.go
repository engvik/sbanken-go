package sbanken

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	HTTP    *http.Client
	config  *Config
	auth    *auth
	baseURL string
}

type auth struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
	ExpiresIn   int    `json:"expires_in"`
	expires     time.Time
}

func NewClient(cfg *Config, httpClient *http.Client) (*Client, error) {
	if err := cfg.validate(); err != nil {
		return nil, err
	}

	c := &Client{}
	c.setHTTPClient(httpClient)
	c.config = cfg
	c.baseURL = "https://api.sbanken.no/exec.bank/api"

	if err := c.authorize(cfg); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Client) setHTTPClient(httpClient *http.Client) {
	if httpClient == nil {
		c.HTTP = http.DefaultClient
		return
	}

	c.HTTP = httpClient
}

func (c *Client) authorize(cfg *Config) error {
	authURL := "https://auth.sbanken.no/identityserver/connect/token"
	payload := []byte("grant_type=client_credentials")

	req, err := http.NewRequest(http.MethodPost, authURL, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	req.SetBasicAuth(url.QueryEscape(cfg.ClientID), url.QueryEscape(cfg.ClientSecret))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")
	req.Header.Set("Accept", "application/json")

	res, err := c.HTTP.Do(req)
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

func (c *Client) getToken() (string, error) {
	if time.Now().After(c.auth.expires) {
		err := c.authorize(c.config)
		if err != nil {
			return "", fmt.Errorf("error renewing token: %w", err)
		}
	}

	return c.auth.AccessToken, nil
}

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
