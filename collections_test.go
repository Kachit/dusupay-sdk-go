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

type CollectionsTestSuite struct {
	suite.Suite
}

func (suite *CollectionsTestSuite) TestCollectionRequestIsValidSuccess() {
	request := CollectionRequest{
		Currency:          CurrencyCodeKES,
		Amount:            100,
		TransactionMethod: TransactionMethodBank,
		ProviderCode:      "provider_code",
		MerchantReference: "merchant_reference",
		Description:       "test description",
	}
	assert.NoError(suite.T(), request.isValid())
}

func (suite *CollectionsTestSuite) TestCollectionRequestIsValidEmptyCurrency() {
	request := CollectionRequest{
		Amount:            100,
		TransactionMethod: TransactionMethodBank,
		ProviderCode:      "provider_code",
		MerchantReference: "merchant_reference",
		Description:       "test description",
	}
	assert.EqualError(suite.T(), request.isValid(), `parameter "currency" is empty`)
}

func (suite *CollectionsTestSuite) TestCollectionRequestIsValidEmptyAmount() {
	request := CollectionRequest{
		Currency:          CurrencyCodeKES,
		TransactionMethod: TransactionMethodBank,
		ProviderCode:      "provider_code",
		MerchantReference: "merchant_reference",
		Description:       "test description",
	}
	assert.EqualError(suite.T(), request.isValid(), `parameter "amount" is empty`)
}

func (suite *CollectionsTestSuite) TestCollectionRequestIsValidEmptyTransactionMethod() {
	request := CollectionRequest{
		Currency:          CurrencyCodeKES,
		Amount:            100,
		ProviderCode:      "provider_code",
		MerchantReference: "merchant_reference",
		Description:       "test description",
	}
	assert.EqualError(suite.T(), request.isValid(), `parameter "transaction_method" is empty`)
}

func (suite *CollectionsTestSuite) TestCollectionRequestIsValidEmptyProviderCode() {
	request := CollectionRequest{
		Currency:          CurrencyCodeKES,
		Amount:            100,
		TransactionMethod: TransactionMethodBank,
		MerchantReference: "merchant_reference",
		Description:       "test description",
	}
	assert.EqualError(suite.T(), request.isValid(), `parameter "provider_code" is empty`)
}

func (suite *CollectionsTestSuite) TestCollectionRequestIsValidEmptyMerchantReference() {
	request := CollectionRequest{
		Currency:          CurrencyCodeKES,
		Amount:            100,
		TransactionMethod: TransactionMethodBank,
		ProviderCode:      "provider_code",
		Description:       "test description",
	}
	assert.EqualError(suite.T(), request.isValid(), `parameter "merchant_reference" is empty`)
}

func (suite *CollectionsTestSuite) TestCollectionRequestIsValidEmptyDescription() {
	request := CollectionRequest{
		Currency:          CurrencyCodeKES,
		Amount:            100,
		TransactionMethod: TransactionMethodBank,
		ProviderCode:      "provider_code",
		MerchantReference: "merchant_reference",
	}
	assert.EqualError(suite.T(), request.isValid(), `parameter "description" is empty`)
}

func (suite *CollectionsTestSuite) TestCollectionRequestIsValidEmptyMsisdnMobileMoney() {
	request := CollectionRequest{
		Currency:          CurrencyCodeKES,
		Amount:            100,
		TransactionMethod: TransactionMethodMobileMoney,
		ProviderCode:      "provider_code",
		MerchantReference: "merchant_reference",
		Description:       "test description",
	}
	assert.EqualError(suite.T(), request.isValid(), `parameter "msisdn" is empty`)
}

func (suite *CollectionsTestSuite) TestCollectionRequestIsValidMobileMoneyWithMsisdn() {
	request := CollectionRequest{
		Currency:          CurrencyCodeKES,
		Amount:            100,
		TransactionMethod: TransactionMethodMobileMoney,
		ProviderCode:      "provider_code",
		MerchantReference: "merchant_reference",
		Description:       "test description",
		Msisdn:            "256777000001",
	}
	assert.NoError(suite.T(), request.isValid())
}

func TestCollectionsTestSuite(t *testing.T) {
	suite.Run(t, new(CollectionsTestSuite))
}

type CollectionsResourceTestSuite struct {
	suite.Suite
	cfg      *Config
	ctx      context.Context
	testable *CollectionsResource
}

func (suite *CollectionsResourceTestSuite) SetupTest() {
	cfg := BuildStubConfig()
	transport := BuildStubHTTPTransport()
	suite.cfg = cfg
	suite.ctx = context.Background()
	suite.testable = &CollectionsResource{NewResourceAbstract(transport, cfg)}
	httpmock.Activate()
}

func (suite *CollectionsResourceTestSuite) TearDownTest() {
	httpmock.DeactivateAndReset()
}

func (suite *CollectionsResourceTestSuite) TestCreateSuccess() {
	body, _ := LoadStubResponseData("stubs/collections/initialize/success.json")
	httpmock.RegisterResponder(http.MethodPost, suite.cfg.URI+"/collections/initialize", httpmock.NewBytesResponder(http.StatusAccepted, body))

	request := &CollectionRequest{
		Currency:          CurrencyCodeKES,
		Amount:            100,
		TransactionMethod: TransactionMethodBank,
		ProviderCode:      "bank_ke",
		MerchantReference: "MCTREF9WVCPGFRHYCWSK",
		Description:       "Test collection",
	}
	result, resp, err := suite.testable.Create(suite.ctx, request)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), resp)
	assert.NotEmpty(suite.T(), result)
	assert.True(suite.T(), result.IsSuccess())
	assert.Equal(suite.T(), http.StatusAccepted, result.Code)
	assert.Equal(suite.T(), "accepted", result.Status)
	assert.Equal(suite.T(), "Request Accepted", result.Message)
	assert.Equal(suite.T(), "DUSUPAYSTPDKZ6MLY5WPQ", result.Data.InternalReference)
	assert.Equal(suite.T(), "MCTREF9WVCPGFRHYCWSK", result.Data.MerchantReference)
	assert.Equal(suite.T(), "https://sandboxapi.dusupay.com/checkout/DUSUPAYSTPDKZ6MLY5WPQ", result.Data.PaymentURL)

	defer func() { _ = resp.Body.Close() }()
	bodyRsp, _ := io.ReadAll(resp.Body)
	assert.Equal(suite.T(), body, bodyRsp)
}

func (suite *CollectionsResourceTestSuite) TestCreateJsonError() {
	body, _ := LoadStubResponseData("stubs/errors/401.json")
	httpmock.RegisterResponder(http.MethodPost, suite.cfg.URI+"/collections/initialize", httpmock.NewBytesResponder(http.StatusUnauthorized, body))

	request := &CollectionRequest{
		Currency:          CurrencyCodeKES,
		Amount:            100,
		TransactionMethod: TransactionMethodBank,
		ProviderCode:      "bank_ke",
		MerchantReference: "MCTREF9WVCPGFRHYCWSK",
		Description:       "Test collection",
	}
	result, resp, err := suite.testable.Create(suite.ctx, request)
	assert.Error(suite.T(), err)
	assert.NotEmpty(suite.T(), resp)
	assert.NotEmpty(suite.T(), result)
	assert.False(suite.T(), result.IsSuccess())
	assert.Equal(suite.T(), "Unauthorized API access. Unknown Merchant", err.Error())
}

func (suite *CollectionsResourceTestSuite) TestCreateNonJsonError() {
	body, _ := LoadStubResponseData("stubs/errors/500.html")
	httpmock.RegisterResponder(http.MethodPost, suite.cfg.URI+"/collections/initialize", httpmock.NewBytesResponder(http.StatusInternalServerError, body))

	request := &CollectionRequest{
		Currency:          CurrencyCodeKES,
		Amount:            100,
		TransactionMethod: TransactionMethodBank,
		ProviderCode:      "bank_ke",
		MerchantReference: "MCTREF9WVCPGFRHYCWSK",
		Description:       "Test collection",
	}
	result, resp, err := suite.testable.Create(suite.ctx, request)
	assert.Error(suite.T(), err)
	assert.NotEmpty(suite.T(), resp)
	assert.Empty(suite.T(), result)
}

func (suite *CollectionsResourceTestSuite) TestCreateInvalidRequest() {
	result, rsp, err := suite.testable.Create(suite.ctx, &CollectionRequest{})
	assert.Nil(suite.T(), rsp)
	assert.Nil(suite.T(), result)
	assert.Error(suite.T(), err)
}

func (suite *CollectionsResourceTestSuite) TestAuthorizeSuccess() {
	body, _ := LoadStubResponseData("stubs/collections/authorize/success.json")
	httpmock.RegisterResponder(http.MethodPost, suite.cfg.URI+"/collections/authorize", httpmock.NewBytesResponder(http.StatusAccepted, body))

	request := &CollectionAuthorizeRequest{
		InternalReference: "DUSUPAYSTPDKZ6MLY5WPQ",
		OTP:               "123456",
	}
	result, resp, err := suite.testable.Authorize(suite.ctx, request)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), resp)
	assert.NotEmpty(suite.T(), result)
	assert.True(suite.T(), result.IsSuccess())
	assert.Equal(suite.T(), "DUSUPAYSTPDKZ6MLY5WPQ", result.Data.InternalReference)
	assert.Equal(suite.T(), "MCTREF9WVCPGFRHYCWSK", result.Data.MerchantReference)
}

func (suite *CollectionsResourceTestSuite) TestAuthorizeInvalidRequest() {
	result, rsp, err := suite.testable.Authorize(suite.ctx, &CollectionAuthorizeRequest{})
	assert.Nil(suite.T(), rsp)
	assert.Nil(suite.T(), result)
	assert.Error(suite.T(), err)
}

func TestCollectionsResourceTestSuite(t *testing.T) {
	suite.Run(t, new(CollectionsResourceTestSuite))
}
