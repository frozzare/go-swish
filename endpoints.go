package swish

import (
	"errors"
	"strings"
)

var (
	// ErrNoLocationHeader is the error when Swish API don't return a location header.
	ErrNoLocationHeader = errors.New("No location header from Swish API")
)

// PaymentData represents the payment object from Swish API.
type PaymentData struct {
	AdditionalInformation string      `json:"additionalInformation,omitempty"`
	Amount                interface{} `json:"amount,omitempty"`
	CallbackURL           string      `json:"callbackUrl,omitempty"`
	Currency              string      `json:"currency,omitempty"`
	DateCreated           string      `json:"dateCreated,omitempty"`
	DatePaid              string      `json:"datePaid,omitempty"`
	ErrorCode             string      `json:"errorCode,omitempty"`
	ErrorMessage          string      `json:"errorMessage,omitempty"`
	ID                    string      `json:"id,omitempty"`
	Message               string      `json:"message,omitempty"`
	PayeeAlias            string      `json:"payeeAlias,omitempty"`
	PayeePaymentReference string      `json:"payeePaymentReference,omitempty"`
	PayerPaymentReference string      `json:"payerPaymentReference,omitempty"`
	PayerAlias            string      `json:"payerAlias,omitempty"`
	PaymentReference      string      `json:"paymentReference,omitempty"`
	Status                string      `json:"status,omitempty"`
}

// CreatePayment will create a payment request to Swish and return a payment
// request containing the ID of the request and the data sent to Swish or a error.
func (c *Client) CreatePayment(p *PaymentData) (*PaymentData, error) {
	res, err := c.createRequest("POST", "/paymentrequests", p)

	if err != nil {
		return nil, err
	}

	if len(res.Header.Get("Location")) == 0 {
		return nil, ErrNoLocationHeader
	}

	p.ID = strings.Replace(res.Header.Get("Location"), c.URL()+"/paymentrequests/", "", -1)

	return p, nil
}

// Payment will return a payment request or a error for the given id.
func (c *Client) Payment(id string) (*PaymentData, error) {
	res, err := c.createRequest("GET", "/paymentrequests/"+id, nil)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var p *PaymentData

	if err := readJSON(res, &p); err != nil {
		return nil, err
	}

	return p, nil
}

// RefundData represents the refund object from Swish API.
type RefundData struct {
	AdditionalInformation    string  `json:"additionalInformation,omitempty"`
	Amount                   float64 `json:"amount,omitempty"`
	CallbackURL              string  `json:"callbackUrl,omitempty"`
	Currency                 string  `json:"currency,omitempty"`
	DateCreated              string  `json:"dateCreated,omitempty"`
	DatePaid                 string  `json:"datePaid,omitempty"`
	ErrorCode                string  `json:"errorCode,omitempty"`
	ErrorMessage             string  `json:"errorMessage,omitempty"`
	ID                       string  `json:"id,omitempty"`
	Message                  string  `json:"message,omitempty"`
	PayeeAlias               string  `json:"payeeAlias,omitempty"`
	PayerPaymentReference    string  `json:"payerPaymentReference,omitempty"`
	PayerAlias               string  `json:"payerAlias,omitempty"`
	PaymentReference         string  `json:"paymentReference,omitempty"`
	OriginalPaymentReference string  `json:"originalPaymentReference,omitempty"`
	Status                   string  `json:"status,omitempty"`
}

// CreateRefund will create a refund request to Swish and return a refund
// request containing the ID of the request and the data sent to Swish or a error.
func (c *Client) CreateRefund(r *RefundData) (*RefundData, error) {
	res, err := c.createRequest("POST", "/refunds", r)

	if err != nil {
		return nil, err
	}

	if len(res.Header.Get("Location")) == 0 {
		return nil, ErrNoLocationHeader
	}

	r.ID = strings.Replace(res.Header.Get("Location"), c.URL()+"/refunds/", "", -1)

	return r, nil
}

// Refund will return a payment request or a error for the given id.
func (c *Client) Refund(id string) (*RefundData, error) {
	res, err := c.createRequest("GET", "/refunds/"+id, nil)

	if err != nil {
		return nil, err
	}

	var r *RefundData

	if err := readJSON(res, &r); err != nil {
		return nil, err
	}

	return r, nil
}
