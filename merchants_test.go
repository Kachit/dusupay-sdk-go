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

type MerchantsResourceTestSuite struct {
	suite.Suite
	cfg      *Config
	ctx      context.Context
	testable *MerchantsResource
}

func (suite *MerchantsResourceTestSuite) SetupTest() {
	cfg := BuildStubConfig()
	transport := BuildStubHTTPTransport()
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
	httpmock.RegisterResponder(http.MethodGet, suite.cfg.URI+"/data/wallet-balances", httpmock.NewBytesResponder(http.StatusOK, body))

	result, resp, err := suite.testable.GetBalances(suite.ctx, nil)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), resp)
	assert.NotEmpty(suite.T(), result)
	assert.True(suite.T(), result.IsSuccess())
	assert.Equal(suite.T(), http.StatusOK, result.Code)
	assert.Equal(suite.T(), 2, len(result.Data.MerchantWallets))
	assert.Equal(suite.T(), "UGX", result.Data.MerchantWallets[0].Currency)
	assert.Equal(suite.T(), 5475.816, result.Data.MerchantWallets[0].AvailableBalance)
	assert.Equal(suite.T(), "MERCHANT_COLLECTION", result.Data.MerchantWallets[0].WalletType)
	assert.Equal(suite.T(), "USD", result.Data.MerchantWallets[1].Currency)
	assert.Equal(suite.T(), float64(12), result.Data.MerchantWallets[1].AvailableBalance)

	defer func() { _ = resp.Body.Close() }()
	bodyRsp, _ := io.ReadAll(resp.Body)
	assert.Equal(suite.T(), body, bodyRsp)
}

func (suite *MerchantsResourceTestSuite) TestGetBalancesJsonError() {
	body, _ := LoadStubResponseData("stubs/errors/401.json")
	httpmock.RegisterResponder(http.MethodGet, suite.cfg.URI+"/data/wallet-balances", httpmock.NewBytesResponder(http.StatusUnauthorized, body))

	result, resp, err := suite.testable.GetBalances(suite.ctx, nil)
	defer func() { _ = resp.Body.Close() }()
	assert.Error(suite.T(), err)
	assert.NotEmpty(suite.T(), resp)
	assert.NotEmpty(suite.T(), result)
	assert.Equal(suite.T(), "Unauthorized API access. Unknown Merchant", err.Error())
}

func (suite *MerchantsResourceTestSuite) TestGetBalancesNonJsonError() {
	body, _ := LoadStubResponseData("stubs/errors/500.html")
	httpmock.RegisterResponder(http.MethodGet, suite.cfg.URI+"/data/wallet-balances", httpmock.NewBytesResponder(http.StatusInternalServerError, body))

	result, resp, err := suite.testable.GetBalances(suite.ctx, nil)
	defer func() { _ = resp.Body.Close() }()
	assert.Error(suite.T(), err)
	assert.NotEmpty(suite.T(), resp)
	assert.Empty(suite.T(), result)
}

func (suite *MerchantsResourceTestSuite) TestVerifySuccess() {
	body, _ := LoadStubResponseData("stubs/merchants/verify/success.json")
	httpmock.RegisterResponder(http.MethodGet, suite.cfg.URI+"/data/transaction/verify/MCTREFT2WMNWZ23SBN6Y", httpmock.NewBytesResponder(http.StatusOK, body))

	result, resp, err := suite.testable.Verify(suite.ctx, "MCTREFT2WMNWZ23SBN6Y")
	defer func() { _ = resp.Body.Close() }()
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), resp)
	assert.NotEmpty(suite.T(), result)
	assert.True(suite.T(), result.IsSuccess())
	assert.Equal(suite.T(), int64(226), result.Data.ID)
	assert.Equal(suite.T(), "MCTREFT2WMNWZ23SBN6Y", result.Data.MerchantReference)
	assert.Equal(suite.T(), "DUSUPAYRMGRXNNYBWATKJ", result.Data.InternalReference)
	assert.Equal(suite.T(), "COMPLETED", result.Data.TransactionStatus)
	assert.Equal(suite.T(), float64(10000), result.Data.TransactionAmount)
	assert.Equal(suite.T(), "mtn_ug", result.Data.ProviderCode)
	assert.Equal(suite.T(), "COLLECTION", result.Data.TransactionType)
	assert.Equal(suite.T(), "UGX", result.Data.RequestCurrency)
	assert.Equal(suite.T(), float64(10000), result.Data.RequestAmount)
	assert.Equal(suite.T(), "UGX", result.Data.TransactionCurrency)
	assert.Equal(suite.T(), "256787008803", result.Data.TransactionAccount)
	assert.Equal(suite.T(), float64(9700), result.Data.TotalCredit)
}

func (suite *MerchantsResourceTestSuite) TestGetBalancesWithCurrencyFilter() {
	body, _ := LoadStubResponseData("stubs/merchants/balance/success.json")
	httpmock.RegisterResponderWithQuery(http.MethodGet, suite.cfg.URI+"/data/wallet-balances", "currency=UGX", httpmock.NewBytesResponder(http.StatusOK, body))

	result, resp, err := suite.testable.GetBalances(suite.ctx, &BalancesFilter{Currency: CurrencyCodeUGX})
	defer func() { _ = resp.Body.Close() }()
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), result.IsSuccess())
}

func (suite *MerchantsResourceTestSuite) TestVerifyEmptyReference() {
	result, rsp, err := suite.testable.Verify(suite.ctx, "") //nolint:bodyclose // returns nil response on validation error
	assert.Nil(suite.T(), rsp)
	assert.Nil(suite.T(), result)
	assert.EqualError(suite.T(), err, `MerchantsResource.Verify error: parameter "merchant_reference" is empty`)
}

func TestMerchantsResourceTestSuite(t *testing.T) {
	suite.Run(t, new(MerchantsResourceTestSuite))
}
