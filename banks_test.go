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

type BanksTestSuite struct {
	suite.Suite
}

func (suite *BanksTestSuite) TestBanksFilterIsValidSuccess() {
	filter := BanksFilter{Country: CountryCodeKenya, TransactionType: TransactionTypeCollection}
	assert.Nil(suite.T(), filter.isValid())
	assert.NoError(suite.T(), filter.isValid())
}

func (suite *BanksTestSuite) TestBanksFilterIsValidEmptyCountryCode() {
	filter := BanksFilter{TransactionType: "qwerty"}
	result := filter.isValid()
	assert.Error(suite.T(), result)
	assert.Equal(suite.T(), `parameter "country_code" is empty`, result.Error())
}

func (suite *BanksTestSuite) TestBanksFilterIsValidEmptyMethod() {
	filter := BanksFilter{Country: CountryCodeKenya}
	result := filter.isValid()
	assert.Error(suite.T(), result)
	assert.Equal(suite.T(), `parameter "transaction_type" is empty`, result.Error())
}

func (suite *BanksTestSuite) TestBanksFilterBuildPath() {
	filter := BanksFilter{Country: CountryCodeKenya, TransactionType: TransactionTypeCollection}
	result := filter.buildPath()
	assert.Equal(suite.T(), `collection/bank/ke`, result)
}

func (suite *BanksTestSuite) TestBanksBranchesFilterIsValidSuccess() {
	filter := BanksBranchesFilter{Country: CountryCodeKenya, Bank: "qwerty"}
	assert.Nil(suite.T(), filter.isValid())
	assert.NoError(suite.T(), filter.isValid())
}

func (suite *BanksTestSuite) TestBanksBranchesFilterIsValidEmptyCountryCode() {
	filter := BanksBranchesFilter{Country: CountryCodeKenya, Bank: "qwerty"}
	assert.Nil(suite.T(), filter.isValid())
	assert.NoError(suite.T(), filter.isValid())
}

func (suite *BanksTestSuite) TestBanksBranchesFilterIsValidEmptyBankCode() {
	filter := BanksBranchesFilter{Country: CountryCodeKenya}
	result := filter.isValid()
	assert.Error(suite.T(), result)
	assert.Equal(suite.T(), `parameter "bank_code" is empty`, result.Error())
}

func (suite *BanksTestSuite) TestBanksBranchesFilterBuildPath() {
	filter := BanksBranchesFilter{Country: CountryCodeKenya, Bank: "qwerty"}
	result := filter.buildPath()
	assert.Equal(suite.T(), `ke/branches/qwerty`, result)
}

func TestBanksTestSuite(t *testing.T) {
	suite.Run(t, new(BanksTestSuite))
}

type BanksResourceTestSuite struct {
	suite.Suite
	cfg      *Config
	ctx      context.Context
	testable *BanksResource
}

func (suite *BanksResourceTestSuite) SetupTest() {
	cfg := BuildStubConfig()
	transport := BuildStubHttpTransport()
	suite.cfg = cfg
	suite.ctx = context.Background()
	suite.testable = &BanksResource{NewResourceAbstract(transport, cfg)}
	httpmock.Activate()
}

func (suite *BanksResourceTestSuite) TearDownTest() {
	httpmock.DeactivateAndReset()
}

func (suite *BanksResourceTestSuite) TestGetListInvalidFilter() {
	filter := &BanksFilter{}
	result, rsp, err := suite.testable.GetList(suite.ctx, filter)
	assert.Nil(suite.T(), result)
	assert.Nil(suite.T(), rsp)
	assert.Error(suite.T(), err)
}

func (suite *BanksResourceTestSuite) TestGetListSuccess() {
	body, _ := LoadStubResponseData("stubs/banks/list/success.json")
	httpmock.RegisterResponder(http.MethodGet, suite.cfg.Uri+"/v1/payment-options/payout/bank/ke", httpmock.NewBytesResponder(http.StatusOK, body))

	filter := &BanksFilter{Country: CountryCodeKenya, TransactionType: TransactionTypePayout}
	result, resp, err := suite.testable.GetList(suite.ctx, filter)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), resp)
	assert.NotEmpty(suite.T(), result)
	//result
	assert.True(suite.T(), result.IsSuccess())
	assert.Equal(suite.T(), http.StatusOK, result.Code)
	assert.Equal(suite.T(), "success", result.Status)
	assert.Equal(suite.T(), "Request completed successfully.", result.Message)
	assert.Equal(suite.T(), "access_bank", (*result.Data)[0].BankCode)
	assert.Equal(suite.T(), "Access Bank", (*result.Data)[0].Name)
	assert.Equal(suite.T(), "NGN", (*result.Data)[0].TransactionCurrency)
	assert.Equal(suite.T(), float64(1000), (*result.Data)[0].MinAmount)
	assert.Equal(suite.T(), float64(380000), (*result.Data)[0].MaxAmount)
	assert.Equal(suite.T(), true, (*result.Data)[0].Available)
	assert.Empty(suite.T(), (*result.Data)[0].SandboxTestAccounts)
	//response
	defer resp.Body.Close()
	bodyRsp, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(suite.T(), body, bodyRsp)
}

func (suite *BanksResourceTestSuite) TestGetListSuccessSandbox() {
	body, _ := LoadStubResponseData("stubs/banks/list/success-sandbox.json")
	httpmock.RegisterResponder(http.MethodGet, suite.cfg.Uri+"/v1/payment-options/payout/bank/ke", httpmock.NewBytesResponder(http.StatusOK, body))

	filter := &BanksFilter{Country: CountryCodeKenya, TransactionType: TransactionTypePayout}
	result, resp, err := suite.testable.GetList(suite.ctx, filter)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), resp)
	assert.NotEmpty(suite.T(), result)
	//result
	assert.True(suite.T(), result.IsSuccess())
	assert.Equal(suite.T(), http.StatusOK, result.Code)
	assert.Equal(suite.T(), "success", result.Status)
	assert.Equal(suite.T(), "Request completed successfully.", result.Message)
	assert.Equal(suite.T(), "access_bank", (*result.Data)[0].BankCode)
	assert.Equal(suite.T(), "Access Bank", (*result.Data)[0].Name)
	assert.Equal(suite.T(), "NGN", (*result.Data)[0].TransactionCurrency)
	assert.Equal(suite.T(), float64(1000), (*result.Data)[0].MinAmount)
	assert.Equal(suite.T(), float64(380000), (*result.Data)[0].MaxAmount)
	assert.Equal(suite.T(), true, (*result.Data)[0].Available)
	assert.Equal(suite.T(), "256777000456", (*result.Data)[0].SandboxTestAccounts.Failure)
	assert.Equal(suite.T(), "256777000123", (*result.Data)[0].SandboxTestAccounts.Success)
	//response
	defer resp.Body.Close()
	bodyRsp, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(suite.T(), body, bodyRsp)
}

func (suite *BanksResourceTestSuite) TestGetListJsonError() {
	body, _ := LoadStubResponseData("stubs/errors/401.json")
	httpmock.RegisterResponder(http.MethodGet, suite.cfg.Uri+"/v1/payment-options/payout/bank/ke", httpmock.NewBytesResponder(http.StatusOK, body))

	filter := &BanksFilter{Country: CountryCodeKenya, TransactionType: TransactionTypePayout}
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

func (suite *BanksResourceTestSuite) TestGetListNonJsonError() {
	body, _ := LoadStubResponseData("stubs/errors/500.html")
	httpmock.RegisterResponder(http.MethodGet, suite.cfg.Uri+"/v1/payment-options/payout/bank/ke", httpmock.NewBytesResponder(http.StatusOK, body))

	filter := &BanksFilter{Country: CountryCodeKenya, TransactionType: TransactionTypePayout}
	result, resp, err := suite.testable.GetList(suite.ctx, filter)
	assert.Error(suite.T(), err)
	assert.NotEmpty(suite.T(), resp)
	assert.Empty(suite.T(), result)
	//response
	defer resp.Body.Close()
	bodyRsp, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(suite.T(), body, bodyRsp)
}

func (suite *BanksResourceTestSuite) TestGetBranchesListInvalidFilter() {
	filter := &BanksBranchesFilter{}
	result, rsp, err := suite.testable.GetBranchesList(suite.ctx, filter)
	assert.Nil(suite.T(), result)
	assert.Nil(suite.T(), rsp)
	assert.Error(suite.T(), err)
}

func (suite *BanksResourceTestSuite) TestGetBranchesListSuccess() {
	body, _ := LoadStubResponseData("stubs/banks/branches/success.json")
	httpmock.RegisterResponder(http.MethodGet, suite.cfg.Uri+"/v1/bank/ke/branches/qwerty", httpmock.NewBytesResponder(http.StatusOK, body))

	filter := &BanksBranchesFilter{Country: CountryCodeKenya, Bank: "qwerty"}
	result, resp, err := suite.testable.GetBranchesList(suite.ctx, filter)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), resp)
	assert.NotEmpty(suite.T(), result)
	//result
	assert.True(suite.T(), result.IsSuccess())
	assert.Equal(suite.T(), http.StatusOK, result.Code)
	assert.Equal(suite.T(), "success", result.Status)
	assert.Equal(suite.T(), "Request completed successfully.", result.Message)
	assert.Equal(suite.T(), "GH030243", (*result.Data)[0].Code)
	assert.Equal(suite.T(), "BARCLAYS BANK(GH) LTD-NKAWKAW", (*result.Data)[0].Name)
	//response
	defer resp.Body.Close()
	bodyRsp, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(suite.T(), body, bodyRsp)
}

func (suite *BanksResourceTestSuite) TestGetBranchesListJsonError() {
	body, _ := LoadStubResponseData("stubs/errors/401.json")
	httpmock.RegisterResponder(http.MethodGet, suite.cfg.Uri+"/v1/bank/ke/branches/qwerty", httpmock.NewBytesResponder(http.StatusOK, body))

	filter := &BanksBranchesFilter{Country: CountryCodeKenya, Bank: "qwerty"}
	result, resp, err := suite.testable.GetBranchesList(suite.ctx, filter)
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

func (suite *BanksResourceTestSuite) TestGetBranchesListNonJsonError() {
	body, _ := LoadStubResponseData("stubs/errors/500.html")
	httpmock.RegisterResponder(http.MethodGet, suite.cfg.Uri+"/v1/bank/ke/branches/qwerty", httpmock.NewBytesResponder(http.StatusOK, body))

	filter := &BanksBranchesFilter{Country: CountryCodeKenya, Bank: "qwerty"}
	result, resp, err := suite.testable.GetBranchesList(suite.ctx, filter)
	assert.Error(suite.T(), err)
	assert.NotEmpty(suite.T(), resp)
	assert.Empty(suite.T(), result)
	//response
	defer resp.Body.Close()
	bodyRsp, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(suite.T(), body, bodyRsp)
}

func TestBanksResourceTestSuite(t *testing.T) {
	suite.Run(t, new(BanksResourceTestSuite))
}
