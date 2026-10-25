package dusupay

import (
	"encoding/json"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type WebhooksTestSuite struct {
	suite.Suite
}

func (suite *WebhooksTestSuite) TestCollectionCallbackUnmarshalSuccess() {
	var callback CollectionCallback
	body, _ := LoadStubResponseData("stubs/webhooks/callback/collection-success.json")
	err := json.Unmarshal(body, &callback)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), CallbackEventTransactionCompleted, callback.Event)
	payload := callback.Payload
	assert.Equal(suite.T(), int64(226), payload.ID)
	assert.Equal(suite.T(), "76859aae-f148-48c5-9901-2e474cf19b71", payload.MerchantReference)
	assert.Equal(suite.T(), "DUSUPAY405GZM1G5JXGA71IK", payload.InternalReference)
	assert.Equal(suite.T(), "COLLECTION", payload.TransactionType)
	assert.Equal(suite.T(), "COMPLETED", payload.TransactionStatus)
	assert.Equal(suite.T(), 0.2, payload.RequestAmount)
	assert.Equal(suite.T(), "USD", payload.RequestCurrency)
	assert.Equal(suite.T(), 737.9934, payload.TransactionAmount)
	assert.Equal(suite.T(), "UGX", payload.TransactionCurrency)
	assert.Equal(suite.T(), "256777111786", payload.TransactionAccount)
	assert.Equal(suite.T(), 716.5916, payload.TotalCredit)
	assert.Equal(suite.T(), "mtn_ug", payload.ProviderCode)
	assert.Equal(suite.T(), "JOHN DOE", payload.CustomerName)
	assert.Equal(suite.T(), false, payload.ChargeCustomer)
	assert.Nil(suite.T(), payload.Authorization)
}

func (suite *WebhooksTestSuite) TestCollectionCallbackProcessingOTP() {
	var callback CollectionCallback
	body, _ := LoadStubResponseData("stubs/webhooks/callback/collection-processing-otp.json")
	err := json.Unmarshal(body, &callback)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), CallbackEventTransactionProcessing, callback.Event)
	assert.Equal(suite.T(), string(TransactionStatusProcessing), callback.Payload.TransactionStatus)
	assert.Equal(suite.T(), AuthorizationModeOTP, callback.Payload.Authorization.Mode)
	assert.Equal(suite.T(), "otp", callback.Payload.Authorization.RequiredParameters[0].ParameterName)
	assert.Contains(suite.T(), callback.Payload.PaymentInstructions, "Authorize Payment using the OTP")
}

func (suite *WebhooksTestSuite) TestCollectionCallbackProcessingBank() {
	var callback CollectionCallback
	body, _ := LoadStubResponseData("stubs/webhooks/callback/collection-processing-bank.json")
	err := json.Unmarshal(body, &callback)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "2121562123", callback.Payload.SubAccountDetails.AccountNo)
	assert.Equal(suite.T(), "WEMA", callback.Payload.SubAccountDetails.BankName)
	assert.Equal(suite.T(), float64(4000), callback.Payload.SubAccountDetails.Amount)
	assert.Equal(suite.T(), AuthorizationModeCode(""), callback.Payload.Authorization.Mode)
}

func (suite *WebhooksTestSuite) TestCollectionCallbackBuildPayloadString() {
	var callback CollectionCallback
	body, _ := LoadStubResponseData("stubs/webhooks/callback/collection-success.json")
	_ = json.Unmarshal(body, &callback)
	result := callback.BuildPayloadString()
	assert.Equal(suite.T(), "transaction.completed:76859aae-f148-48c5-9901-2e474cf19b71:DUSUPAY405GZM1G5JXGA71IK:COLLECTION:COMPLETED", result)
}

func (suite *WebhooksTestSuite) TestCollectionCallbackBuildPayloadStringNilPayload() {
	callback := &CollectionCallback{Event: CallbackEventRequestFailed}
	assert.Equal(suite.T(), "request.failed::::", callback.BuildPayloadString())
}

func (suite *WebhooksTestSuite) TestPayoutCallbackUnmarshalSuccess() {
	var callback PayoutCallback
	body, _ := LoadStubResponseData("stubs/webhooks/callback/payout-success.json")
	err := json.Unmarshal(body, &callback)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), CallbackEventTransactionCompleted, callback.Event)
	payload := callback.Payload
	assert.Equal(suite.T(), int64(227), payload.ID)
	assert.Equal(suite.T(), "76859aae-f148-48c5-9901-2e474cf19b72", payload.MerchantReference)
	assert.Equal(suite.T(), "DUSUPAY405GZM1G5JXGA71IZ", payload.InternalReference)
	assert.Equal(suite.T(), "PAYOUT", payload.TransactionType)
	assert.Equal(suite.T(), "COMPLETED", payload.TransactionStatus)
	assert.Equal(suite.T(), float64(10000), payload.RequestAmount)
	assert.Equal(suite.T(), "UGX", payload.RequestCurrency)
	assert.Equal(suite.T(), "256777111786", payload.TransactionAccount)
	assert.Equal(suite.T(), float64(10300), payload.TotalDebit)
}

func (suite *WebhooksTestSuite) TestPayoutCallbackBuildPayloadString() {
	var callback PayoutCallback
	body, _ := LoadStubResponseData("stubs/webhooks/callback/payout-success.json")
	_ = json.Unmarshal(body, &callback)
	result := callback.BuildPayloadString()
	assert.Equal(suite.T(), "transaction.completed:76859aae-f148-48c5-9901-2e474cf19b72:DUSUPAY405GZM1G5JXGA71IZ:PAYOUT:COMPLETED", result)
}

func (suite *WebhooksTestSuite) TestPayoutCallbackBuildPayloadStringNilPayload() {
	callback := &PayoutCallback{Event: CallbackEventRequestFailed}
	assert.Equal(suite.T(), "request.failed::::", callback.BuildPayloadString())
}

func (suite *WebhooksTestSuite) TestNewRedirectPayload() {
	query, _ := url.ParseQuery("id=34639&event=transaction.failed&merchant_reference=MCTREF8DQRHYTNCQTACH" +
		"&internal_reference=DUSUPAY8A54ZUJFGSMENE&transaction_type=COLLECTION&transaction_status=FAILED" +
		"&status_message=Balance%20Insufficient%20for%20the%20transaction&rsa_signature=yvODC%2BVi" +
		"&hmac_signature=t%3D1721967214552%2Cs%3D0a3505923c")
	payload := NewRedirectPayload(query)
	assert.Equal(suite.T(), int64(34639), payload.ID)
	assert.Equal(suite.T(), CallbackEventTransactionFailed, payload.Event)
	assert.Equal(suite.T(), "Balance Insufficient for the transaction", payload.StatusMessage)
	assert.Equal(suite.T(), "yvODC+Vi", payload.RSASignature)
	assert.Equal(suite.T(), "t=1721967214552,s=0a3505923c", payload.HMACSignature)
	assert.Equal(suite.T(), "transaction.failed:MCTREF8DQRHYTNCQTACH:DUSUPAY8A54ZUJFGSMENE:COLLECTION:FAILED", payload.BuildPayloadString())
}

func TestWebhooksTestSuite(t *testing.T) {
	suite.Run(t, new(WebhooksTestSuite))
}
