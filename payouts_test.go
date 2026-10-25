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

type PayoutsTestSuite struct {
	suite.Suite
}

func (suite *PayoutsTestSuite) TestPayoutRequestIsValidSuccess() {
	request := PayoutRequest{
		Currency:          CurrencyCodeUGX,
		Amount:            10000,
		TransactionMethod: TransactionMethodMobileMoney,
		ProviderCode:      "mtn_ug",
		AccountNumber:     "256777000001",
		CustomerName:      "John Doe",
		MerchantReference: "MCTREFT2WMNWZ23SBN6Y",
		Description:       "Test payout",
	}
	assert.NoError(suite.T(), request.isValid())
}

func (suite *PayoutsTestSuite) TestPayoutRequestIsValidEmptyCurrency() {
	request := PayoutRequest{Amount: 100, TransactionMethod: TransactionMethodMobileMoney, ProviderCode: "mtn_ug", AccountNumber: "256777000001", CustomerName: "John Doe", MerchantReference: "ref", Description: "desc"}
	assert.EqualError(suite.T(), request.isValid(), `parameter "currency" is empty`)
}

func (suite *PayoutsTestSuite) TestPayoutRequestIsValidEmptyAmount() {
	request := PayoutRequest{Currency: CurrencyCodeUGX, TransactionMethod: TransactionMethodMobileMoney, ProviderCode: "mtn_ug", AccountNumber: "256777000001", CustomerName: "John Doe", MerchantReference: "ref", Description: "desc"}
	assert.EqualError(suite.T(), request.isValid(), `parameter "amount" is empty`)
}

func (suite *PayoutsTestSuite) TestPayoutRequestIsValidEmptyTransactionMethod() {
	request := PayoutRequest{Currency: CurrencyCodeUGX, Amount: 100, ProviderCode: "mtn_ug", AccountNumber: "256777000001", CustomerName: "John Doe", MerchantReference: "ref", Description: "desc"}
	assert.EqualError(suite.T(), request.isValid(), `parameter "transaction_method" is empty`)
}

func (suite *PayoutsTestSuite) TestPayoutRequestIsValidEmptyProviderCode() {
	request := PayoutRequest{Currency: CurrencyCodeUGX, Amount: 100, TransactionMethod: TransactionMethodMobileMoney, AccountNumber: "256777000001", CustomerName: "John Doe", MerchantReference: "ref", Description: "desc"}
	assert.EqualError(suite.T(), request.isValid(), `parameter "provider_code" is empty`)
}

func (suite *PayoutsTestSuite) TestPayoutRequestIsValidEmptyMerchantReference() {
	request := PayoutRequest{Currency: CurrencyCodeUGX, Amount: 100, TransactionMethod: TransactionMethodMobileMoney, ProviderCode: "mtn_ug", AccountNumber: "256777000001", CustomerName: "John Doe", Description: "desc"}
	assert.EqualError(suite.T(), request.isValid(), `parameter "merchant_reference" is empty`)
}

func (suite *PayoutsTestSuite) TestPayoutRequestIsValidEmptyDescription() {
	request := PayoutRequest{Currency: CurrencyCodeUGX, Amount: 100, TransactionMethod: TransactionMethodMobileMoney, ProviderCode: "mtn_ug", AccountNumber: "256777000001", CustomerName: "John Doe", MerchantReference: "ref"}
	assert.EqualError(suite.T(), request.isValid(), `parameter "description" is empty`)
}

func (suite *PayoutsTestSuite) TestPayoutRequestIsValidEmptyAccountNumber() {
	request := PayoutRequest{Currency: CurrencyCodeUGX, Amount: 100, TransactionMethod: TransactionMethodMobileMoney, ProviderCode: "mtn_ug", CustomerName: "John Doe", MerchantReference: "ref", Description: "desc"}
	assert.EqualError(suite.T(), request.isValid(), `parameter "account_number" is empty`)
}

func (suite *PayoutsTestSuite) TestPayoutRequestIsValidEmptyCustomerName() {
	request := PayoutRequest{Currency: CurrencyCodeUGX, Amount: 100, TransactionMethod: TransactionMethodMobileMoney, ProviderCode: "mtn_ug", AccountNumber: "256777000001", MerchantReference: "ref", Description: "desc"}
	assert.EqualError(suite.T(), request.isValid(), `parameter "customer_name" is empty`)
}

func (suite *PayoutsTestSuite) TestPayoutRequestIsValidBankEmptyBankCode() {
	request := PayoutRequest{
		Currency:          CurrencyCodeNGN,
		Amount:            4000,
		TransactionMethod: TransactionMethodBank,
		ProviderCode:      "bank_ng",
		AccountNumber:     "2121562123",
		CustomerName:      "John Doe",
		MerchantReference: "MCTREFT2WMNWZ23SBN6Y",
		Description:       "Test payout",
	}
	assert.EqualError(suite.T(), request.isValid(), `parameter "extra_params.bank_code" is empty`)

	request.ExtraParams = &PayoutExtraParams{BankCode: "access_bank_ng", BranchCode: "046"}
	assert.NoError(suite.T(), request.isValid())

	result, _ := transformStructToMap(&request)
	extra := result["extra_params"].(map[string]interface{})
	assert.Equal(suite.T(), "access_bank_ng", extra["bank_code"])
	assert.Equal(suite.T(), "046", extra["branch_code"])
	assert.NotContains(suite.T(), extra, "swift_code")
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
	transport := BuildStubHTTPTransport()
	suite.cfg = cfg
	suite.ctx = context.Background()
	suite.testable = &PayoutsResource{NewResourceAbstract(transport, cfg)}
	httpmock.Activate()
}

func (suite *PayoutsResourceTestSuite) TearDownTest() {
	httpmock.DeactivateAndReset()
}

func (suite *PayoutsResourceTestSuite) TestCreateSuccess() {
	body, _ := LoadStubResponseData("stubs/payouts/send-funds/success.json")
	httpmock.RegisterResponder(http.MethodPost, suite.cfg.URI+"/payout/send-funds", httpmock.NewBytesResponder(http.StatusAccepted, body))

	request := &PayoutRequest{
		Currency:          CurrencyCodeUGX,
		Amount:            10000,
		TransactionMethod: TransactionMethodMobileMoney,
		ProviderCode:      "mtn_ug",
		AccountNumber:     "256777000001",
		CustomerName:      "John Doe",
		MerchantReference: "MCTREFT2WMNWZ23SBN6Y",
		Description:       "Test payout",
	}
	result, resp, err := suite.testable.Create(suite.ctx, request)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), resp)
	assert.NotEmpty(suite.T(), result)
	assert.True(suite.T(), result.IsSuccess())
	assert.Equal(suite.T(), http.StatusAccepted, result.Code)
	assert.Equal(suite.T(), "DUSUPAYRMGRXNNYBWATKJ", result.Data.InternalReference)
	assert.Equal(suite.T(), "MCTREFT2WMNWZ23SBN6Y", result.Data.MerchantReference)

	defer func() { _ = resp.Body.Close() }()
	bodyRsp, _ := io.ReadAll(resp.Body)
	assert.Equal(suite.T(), body, bodyRsp)
}

func (suite *PayoutsResourceTestSuite) TestCreateJsonError() {
	body, _ := LoadStubResponseData("stubs/errors/401.json")
	httpmock.RegisterResponder(http.MethodPost, suite.cfg.URI+"/payout/send-funds", httpmock.NewBytesResponder(http.StatusUnauthorized, body))

	request := &PayoutRequest{Currency: CurrencyCodeUGX, Amount: 100, TransactionMethod: TransactionMethodMobileMoney, ProviderCode: "mtn_ug", AccountNumber: "256777000001", CustomerName: "John Doe", MerchantReference: "ref", Description: "desc"}
	result, resp, err := suite.testable.Create(suite.ctx, request)
	defer func() { _ = resp.Body.Close() }()
	assert.Error(suite.T(), err)
	assert.NotEmpty(suite.T(), resp)
	assert.NotEmpty(suite.T(), result)
	assert.Equal(suite.T(), "Unauthorized API access. Unknown Merchant", err.Error())
}

func (suite *PayoutsResourceTestSuite) TestCreateNonJsonError() {
	body, _ := LoadStubResponseData("stubs/errors/500.html")
	httpmock.RegisterResponder(http.MethodPost, suite.cfg.URI+"/payout/send-funds", httpmock.NewBytesResponder(http.StatusInternalServerError, body))

	request := &PayoutRequest{Currency: CurrencyCodeUGX, Amount: 100, TransactionMethod: TransactionMethodMobileMoney, ProviderCode: "mtn_ug", AccountNumber: "256777000001", CustomerName: "John Doe", MerchantReference: "ref", Description: "desc"}
	result, resp, err := suite.testable.Create(suite.ctx, request)
	defer func() { _ = resp.Body.Close() }()
	assert.Error(suite.T(), err)
	assert.NotEmpty(suite.T(), resp)
	assert.Empty(suite.T(), result)
}

func (suite *PayoutsResourceTestSuite) TestCreateInvalidRequest() {
	result, rsp, err := suite.testable.Create(suite.ctx, &PayoutRequest{}) //nolint:bodyclose // returns nil response on validation error
	assert.Nil(suite.T(), rsp)
	assert.Nil(suite.T(), result)
	assert.Error(suite.T(), err)
}

func TestPayoutsResourceTestSuite(t *testing.T) {
	suite.Run(t, new(PayoutsResourceTestSuite))
}
