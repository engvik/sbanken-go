package sbanken

type Client struct{}

func NewClient(*Config) *Client {
	return &Client{}
}
