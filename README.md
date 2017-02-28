# Swish [![Build Status](https://travis-ci.org/frozzare/go-swish.svg?branch=master)](https://travis-ci.org/frozzare/go-swish) [![GoDoc](https://godoc.org/github.com/frozzare/go-swish?status.svg)](https://godoc.org/github.com/frozzare/go-swish) [![Go Report Card](https://goreportcard.com/badge/github.com/frozzare/go-swish)](https://goreportcard.com/report/github.com/frozzare/go-swish)

> Work in progress

Go package for deailing with [Swish](https://www.getswish.se/) merchant API.

## Installation

```
$ go get github.com/frozzare/go-swish
```

## Usage

Please read the Swish [documentation](https://www.getswish.se/manualer/) first so you know what you need and what the different fields means.

Begin by obtaining the SSL certificates required by Swish. The Swish server itself uses a self-signed root certificated so a CA-bundle to verify its origin is needed.
You will also need a client certificate and corresponding private key so the Swish server can identify you.

Certificates in `certs` directory is the test certificates from Swish and cannot be used in production.

```go
package main

import (
	"log"

	"github.com/frozzare/go-swish"
)

func main() {
	client, err := swish.NewClient(&swish.Options{
		Env:        "test",
		Passphrase: "swish",
		P12:        "./certs/test.p12",
		Root:       "./certs/root.pem",
	})

	if err != nil {
		log.Fatal(err)
	}

	res, err := client.CreatePayment(&swish.PaymentData{
		CallbackURL:           "https://example.com/api/swishcb/paymentrequests",
		PayeePaymentReference: "0123456789",
		PayeeAlias:            "1231181189",
		Amount:                100.00,
		Currency:              "SEK",
		Message:               "Kingston USB Flash Drive 8 GB",
	})

	if err != nil {
		log.Fatal(err)
	}

	log.Println(res.ID)

	res, err = client.Payment(res.ID)

	if err != nil {
		log.Fatal(err)
	}

	log.Println(res.Status)
}
```

# License

MIT © [Fredrik Forsmo](https://github.com/frozzare)