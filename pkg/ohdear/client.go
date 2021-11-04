package ohdear

import (
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
)

type Client struct {
	*resty.Client
}

func NewClient(baseURL, token string) *Client {
	client := resty.New()
	client.SetBaseURL(baseURL)
	client.SetAuthToken(token)
	client.SetHeaders(map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	})
	client.SetError(&Error{})
	client.OnAfterResponse(errorFromResponse)

	client.SetRetryCount(3)
	client.SetRetryWaitTime(5 * time.Second)
	client.SetRetryMaxWaitTime(20 * time.Second)
	client.AddRetryCondition(func(r *resty.Response, err error) bool {
		return r.StatusCode() == http.StatusTooManyRequests
	})

	client.SetDebug(true)
	client.SetLogger(&TerraformLogger{})

	return &Client{client}
}

func (c *Client) SetUserAgent(ua string) *Client {
	c.Header.Set("User-Agent", ua)
	return c
}
