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

type BanksTestSuite struct {
	suite.Suite
}

func (suite *BanksTestSuite) TestBanksFilterIsValidSuccess() {
	filter := BanksFilter{ProviderCode: "bank_ng"}
	assert.NoError(suite.T(), filter.isValid())
}

func (suite *BanksTestSuite) TestBanksFilterIsValidEmptyProviderCode() {
	filter := BanksFilter{}
	assert.EqualError(suite.T(), filter.isValid(), `parameter "provider_code" is empty`)
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
	transport := BuildStubHTTPTransport()
	suite.cfg = cfg
	suite.ctx = context.Background()
	suite.testable = &BanksResource{NewResourceAbstract(transport, cfg)}
	httpmock.Activate()
}

func (suite *BanksResourceTestSuite) TearDownTest() {
	httpmock.DeactivateAndReset()
}

func (suite *BanksResourceTestSuite) TestGetListSuccess() {
	body, _ := LoadStubResponseData("stubs/banks/list/success.json")
	httpmock.RegisterResponder(http.MethodGet, suite.cfg.URI+"/data/payout-bank-codes", httpmock.NewBytesResponder(http.StatusOK, body))

	filter := &BanksFilter{ProviderCode: "bank_ng"}
	result, resp, err := suite.testable.GetList(suite.ctx, filter)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), resp)
	assert.NotEmpty(suite.T(), result)
	assert.True(suite.T(), result.IsSuccess())
	assert.Equal(suite.T(), 2, len(result.Data.PayoutBanks))
	assert.Equal(suite.T(), "Access Bank Nigeria", result.Data.PayoutBanks[0].BankName)
	assert.Equal(suite.T(), "access_bank_ng", result.Data.PayoutBanks[0].BankCode)
	assert.Equal(suite.T(), true, result.Data.PayoutBanks[0].IsActive)

	defer func() { _ = resp.Body.Close() }()
	bodyRsp, _ := io.ReadAll(resp.Body)
	assert.Equal(suite.T(), body, bodyRsp)
}

func (suite *BanksResourceTestSuite) TestGetListJsonError() {
	body, _ := LoadStubResponseData("stubs/errors/401.json")
	httpmock.RegisterResponder(http.MethodGet, suite.cfg.URI+"/data/payout-bank-codes", httpmock.NewBytesResponder(http.StatusUnauthorized, body))

	filter := &BanksFilter{ProviderCode: "bank_ng"}
	result, resp, err := suite.testable.GetList(suite.ctx, filter)
	defer func() { _ = resp.Body.Close() }()
	assert.Error(suite.T(), err)
	assert.NotEmpty(suite.T(), resp)
	assert.NotEmpty(suite.T(), result)
	assert.Equal(suite.T(), "Unauthorized API access. Unknown Merchant", err.Error())
}

func (suite *BanksResourceTestSuite) TestGetListNonJsonError() {
	body, _ := LoadStubResponseData("stubs/errors/500.html")
	httpmock.RegisterResponder(http.MethodGet, suite.cfg.URI+"/data/payout-bank-codes", httpmock.NewBytesResponder(http.StatusInternalServerError, body))

	filter := &BanksFilter{ProviderCode: "bank_ng"}
	result, resp, err := suite.testable.GetList(suite.ctx, filter)
	defer func() { _ = resp.Body.Close() }()
	assert.Error(suite.T(), err)
	assert.NotEmpty(suite.T(), resp)
	assert.Empty(suite.T(), result)
}

func (suite *BanksResourceTestSuite) TestGetListInvalidFilter() {
	result, rsp, err := suite.testable.GetList(suite.ctx, &BanksFilter{}) //nolint:bodyclose // returns nil response on validation error
	assert.Nil(suite.T(), rsp)
	assert.Nil(suite.T(), result)
	assert.Error(suite.T(), err)
}

func TestBanksResourceTestSuite(t *testing.T) {
	suite.Run(t, new(BanksResourceTestSuite))
}
