package dusupay

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

//NewHttpTransport create new http transport
func NewHttpTransport(config *Config, h *http.Client) *Transport {
	rb := &RequestBuilder{cfg: config}
	return &Transport{http: h, rb: rb}
}

//Transport wrapper
type Transport struct {
	http *http.Client
	rb   *RequestBuilder
}

//SendRequest Send request method
func (tr *Transport) SendRequest(ctx context.Context, method string, path string, query map[string]interface{}, body map[string]interface{}) (resp *http.Response, err error) {
	req, err := tr.rb.BuildRequest(ctx, method, path, query, body)
	if err != nil {
		return nil, fmt.Errorf("transport.SendRequest: %v", err)
	}
	return tr.http.Do(req)
}

//Get method
func (tr *Transport) Get(ctx context.Context, path string, query map[string]interface{}) (resp *http.Response, err error) {
	return tr.SendRequest(ctx, http.MethodGet, path, query, nil)
}

//Post method
func (tr *Transport) Post(ctx context.Context, path string, body map[string]interface{}, query map[string]interface{}) (resp *http.Response, err error) {
	return tr.SendRequest(ctx, http.MethodPost, path, query, body)
}

//RequestBuilder handler
type RequestBuilder struct {
	cfg *Config
}

//BuildUri method
func (rb *RequestBuilder) buildUri(path string, query map[string]interface{}) (uri *url.URL, err error) {
	u, err := url.Parse(rb.cfg.Uri)
	if err != nil {
		return nil, fmt.Errorf("RequestBuilder.buildUri parse: %v", err)
	}
	u.Path = "/" + path
	u.RawQuery = rb.buildQueryParams(query)
	return u, err
}

//BuildQueryParams method
func (rb *RequestBuilder) buildQueryParams(query map[string]interface{}) string {
	q := url.Values{}
	if query != nil {
		for k, v := range query {
			q.Set(k, fmt.Sprintf("%v", v))
		}
	}
	return q.Encode()
}

//BuildHeaders method
func (rb *RequestBuilder) buildHeaders() http.Header {
	headers := http.Header{}
	headers.Set("Content-Type", "application/json")
	headers.Set("secret-key", rb.cfg.SecretKey)
	return headers
}

//BuildAuthParams method
func (rb *RequestBuilder) buildAuthParams(params map[string]interface{}) map[string]interface{} {
	if params == nil {
		params = make(map[string]interface{})
	}
	params["api_key"] = rb.cfg.PublicKey
	return params
}

//BuildBody method
func (rb *RequestBuilder) buildBody(data map[string]interface{}) (io.Reader, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("RequestBuilder.buildBody json convert: %v", err)
	}
	return bytes.NewBuffer(b), nil
}

//BuildRequest method
func (rb *RequestBuilder) BuildRequest(ctx context.Context, method string, path string, query map[string]interface{}, body map[string]interface{}) (req *http.Request, err error) {
	method = strings.ToUpper(method)
	//build body
	var bodyReader io.Reader
	if method == http.MethodPost {
		body = rb.buildAuthParams(body)
		bodyReader, err = rb.buildBody(body)
		if err != nil {
			return nil, fmt.Errorf("transport.request build request body: %v", err)
		}
	} else {
		query = rb.buildAuthParams(query)
	}
	//build uri
	uri, err := rb.buildUri(path, query)
	if err != nil {
		return nil, fmt.Errorf("transport.request build uri: %v", err)
	}
	//build request
	req, err = http.NewRequestWithContext(ctx, method, uri.String(), bodyReader)
	if err != nil {
		return nil, fmt.Errorf("transport.request new request error: %v", err)
	}
	//build headers
	req.Header = rb.buildHeaders()
	return req, nil
}

//ResponseBody struct
type ResponseBody struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

//IsSuccess method
func (r *ResponseBody) IsSuccess() bool {
	return r.Code < http.StatusMultipleChoices
}

//UnmarshalResponse func
func unmarshalResponse(resp *http.Response, v interface{}) error {
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Response.Unmarshal read body: %v", err)
	}
	//reset the response body to the original unread state
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	return json.Unmarshal(bodyBytes, &v)
}

//IsEmptyObjectResponseData func
func isEmptyObjectResponseData(data []byte) bool {
	return data[0] == 123 && data[1] == 125
}
