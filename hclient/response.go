package hclient

import "encoding/json"

func (r *Request) ParseJSON(v interface{}) *Request {
	return r.do(func(c *Request) {
		if !c.done {
			c.mu.Unlock()
			c.Do()
			c.mu.Lock()
		}
		c.setError(json.Unmarshal(c.response, v))
	})
}
