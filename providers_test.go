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

type ProvidersTestSuite struct {
	suite.Suite
}

func (suite *ProvidersTestSuite) TestProvidersFilterIsValidSuccess() {
	filter := ProvidersFilter{Country: CountryCodeKenya, Method: TransactionMethodCard, TransactionType: TransactionTypeCollection}
	assert.Nil(suite.T(), filter.isValid())
	assert.NoError(suite.T(), filter.isValid())
}

func (suite *ProvidersTestSuite) TestProvidersFilterIsValidEmptyCountryCode() {
	filter := ProvidersFilter{Method: TransactionMethodCard}
	result := filter.isValid()
	assert.Error(suite.T(), result)
	assert.Equal(suite.T(), `parameter "country_code" is empty`, result.Error())
}

func (suite *ProvidersTestSuite) TestProvidersFilterIsValidEmptyMethod() {
	filter := ProvidersFilter{Country: CountryCodeKenya}
	result := filter.isValid()
	assert.Error(suite.T(), result)
	assert.Equal(suite.T(), `parameter "method" is empty`, result.Error())
}

func (suite *ProvidersTestSuite) TestProvidersFilterIsValidEmptyTransactionType() {
	filter := ProvidersFilter{Country: CountryCodeKenya, Method: TransactionMethodCard}
	result := filter.isValid()
	assert.Error(suite.T(), result)
	assert.Equal(suite.T(), `parameter "transaction_type" is empty`, result.Error())
}

func (suite *ProvidersTestSuite) TestProvidersFilterBuildPath() {
	filter := ProvidersFilter{Country: CountryCodeKenya, Method: TransactionMethodCard, TransactionType: TransactionTypeCollection}
	result := filter.buildPath()
	assert.Equal(suite.T(), `collection/card/ke`, result)
}

func TestProvidersTestSuite(t *testing.T) {
	suite.Run(t, new(ProvidersTestSuite))
}

type ProvidersResourceTestSuite struct {
	suite.Suite
	cfg      *Config
	ctx      context.Context
	testable *ProvidersResource
}

func (suite *ProvidersResourceTestSuite) SetupTest() {
	cfg := BuildStubConfig()
	transport := BuildStubHttpTransport()
	suite.cfg = cfg
	suite.ctx = context.Background()
	suite.testable = &ProvidersResource{NewResourceAbstract(transport, cfg)}
	httpmock.Activate()
}

func (suite *ProvidersResourceTestSuite) TearDownTest() {
	httpmock.DeactivateAndReset()
}

func (suite *ProvidersResourceTestSuite) TestGetListSuccess() {
	body, _ := LoadStubResponseData("stubs/providers/payment-options/success.json")
	httpmock.RegisterResponder(http.MethodGet, suite.cfg.Uri+"/v1/payment-options/collection/card/ke", httpmock.NewBytesResponder(http.StatusOK, body))

	filter := &ProvidersFilter{Country: CountryCodeKenya, Method: TransactionMethodCard, TransactionType: TransactionTypeCollection}
	result, resp, err := suite.testable.GetList(suite.ctx, filter)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), resp)
	assert.NotEmpty(suite.T(), result)
	//result
	assert.True(suite.T(), result.IsSuccess())
	assert.Equal(suite.T(), http.StatusOK, result.Code)
	assert.Equal(suite.T(), "success", result.Status)
	assert.Equal(suite.T(), "Request completed successfully.", result.Message)
	assert.Equal(suite.T(), "mtn_ug", (*result.Data)[0].ID)
	assert.Equal(suite.T(), "MTN Mobile Money", (*result.Data)[0].Name)
	assert.Equal(suite.T(), "UGX", (*result.Data)[0].TransactionCurrency)
	assert.Equal(suite.T(), float64(3000), (*result.Data)[0].MinAmount)
	assert.Equal(suite.T(), float64(5000000), (*result.Data)[0].MaxAmount)
	assert.Equal(suite.T(), true, (*result.Data)[0].Available)
	assert.Empty(suite.T(), (*result.Data)[0].SandboxTestAccounts)
	//response
	defer resp.Body.Close()
	bodyRsp, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(suite.T(), body, bodyRsp)
}

func (suite *ProvidersResourceTestSuite) TestGetListSuccessSandbox() {
	body, _ := LoadStubResponseData("stubs/providers/payment-options/success-sandbox.json")
	httpmock.RegisterResponder(http.MethodGet, suite.cfg.Uri+"/v1/payment-options/collection/card/ke", httpmock.NewBytesResponder(http.StatusOK, body))

	filter := &ProvidersFilter{Country: CountryCodeKenya, Method: TransactionMethodCard, TransactionType: TransactionTypeCollection}
	result, resp, err := suite.testable.GetList(suite.ctx, filter)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), resp)
	assert.NotEmpty(suite.T(), result)
	//result
	assert.True(suite.T(), result.IsSuccess())
	assert.Equal(suite.T(), http.StatusOK, result.Code)
	assert.Equal(suite.T(), "success", result.Status)
	assert.Equal(suite.T(), "Request completed successfully.", result.Message)
	assert.Equal(suite.T(), "mtn_ug", (*result.Data)[0].ID)
	assert.Equal(suite.T(), "MTN Mobile Money", (*result.Data)[0].Name)
	assert.Equal(suite.T(), "UGX", (*result.Data)[0].TransactionCurrency)
	assert.Equal(suite.T(), float64(3000), (*result.Data)[0].MinAmount)
	assert.Equal(suite.T(), float64(5000000), (*result.Data)[0].MaxAmount)
	assert.Equal(suite.T(), true, (*result.Data)[0].Available)
	assert.Equal(suite.T(), "256777000456", (*result.Data)[0].SandboxTestAccounts.Failure)
	assert.Equal(suite.T(), "256777000123", (*result.Data)[0].SandboxTestAccounts.Success)
	//response
	defer resp.Body.Close()
	bodyRsp, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(suite.T(), body, bodyRsp)
}

func (suite *ProvidersResourceTestSuite) TestGetListJsonError() {
	body, _ := LoadStubResponseData("stubs/errors/401.json")
	httpmock.RegisterResponder(http.MethodGet, suite.cfg.Uri+"/v1/payment-options/collection/card/ke", httpmock.NewBytesResponder(http.StatusOK, body))

	filter := &ProvidersFilter{Country: CountryCodeKenya, Method: TransactionMethodCard, TransactionType: TransactionTypeCollection}
	result, resp, err := suite.testable.GetList(suite.ctx, filter)
	assert.Error(suite.T(), err)
	assert.NotEmpty(suite.T(), resp)
	assert.NotEmpty(suite.T(), result)
	//result
	assert.False(suite.T(), result.IsSuccess())
	assert.Equal(suite.T(), http.StatusUnauthorized, result.Code)
	assert.Equal(suite.T(), "error", result.Status)
	assert.Equal(suite.T(), "Unauthorized API access. Unknown Merchant", result.Message)
	assert.Empty(suite.T(), result.Data)
	//response
	defer resp.Body.Close()
	bodyRsp, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(suite.T(), body, bodyRsp)
	//error
	assert.Equal(suite.T(), "Unauthorized API access. Unknown Merchant", err.Error())
}

func (suite *ProvidersResourceTestSuite) TestGetListNonError() {
	body, _ := LoadStubResponseData("stubs/errors/500.html")
	httpmock.RegisterResponder(http.MethodGet, suite.cfg.Uri+"/v1/payment-options/collection/card/ke", httpmock.NewBytesResponder(http.StatusOK, body))

	filter := &ProvidersFilter{Country: CountryCodeKenya, Method: TransactionMethodCard, TransactionType: TransactionTypeCollection}
	result, resp, err := suite.testable.GetList(suite.ctx, filter)
	assert.Error(suite.T(), err)
	assert.NotEmpty(suite.T(), resp)
	assert.Empty(suite.T(), result)
	//response
	defer resp.Body.Close()
	bodyRsp, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(suite.T(), body, bodyRsp)
}

func (suite *ProvidersResourceTestSuite) TestGetListInvalidFilter() {
	filter := &ProvidersFilter{}
	result, rsp, err := suite.testable.GetList(suite.ctx, filter)
	assert.Nil(suite.T(), rsp)
	assert.Nil(suite.T(), result)
	assert.Error(suite.T(), err)
}

func TestProvidersResourceTestSuite(t *testing.T) {
	suite.Run(t, new(ProvidersResourceTestSuite))
}
