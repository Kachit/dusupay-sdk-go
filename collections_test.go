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

type CollectionsTestSuite struct {
	suite.Suite
}

func (suite *CollectionsTestSuite) TestCollectionRequestIsValidSuccess() {
	request := CollectionRequest{
		Currency:          CurrencyCodeKES,
		Amount:            100,
		Method:            TransactionMethodBank,
		ProviderId:        "provider_id",
		MerchantReference: "merchant_reference",
		RedirectUrl:       "redirect_url",
		Narration:         "narration",
	}
	assert.Nil(suite.T(), request.isValid())
	assert.NoError(suite.T(), request.isValid())
}

func (suite *CollectionsTestSuite) TestCollectionRequestIsValidEmptyCurrency() {
	request := CollectionRequest{
		Amount:            100,
		Method:            TransactionMethodBank,
		ProviderId:        "provider_id",
		MerchantReference: "merchant_reference",
		RedirectUrl:       "redirect_url",
		Narration:         "narration",
	}
	result := request.isValid()
	assert.Error(suite.T(), result)
	assert.Equal(suite.T(), `parameter "currency" is empty`, result.Error())
}

func (suite *CollectionsTestSuite) TestCollectionRequestIsValidEmptyAmount() {
	request := CollectionRequest{
		Currency:          CurrencyCodeKES,
		Method:            TransactionMethodBank,
		ProviderId:        "provider_id",
		MerchantReference: "merchant_reference",
		RedirectUrl:       "redirect_url",
		Narration:         "narration",
	}
	result := request.isValid()
	assert.Error(suite.T(), result)
	assert.Equal(suite.T(), `parameter "amount" is empty`, result.Error())
}

func (suite *CollectionsTestSuite) TestCollectionRequestIsValidEmptyMethod() {
	request := CollectionRequest{
		Currency:          CurrencyCodeKES,
		Amount:            100,
		ProviderId:        "provider_id",
		MerchantReference: "merchant_reference",
		RedirectUrl:       "redirect_url",
		Narration:         "narration",
	}
	result := request.isValid()
	assert.Error(suite.T(), result)
	assert.Equal(suite.T(), `parameter "method" is empty`, result.Error())
}

func (suite *CollectionsTestSuite) TestCollectionRequestIsValidEmptyProviderId() {
	request := CollectionRequest{
		Currency:          CurrencyCodeKES,
		Amount:            100,
		Method:            TransactionMethodBank,
		MerchantReference: "merchant_reference",
		RedirectUrl:       "redirect_url",
		Narration:         "narration",
	}
	result := request.isValid()
	assert.Error(suite.T(), result)
	assert.Equal(suite.T(), `parameter "provider_id" is empty`, result.Error())
}

func (suite *CollectionsTestSuite) TestCollectionRequestIsValidEmptyMerchantReference() {
	request := CollectionRequest{
		Currency:    CurrencyCodeKES,
		Amount:      100,
		Method:      TransactionMethodBank,
		ProviderId:  "provider_id",
		RedirectUrl: "redirect_url",
		Narration:   "narration",
	}
	result := request.isValid()
	assert.Error(suite.T(), result)
	assert.Equal(suite.T(), `parameter "merchant_reference" is empty`, result.Error())
}

func (suite *CollectionsTestSuite) TestCollectionRequestIsValidEmptyNarration() {
	request := CollectionRequest{
		Currency:          CurrencyCodeKES,
		Amount:            100,
		Method:            TransactionMethodBank,
		ProviderId:        "provider_id",
		MerchantReference: "merchant_reference",
		RedirectUrl:       "redirect_url",
	}
	result := request.isValid()
	assert.Error(suite.T(), result)
	assert.Equal(suite.T(), `parameter "narration" is empty`, result.Error())
}

func (suite *CollectionsTestSuite) TestCollectionRequestIsValidEmptyRedirectUrlByDefault() {
	request := CollectionRequest{
		Currency:          CurrencyCodeKES,
		Amount:            100,
		Method:            TransactionMethodBank,
		ProviderId:        "provider_id",
		MerchantReference: "merchant_reference",
		Narration:         "narration",
	}
	result := request.isValid()
	assert.Error(suite.T(), result)
	assert.Equal(suite.T(), `parameter "redirect_url" is empty`, result.Error())
}

func (suite *CollectionsTestSuite) TestCollectionRequestIsValidEmptyRedirectUrlMobileMoney() {
	request := CollectionRequest{
		Currency:          CurrencyCodeKES,
		Amount:            100,
		Method:            TransactionMethodMobileMoney,
		ProviderId:        "provider_id",
		MerchantReference: "merchant_reference",
		Narration:         "narration",
		AccountNumber:     "account_number",
	}
	result := request.isValid()
	assert.Nil(suite.T(), result)
	assert.NoError(suite.T(), result)
}

func (suite *CollectionsTestSuite) TestCollectionRequestIsValidEmptyRedirectUrlMobileMoneyWithoutAccountNumber() {
	request := CollectionRequest{
		Currency:          CurrencyCodeKES,
		Amount:            100,
		Method:            TransactionMethodMobileMoney,
		ProviderId:        "provider_id",
		MerchantReference: "merchant_reference",
		Narration:         "narration",
	}
	result := request.isValid()
	assert.Error(suite.T(), result)
	assert.Equal(suite.T(), `parameter "account_number" is empty`, result.Error())
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
	transport := BuildStubHttpTransport()
	suite.cfg = cfg
	suite.ctx = context.Background()
	suite.testable = &CollectionsResource{NewResourceAbstract(transport, cfg)}
	httpmock.Activate()
}

func (suite *CollectionsResourceTestSuite) TearDownTest() {
	httpmock.DeactivateAndReset()
}

func (suite *CollectionsResourceTestSuite) TestCreateSuccess() {
	body, _ := LoadStubResponseData("stubs/collections/create/success.json")
	httpmock.RegisterResponder(http.MethodPost, suite.cfg.Uri+"/v1/collections", httpmock.NewBytesResponder(http.StatusOK, body))

	request := &CollectionRequest{
		Currency:          CurrencyCodeKES,
		Amount:            100,
		Method:            TransactionMethodBank,
		ProviderId:        "provider_id",
		MerchantReference: "merchant_reference",
		RedirectUrl:       "redirect_url",
		Narration:         "narration",
	}
	result, resp, err := suite.testable.Create(suite.ctx, request)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), resp)
	assert.NotEmpty(suite.T(), result)
	//result
	assert.True(suite.T(), result.IsSuccess())
	assert.Equal(suite.T(), http.StatusAccepted, result.Code)
	assert.Equal(suite.T(), "accepted", result.Status)
	assert.Equal(suite.T(), "Transaction Initiated", result.Message)
	assert.Equal(suite.T(), int64(226), result.Data.ID)
	assert.Equal(suite.T(), 0.2, result.Data.RequestAmount)
	assert.Equal(suite.T(), "USD", result.Data.RequestCurrency)
	assert.Equal(suite.T(), 737.9934, result.Data.AccountAmount)
	assert.Equal(suite.T(), "UGX", result.Data.AccountCurrency)
	assert.Equal(suite.T(), 21.4018, result.Data.TransactionFee)
	assert.Equal(suite.T(), 716.5916, result.Data.TotalCredit)
	assert.Equal(suite.T(), "mtn_ug", result.Data.ProviderID)
	assert.Equal(suite.T(), "76859aae-f148-48c5-9901-2e474cf19b71", result.Data.MerchantReference)
	assert.Equal(suite.T(), "DUSUPAY405GZM1G5JXGA71IK", result.Data.InternalReference)
	assert.Equal(suite.T(), "PENDING", result.Data.TransactionStatus)
	assert.Equal(suite.T(), "collection", result.Data.TransactionType)
	assert.Equal(suite.T(), "Transaction Initiated", result.Data.Message)
	assert.Equal(suite.T(), false, result.Data.CustomerCharged)
	assert.Equal(suite.T(), "https://sandbox.dusupay.com/v1/complete-payment/DUSUPAY405GZM1G5JXGA71IK", result.Data.PaymentURL)
	assert.Equal(suite.T(), "Ensure that you have sufficient balance on your MTN Mobile Money account", result.Data.Instructions[0].Description)
	assert.Equal(suite.T(), "1", result.Data.Instructions[0].StepNo)
	//response
	defer resp.Body.Close()
	bodyRsp, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(suite.T(), body, bodyRsp)
}

func (suite *CollectionsResourceTestSuite) TestCreateJsonError() {
	body, _ := LoadStubResponseData("stubs/errors/401.json")
	httpmock.RegisterResponder(http.MethodPost, suite.cfg.Uri+"/v1/collections", httpmock.NewBytesResponder(http.StatusOK, body))

	request := &CollectionRequest{
		Currency:          CurrencyCodeKES,
		Amount:            100,
		Method:            TransactionMethodBank,
		ProviderId:        "provider_id",
		MerchantReference: "merchant_reference",
		RedirectUrl:       "redirect_url",
		Narration:         "narration",
	}
	result, resp, err := suite.testable.Create(suite.ctx, request)
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

func (suite *CollectionsResourceTestSuite) TestCreateNonJsonError() {
	body, _ := LoadStubResponseData("stubs/errors/500.html")
	httpmock.RegisterResponder(http.MethodPost, suite.cfg.Uri+"/v1/collections", httpmock.NewBytesResponder(http.StatusOK, body))

	request := &CollectionRequest{
		Currency:          CurrencyCodeKES,
		Amount:            100,
		Method:            TransactionMethodBank,
		ProviderId:        "provider_id",
		MerchantReference: "merchant_reference",
		RedirectUrl:       "redirect_url",
		Narration:         "narration",
	}
	result, resp, err := suite.testable.Create(suite.ctx, request)
	assert.Error(suite.T(), err)
	assert.NotEmpty(suite.T(), resp)
	assert.Empty(suite.T(), result)
	//response
	defer resp.Body.Close()
	bodyRsp, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(suite.T(), body, bodyRsp)
}

func (suite *CollectionsResourceTestSuite) TestCreateInvalidRequest() {
	req := &CollectionRequest{}
	result, rsp, err := suite.testable.Create(suite.ctx, req)
	assert.Nil(suite.T(), rsp)
	assert.Nil(suite.T(), result)
	assert.Error(suite.T(), err)
}

func TestCollectionsResourceTestSuite(t *testing.T) {
	suite.Run(t, new(CollectionsResourceTestSuite))
}
