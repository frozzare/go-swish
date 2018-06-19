package swish

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/frozzare/go-assert"

	"gopkg.in/jarcoal/httpmock.v1"
)

func TestCreatePaymentRequest(t *testing.T) {
	httpmock.Activate()

	defer httpmock.DeactivateAndReset()

	tests := []struct {
		description    string
		responder      func(req *http.Request) (*http.Response, error)
		expectedResult *PaymentRequest
		expectedError  error
	}{
		{
			description: "create payment request success",
			responder: func(req *http.Request) (*http.Response, error) {
				resp := httpmock.NewStringResponse(200, "")

				resp.Header.Set("Location", "https://mss.swicpc.bankgirot.se/swish-cpcapi/api/v1/paymentrequests/AB23D7406ECE4542A80152D909EF9F6B")

				return resp, nil
			},
			expectedResult: &PaymentRequest{
				ID: "AB23D7406ECE4542A80152D909EF9F6B",
				PayeePaymentReference: "0123456789",
				CallbackURL:           "https://example.com/api/swishcb/paymentrequests",
				PayerAlias:            "46701234567",
				PayeeAlias:            "1234760039",
				Amount:                "100",
				Currency:              "SEK",
				Message:               "Kingston USB Flash Drive 8 GB",
			},
			expectedError: nil,
		},
		{
			description: "create payment request failed - no location header",
			responder: func(req *http.Request) (*http.Response, error) {
				return httpmock.NewStringResponse(200, ""), nil
			},
			expectedResult: nil,
			expectedError:  ErrNoLocationHeader,
		},
		{
			description: "create payment request failed - bad status code",
			responder: func(req *http.Request) (*http.Response, error) {
				resp := httpmock.NewStringResponse(500, "")

				resp.Header.Set("Location", "https://mss.swicpc.bankgirot.se/swish-cpcapi/api/v1/paymentrequests/AB23D7406ECE4542A80152D909EF9F6B")

				return resp, nil
			},
			expectedResult: nil,
			expectedError:  errors.New("Bad status code from Swish API: 500"),
		},
	}

	client, err := NewClient(&Options{
		Env:        "test",
		Passphrase: "swish",
		P12:        "./certs/test.p12",
		Root:       "./certs/root.pem",
	})

	assert.Nil(t, err)

	for _, test := range tests {
		httpmock.RegisterResponder("POST", "https://mss.swicpc.bankgirot.se/swish-cpcapi/api/v1/paymentrequests", test.responder)

		res, err := client.CreatePaymentRequest(context.Background(), &PaymentRequest{
			PayeePaymentReference: "0123456789",
			CallbackURL:           "https://example.com/api/swishcb/paymentrequests",
			PayerAlias:            "46701234567",
			PayeeAlias:            "1234760039",
			Amount:                "100",
			Currency:              "SEK",
			Message:               "Kingston USB Flash Drive 8 GB",
		})

		assert.Equal(t, res, test.expectedResult, test.description)
		assert.Equal(t, err, test.expectedError, test.description)

		httpmock.Reset()
	}
}

func TestPaymentRequest(t *testing.T) {
	httpmock.Activate()

	defer httpmock.DeactivateAndReset()

	tests := []struct {
		description    string
		responder      func(req *http.Request) (*http.Response, error)
		expectedResult *PaymentRequest
		expectedError  error
	}{
		{
			description: "payment request success",
			responder: func(req *http.Request) (*http.Response, error) {
				return httpmock.NewStringResponse(200, `
                    {
                        "id": "AB23D7406ECE4542A80152D909EF9F6B",
                        "payeePaymentReference": "0123456789",
                        "paymentReference": "6D6CD7406ECE4542A80152D909EF9F6B",
                        "callbackUrl": "https://example.com/api/swishcb/paymentrequests",
                        "payerAlias": "46701234567",
                        "payeeAlias": "1234760039",
                        "amount": "100",
                        "currency": "SEK",
                        "message": "Kingston USB Flash Drive 8 GB",
                        "status": "PAID",
                        "dateCreated": "2015-02-19T22:01:53+01:00",
                        "datePaid": "2015-02-19T22:03:53+01:00"
                    }
                `), nil
			},
			expectedResult: &PaymentRequest{
				ID: "AB23D7406ECE4542A80152D909EF9F6B",
				PayeePaymentReference: "0123456789",
				PaymentReference:      "6D6CD7406ECE4542A80152D909EF9F6B",
				CallbackURL:           "https://example.com/api/swishcb/paymentrequests",
				PayerAlias:            "46701234567",
				PayeeAlias:            "1234760039",
				Amount:                "100",
				Currency:              "SEK",
				Message:               "Kingston USB Flash Drive 8 GB",
				Status:                "PAID",
				DateCreated:           "2015-02-19T22:01:53+01:00",
				DatePaid:              "2015-02-19T22:03:53+01:00",
			},
			expectedError: nil,
		},
		{
			description: "payment request failed - bad status code",
			responder: func(req *http.Request) (*http.Response, error) {
				return httpmock.NewStringResponse(500, ""), nil
			},
			expectedResult: nil,
			expectedError:  errors.New("Bad status code from Swish API: 500"),
		},
	}

	client, err := NewClient(&Options{
		Env:        "test",
		Passphrase: "swish",
		P12:        "./certs/test.p12",
		Root:       "./certs/root.pem",
	})

	assert.Nil(t, err)

	for _, test := range tests {
		httpmock.RegisterResponder("GET", "https://mss.swicpc.bankgirot.se/swish-cpcapi/api/v1/paymentrequests/AB23D7406ECE4542A80152D909EF9F6B", test.responder)

		res, err := client.PaymentRequest(context.Background(), "AB23D7406ECE4542A80152D909EF9F6B")

		assert.Equal(t, res, test.expectedResult, test.description)
		assert.Equal(t, err, test.expectedError, test.description)

		httpmock.Reset()
	}
}

func TestCreateRefundRequest(t *testing.T) {
	httpmock.Activate()

	defer httpmock.DeactivateAndReset()

	tests := []struct {
		description    string
		responder      func(req *http.Request) (*http.Response, error)
		expectedResult *PaymentRequest
		expectedError  error
	}{
		{
			description: "create refund request success",
			responder: func(req *http.Request) (*http.Response, error) {
				resp := httpmock.NewStringResponse(200, "")

				resp.Header.Set("Location", "https://mss.swicpc.bankgirot.se/swish-cpcapi/api/v1/refunds/AB23D7406ECE4542A80152D909EF9F6B")

				return resp, nil
			},
			expectedResult: &PaymentRequest{
				ID: "AB23D7406ECE4542A80152D909EF9F6B",
				OriginalPaymentReference: "AB23D7406ECE4542A80152D909EF9F6B",
				PayerPaymentReference:    "0123456789",
				CallbackURL:              "https://example.com/api/swishcb/paymentrequests",
				PayerAlias:               "46701234567",
				Amount:                   "100",
				Currency:                 "SEK",
				Message:                  "Refund for Kingston USB Flash Drive 8 GB",
			},
			expectedError: nil,
		},
		{
			description: "create refund request failed - no location header",
			responder: func(req *http.Request) (*http.Response, error) {
				return httpmock.NewStringResponse(200, ""), nil
			},
			expectedResult: nil,
			expectedError:  ErrNoLocationHeader,
		},
		{
			description: "create refund request failed - swish errors",
			responder: func(req *http.Request) (*http.Response, error) {
				res := httpmock.NewStringResponse(422, `[{"errorCode":"RF02","errorMessage":"Original Payment not found or original payment is more than than 13 months old","additionalInformation":null}]`)
				return res, nil
			},
			expectedResult: nil,
			expectedError:  errors.New("Original Payment not found or original payment is more than than 13 months old"),
		},
		{
			description: "create refund request failed - bad status code",
			responder: func(req *http.Request) (*http.Response, error) {
				resp := httpmock.NewStringResponse(500, "")

				resp.Header.Set("Location", "https://mss.swicpc.bankgirot.se/swish-cpcapi/api/v1/paymentrequests/AB23D7406ECE4542A80152D909EF9F6B")

				return resp, nil
			},
			expectedResult: nil,
			expectedError:  errors.New("Bad status code from Swish API: 500"),
		},
	}

	client, err := NewClient(&Options{
		Env:        "test",
		Passphrase: "swish",
		P12:        "./certs/test.p12",
		Root:       "./certs/root.pem",
	})

	assert.Nil(t, err)

	for _, test := range tests {
		httpmock.RegisterResponder("POST", "https://mss.swicpc.bankgirot.se/swish-cpcapi/api/v1/refunds", test.responder)

		res, err := client.CreateRefundRequest(context.Background(), &PaymentRequest{
			OriginalPaymentReference: "AB23D7406ECE4542A80152D909EF9F6B",
			PayerPaymentReference:    "0123456789",
			CallbackURL:              "https://example.com/api/swishcb/paymentrequests",
			PayerAlias:               "46701234567",
			Amount:                   "100",
			Currency:                 "SEK",
			Message:                  "Refund for Kingston USB Flash Drive 8 GB",
		})

		assert.Equal(t, res, test.expectedResult, test.description)
		assert.Equal(t, err, test.expectedError, test.description)

		httpmock.Reset()
	}
}

func TestRefundRequest(t *testing.T) {
	httpmock.Activate()

	defer httpmock.DeactivateAndReset()

	tests := []struct {
		description    string
		responder      func(req *http.Request) (*http.Response, error)
		expectedResult *PaymentRequest
		expectedError  error
	}{
		{
			description: "refund request success",
			responder: func(req *http.Request) (*http.Response, error) {
				return httpmock.NewStringResponse(200, `
                    {
                        "id": "AB23D7406ECE4542A80152D909EF9F6B",
                        "payeePaymentReference": "0123456789",
                        "paymentReference": "6D6CD7406ECE4542A80152D909EF9F6B",
                        "callbackUrl": "https://example.com/api/swishcb/refunds",
                        "payerAlias": "46701234567",
                        "payeeAlias": "1234760039",
                        "amount": "100",
                        "currency": "SEK",
                        "message": "Refund for Kingston USB Flash Drive 8 GB",
                        "status": "PAID",
                        "dateCreated": "2015-02-19T22:01:53+01:00",
                        "datePaid": "2015-02-19T22:03:53+01:00"
                    }
                `), nil
			},
			expectedResult: &PaymentRequest{
				ID: "AB23D7406ECE4542A80152D909EF9F6B",
				PayeePaymentReference: "0123456789",
				PaymentReference:      "6D6CD7406ECE4542A80152D909EF9F6B",
				CallbackURL:           "https://example.com/api/swishcb/refunds",
				PayerAlias:            "46701234567",
				PayeeAlias:            "1234760039",
				Amount:                "100",
				Currency:              "SEK",
				Message:               "Refund for Kingston USB Flash Drive 8 GB",
				Status:                "PAID",
				DateCreated:           "2015-02-19T22:01:53+01:00",
				DatePaid:              "2015-02-19T22:03:53+01:00",
			},
			expectedError: nil,
		},
		{
			description: "refund request failed - bad status code",
			responder: func(req *http.Request) (*http.Response, error) {
				return httpmock.NewStringResponse(500, ""), nil
			},
			expectedResult: nil,
			expectedError:  errors.New("Bad status code from Swish API: 500"),
		},
	}

	client, err := NewClient(&Options{
		Env:        "test",
		Passphrase: "swish",
		P12:        "./certs/test.p12",
		Root:       "./certs/root.pem",
	})

	assert.Nil(t, err)

	for _, test := range tests {
		httpmock.RegisterResponder("GET", "https://mss.swicpc.bankgirot.se/swish-cpcapi/api/v1/refunds/AB23D7406ECE4542A80152D909EF9F6B", test.responder)

		res, err := client.RefundRequest(context.Background(), "AB23D7406ECE4542A80152D909EF9F6B")

		assert.Equal(t, res, test.expectedResult, test.description)
		assert.Equal(t, err, test.expectedError, test.description)

		httpmock.Reset()
	}
}

func TestConfigWithCertificateData(t *testing.T) {
	p12File := "./certs/test.p12"
	certFile := "./certs/root.pem"

	fileConfig, err := createTLSConfig(&Options{
		Env:        "test",
		Passphrase: "swish",
		P12:        p12File,
		Root:       certFile,
	})
	assert.Nil(t, err)

	p12, err := ioutil.ReadFile(p12File)
	assert.Nil(t, err)

	root, _ := ioutil.ReadFile(certFile)
	assert.Nil(t, err)

	dataConfig, err := createTLSConfig(&Options{
		Env:        "test",
		Passphrase: "swish",
		P12Data:    p12,
		RootData:   root,
	})
	assert.Nil(t, err)

	// Assert both methods result in the same TLS config
	assert.Equal(t, fileConfig, dataConfig)
}
