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

type MerchantsResourceTestSuite struct {
	suite.Suite
	cfg      *Config
	ctx      context.Context
	testable *MerchantsResource
}

func (suite *MerchantsResourceTestSuite) SetupTest() {
	cfg := BuildStubConfig()
	transport := BuildStubHttpTransport()
	suite.cfg = cfg
	suite.ctx = context.Background()
	suite.testable = &MerchantsResource{NewResourceAbstract(transport, cfg)}
	httpmock.Activate()
}

func (suite *MerchantsResourceTestSuite) TearDownTest() {
	httpmock.DeactivateAndReset()
}

func (suite *MerchantsResourceTestSuite) TestGetBalancesSuccess() {
	body, _ := LoadStubResponseData("stubs/merchants/balance/success.json")
	httpmock.RegisterResponder(http.MethodGet, suite.cfg.Uri+"/v1/merchants/balance", httpmock.NewBytesResponder(http.StatusOK, body))

	result, resp, err := suite.testable.GetBalances(suite.ctx)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), resp)
	assert.NotEmpty(suite.T(), result)
	//result
	assert.True(suite.T(), result.IsSuccess())
	assert.Equal(suite.T(), http.StatusOK, result.Code)
	assert.Equal(suite.T(), "success", result.Status)
	assert.Equal(suite.T(), "Request completed successfully.", result.Message)
	assert.Equal(suite.T(), "UGX", (*result.Data)[0].Currency)
	assert.Equal(suite.T(), 5475.816, (*result.Data)[0].Balance)
	assert.Equal(suite.T(), "USD", (*result.Data)[1].Currency)
	assert.Equal(suite.T(), float64(12), (*result.Data)[1].Balance)
	//response
	defer resp.Body.Close()
	bodyRsp, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(suite.T(), body, bodyRsp)
}

func (suite *MerchantsResourceTestSuite) TestGetBalancesJsonError() {
	body, _ := LoadStubResponseData("stubs/errors/401.json")
	httpmock.RegisterResponder(http.MethodGet, suite.cfg.Uri+"/v1/merchants/balance", httpmock.NewBytesResponder(http.StatusOK, body))

	result, resp, err := suite.testable.GetBalances(suite.ctx)
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

func (suite *MerchantsResourceTestSuite) TestGetBalancesNonJsonError() {
	body, _ := LoadStubResponseData("stubs/errors/500.html")
	httpmock.RegisterResponder(http.MethodGet, suite.cfg.Uri+"/v1/merchants/balance", httpmock.NewBytesResponder(http.StatusOK, body))

	result, resp, err := suite.testable.GetBalances(suite.ctx)
	assert.Error(suite.T(), err)
	assert.NotEmpty(suite.T(), resp)
	assert.Empty(suite.T(), result)
	//response
	defer resp.Body.Close()
	bodyRsp, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(suite.T(), body, bodyRsp)
}

func TestMerchantsResourceTestSuite(t *testing.T) {
	suite.Run(t, new(MerchantsResourceTestSuite))
}
