package requests

import (
	"net/http"
	"time"
)

// NewTransport returns default HTTP transport which
// should be reused as it caches underlying TCP connections.
// If connections pooling is not needed consider to set
// DisableKeepAlives=false and MaxIdleConnsPerHost=-1.
func NewTransport(dialFunc DialContext) *http.Transport {
	return &http.Transport{
		DialContext:           dialFunc,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   30 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		DisableKeepAlives:     true,
	}
}
