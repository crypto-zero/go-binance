package binance

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type SecType int

const (
	SecTypeNone SecType = iota
	SecTypeAPIKey
	SecTypeSigned // if the 'timestamp' parameter is required
)

type Params map[string]interface{}

// Request define an API Request
type Request struct {
	id         uint64
	method     string
	endpoint   string
	query      url.Values
	form       url.Values
	recvWindow int64
	secType    SecType
	header     http.Header
	body       io.Reader
	fullURL    string
}

// AddQuery add param with key/value to query string
func (r *Request) AddQuery(key string, value interface{}) *Request {
	if r.query == nil {
		r.query = url.Values{}
	}
	r.query.Add(key, fmt.Sprintf("%v", value))
	return r
}

// SetQuery set param with key/value to query string
func (r *Request) SetQuery(key string, value interface{}) *Request {
	if r.query == nil {
		r.query = url.Values{}
	}
	r.query.Set(key, fmt.Sprintf("%v", value))
	return r
}

// SetQueryParams set Params with key/values to query string
func (r *Request) SetQueryParams(m Params) *Request {
	for k, v := range m {
		r.SetQuery(k, v)
	}
	return r
}

// SetForm set param with key/value to Request form body
func (r *Request) SetForm(key string, value interface{}) *Request {
	if r.form == nil {
		r.form = url.Values{}
	}
	r.form.Set(key, fmt.Sprintf("%v", value))
	return r
}

// SetFormParams set Params with key/values to Request form body
func (r *Request) SetFormParams(m Params) *Request {
	for k, v := range m {
		r.SetForm(k, v)
	}
	return r
}

func (r *Request) Validate() (err error) {
	if r.query == nil {
		r.query = url.Values{}
	}
	if r.form == nil {
		r.form = url.Values{}
	}
	return nil
}

// RequestOption define option type for Request
type RequestOption func(*Request)

// WithRecvWindow set recvWindow param for the Request
func WithRecvWindow(recvWindow int64) RequestOption {
	return func(r *Request) {
		r.recvWindow = recvWindow
	}
}

// WithHeader set or add a header value to the Request
func WithHeader(key, value string, replace bool) RequestOption {
	return func(r *Request) {
		if r.header == nil {
			r.header = http.Header{}
		}
		if replace {
			r.header.Set(key, value)
		} else {
			r.header.Add(key, value)
		}
	}
}

// WithHeaders set or replace the headers of the Request
func WithHeaders(header http.Header) RequestOption {
	return func(r *Request) {
		r.header = header.Clone()
	}
}
