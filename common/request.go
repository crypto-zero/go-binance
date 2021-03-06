package common

import (
	"fmt"
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
	ID         uint64
	Method     string
	Endpoint   string
	SecType    SecType
	Query      url.Values
	Form       url.Values
	RecvWindow int64
	Header     http.Header
}

// AddQuery add param with key/value to query string
func (r *Request) AddQuery(key string, value interface{}) *Request {
	r.Query.Add(key, fmt.Sprintf("%v", value))
	return r
}

// SetQuery set param with key/value to query string
func (r *Request) SetQuery(key string, value interface{}) *Request {
	r.Query.Set(key, fmt.Sprintf("%v", value))
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
	r.Form.Set(key, fmt.Sprintf("%v", value))
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
	if r.Query == nil {
		r.Query = url.Values{}
	}
	if r.Form == nil {
		r.Form = url.Values{}
	}
	return nil
}

// RequestOption define option type for Request
type RequestOption func(*Request)

// WithRecvWindow set recvWindow param for the Request
func WithRecvWindow(recvWindow int64) RequestOption {
	return func(r *Request) {
		r.RecvWindow = recvWindow
	}
}

// WithHeader set or add a Header value to the Request
func WithHeader(key, value string, replace bool) RequestOption {
	return func(r *Request) {
		if r.Header == nil {
			r.Header = http.Header{}
		}
		if replace {
			r.Header.Set(key, value)
		} else {
			r.Header.Add(key, value)
		}
	}
}

// WithHeaders set or replace the headers of the Request
func WithHeaders(header http.Header) RequestOption {
	return func(r *Request) {
		r.Header = header.Clone()
	}
}

func NewRequest(method, endpoint string, secType SecType) *Request {
	return &Request{
		Method:   method,
		Endpoint: endpoint,
		SecType:  secType,
		Query:    url.Values{},
		Form:     url.Values{},
	}
}

func NewGetRequest(endpoint string, secType SecType) *Request {
	return NewRequest(http.MethodGet, endpoint, secType)
}

func NewGetRequestPublic(endpoint string) *Request {
	return NewGetRequest(endpoint, SecTypeNone)
}

func NewGetRequestAPIKey(endpoint string) *Request {
	return NewGetRequest(endpoint, SecTypeAPIKey)
}

func NewGetRequestSigned(endpoint string) *Request {
	return NewGetRequest(endpoint, SecTypeSigned)
}

func NewPostRequest(endpoint string, secType SecType) *Request {
	return NewRequest(http.MethodPost, endpoint, secType)
}

func NewPostRequestSigned(endpoint string) *Request {
	return NewPostRequest(endpoint, SecTypeSigned)
}

func NewPostRequestAPIKey(endpoint string) *Request {
	return NewPostRequest(endpoint, SecTypeAPIKey)
}

func NewPutRequest(endpoint string, secType SecType) *Request {
	return NewRequest(http.MethodPut, endpoint, secType)
}

func NewPutRequestAPIKey(endpoint string) *Request {
	return NewPutRequest(endpoint, SecTypeAPIKey)
}

func NewPutRequestSigned(endpoint string) *Request {
	return NewPutRequest(endpoint, SecTypeSigned)
}

func NewDeleteRequest(endpoint string, secType SecType) *Request {
	return NewRequest(http.MethodDelete, endpoint, secType)
}

func NewDeleteRequestAPIKey(endpoint string) *Request {
	return NewDeleteRequest(endpoint, SecTypeAPIKey)
}

func NewDeleteRequestSigned(endpoint string) *Request {
	return NewDeleteRequest(endpoint, SecTypeSigned)
}
