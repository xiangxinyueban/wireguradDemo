package requests

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"
)

const (
	// DefaultTimeout is a default HTTP client timeout.
	DefaultTimeout = 20 * time.Second
)

// NewHTTPClientWithTransport creates a new HTTP client with custom transport.
func NewHTTPClientWithTransport(transport *http.Transport, timeout time.Duration) *HTTPClient {
	c := &HTTPClient{
		clientFactory: func(proxyPort int) *http.Client {
			t := transport.Clone()

			if proxyPort > 0 {
				url, _ := url.Parse(fmt.Sprintf("http://localhost:%d", proxyPort))
				t.Proxy = http.ProxyURL(url)
			}

			return &http.Client{
				Timeout:   timeout,
				Transport: setUserAgent(t, getUserAgent()),
			}
		},
	}
	// Create initial clean before any HTTP request is made.
	c.client = c.clientFactory(0)

	return c
}

func setUserAgent(transport http.RoundTripper, userAgent string) http.RoundTripper {
	return &userAgenter{
		transport: transport,
		Agent:     userAgent,
	}
}

func getUserAgent() string {
	return fmt.Sprintf("Ladder %v", "0.1")
}

type userAgenter struct {
	transport http.RoundTripper
	Agent     string
}

func (ua *userAgenter) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Set("User-Agent", ua.Agent)
	return ua.transport.RoundTrip(r)
}

// NewHTTPClient creates a new HTTP client.
func NewHTTPClient(srcIP string, timeout time.Duration) *HTTPClient {
	return NewHTTPClientWithTransport(NewTransport(NewDialer(srcIP).DialContext), timeout)
}

// HTTPClient describes a client for performing HTTP requests.
type HTTPClient struct {
	client        *http.Client
	clientMu      sync.Mutex
	clientFactory func(proxyPort int) *http.Client
}

// Do send an HTTP request and returns an HTTP response.
func (c *HTTPClient) Do(req *http.Request) (*http.Response, error) {
	return c.resolveClient().Do(req)
}

// DoViaProxy send an HTTP request via proxy and returns an HTTP response.
func (c *HTTPClient) DoViaProxy(req *http.Request, proxyPort int) (*http.Response, error) {
	return c.clientFactory(proxyPort).Do(req)
}

// DoRequest performs HTTP requests and parses error without returning response.
func (c *HTTPClient) DoRequest(req *http.Request) error {
	response, err := c.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	return ParseResponseError(response)
}

// DoRequestViaProxy performs HTTP requests via proxy and parses error without returning response.
func (c *HTTPClient) DoRequestViaProxy(req *http.Request, proxyPort int) error {
	response, err := c.DoViaProxy(req, proxyPort)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	return ParseResponseError(response)
}

// DoRequestAndParseResponse performs HTTP requests and response from JSON.
func (c *HTTPClient) DoRequestAndParseResponse(req *http.Request, resp interface{}) error {
	response, err := c.Do(req)
	log.Error().Err(err)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	//httptrace.TraceRequestResponse(req, response)
	blob, _ := io.ReadAll(response.Body)
	log.Error().Msg(response.Status)
	log.Error().Msg(string(blob))
	err = ParseResponseError(response)
	if err != nil {
		return err
	}

	return ParseResponseJSON(response, &resp)
}

// DoRequestViaProxyAndParseResponse performs HTTP requests and response from JSON.
func (c *HTTPClient) DoRequestViaProxyAndParseResponse(req *http.Request, resp interface{}, proxyPort int) error {
	response, err := c.DoViaProxy(req, proxyPort)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	err = ParseResponseError(response)
	if err != nil {
		return err
	}

	return ParseResponseJSON(response, &resp)
}

func (c *HTTPClient) resolveClient() *http.Client {
	c.clientMu.Lock()
	defer c.clientMu.Unlock()
	if c.client != nil {
		return c.client
	}
	c.client = c.clientFactory(0)
	return c.client
}

// ParseResponseJSON parses http.Response into given struct.
func ParseResponseJSON(response *http.Response, dto interface{}) error {
	err := json.NewDecoder(response.Body).Decode(dto)
	if err == io.EOF {
		return nil
	}
	return err
}

// ParseResponseError parses http.Response error.
func ParseResponseError(response *http.Response) error {
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		blob, _ := io.ReadAll(response.Body)
		return errors.New(string(blob))
	}
	return nil
}
