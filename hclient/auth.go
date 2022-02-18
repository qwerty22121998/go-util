package hclient

func (r *Request) BasicAuth(username string, password string) *Request {
	return r.do(func(c *Request) {
		c.request.SetBasicAuth(username, password)
	})
}
