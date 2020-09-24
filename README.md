# sbanken-go
![tests](https://github.com/engvik/sbanken-go/workflows/master/badge.svg)
[![GoDoc](https://godoc.org/github.com/engvik/sbanken-go?status.svg)](https://godoc.org/github.com/engvik/sbanken-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/engvik/sbanken-go)](https://goreportcard.com/report/github.com/engvik/sbanken-go)

A Go client for the Sbanken API

Exernal information:
* [Sbanken API information](https://sbanken.no/bruke/utviklerportalen/)
* [Sbanken API documentation](https://api.sbanken.no/exec.bank/swagger/index.html)

## Get access to the API

See [this page](https://sbanken.no/bruke/utviklerportalen/) on how to get access to Sbankens API.

## Quick example

```go
ctx := context.Background()
cfg := sbanken.Config{
    ClientID:     os.Getenv("CLIENT_ID"),
    ClientSecret: os.Getenv("CLIENT_SECRET"),
    CustomerID:   os.Getenv("CUSTOMER_ID"),
}

c, err := sbanken.NewClient(ctx, &cfg, nil)
if err != nil {
    log.Fatal(err)
}


accounts, err := c.ListAccounts(ctx)
if err != nil {
    log.Fatal(err)
}

log.Println(accounts)

```

## Documentation

See [GoDoc](https://godoc.org/github.com/engvik/sbanken-go).
