package dusupay

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type WebhooksTestSuite struct {
	suite.Suite
}

func (suite *WebhooksTestSuite) TestCollectionCallbackPayloadUnmarshalSuccess() {
	var payload CollectionCallbackPayload
	body, _ := LoadStubResponseData("stubs/webhooks/callback/collection-success.json")
	err := json.Unmarshal(body, &payload)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "transaction.completed", payload.Event)
	assert.Equal(suite.T(), int64(226), payload.ID)
	assert.Equal(suite.T(), "76859aae-f148-48c5-9901-2e474cf19b71", payload.MerchantReference)
	assert.Equal(suite.T(), "DUSUPAY405GZM1G5JXGA71IK", payload.InternalReference)
	assert.Equal(suite.T(), "COLLECTION", payload.TransactionType)
	assert.Equal(suite.T(), "COMPLETED", payload.TransactionStatus)
	assert.Equal(suite.T(), 0.2, payload.RequestAmount)
	assert.Equal(suite.T(), "USD", payload.RequestCurrency)
	assert.Equal(suite.T(), "mtn_ug", payload.ProviderCode)
	assert.Equal(suite.T(), "JOHN DOE", payload.CustomerName)
	assert.Equal(suite.T(), false, payload.ChargeCustomer)
}

func (suite *WebhooksTestSuite) TestCollectionCallbackPayloadBuildPayloadString() {
	var payload CollectionCallbackPayload
	body, _ := LoadStubResponseData("stubs/webhooks/callback/collection-success.json")
	_ = json.Unmarshal(body, &payload)
	result := payload.BuildPayloadString()
	assert.Equal(suite.T(), "transaction.completed:76859aae-f148-48c5-9901-2e474cf19b71:DUSUPAY405GZM1G5JXGA71IK:COLLECTION:COMPLETED", result)
}

func (suite *WebhooksTestSuite) TestPayoutCallbackPayloadUnmarshalSuccess() {
	var payload PayoutCallbackPayload
	body, _ := LoadStubResponseData("stubs/webhooks/callback/payout-success.json")
	err := json.Unmarshal(body, &payload)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "transaction.completed", payload.Event)
	assert.Equal(suite.T(), int64(227), payload.ID)
	assert.Equal(suite.T(), "76859aae-f148-48c5-9901-2e474cf19b72", payload.MerchantReference)
	assert.Equal(suite.T(), "DUSUPAY405GZM1G5JXGA71IZ", payload.InternalReference)
	assert.Equal(suite.T(), "PAYOUT", payload.TransactionType)
	assert.Equal(suite.T(), "COMPLETED", payload.TransactionStatus)
	assert.Equal(suite.T(), float64(10000), payload.RequestAmount)
	assert.Equal(suite.T(), "UGX", payload.RequestCurrency)
	assert.Equal(suite.T(), "256777111786", payload.AccountNumber)
}

func (suite *WebhooksTestSuite) TestPayoutCallbackPayloadBuildPayloadString() {
	var payload PayoutCallbackPayload
	body, _ := LoadStubResponseData("stubs/webhooks/callback/payout-success.json")
	_ = json.Unmarshal(body, &payload)
	result := payload.BuildPayloadString()
	assert.Equal(suite.T(), "transaction.completed:76859aae-f148-48c5-9901-2e474cf19b72:DUSUPAY405GZM1G5JXGA71IZ:PAYOUT:COMPLETED", result)
}

func TestWebhooksTestSuite(t *testing.T) {
	suite.Run(t, new(WebhooksTestSuite))
}
