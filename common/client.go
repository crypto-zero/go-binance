package common

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync/atomic"
	"time"
)

const (
	timestampKey  = "timestamp"
	signatureKey  = "signature"
	recvWindowKey = "recvWindow"
)

func currentTimestamp() int64 {
	return FormatTimestamp(time.Now())
}

// FormatTimestamp formats a time into Unix timestamp in milliseconds, as requested by Binance.
func FormatTimestamp(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond)
}

type DoFunc func(req *http.Request) (*http.Response, error)

type client struct {
	globalRequestID uint64
	apiKey          string
	secretKey       string
	baseURL         string
	userAgent       string
	httpClient      *http.Client
	logger          Logger
	timeOffset      int64
	do              DoFunc
}

func (c *client) GetTimeOffset() int64 {
	return c.timeOffset
}

func (c *client) UpdateTimeOffset(offset int64) {
	c.timeOffset = offset
}

func (c *client) UpdateDoFunc(f DoFunc) {
	c.do = f
}

func (c *client) UpdateHTTPClient(hc *http.Client) {
	c.httpClient = hc
}

func (c *client) prepareRequest(r *Request, opts ...RequestOption) (bodyString,
	fullURL string, header http.Header, err error,
) {
	// set Request options from user
	for _, opt := range opts {
		opt(r)
	}
	if err = r.Validate(); err != nil {
		return "", "", nil, err
	}

	fURL, err := url.Parse(c.baseURL)
	if err != nil {
		return "", "", nil, err
	}
	fURL.Path = r.Endpoint

	if r.RecvWindow > 0 {
		r.SetQuery(recvWindowKey, r.RecvWindow)
	}
	if r.SecType == SecTypeSigned {
		r.SetQuery(timestampKey, currentTimestamp()-c.timeOffset)
	}

	queryString := r.Query.Encode()
	bodyString = r.Form.Encode()
	header = http.Header{}

	if r.Header != nil {
		header = r.Header.Clone()
	}
	if bodyString != "" {
		header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	if r.SecType == SecTypeAPIKey || r.SecType == SecTypeSigned {
		header.Set("X-MBX-APIKEY", c.apiKey)
	}
	if r.SecType == SecTypeSigned {
		raw := fmt.Sprintf("%s%s", queryString, bodyString)
		mac := hmac.New(sha256.New, []byte(c.secretKey))
		_, err = mac.Write([]byte(raw))
		if err != nil {
			return "", "", nil, err
		}
		signature := fmt.Sprintf("%x", mac.Sum(nil))
		// NOTE: The signature pair MUST be appended to the last of query string.
		if queryString != "" {
			queryString += "&"
		}
		queryString = fmt.Sprintf("%s%s=%s", queryString, signatureKey, signature)
	}

	fURL.RawQuery = queryString
	r.ID = atomic.AddUint64(&c.globalRequestID, 1)
	return bodyString, fURL.String(), header, nil
}

func (c *client) CallAPIBytes(ctx context.Context, r *Request, opts ...RequestOption) (
	data []byte, err error,
) {
	body, fullURL, headers, err := c.prepareRequest(r, opts...)
	if err != nil {
		return nil, err
	}

	var inBody io.Reader
	if body != "" {
		inBody = bytes.NewBufferString(body)
	}
	req, err := http.NewRequest(r.Method, fullURL, inBody)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	req.Header = headers

	c.logger.Debugw("call api prepare", "id", r.ID, "url", fullURL, "body", body)

	f := c.do
	if f == nil {
		f = c.httpClient.Do
	}

	res, err := f(req)
	if err != nil {
		return nil, err
	}

	data, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	defer func() {
		// Only overwrite the returned error if the original error was nil and an
		// error occurred while closing the body.
		if cerr := res.Body.Close(); err == nil && cerr != nil {
			err = cerr
		}
	}()

	c.logger.Debugw("call api reply", "id", r.ID, "status_code", res.StatusCode,
		"response_headers", res.Header, "response_body", string(data))

	if res.StatusCode >= 400 {
		apiErr := &APIError{Status: res.StatusCode}
		if e := json.Unmarshal(data, apiErr); e != nil {
			c.logger.Debugw("call api parse error failed", "id", r.ID, "err", e)
		}
		return nil, apiErr
	}
	return data, nil
}

func (c *client) CallAPI(ctx context.Context, r *Request, result interface{},
	opts ...RequestOption,
) (err error) {
	data, err := c.CallAPIBytes(ctx, r, opts...)
	if err != nil {
		return err
	}
	if result == nil {
		return nil
	}

	f, ok := result.(func(data []byte) error)
	if ok {
		if err = f(data); err != nil {
			return err
		}
		return nil
	}
	if err = json.Unmarshal(data, result); err != nil {
		return err
	}
	return nil
}

type Client interface {
	GetTimeOffset() int64
	UpdateTimeOffset(offset int64)
	UpdateDoFunc(f DoFunc)
	UpdateHTTPClient(hc *http.Client)
	CallAPIBytes(ctx context.Context, r *Request, opts ...RequestOption) (data []byte, err error)
	CallAPI(ctx context.Context, r *Request, result interface{}, opts ...RequestOption) (err error)
}

func NewClient(apiKey, secretKey, baseURL, userAgent string, httpClient *http.Client,
	logger Logger,
) Client {
	return &client{
		globalRequestID: 0,
		apiKey:          apiKey,
		secretKey:       secretKey,
		baseURL:         baseURL,
		userAgent:       userAgent,
		httpClient:      httpClient,
		logger:          logger,
	}
}
