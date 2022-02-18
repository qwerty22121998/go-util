package hclient

import "net/http"

func WithHTTPClient(client *http.Client) RequestOption {
	return func(c *Request) {
		c.client = client
	}
}
