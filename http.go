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

func (tr *Transport) SendRequest(ctx context.Context, method string, path string, query map[string]interface{}, body map[string]interface{}) (resp *http.Response, err error) {
	req, err := tr.rb.BuildRequest(ctx, method, path, query, body)
	if err != nil {
		return nil, fmt.Errorf("transport@SendRequest: %v", err)
	}
	return tr.http.Do(req)
}

//Get method
func (tr *Transport) Get(ctx context.Context, path string, query map[string]interface{}) (resp *http.Response, err error) {
	return tr.SendRequest(ctx, "GET", path, query, nil)
}

//Post method
func (tr *Transport) Post(ctx context.Context, path string, body map[string]interface{}, query map[string]interface{}) (resp *http.Response, err error) {
	return tr.SendRequest(ctx, "POST", path, query, body)
}

//RequestBuilder handler
type RequestBuilder struct {
	cfg *Config
}

//buildUri method
func (rb *RequestBuilder) buildUri(path string, query map[string]interface{}) (uri *url.URL, err error) {
	u, err := url.Parse(rb.cfg.Uri)
	if err != nil {
		return nil, fmt.Errorf("RequestBuilder@buildUri parse: %v", err)
	}
	u.Path = "/" + path
	u.RawQuery = rb.buildQueryParams(query)
	return u, err
}

//buildQueryParams method
func (rb *RequestBuilder) buildQueryParams(query map[string]interface{}) string {
	q := url.Values{}
	if query != nil {
		for k, v := range query {
			q.Set(k, fmt.Sprintf("%v", v))
		}
	}
	return q.Encode()
}

//buildHeaders method
func (rb *RequestBuilder) buildHeaders() http.Header {
	headers := http.Header{}
	headers.Set("Content-Type", "application/json")
	headers.Set("secret-key", rb.cfg.SecretKey)
	return headers
}

//buildBody method
func (rb *RequestBuilder) buildBody(data map[string]interface{}) (io.Reader, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("RequestBuilder@buildBody json convert: %v", err)
	}
	return bytes.NewBuffer(b), nil
}

//BuildRequest method
func (rb *RequestBuilder) BuildRequest(ctx context.Context, method string, path string, query map[string]interface{}, body map[string]interface{}) (req *http.Request, err error) {
	//build uri
	uri, err := rb.buildUri(path, query)
	if err != nil {
		return nil, fmt.Errorf("transport@request build uri: %v", err)
	}
	//build body
	var bodyReader io.Reader
	if method == "POST" {
		bodyReader, err = rb.buildBody(body)
		if err != nil {
			return nil, fmt.Errorf("transport@request build request body: %v", err)
		}
	}
	//build request
	req, err = http.NewRequestWithContext(ctx, strings.ToUpper(method), uri.String(), bodyReader)
	if err != nil {
		return nil, fmt.Errorf("transport@request new request error: %v", err)
	}
	//build headers
	req.Header = rb.buildHeaders()
	return req, nil
}

//Response wrapper
type Response struct {
	raw *http.Response
}

//IsSuccess method
func (r *Response) IsSuccess() bool {
	return r.raw.StatusCode < http.StatusMultipleChoices
}

//GetRawResponse method
func (r *Response) GetRawResponse() *http.Response {
	return r.raw
}

//GetRawBody method
func (r *Response) GetRawBody() string {
	body, _ := r.ReadBody()
	return string(body)
}

//Unmarshal method
func (r *Response) Unmarshal(v interface{}) error {
	data, err := r.ReadBody()
	if err != nil {
		return fmt.Errorf("Response@Unmarshal read body: %v", err)
	}
	return json.Unmarshal(data, &v)
}

//ReadBody method
func (r *Response) ReadBody() ([]byte, error) {
	defer r.raw.Body.Close()
	return ioutil.ReadAll(r.raw.Body)
}

//NewResponse create new response
func NewResponse(raw *http.Response) *Response {
	return &Response{raw: raw}
}

type ResponseBody struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
}
