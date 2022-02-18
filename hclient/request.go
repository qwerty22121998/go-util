package hclient

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/url"
	"os"
	"strings"
)

func (r *Request) Header(key string, value string) *Request {
	r.request.Header.Set(key, value)
	return r
}

func (r *Request) Form(data url.Values) *Request {
	return r.do(func(c *Request) {
		c.request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		c.request.Body = io.NopCloser(strings.NewReader(data.Encode()))
	})
}

func (r *Request) FormFile(data map[string]io.Reader) *Request {
	var err error
	return r.do(func(c *Request) {
		var b bytes.Buffer
		mw := multipart.NewWriter(&b)
		for k, v := range data {
			var iw io.Writer
			if c, ok := v.(io.Closer); ok {
				defer c.Close()
			}

			if f, ok := v.(*os.File); ok {
				if iw, err = mw.CreateFormFile(k, f.Name()); err != nil {
					r.setError(err)
					return
				}
			} else {
				if iw, err = mw.CreateFormField(k); err != nil {
					r.setError(err)
					return
				}
			}

			if _, err := io.Copy(iw, v); err != nil {
				r.setError(err)
				return
			}
		}
		mw.Close()
		c.request.Header.Set("Content-Type", mw.FormDataContentType())
		c.request.Body = io.NopCloser(&b)
	})
}

func (r *Request) JSON(data interface{}) *Request {
	return r.do(func(c *Request) {
		dataByte, err := json.Marshal(data)
		if err != nil {
			c.setError(err)
			return
		}
		c.request.Header.Set("Content-Type", "application/json")
		c.request.Body = io.NopCloser(bytes.NewReader(dataByte))
	})
}
