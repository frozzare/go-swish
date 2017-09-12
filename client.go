package swish

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"golang.org/x/crypto/pkcs12"
)

// Options represents Swish client options.
type Options struct {
	Env        string
	P12        string
	Passphrase string
	Root       string
	Client     *http.Client
}

// Client represents a Swish client.
type Client struct {
	*Options
}

// Error represents a error object from Swish API.
type Error struct {
	AdditionalInformation string `json:"additionalInformation,omitempty"`
	ErrorCode             string `json:"errorCode,omitempty"`
	ErrorMessage          string `json:"errorMessage,omitempty"`
}

// NewClient creats a new Swish client.
func NewClient(opts *Options) (*Client, error) {
	if opts.Client == nil {
		opts.Client = http.DefaultClient
	}

	cfg, err := createTLSConfig(opts)
	if err != nil {
		return nil, err
	}

	if opts.Client.Transport == nil {
		opts.Client.Transport = http.DefaultTransport
	}

	if _, ok := opts.Client.Transport.(*http.Transport); ok {
		opts.Client.Transport.(*http.Transport).TLSClientConfig = cfg
	}

	return &Client{opts}, nil
}

// URL returns the BankID url.
func (s *Client) URL() string {
	switch s.Env {
	case "production":
		return "https://swicpc.bankgirot.se/swish-cpcapi/api/v1"
	default:
		return "https://mss.swicpc.bankgirot.se/swish-cpcapi/api/v1"
	}
}

// createTLSConfig creates a TLSConfig with the certificates that are configured.
func createTLSConfig(opts *Options) (*tls.Config, error) {
	p12, err := ioutil.ReadFile(opts.P12)
	if err != nil {
		return nil, err
	}

	blocks, err := pkcs12.ToPEM(p12, opts.Passphrase)
	if err != nil {
		return nil, err
	}

	var pemData []byte
	for _, b := range blocks {
		pemData = append(pemData, pem.EncodeToMemory(b)...)
	}

	cert, err := tls.X509KeyPair(pemData, pemData)
	if err != nil {
		return nil, err
	}

	// Load CA cert
	caCert, err := ioutil.ReadFile(opts.Root)
	if err != nil {
		return nil, err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	}

	return tlsConfig, nil
}

// createRequest will create a http request with given method to the given endpoint with the given data.
func (s *Client) createRequest(ctx context.Context, method, endpoint string, data interface{}) (*http.Response, error) {
	var body io.Reader

	if data != nil {
		j, err := json.Marshal(data)

		if err != nil {
			return nil, err
		}

		body = bytes.NewBuffer(j)
	}

	req, err := http.NewRequest(method, s.URL()+endpoint, body)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	res, err := s.Client.Do(req.WithContext(ctx))

	if err != nil {
		// If we got an error, and the context has been canceled,
		// the context's error is probably more useful.
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		return nil, err
	}

	if res.StatusCode != 200 && res.StatusCode != 201 {
		var errs []Error

		readChuncked(res, &errs)

		if len(errs) > 0 {
			return res, errors.New(errs[0].ErrorMessage)
		}

		return res, fmt.Errorf("Bad status code from Swish API: %d", res.StatusCode)
	}

	return res, nil
}

func readChuncked(res *http.Response, target interface{}) error {
	defer func() {
		// Drain up to 512 bytes and close the body to let the Transport reuse the connection.
		io.CopyN(ioutil.Discard, res.Body, 512)
		res.Body.Close()
	}()

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(res.Body); err != nil {
		return err
	}

	return json.Unmarshal(buf.Bytes(), &target)
}
