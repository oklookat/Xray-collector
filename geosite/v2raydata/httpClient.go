package v2raydata

import (
	"context"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

func newHttpClient(userAgent string, rps float64, timeout time.Duration) *httpClient {
	if rps <= 0 {
		rps = 1
	}

	limiter := rate.NewLimiter(rate.Limit(rps), 1) // burst = 1

	return &httpClient{
		client:    &http.Client{Timeout: timeout},
		limiter:   limiter,
		userAgent: userAgent,
	}
}

type httpClient struct {
	client    *http.Client
	limiter   *rate.Limiter
	userAgent string
}

func (c *httpClient) Get(url, ifNoneMatch string) (*http.Response, error) {
	err := c.limiter.Wait(context.Background())
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", c.userAgent)
	if len(ifNoneMatch) > 0 {
		req.Header.Set("If-None-Match", ifNoneMatch)
	}

	return c.client.Do(req)
}
