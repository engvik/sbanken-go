package sbanken

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Account struct {
	ID          string  `json:"accountId"`
	Name        string  `json:"name"`
	Type        string  `json:"accountType"`
	Available   float32 `json:"available"`
	Balance     float32 `json:"balance"`
	CreditLimit float32 `json:"creditLimit"`
	Number      int     `json:"accountNumber"`
}

func (c *Client) ListAccounts() ([]Account, error) {
	url := fmt.Sprintf("%s/v1/Accounts", c.baseURL)

	res, sc, err := c.request(url)
	if err != nil {
		return nil, err
	}

	if sc != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", sc)
	}

	data := struct {
		Accounts []Account `json:"items"`
	}{}

	json.Unmarshal(res, &data)

	return data.Accounts, nil
}

func (c *Client) ListAccount(accountID string) (Account, error) {
	url := fmt.Sprintf("%s/v1/Accounts/%s", c.baseURL, accountID)

	res, sc, err := c.request(url)
	if err != nil {
		return Account{}, err
	}

	if sc != http.StatusOK {
		return Account{}, fmt.Errorf("unexpected status code: %d", sc)
	}

	data := struct {
		Account Account `json:"item"`
	}{}

	json.Unmarshal(res, &data)

	return data.Account, nil
}
