package ghttp

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"time"
)

func WithSkipSSLVerify() func(s *http.Client) {
	return func(s *http.Client) {
		s.Transport.(*http.Transport).TLSClientConfig.InsecureSkipVerify = true
	}
}

func WithProxy(proxy string) func(s *http.Client) {
	return func(s *http.Client) {
		s.Transport.(*http.Transport).Proxy = func(req *http.Request) (*url.URL, error) {
			return url.Parse(proxy)
		}
	}
}

func WithEnvironmentProxy() func(s *http.Client) {
	return func(s *http.Client) {
		s.Transport.(*http.Transport).Proxy = http.ProxyFromEnvironment
	}
}

func NewDefaultHttpClient(opts ...func(s *http.Client)) *http.Client {
	s := &http.Client{
		Transport: &http.Transport{
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			MaxConnsPerHost:       100,
			MaxIdleConnsPerHost:   100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func ReadAndCloseHttpResponse(resp *http.Response) ([]byte, error) {
	if resp == nil {
		return nil, errors.New("response is nil")
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}
