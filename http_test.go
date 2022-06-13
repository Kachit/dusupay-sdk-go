package dusupay

import (
	"context"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"net/http"
	"testing"
)

type HttpRequestBuilderTestSuite struct {
	suite.Suite
	cfg      *Config
	ctx      context.Context
	testable *RequestBuilder
}

func (suite *HttpRequestBuilderTestSuite) SetupTest() {
	suite.cfg = BuildStubConfig()
	suite.ctx = context.Background()
	suite.testable = &RequestBuilder{cfg: suite.cfg}
}

func (suite *HttpRequestBuilderTestSuite) TestBuildUriWithoutQueryParams() {
	uri, err := suite.testable.buildUri("qwerty", nil)
	assert.NotEmpty(suite.T(), uri)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), SandboxAPIUrl+"/qwerty", uri.String())
}

func (suite *HttpRequestBuilderTestSuite) TestBuildUriWithQueryParams() {
	data := make(map[string]interface{})
	data["foo"] = "bar"
	data["bar"] = "baz"

	uri, err := suite.testable.buildUri("qwerty", data)
	assert.NotEmpty(suite.T(), uri)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), SandboxAPIUrl+"/qwerty?bar=baz&foo=bar", uri.String())
}

func (suite *HttpRequestBuilderTestSuite) TestBuildHeaders() {
	headers := suite.testable.buildHeaders()
	assert.NotEmpty(suite.T(), headers)
	assert.Equal(suite.T(), "application/json", headers.Get("Content-Type"))
	assert.Equal(suite.T(), suite.cfg.SecretKey, headers.Get("secret-key"))
}

func (suite *HttpRequestBuilderTestSuite) TestBuildBody() {
	data := make(map[string]interface{})
	data["foo"] = "bar"
	data["bar"] = "baz"

	body, _ := suite.testable.buildBody(data)
	assert.NotEmpty(suite.T(), body)
}

func (suite *HttpRequestBuilderTestSuite) TestBuildAuthParams() {
	data := make(map[string]interface{})
	data["foo"] = "bar"
	data["bar"] = "baz"

	result := suite.testable.buildAuthParams(data)
	assert.Equal(suite.T(), data["foo"], result["foo"])
	assert.Equal(suite.T(), data["bar"], result["bar"])
	assert.Equal(suite.T(), suite.cfg.PublicKey, result["api_key"])
}

func (suite *HttpRequestBuilderTestSuite) TestBuildAuthParamsEmpty() {
	result := suite.testable.buildAuthParams(nil)
	assert.Equal(suite.T(), suite.cfg.PublicKey, result["api_key"])
}

func (suite *HttpRequestBuilderTestSuite) TestBuildRequestGET() {
	result, err := suite.testable.BuildRequest(suite.ctx, "get", "foo", map[string]interface{}{"foo": "bar"}, map[string]interface{}{"foo": "bar"})
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.MethodGet, result.Method)
	assert.Equal(suite.T(), "https://sandbox.dusupay.com/foo?api_key=PublicKey&foo=bar", result.URL.String())
	assert.Equal(suite.T(), "application/json", result.Header.Get("Content-Type"))
	assert.Equal(suite.T(), suite.cfg.SecretKey, result.Header.Get("secret-key"))
	assert.Nil(suite.T(), result.Body)
}

func (suite *HttpRequestBuilderTestSuite) TestBuildRequestPOST() {
	result, err := suite.testable.BuildRequest(suite.ctx, "post", "foo", map[string]interface{}{"foo": "bar"}, map[string]interface{}{"foo": "bar"})
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.MethodPost, result.Method)
	assert.Equal(suite.T(), "https://sandbox.dusupay.com/foo?foo=bar", result.URL.String())
	assert.Equal(suite.T(), "application/json", result.Header.Get("Content-Type"))
	assert.Equal(suite.T(), suite.cfg.SecretKey, result.Header.Get("secret-key"))
	assert.NotEmpty(suite.T(), result.Body)
}

func TestHttpRequestBuilderTestSuite(t *testing.T) {
	suite.Run(t, new(HttpRequestBuilderTestSuite))
}

type HttpTestSuite struct {
	suite.Suite
}

func (suite *HttpTestSuite) TestIsEmptyObjectResponseDataEmptyObject() {
	data := "{}"
	assert.True(suite.T(), isEmptyObjectResponseData([]byte(data)))
}

func (suite *HttpTestSuite) TestIsEmptyObjectResponseDataFilledObject() {
	data := `{"foo":"bar"}`
	assert.False(suite.T(), isEmptyObjectResponseData([]byte(data)))
}

func (suite *HttpTestSuite) TestIsEmptyObjectResponseDataEmptyArray() {
	data := "[]"
	assert.False(suite.T(), isEmptyObjectResponseData([]byte(data)))
}

func (suite *HttpTestSuite) TestResponseBodyIsSuccess() {
	rsp := &ResponseBody{Code: http.StatusAccepted}
	assert.True(suite.T(), rsp.IsSuccess())
	rsp.Code = http.StatusMultipleChoices
	assert.False(suite.T(), rsp.IsSuccess())
	rsp.Code = http.StatusBadRequest
	assert.False(suite.T(), rsp.IsSuccess())
}

func TestHttpTestSuite(t *testing.T) {
	suite.Run(t, new(HttpTestSuite))
}

type HttpTransportTestSuite struct {
	suite.Suite
	cfg      *Config
	ctx      context.Context
	testable *Transport
}

func (suite *HttpTransportTestSuite) SetupTest() {
	suite.cfg = BuildStubConfig()
	suite.ctx = context.Background()
	suite.testable = NewHttpTransport(suite.cfg, &http.Client{})
	httpmock.Activate()
}

func (suite *HttpTransportTestSuite) TearDownTest() {
	httpmock.DeactivateAndReset()
}

func (suite *HttpTransportTestSuite) TestSendRequestSuccess() {
	body, _ := LoadStubResponseData("stubs/merchants/balance/success.json")
	httpmock.RegisterResponder(http.MethodGet, suite.cfg.Uri+"/foo?api_key=PublicKey", httpmock.NewBytesResponder(http.StatusOK, body))

	resp, err := suite.testable.SendRequest(suite.ctx, http.MethodGet, "foo", nil, nil)

	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), resp)

	defer resp.Body.Close()
	bodyRsp, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(suite.T(), body, bodyRsp)
}

func (suite *HttpTransportTestSuite) TestGet() {
	body, _ := LoadStubResponseData("stubs/merchants/balance/success.json")
	httpmock.RegisterResponder(http.MethodGet, suite.cfg.Uri+"/foo?api_key=PublicKey", httpmock.NewBytesResponder(http.StatusOK, body))

	resp, err := suite.testable.Get(suite.ctx, "foo", nil)

	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), resp)

	defer resp.Body.Close()
	bodyRsp, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(suite.T(), body, bodyRsp)
}

func (suite *HttpTransportTestSuite) TestPost() {
	body, _ := LoadStubResponseData("stubs/merchants/balance/success.json")
	httpmock.RegisterResponder(http.MethodPost, suite.cfg.Uri+"/foo", httpmock.NewBytesResponder(http.StatusOK, body))

	resp, err := suite.testable.Post(suite.ctx, "foo", nil, nil)

	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), resp)

	defer resp.Body.Close()
	bodyRsp, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(suite.T(), body, bodyRsp)
}

func TestHttpTransportTestSuite(t *testing.T) {
	suite.Run(t, new(HttpTransportTestSuite))
}
