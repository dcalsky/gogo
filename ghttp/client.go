package ghttp

import (
	"errors"
	"io"
	"net/http"
	"time"
)

func NewDefaultHttpClient(opts ...func(s *http.Client)) *http.Client {
	s := &http.Client{}
	for _, opt := range opts {
		opt(s)
	}
	if s.Transport == nil {
		s.Transport = &http.Transport{
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			MaxConnsPerHost:       100,
			MaxIdleConnsPerHost:   100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		}
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
