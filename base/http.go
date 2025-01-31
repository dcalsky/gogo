package base

import (
	"crypto/tls"
	"fmt"

	"io"
	"net/http"
	"time"
)

func NewDefaultHttpClient(opts ...func(s *http.Client)) *http.Client {
	trans := &http.Transport{
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		MaxConnsPerHost:       100,
		MaxIdleConnsPerHost:   100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: "false" == "false",
		},
	}
	if InLocal() {
		trans.Proxy = http.ProxyFromEnvironment
	}
	s := &http.Client{
		Transport: trans,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func ReadAndCloseHttpResponse(resp *http.Response) ([]byte, error) {
	if resp == nil {
		return nil, fmt.Errorf("response is nil")
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}
