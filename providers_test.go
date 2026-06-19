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

type ProvidersTestSuite struct {
	suite.Suite
}

func (suite *ProvidersTestSuite) TestProvidersFilterIsValidSuccess() {
	filter := ProvidersFilter{Currency: CurrencyCodeUGX, TransactionType: TransactionTypeCollection}
	assert.NoError(suite.T(), filter.isValid())
}

func (suite *ProvidersTestSuite) TestProvidersFilterIsValidEmptyCurrency() {
	filter := ProvidersFilter{TransactionType: TransactionTypeCollection}
	assert.EqualError(suite.T(), filter.isValid(), `parameter "currency" is empty`)
}

func (suite *ProvidersTestSuite) TestProvidersFilterIsValidEmptyTransactionType() {
	filter := ProvidersFilter{Currency: CurrencyCodeUGX}
	assert.EqualError(suite.T(), filter.isValid(), `parameter "transaction_type" is empty`)
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
	transport := BuildStubHTTPTransport()
	suite.cfg = cfg
	suite.ctx = context.Background()
	suite.testable = &ProvidersResource{NewResourceAbstract(transport, cfg)}
	httpmock.Activate()
}

func (suite *ProvidersResourceTestSuite) TearDownTest() {
	httpmock.DeactivateAndReset()
}

func (suite *ProvidersResourceTestSuite) TestGetListSuccess() {
	body, _ := LoadStubResponseData("stubs/providers/list/success.json")
	httpmock.RegisterResponder(http.MethodGet, suite.cfg.URI+"/data/payment-providers", httpmock.NewBytesResponder(http.StatusOK, body))

	filter := &ProvidersFilter{Currency: CurrencyCodeUGX, TransactionType: TransactionTypeCollection}
	result, resp, err := suite.testable.GetList(suite.ctx, filter)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), resp)
	assert.NotEmpty(suite.T(), result)
	assert.True(suite.T(), result.IsSuccess())
	assert.Equal(suite.T(), 1, len(result.Data.PaymentProviders))
	assert.Equal(suite.T(), "mtn_ug", result.Data.PaymentProviders[0].ProviderCode)
	assert.Equal(suite.T(), "MTN Mobile Money Uganda", result.Data.PaymentProviders[0].ProviderName)
	assert.Equal(suite.T(), "MOBILE_MONEY", result.Data.PaymentProviders[0].TransactionMethod)
	assert.Equal(suite.T(), "UGX", result.Data.PaymentProviders[0].TransactionCurrency)
	assert.Equal(suite.T(), float64(500), result.Data.PaymentProviders[0].MinAmount)
	assert.Equal(suite.T(), true, result.Data.PaymentProviders[0].IsActive)

	defer func() { _ = resp.Body.Close() }()
	bodyRsp, _ := io.ReadAll(resp.Body)
	assert.Equal(suite.T(), body, bodyRsp)
}

func (suite *ProvidersResourceTestSuite) TestGetListJsonError() {
	body, _ := LoadStubResponseData("stubs/errors/401.json")
	httpmock.RegisterResponder(http.MethodGet, suite.cfg.URI+"/data/payment-providers", httpmock.NewBytesResponder(http.StatusUnauthorized, body))

	filter := &ProvidersFilter{Currency: CurrencyCodeUGX, TransactionType: TransactionTypeCollection}
	result, resp, err := suite.testable.GetList(suite.ctx, filter)
	defer func() { _ = resp.Body.Close() }()
	assert.Error(suite.T(), err)
	assert.NotEmpty(suite.T(), resp)
	assert.NotEmpty(suite.T(), result)
	assert.Equal(suite.T(), "Unauthorized API access. Unknown Merchant", err.Error())
}

func (suite *ProvidersResourceTestSuite) TestGetListNonJsonError() {
	body, _ := LoadStubResponseData("stubs/errors/500.html")
	httpmock.RegisterResponder(http.MethodGet, suite.cfg.URI+"/data/payment-providers", httpmock.NewBytesResponder(http.StatusInternalServerError, body))

	filter := &ProvidersFilter{Currency: CurrencyCodeUGX, TransactionType: TransactionTypeCollection}
	result, resp, err := suite.testable.GetList(suite.ctx, filter)
	defer func() { _ = resp.Body.Close() }()
	assert.Error(suite.T(), err)
	assert.NotEmpty(suite.T(), resp)
	assert.Empty(suite.T(), result)
}

func (suite *ProvidersResourceTestSuite) TestGetListInvalidFilter() {
	result, rsp, err := suite.testable.GetList(suite.ctx, &ProvidersFilter{}) //nolint:bodyclose // returns nil response on validation error
	assert.Nil(suite.T(), rsp)
	assert.Nil(suite.T(), result)
	assert.Error(suite.T(), err)
}

func TestProvidersResourceTestSuite(t *testing.T) {
	suite.Run(t, new(ProvidersResourceTestSuite))
}
