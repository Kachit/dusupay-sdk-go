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

type PayoutsTestSuite struct {
	suite.Suite
}

func (suite *PayoutsTestSuite) TestPayoutsRequestIsValidSuccess() {
	request := PayoutRequest{
		Currency:          CurrencyCodeKES,
		Amount:            100,
		Method:            TransactionMethodBank,
		ProviderId:        "provider_id",
		MerchantReference: "merchant_reference",
		Narration:         "narration",
		AccountNumber:     "account_number",
		AccountName:       "account_name",
	}
	assert.Nil(suite.T(), request.isValid())
	assert.NoError(suite.T(), request.isValid())
}

func (suite *PayoutsTestSuite) TestPayoutsRequestIsValidEmptyCurrency() {
	request := PayoutRequest{
		Amount:            100,
		Method:            TransactionMethodBank,
		ProviderId:        "provider_id",
		MerchantReference: "merchant_reference",
		Narration:         "narration",
	}
	result := request.isValid()
	assert.Error(suite.T(), result)
	assert.Equal(suite.T(), `parameter "currency" is empty`, result.Error())
}

func (suite *PayoutsTestSuite) TestPayoutsRequestIsValidEmptyAmount() {
	request := PayoutRequest{
		Currency:          CurrencyCodeKES,
		Method:            TransactionMethodBank,
		ProviderId:        "provider_id",
		MerchantReference: "merchant_reference",
		Narration:         "narration",
	}
	result := request.isValid()
	assert.Error(suite.T(), result)
	assert.Equal(suite.T(), `parameter "amount" is empty`, result.Error())
}

func (suite *PayoutsTestSuite) TestPayoutsRequestIsValidEmptyMethod() {
	request := PayoutRequest{
		Currency:          CurrencyCodeKES,
		Amount:            100,
		ProviderId:        "provider_id",
		MerchantReference: "merchant_reference",
		Narration:         "narration",
	}
	result := request.isValid()
	assert.Error(suite.T(), result)
	assert.Equal(suite.T(), `parameter "method" is empty`, result.Error())
}

func (suite *PayoutsTestSuite) TestPayoutsRequestIsValidEmptyProviderId() {
	request := PayoutRequest{
		Currency:          CurrencyCodeKES,
		Amount:            100,
		Method:            TransactionMethodBank,
		MerchantReference: "merchant_reference",
		Narration:         "narration",
	}
	result := request.isValid()
	assert.Error(suite.T(), result)
	assert.Equal(suite.T(), `parameter "provider_id" is empty`, result.Error())
}

func (suite *PayoutsTestSuite) TestPayoutsRequestIsValidEmptyMerchantReference() {
	request := PayoutRequest{
		Currency:   CurrencyCodeKES,
		Amount:     100,
		Method:     TransactionMethodBank,
		ProviderId: "provider_id",
		Narration:  "narration",
	}
	result := request.isValid()
	assert.Error(suite.T(), result)
	assert.Equal(suite.T(), `parameter "merchant_reference" is empty`, result.Error())
}

func (suite *PayoutsTestSuite) TestPayoutsRequestIsValidEmptyNarration() {
	request := PayoutRequest{
		Currency:          CurrencyCodeKES,
		Amount:            100,
		Method:            TransactionMethodBank,
		ProviderId:        "provider_id",
		MerchantReference: "merchant_reference",
	}
	result := request.isValid()
	assert.Error(suite.T(), result)
	assert.Equal(suite.T(), `parameter "narration" is empty`, result.Error())
}

func (suite *PayoutsTestSuite) TestPayoutsRequestIsValidEmptyAccountNumber() {
	request := PayoutRequest{
		Currency:          CurrencyCodeKES,
		Amount:            100,
		Method:            TransactionMethodBank,
		ProviderId:        "provider_id",
		MerchantReference: "merchant_reference",
		Narration:         "narration",
	}
	result := request.isValid()
	assert.Error(suite.T(), result)
	assert.Equal(suite.T(), `parameter "account_number" is empty`, result.Error())
}

func (suite *PayoutsTestSuite) TestPayoutsRequestIsValidEmptyAccountName() {
	request := PayoutRequest{
		Currency:          CurrencyCodeKES,
		Amount:            100,
		Method:            TransactionMethodBank,
		ProviderId:        "provider_id",
		MerchantReference: "merchant_reference",
		Narration:         "narration",
		AccountNumber:     "account_number",
	}
	result := request.isValid()
	assert.Error(suite.T(), result)
	assert.Equal(suite.T(), `parameter "account_name" is empty`, result.Error())
}

func TestPayoutsTestSuite(t *testing.T) {
	suite.Run(t, new(PayoutsTestSuite))
}

type PayoutsResourceTestSuite struct {
	suite.Suite
	cfg      *Config
	ctx      context.Context
	testable *PayoutsResource
}

func (suite *PayoutsResourceTestSuite) SetupTest() {
	cfg := BuildStubConfig()
	transport := BuildStubHttpTransport()
	suite.cfg = cfg
	suite.ctx = context.Background()
	suite.testable = &PayoutsResource{NewResourceAbstract(transport, cfg)}
	httpmock.Activate()
}

func (suite *PayoutsResourceTestSuite) TearDownTest() {
	httpmock.DeactivateAndReset()
}

func (suite *PayoutsResourceTestSuite) TestCreateSuccess() {
	body, _ := LoadStubResponseData("stubs/payouts/create/success.json")
	httpmock.RegisterResponder(http.MethodPost, suite.cfg.Uri+"/v1/payouts", httpmock.NewBytesResponder(http.StatusOK, body))

	request := &PayoutRequest{
		Currency:          CurrencyCodeKES,
		Amount:            100,
		Method:            TransactionMethodBank,
		ProviderId:        "provider_id",
		MerchantReference: "merchant_reference",
		Narration:         "narration",
		AccountNumber:     "account_number",
		AccountName:       "account_name",
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
	assert.Equal(suite.T(), int64(124468), result.Data.ID)
	assert.Equal(suite.T(), float64(700), result.Data.RequestAmount)
	assert.Equal(suite.T(), "UGX", result.Data.RequestCurrency)
	assert.Equal(suite.T(), float64(700), result.Data.AccountAmount)
	assert.Equal(suite.T(), "UGX", result.Data.AccountCurrency)
	assert.Equal(suite.T(), float64(1500), result.Data.TransactionFee)
	assert.Equal(suite.T(), float64(2200), result.Data.TotalDebit)
	assert.Equal(suite.T(), "mtn_ug", result.Data.ProviderID)
	assert.Equal(suite.T(), "payout-1005", result.Data.MerchantReference)
	assert.Equal(suite.T(), "DUSUPAY405GZMDVTKASJL8UQ", result.Data.InternalReference)
	assert.Equal(suite.T(), "PENDING", result.Data.TransactionStatus)
	assert.Equal(suite.T(), "payout", result.Data.TransactionType)
	assert.Equal(suite.T(), "Transaction Initiated", result.Data.Message)
	//response
	defer resp.Body.Close()
	bodyRsp, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(suite.T(), body, bodyRsp)
}

func (suite *PayoutsResourceTestSuite) TestCreateJsonError() {
	body, _ := LoadStubResponseData("stubs/errors/401.json")
	httpmock.RegisterResponder(http.MethodPost, suite.cfg.Uri+"/v1/payouts", httpmock.NewBytesResponder(http.StatusOK, body))

	request := &PayoutRequest{
		Currency:          CurrencyCodeKES,
		Amount:            100,
		Method:            TransactionMethodBank,
		ProviderId:        "provider_id",
		MerchantReference: "merchant_reference",
		Narration:         "narration",
		AccountNumber:     "account_number",
		AccountName:       "account_name",
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

func (suite *PayoutsResourceTestSuite) TestCreateNonJsonError() {
	body, _ := LoadStubResponseData("stubs/errors/500.html")
	httpmock.RegisterResponder(http.MethodPost, suite.cfg.Uri+"/v1/payouts", httpmock.NewBytesResponder(http.StatusOK, body))

	request := &PayoutRequest{
		Currency:          CurrencyCodeKES,
		Amount:            100,
		Method:            TransactionMethodBank,
		ProviderId:        "provider_id",
		MerchantReference: "merchant_reference",
		Narration:         "narration",
		AccountNumber:     "account_number",
		AccountName:       "account_name",
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

func (suite *PayoutsResourceTestSuite) TestCreateInvalidRequest() {
	req := &PayoutRequest{}
	result, rsp, err := suite.testable.Create(suite.ctx, req)
	assert.Nil(suite.T(), rsp)
	assert.Nil(suite.T(), result)
	assert.Error(suite.T(), err)
}

func TestPayoutsResourceTestSuite(t *testing.T) {
	suite.Run(t, new(PayoutsResourceTestSuite))
}
