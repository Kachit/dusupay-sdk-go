package dusupay

import (
	"context"
	"github.com/stretchr/testify/assert"
	"net/http"

	//"net/http"
	"testing"
)

func Test_HTTP_RequestBuilder_BuildUriWithoutQueryParams(t *testing.T) {
	cfg := BuildStubConfig()
	builder := RequestBuilder{cfg: cfg}
	uri, err := builder.buildUri("qwerty", nil)
	assert.NotEmpty(t, uri)
	assert.Equal(t, SandboxAPIUrl+"/qwerty", uri.String())
	assert.Nil(t, err)
}

func Test_HTTP_RequestBuilder_BuildUriWithQueryParams(t *testing.T) {
	cfg := BuildStubConfig()
	builder := RequestBuilder{cfg: cfg}

	data := make(map[string]interface{})
	data["foo"] = "bar"
	data["bar"] = "baz"

	uri, err := builder.buildUri("qwerty", data)
	assert.NotEmpty(t, uri)
	assert.Equal(t, SandboxAPIUrl+"/qwerty?bar=baz&foo=bar", uri.String())
	assert.Nil(t, err)
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

func Test_HTTP_NewHttpTransport(t *testing.T) {
	cfg := BuildStubConfig()
	transport := NewHttpTransport(cfg, nil)
	assert.NotEmpty(t, transport)
}
