package dusupay

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// NewHTTPTransport create new http transport
func NewHTTPTransport(config *Config, h *http.Client) *Transport {
	rb := &RequestBuilder{cfg: config}
	return &Transport{http: h, rb: rb}
}

// Transport wrapper
type Transport struct {
	http *http.Client
	rb   *RequestBuilder
}

// SendRequest Send request method
func (tr *Transport) SendRequest(ctx context.Context, method string, path string, query map[string]interface{}, body map[string]interface{}) (resp *http.Response, err error) {
	req, err := tr.rb.BuildRequest(ctx, method, path, query, body)
	if err != nil {
		return nil, fmt.Errorf("transport.SendRequest: %v", err)
	}
	return tr.http.Do(req)
}

// Get method
func (tr *Transport) Get(ctx context.Context, path string, query map[string]interface{}) (resp *http.Response, err error) {
	return tr.SendRequest(ctx, http.MethodGet, path, query, nil)
}

// Post method
func (tr *Transport) Post(ctx context.Context, path string, body map[string]interface{}, query map[string]interface{}) (resp *http.Response, err error) {
	return tr.SendRequest(ctx, http.MethodPost, path, query, body)
}

// RequestBuilder handler
type RequestBuilder struct {
	cfg *Config
}

// buildURI method
func (rb *RequestBuilder) buildURI(path string, query map[string]interface{}) (uri *url.URL, err error) {
	u, err := url.Parse(rb.cfg.URI)
	if err != nil {
		return nil, fmt.Errorf("RequestBuilder.buildURI parse: %v", err)
	}
	u.Path = "/" + path
	u.RawQuery = rb.buildQueryParams(query)
	return u, err
}

// buildQueryParams method
func (rb *RequestBuilder) buildQueryParams(query map[string]interface{}) string {
	q := url.Values{}
	for k, v := range query {
		q.Set(k, fmt.Sprintf("%v", v))
	}
	return q.Encode()
}

// buildHeaders method
func (rb *RequestBuilder) buildHeaders() http.Header {
	headers := http.Header{}
	headers.Set("Content-Type", "application/json")
	headers.Set("x-api-version", "1")
	headers.Set("public-key", rb.cfg.PublicKey)
	headers.Set("secret-key", rb.cfg.SecretKey)
	return headers
}

// buildBody method
func (rb *RequestBuilder) buildBody(data map[string]interface{}) (io.Reader, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("RequestBuilder.buildBody json convert: %v", err)
	}
	return bytes.NewBuffer(b), nil
}

// BuildRequest method
func (rb *RequestBuilder) BuildRequest(ctx context.Context, method string, path string, query map[string]interface{}, body map[string]interface{}) (req *http.Request, err error) {
	method = strings.ToUpper(method)
	var bodyReader io.Reader
	if method == http.MethodPost {
		bodyReader, err = rb.buildBody(body)
		if err != nil {
			return nil, fmt.Errorf("transport.request build request body: %v", err)
		}
	}
	uri, err := rb.buildURI(path, query)
	if err != nil {
		return nil, fmt.Errorf("transport.request build uri: %v", err)
	}
	req, err = http.NewRequestWithContext(ctx, method, uri.String(), bodyReader)
	if err != nil {
		return nil, fmt.Errorf("transport.request new request error: %v", err)
	}
	req.Header = rb.buildHeaders()
	return req, nil
}

// ResponseBody struct
type ResponseBody struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

// IsSuccess method
func (r *ResponseBody) IsSuccess() bool {
	return r.Code < http.StatusMultipleChoices
}

// unmarshalResponse func
func unmarshalResponse(resp *http.Response, v interface{}) error {
	defer func() { _ = resp.Body.Close() }()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Response.Unmarshal read body: %v", err)
	}
	resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	return json.Unmarshal(bodyBytes, &v)
}

// isEmptyObjectResponseData func
func isEmptyObjectResponseData(data []byte) bool {
	return data[0] == 123 && data[1] == 125
}
