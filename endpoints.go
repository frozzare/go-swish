package swish

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"

	"strings"
)

var (
	// ErrNoLocationHeader is the error when no location header exists in the response from Swish API.
	ErrNoLocationHeader = errors.New("Error: No location header from Swish API")
)

// PaymentRequest represents a payment request from Swish API.
type PaymentRequest struct {
	AdditionalInformation    string `json:"additionalInformation,omitempty"`
	Amount                   string `json:"amount,omitempty"`
	CallbackURL              string `json:"callbackUrl,omitempty"`
	Currency                 string `json:"currency,omitempty"`
	DateCreated              string `json:"dateCreated,omitempty"`
	DatePaid                 string `json:"datePaid,omitempty"`
	ErrorCode                string `json:"errorCode,omitempty"`
	ErrorMessage             string `json:"errorMessage,omitempty"`
	ID                       string `json:"id,omitempty"`
	Message                  string `json:"message,omitempty"`
	PayeeAlias               string `json:"payeeAlias,omitempty"`
	PayeePaymentReference    string `json:"payeePaymentReference,omitempty"`
	PayerPaymentReference    string `json:"payerPaymentReference,omitempty"`
	PayerAlias               string `json:"payerAlias,omitempty"`
	PaymentReference         string `json:"paymentReference,omitempty"`
	OriginalPaymentReference string `json:"originalPaymentReference,omitempty"`
	Status                   string `json:"status,omitempty"`
}

// CreatePaymentRequest will create a payment request to Swish and return a payment
// request containing the ID of the request and the data sent to Swish or a error.
func (c *Client) CreatePaymentRequest(ctx context.Context, req *PaymentRequest) (*PaymentRequest, error) {
	res, err := c.createRequest(ctx, "POST", "/paymentrequests", req)

	if err != nil {
		return nil, err
	}

	if len(res.Header.Get("Location")) == 0 {
		return nil, ErrNoLocationHeader
	}

	req.ID = strings.Replace(res.Header.Get("Location"), c.URL()+"/paymentrequests/", "", -1)

	return req, nil
}

// PaymentRequest will return a payment request or a error for the given id.
func (c *Client) PaymentRequest(ctx context.Context, id string) (*PaymentRequest, error) {
	res, err := c.createRequest(ctx, "GET", "/paymentrequests/"+id, nil)

	if err != nil {
		return nil, err
	}

	defer func() {
		// Drain up to 512 bytes and close the body to let the Transport reuse the connection.
		io.CopyN(ioutil.Discard, res.Body, 512)
		res.Body.Close()
	}()

	var paymentRequest *PaymentRequest

	if err := readChuncked(res, &paymentRequest); err != nil {
		return nil, err
	}

	return paymentRequest, nil
}

// CreateRefundRequest will create a refund request to Swish and return a refund
// request containing the ID of the request and the data sent to Swish or a error.
func (c *Client) CreateRefundRequest(ctx context.Context, req *PaymentRequest) (*PaymentRequest, error) {
	res, err := c.createRequest(ctx, "POST", "/refunds", req)

	if err != nil {
		return nil, err
	}

	if len(res.Header.Get("Location")) == 0 {
		return nil, ErrNoLocationHeader
	}

	req.ID = strings.Replace(res.Header.Get("Location"), c.URL()+"/refunds/", "", -1)

	return req, nil
}

// RefundRequest will return a payment request or a error for the given id.
func (c *Client) RefundRequest(ctx context.Context, id string) (*PaymentRequest, error) {
	res, err := c.createRequest(ctx, "GET", "/refunds/"+id, nil)

	if err != nil {
		return nil, err
	}

	defer func() {
		// Drain up to 512 bytes and close the body to let the Transport reuse the connection.
		io.CopyN(ioutil.Discard, res.Body, 512)
		res.Body.Close()
	}()

	var paymentRequest *PaymentRequest

	json.NewDecoder(res.Body).Decode(&paymentRequest)

	return paymentRequest, nil
}
