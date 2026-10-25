package dusupay

import (
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type HTTPRequestBuilderTestSuite struct {
	suite.Suite
	cfg      *Config
	ctx      context.Context
	testable *RequestBuilder
}

func (suite *HTTPRequestBuilderTestSuite) SetupTest() {
	suite.cfg = BuildStubConfig()
	suite.ctx = context.Background()
	suite.testable = &RequestBuilder{cfg: suite.cfg}
}

func (suite *HTTPRequestBuilderTestSuite) TestBuildURIWithoutQueryParams() {
	uri, err := suite.testable.buildURI("qwerty", nil)
	assert.NotEmpty(suite.T(), uri)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), SandboxAPIUrl+"/qwerty", uri.String())
}

func (suite *HTTPRequestBuilderTestSuite) TestBuildURIWithQueryParams() {
	data := make(map[string]interface{})
	data["foo"] = "bar"
	data["bar"] = "baz"

	uri, err := suite.testable.buildURI("qwerty", data)
	assert.NotEmpty(suite.T(), uri)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), SandboxAPIUrl+"/qwerty?bar=baz&foo=bar", uri.String())
}

func (suite *HTTPRequestBuilderTestSuite) TestBuildHeaders() {
	headers := suite.testable.buildHeaders()
	assert.NotEmpty(suite.T(), headers)
	assert.Equal(suite.T(), "application/json", headers.Get("Content-Type"))
	assert.Equal(suite.T(), "1", headers.Get("x-api-version"))
	assert.Equal(suite.T(), suite.cfg.PublicKey, headers.Get("public-key"))
	assert.Equal(suite.T(), suite.cfg.SecretKey, headers.Get("secret-key"))
	assert.Empty(suite.T(), headers.Get("Accept-Language"))
}

func (suite *HTTPRequestBuilderTestSuite) TestBuildHeadersWithLanguage() {
	suite.cfg.Language = LanguageFrench
	headers := suite.testable.buildHeaders()
	assert.Equal(suite.T(), "fr", headers.Get("Accept-Language"))
}

func (suite *HTTPRequestBuilderTestSuite) TestBuildBody() {
	data := make(map[string]interface{})
	data["foo"] = "bar"
	data["bar"] = "baz"

	body, _ := suite.testable.buildBody(data)
	assert.NotEmpty(suite.T(), body)
}

func (suite *HTTPRequestBuilderTestSuite) TestBuildRequestGET() {
	result, err := suite.testable.BuildRequest(suite.ctx, "get", "foo", map[string]interface{}{"foo": "bar"}, map[string]interface{}{"foo": "bar"})
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.MethodGet, result.Method)
	assert.Equal(suite.T(), SandboxAPIUrl+"/foo?foo=bar", result.URL.String())
	assert.Equal(suite.T(), "application/json", result.Header.Get("Content-Type"))
	assert.Equal(suite.T(), "1", result.Header.Get("x-api-version"))
	assert.Equal(suite.T(), suite.cfg.PublicKey, result.Header.Get("public-key"))
	assert.Equal(suite.T(), suite.cfg.SecretKey, result.Header.Get("secret-key"))
	assert.Nil(suite.T(), result.Body)
}

func (suite *HTTPRequestBuilderTestSuite) TestBuildRequestPOST() {
	result, err := suite.testable.BuildRequest(suite.ctx, "post", "foo", map[string]interface{}{"foo": "bar"}, map[string]interface{}{"foo": "bar"})
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.MethodPost, result.Method)
	assert.Equal(suite.T(), SandboxAPIUrl+"/foo?foo=bar", result.URL.String())
	assert.Equal(suite.T(), "application/json", result.Header.Get("Content-Type"))
	assert.Equal(suite.T(), "1", result.Header.Get("x-api-version"))
	assert.Equal(suite.T(), suite.cfg.PublicKey, result.Header.Get("public-key"))
	assert.Equal(suite.T(), suite.cfg.SecretKey, result.Header.Get("secret-key"))
	assert.NotEmpty(suite.T(), result.Body)
}

func TestHTTPRequestBuilderTestSuite(t *testing.T) {
	suite.Run(t, new(HTTPRequestBuilderTestSuite))
}

type HTTPTestSuite struct {
	suite.Suite
}

func (suite *HTTPTestSuite) TestIsEmptyObjectResponseDataEmptyObject() {
	data := "{}"
	assert.True(suite.T(), isEmptyObjectResponseData([]byte(data)))
}

func (suite *HTTPTestSuite) TestIsEmptyObjectResponseDataFilledObject() {
	data := `{"foo":"bar"}`
	assert.False(suite.T(), isEmptyObjectResponseData([]byte(data)))
}

func (suite *HTTPTestSuite) TestIsEmptyObjectResponseDataEmptyArray() {
	data := "[]"
	assert.False(suite.T(), isEmptyObjectResponseData([]byte(data)))
}

func (suite *HTTPTestSuite) TestResponseBodyIsSuccess() {
	rsp := &ResponseBody{Code: http.StatusAccepted}
	assert.True(suite.T(), rsp.IsSuccess())
	rsp.Code = http.StatusMultipleChoices
	assert.False(suite.T(), rsp.IsSuccess())
	rsp.Code = http.StatusBadRequest
	assert.False(suite.T(), rsp.IsSuccess())
}

func TestHTTPTestSuite(t *testing.T) {
	suite.Run(t, new(HTTPTestSuite))
}

type HTTPTransportTestSuite struct {
	suite.Suite
	cfg      *Config
	ctx      context.Context
	testable *Transport
}

func (suite *HTTPTransportTestSuite) SetupTest() {
	suite.cfg = BuildStubConfig()
	suite.ctx = context.Background()
	suite.testable = NewHTTPTransport(suite.cfg, &http.Client{})
	httpmock.Activate()
}

func (suite *HTTPTransportTestSuite) TearDownTest() {
	httpmock.DeactivateAndReset()
}

func (suite *HTTPTransportTestSuite) TestSendRequestSuccess() {
	body, _ := LoadStubResponseData("stubs/merchants/balance/success.json")
	httpmock.RegisterResponder(http.MethodGet, suite.cfg.URI+"/foo", httpmock.NewBytesResponder(http.StatusOK, body))

	resp, err := suite.testable.SendRequest(suite.ctx, http.MethodGet, "foo", nil, nil)

	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), resp)

	defer func() { _ = resp.Body.Close() }()
	bodyRsp, _ := io.ReadAll(resp.Body)
	assert.Equal(suite.T(), body, bodyRsp)
}

func (suite *HTTPTransportTestSuite) TestGet() {
	body, _ := LoadStubResponseData("stubs/merchants/balance/success.json")
	httpmock.RegisterResponder(http.MethodGet, suite.cfg.URI+"/foo", httpmock.NewBytesResponder(http.StatusOK, body))

	resp, err := suite.testable.Get(suite.ctx, "foo", nil)

	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), resp)

	defer func() { _ = resp.Body.Close() }()
	bodyRsp, _ := io.ReadAll(resp.Body)
	assert.Equal(suite.T(), body, bodyRsp)
}

func (suite *HTTPTransportTestSuite) TestPost() {
	body, _ := LoadStubResponseData("stubs/merchants/balance/success.json")
	httpmock.RegisterResponder(http.MethodPost, suite.cfg.URI+"/foo", httpmock.NewBytesResponder(http.StatusOK, body))

	resp, err := suite.testable.Post(suite.ctx, "foo", nil, nil)

	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), resp)

	defer func() { _ = resp.Body.Close() }()
	bodyRsp, _ := io.ReadAll(resp.Body)
	assert.Equal(suite.T(), body, bodyRsp)
}

func TestHTTPTransportTestSuite(t *testing.T) {
	suite.Run(t, new(HTTPTransportTestSuite))
}
