package main

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/frozzare/go-swish"
)

var client *swish.Client

func init() {
	var err error

	client, err = swish.NewClient(&swish.Options{
		Env:        "test",
		Passphrase: "swish",
		P12:        "./certs/test.p12",
		Root:       "./certs/root.pem",
	})

	if err != nil {
		log.Fatal(err)
	}
}

func create(w http.ResponseWriter, r *http.Request) {
	res, err := client.CreatePaymentRequest(context.Background(), &swish.PaymentRequest{
		CallbackURL:           "https://c06610e4.ngrok.io",
		PayeePaymentReference: "0123456789",
		PayeeAlias:            "1231181189",
		Amount:                "100.00",
		Currency:              "SEK",
		Message:               "Kingston USB Flash Drive 8 GB",
	})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	writeJSON(w, res)
}

func status(w http.ResponseWriter, r *http.Request) {
	res, err := client.PaymentRequest(context.Background(), r.URL.Query().Get("id"))

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	writeJSON(w, res)
}

func callback(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	log.Println(string(body))

	io.WriteString(w, "ok")
}

func writeJSON(w http.ResponseWriter, data interface{}) {
	j, err := json.Marshal(data)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, string(j))
}

func main() {
	http.HandleFunc("/create", create)
	http.HandleFunc("/status", status)
	http.HandleFunc("/callback", callback)
	http.ListenAndServe(":5000", nil)
}
