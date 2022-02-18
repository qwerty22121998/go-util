package hclient

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
)

type Method string

type RequestOption func(c *Request)

type Request struct {
	mu       sync.Mutex
	client   *http.Client
	request  *http.Request
	response []byte
	done     bool
	method   string
	url      string
	err      error
}

func New(method string, url string, ops ...RequestOption) *Request {
	instance := &Request{
		client: &http.Client{},
		method: method,
		url:    url,
	}

	request, err := http.NewRequest(method, url, nil)
	if err != nil {
		return instance.setError(err)
	}

	instance.request = request

	for _, op := range ops {
		op(instance)
	}

	return instance
}

func (r *Request) Close() {
	r.client.CloseIdleConnections()
}

func (r *Request) setError(err error) *Request {
	r.err = err
	return r
}

func (r *Request) do(op RequestOption) *Request {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.err != nil {
		return r
	}
	op(r)
	return r
}

func (r *Request) Do() *Request {
	return r.do(func(c *Request) {
		defer func() {
			c.done = true
		}()
		resp, err := c.client.Do(c.request)
		if err != nil {
			c.setError(err)
			return
		}
		defer resp.Body.Close()
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			c.setError(err)
			return
		}
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			c.setError(errors.New(fmt.Sprintf("error with status code=%v, body=\"%v\"", resp.StatusCode, string(data))))
		}
		c.response = data
	})
}

func (r *Request) Error() error {
	if !r.done {
		r.Do()
	}
	return r.err
}
