package dusupay

import (
	"context"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

func Test_HTTP_RequestBuilder_BuildUriWithoutQueryParams(t *testing.T) {
	cfg := BuildStubConfig()
	builder := RequestBuilder{cfg: cfg}
	uri, err := builder.buildUri("qwerty", nil)
	assert.NotEmpty(t, uri)
	assert.Nil(t, err)
	assert.Equal(t, SandboxAPIUrl+"/qwerty", uri.String())
}

func Test_HTTP_RequestBuilder_BuildUriWithQueryParams(t *testing.T) {
	cfg := BuildStubConfig()
	builder := RequestBuilder{cfg: cfg}

	data := make(map[string]interface{})
	data["foo"] = "bar"
	data["bar"] = "baz"

	uri, err := builder.buildUri("qwerty", data)
	assert.NotEmpty(t, uri)
	assert.Nil(t, err)
	assert.Equal(t, SandboxAPIUrl+"/qwerty?bar=baz&foo=bar", uri.String())
}

func Test_HTTP_RequestBuilder_BuildHeaders(t *testing.T) {
	cfg := BuildStubConfig()
	builder := RequestBuilder{cfg: cfg}

	headers := builder.buildHeaders()
	assert.NotEmpty(t, headers)
	assert.Equal(t, "application/json", headers.Get("Content-Type"))
	assert.Equal(t, cfg.SecretKey, headers.Get("secret-key"))
}

func Test_HTTP_RequestBuilder_BuildBody(t *testing.T) {
	cfg := BuildStubConfig()
	builder := RequestBuilder{cfg: cfg}

	data := make(map[string]interface{})
	data["foo"] = "bar"
	data["bar"] = "baz"

	body, _ := builder.buildBody(data)
	assert.NotEmpty(t, body)
}

func Test_HTTP_RequestBuilder_BuildAuthParams(t *testing.T) {
	cfg := BuildStubConfig()
	builder := RequestBuilder{cfg: cfg}

	data := make(map[string]interface{})
	data["foo"] = "bar"
	data["bar"] = "baz"

	result := builder.buildAuthParams(data)
	assert.Equal(t, data["foo"], result["foo"])
	assert.Equal(t, data["bar"], result["bar"])
	assert.Equal(t, cfg.PublicKey, result["api_key"])
}

func Test_HTTP_RequestBuilder_BuildAuthParamsEmpty(t *testing.T) {
	cfg := BuildStubConfig()
	builder := RequestBuilder{cfg: cfg}

	result := builder.buildAuthParams(nil)
	assert.Equal(t, cfg.PublicKey, result["api_key"])
}

func Test_HTTP_RequestBuilder_BuildRequestGET(t *testing.T) {
	cfg := BuildStubConfig()
	builder := RequestBuilder{cfg: cfg}

	ctx := context.Background()
	result, err := builder.BuildRequest(ctx, "get", "foo", map[string]interface{}{"foo": "bar"}, map[string]interface{}{"foo": "bar"})
	assert.NoError(t, err)
	assert.Equal(t, http.MethodGet, result.Method)
	assert.Equal(t, "https://sandbox.dusupay.com/foo?api_key=PublicKey&foo=bar", result.URL.String())
	assert.Equal(t, "application/json", result.Header.Get("Content-Type"))
	assert.Equal(t, cfg.SecretKey, result.Header.Get("secret-key"))
	assert.Nil(t, result.Body)
}

func Test_HTTP_RequestBuilder_BuildRequestPOST(t *testing.T) {
	cfg := BuildStubConfig()
	builder := RequestBuilder{cfg: cfg}

	ctx := context.Background()
	result, err := builder.BuildRequest(ctx, "post", "foo", map[string]interface{}{"foo": "bar"}, map[string]interface{}{"foo": "bar"})
	assert.NoError(t, err)
	assert.Equal(t, http.MethodPost, result.Method)
	assert.Equal(t, "https://sandbox.dusupay.com/foo?foo=bar", result.URL.String())
	assert.Equal(t, "application/json", result.Header.Get("Content-Type"))
	assert.Equal(t, cfg.SecretKey, result.Header.Get("secret-key"))
	assert.NotEmpty(t, result.Body)
}

func Test_HTTP_IsEmptyObjectResponseDataEmptyObject(t *testing.T) {
	data := "{}"
	assert.True(t, isEmptyObjectResponseData([]byte(data)))
}

func Test_HTTP_IsEmptyObjectResponseDataFilledObject(t *testing.T) {
	data := `{"foo":"bar"}`
	assert.False(t, isEmptyObjectResponseData([]byte(data)))
}

func Test_HTTP_IsEmptyObjectResponseDataEmptyArray(t *testing.T) {
	data := "[]"
	assert.False(t, isEmptyObjectResponseData([]byte(data)))
}

func Test_HTTP_NewHttpTransport(t *testing.T) {
	cfg := BuildStubConfig()
	transport := NewHttpTransport(cfg, nil)
	assert.NotEmpty(t, transport)
}

func Test_HTTP_Transport_SendRequestSuccess(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	cfg := BuildStubConfig()
	transport := BuildStubHttpTransport()

	body, _ := LoadStubResponseData("stubs/merchants/balance/success.json")
	httpmock.RegisterResponder(http.MethodGet, cfg.Uri+"/foo?api_key=PublicKey", httpmock.NewBytesResponder(http.StatusOK, body))

	ctx := context.Background()
	resp, err := transport.SendRequest(ctx, http.MethodGet, "foo", nil, nil)

	assert.NoError(t, err)
	assert.NotEmpty(t, resp)

	defer resp.Body.Close()
	bodyRsp, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, body, bodyRsp)
}

func Test_HTTP_Transport_Get(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	cfg := BuildStubConfig()
	transport := BuildStubHttpTransport()

	body, _ := LoadStubResponseData("stubs/merchants/balance/success.json")
	httpmock.RegisterResponder(http.MethodGet, cfg.Uri+"/foo?api_key=PublicKey", httpmock.NewBytesResponder(http.StatusOK, body))

	ctx := context.Background()
	resp, err := transport.Get(ctx, "foo", nil)

	assert.NoError(t, err)
	assert.NotEmpty(t, resp)

	defer resp.Body.Close()
	bodyRsp, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, body, bodyRsp)
}

func Test_HTTP_Transport_Post(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	cfg := BuildStubConfig()
	transport := BuildStubHttpTransport()

	body, _ := LoadStubResponseData("stubs/merchants/balance/success.json")
	httpmock.RegisterResponder(http.MethodPost, cfg.Uri+"/foo", httpmock.NewBytesResponder(http.StatusOK, body))

	ctx := context.Background()
	resp, err := transport.Post(ctx, "foo", nil, nil)

	assert.NoError(t, err)
	assert.NotEmpty(t, resp)

	defer resp.Body.Close()
	bodyRsp, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, body, bodyRsp)
}

func Test_HTTP_ResponseBody_IsSuccess(t *testing.T) {
	rsp := &ResponseBody{Code: http.StatusAccepted}
	assert.True(t, rsp.IsSuccess())
	rsp.Code = http.StatusMultipleChoices
	assert.False(t, rsp.IsSuccess())
	rsp.Code = http.StatusBadRequest
	assert.False(t, rsp.IsSuccess())
}
