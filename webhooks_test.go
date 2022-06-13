package dusupay

import (
	"context"
	"encoding/json"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"net/http"
	"testing"
)

type WebhooksTestSuite struct {
	suite.Suite
}

func (suite *WebhooksTestSuite) TestCollectionWebhook_UnmarshalSuccess() {
	var webhook CollectionWebhook
	body, _ := LoadStubResponseData("stubs/webhooks/request/collection-success.json")
	err := json.Unmarshal(body, &webhook)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(226), webhook.ID)
	assert.Equal(suite.T(), 0.2, webhook.RequestAmount)
	assert.Equal(suite.T(), "USD", webhook.RequestCurrency)
	assert.Equal(suite.T(), 737.9934, webhook.AccountAmount)
	assert.Equal(suite.T(), "UGX", webhook.AccountCurrency)
	assert.Equal(suite.T(), 21.4018, webhook.TransactionFee)
	assert.Equal(suite.T(), 716.5916, webhook.TotalCredit)
	assert.Equal(suite.T(), false, webhook.CustomerCharged)
	assert.Equal(suite.T(), "mtn_ug", webhook.ProviderID)
	assert.Equal(suite.T(), "76859aae-f148-48c5-9901-2e474cf19b71", webhook.MerchantReference)
	assert.Equal(suite.T(), "DUSUPAY405GZM1G5JXGA71IK", webhook.InternalReference)
	assert.Equal(suite.T(), "COMPLETED", webhook.TransactionStatus)
	assert.Equal(suite.T(), "collection", webhook.TransactionType)
	assert.Equal(suite.T(), "Transaction Completed Successfully", webhook.Message)
	assert.Equal(suite.T(), "256777111786 - Optional", webhook.AccountNumber)
	assert.Equal(suite.T(), "- Optional", webhook.AccountName)
	assert.Equal(suite.T(), "MTN Mobile Money - Optional", webhook.InstitutionName)
}

func (suite *WebhooksTestSuite) TestPayoutWebhook_UnmarshalSuccess() {
	var webhook PayoutWebhook
	body, _ := LoadStubResponseData("stubs/webhooks/request/payout-success.json")
	err := json.Unmarshal(body, &webhook)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(226), webhook.ID)
	assert.Equal(suite.T(), 0.2, webhook.RequestAmount)
	assert.Equal(suite.T(), "USD", webhook.RequestCurrency)
	assert.Equal(suite.T(), 737.9934, webhook.AccountAmount)
	assert.Equal(suite.T(), "UGX", webhook.AccountCurrency)
	assert.Equal(suite.T(), 21.4018, webhook.TransactionFee)
	assert.Equal(suite.T(), 716.5916, webhook.TotalDebit)
	assert.Equal(suite.T(), "mtn_ug", webhook.ProviderID)
	assert.Equal(suite.T(), "76859aae-f148-48c5-9901-2e474cf19b71", webhook.MerchantReference)
	assert.Equal(suite.T(), "DUSUPAY405GZM1G5JXGA71IK", webhook.InternalReference)
	assert.Equal(suite.T(), "COMPLETED", webhook.TransactionStatus)
	assert.Equal(suite.T(), "payout", webhook.TransactionType)
	assert.Equal(suite.T(), "Transaction Completed Successfully", webhook.Message)
	assert.Equal(suite.T(), "256777111786 - Optional", webhook.AccountNumber)
	assert.Equal(suite.T(), "- Optional", webhook.AccountName)
	assert.Equal(suite.T(), "MTN Mobile Money - Optional", webhook.InstitutionName)
}

func (suite *WebhooksTestSuite) TestRefundWebhook_UnmarshalSuccess() {
	var webhook RefundWebhook
	body, _ := LoadStubResponseData("stubs/webhooks/request/refund-success.json")
	err := json.Unmarshal(body, &webhook)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(65205), webhook.ID)
	assert.Equal(suite.T(), float64(1054), webhook.RefundAmount)
	assert.Equal(suite.T(), "UGX", webhook.RefundCurrency)
	assert.Equal(suite.T(), float64(0), webhook.TransactionFee)
	assert.Equal(suite.T(), float64(1054), webhook.TotalDebit)
	assert.Equal(suite.T(), "international_ugx", webhook.ProviderID)
	assert.Equal(suite.T(), "DUSUPAYXYXYXYXYXYXYXYXYX", webhook.CollectionReference)
	assert.Equal(suite.T(), "RFD-DUSUPAYXYXYXYXYXYXYXYXYX-3486003", webhook.InternalReference)
	assert.Equal(suite.T(), "COMPLETED", webhook.TransactionStatus)
	assert.Equal(suite.T(), "refund", webhook.TransactionType)
	assert.Equal(suite.T(), "Refund Processed Successfully", webhook.Message)
	assert.Equal(suite.T(), "4860610032773134", webhook.AccountNumber)
}

func TestWebhooksTestSuite(t *testing.T) {
	suite.Run(t, new(WebhooksTestSuite))
}

type WebhooksResourceTestSuite struct {
	suite.Suite
	cfg      *Config
	ctx      context.Context
	testable *WebhooksResource
}

func (suite *WebhooksResourceTestSuite) SetupTest() {
	cfg := BuildStubConfig()
	transport := BuildStubHttpTransport()
	suite.cfg = cfg
	suite.ctx = context.Background()
	suite.testable = &WebhooksResource{NewResourceAbstract(transport, cfg)}
	httpmock.Activate()
}

func (suite *WebhooksResourceTestSuite) TearDownTest() {
	httpmock.DeactivateAndReset()
}

func (suite *WebhooksResourceTestSuite) TestSendCallbackSuccess() {
	body, _ := LoadStubResponseData("stubs/webhooks/send-callback/success.json")
	httpmock.RegisterResponder(http.MethodGet, suite.cfg.Uri+"/v1/send-callback/qwerty", httpmock.NewBytesResponder(http.StatusOK, body))

	result, resp, err := suite.testable.SendCallback(suite.ctx, "qwerty")
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), resp)
	assert.NotEmpty(suite.T(), result)
	//result
	assert.True(suite.T(), result.IsSuccess())
	assert.Equal(suite.T(), 200, result.Code)
	assert.Equal(suite.T(), "success", result.Status)
	assert.Equal(suite.T(), "Callback Initiated Successfully", result.Message)
	assert.Equal(suite.T(), int64(613589), result.Data.Payload.ID)
	assert.Equal(suite.T(), float64(520000), result.Data.Payload.RequestAmount)
	assert.Equal(suite.T(), "XAF", result.Data.Payload.RequestCurrency)
	assert.Equal(suite.T(), 791.44, result.Data.Payload.AccountAmount)
	assert.Equal(suite.T(), "EUR", result.Data.Payload.AccountCurrency)
	assert.Equal(suite.T(), 38.7806, result.Data.Payload.TransactionFee)
	assert.Equal(suite.T(), "international_eur", result.Data.Payload.ProviderID)
	assert.Equal(suite.T(), "123456789", result.Data.Payload.MerchantReference)
	assert.Equal(suite.T(), "DUSUPAY5FNZCVUKZ8C0KZE", result.Data.Payload.InternalReference)
	assert.Equal(suite.T(), "COMPLETED", result.Data.Payload.TransactionStatus)
	assert.Equal(suite.T(), "collection", result.Data.Payload.TransactionType)
	assert.Equal(suite.T(), "Transaction Completed Successfully", result.Data.Payload.Message)
	//response
	defer resp.Body.Close()
	bodyRsp, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(suite.T(), body, bodyRsp)
}

func (suite *WebhooksResourceTestSuite) TestSendCallbackJsonError() {
	body, _ := LoadStubResponseData("stubs/errors/401.json")
	httpmock.RegisterResponder(http.MethodGet, suite.cfg.Uri+"/v1/send-callback/qwerty", httpmock.NewBytesResponder(http.StatusOK, body))

	result, resp, err := suite.testable.SendCallback(suite.ctx, "qwerty")
	assert.Error(suite.T(), err)
	assert.NotEmpty(suite.T(), resp)
	assert.NotEmpty(suite.T(), result)
	//result
	assert.False(suite.T(), result.IsSuccess())
	assert.Equal(suite.T(), 401, result.Code)
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

func (suite *WebhooksResourceTestSuite) TestSendCallbackNonJsonError() {
	body, _ := LoadStubResponseData("stubs/errors/500.html")
	httpmock.RegisterResponder(http.MethodGet, suite.cfg.Uri+"/v1/send-callback/qwerty", httpmock.NewBytesResponder(http.StatusOK, body))

	result, resp, err := suite.testable.SendCallback(suite.ctx, "qwerty")
	assert.Error(suite.T(), err)
	assert.NotEmpty(suite.T(), resp)
	assert.Empty(suite.T(), result)
	//response
	defer resp.Body.Close()
	bodyRsp, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(suite.T(), body, bodyRsp)
}

func TestWebhooksResourceTestSuite(t *testing.T) {
	suite.Run(t, new(WebhooksResourceTestSuite))
}
