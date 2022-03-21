package dusupay

import (
	"context"
	"encoding/json"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

func Test_Webhooks_CollectionWebhook_UnmarshalSuccess(t *testing.T) {
	var webhook CollectionWebhook
	body, _ := LoadStubResponseData("stubs/webhooks/request/collection-success.json")
	err := json.Unmarshal(body, &webhook)
	assert.NoError(t, err)
	assert.Equal(t, int64(226), webhook.ID)
	assert.Equal(t, 0.2, webhook.RequestAmount)
	assert.Equal(t, "USD", webhook.RequestCurrency)
	assert.Equal(t, 737.9934, webhook.AccountAmount)
	assert.Equal(t, "UGX", webhook.AccountCurrency)
	assert.Equal(t, 21.4018, webhook.TransactionFee)
	assert.Equal(t, 716.5916, webhook.TotalCredit)
	assert.Equal(t, false, webhook.CustomerCharged)
	assert.Equal(t, "mtn_ug", webhook.ProviderID)
	assert.Equal(t, "76859aae-f148-48c5-9901-2e474cf19b71", webhook.MerchantReference)
	assert.Equal(t, "DUSUPAY405GZM1G5JXGA71IK", webhook.InternalReference)
	assert.Equal(t, "COMPLETED", webhook.TransactionStatus)
	assert.Equal(t, "collection", webhook.TransactionType)
	assert.Equal(t, "Transaction Completed Successfully", webhook.Message)
	assert.Equal(t, "256777111786 - Optional", webhook.AccountNumber)
	assert.Equal(t, "- Optional", webhook.AccountName)
	assert.Equal(t, "MTN Mobile Money - Optional", webhook.InstitutionName)
}

func Test_Webhooks_PayoutWebhook_UnmarshalSuccess(t *testing.T) {
	var webhook PayoutWebhook
	body, _ := LoadStubResponseData("stubs/webhooks/request/payout-success.json")
	err := json.Unmarshal(body, &webhook)
	assert.NoError(t, err)
	assert.Equal(t, int64(226), webhook.ID)
	assert.Equal(t, 0.2, webhook.RequestAmount)
	assert.Equal(t, "USD", webhook.RequestCurrency)
	assert.Equal(t, 737.9934, webhook.AccountAmount)
	assert.Equal(t, "UGX", webhook.AccountCurrency)
	assert.Equal(t, 21.4018, webhook.TransactionFee)
	assert.Equal(t, 716.5916, webhook.TotalDebit)
	assert.Equal(t, "mtn_ug", webhook.ProviderID)
	assert.Equal(t, "76859aae-f148-48c5-9901-2e474cf19b71", webhook.MerchantReference)
	assert.Equal(t, "DUSUPAY405GZM1G5JXGA71IK", webhook.InternalReference)
	assert.Equal(t, "COMPLETED", webhook.TransactionStatus)
	assert.Equal(t, "payout", webhook.TransactionType)
	assert.Equal(t, "Transaction Completed Successfully", webhook.Message)
	assert.Equal(t, "256777111786 - Optional", webhook.AccountNumber)
	assert.Equal(t, "- Optional", webhook.AccountName)
	assert.Equal(t, "MTN Mobile Money - Optional", webhook.InstitutionName)
}

func Test_Webhooks_RefundWebhook_UnmarshalSuccess(t *testing.T) {
	var webhook RefundWebhook
	body, _ := LoadStubResponseData("stubs/webhooks/request/refund-success.json")
	err := json.Unmarshal(body, &webhook)
	assert.NoError(t, err)
	assert.Equal(t, int64(65205), webhook.ID)
	assert.Equal(t, float64(1054), webhook.RefundAmount)
	assert.Equal(t, "UGX", webhook.RefundCurrency)
	assert.Equal(t, float64(0), webhook.TransactionFee)
	assert.Equal(t, float64(1054), webhook.TotalDebit)
	assert.Equal(t, "international_ugx", webhook.ProviderID)
	assert.Equal(t, "DUSUPAYXYXYXYXYXYXYXYXYX", webhook.CollectionReference)
	assert.Equal(t, "RFD-DUSUPAYXYXYXYXYXYXYXYXYX-3486003", webhook.InternalReference)
	assert.Equal(t, "COMPLETED", webhook.TransactionStatus)
	assert.Equal(t, "refund", webhook.TransactionType)
	assert.Equal(t, "Refund Processed Successfully", webhook.Message)
	assert.Equal(t, "4860610032773134", webhook.AccountNumber)
}

func Test_Webhooks_WebhooksResource_SendCallbackSuccess(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	cfg := BuildStubConfig()
	transport := BuildStubHttpTransport()

	body, _ := LoadStubResponseData("stubs/webhooks/send-callback/success.json")
	httpmock.RegisterResponder(http.MethodGet, cfg.Uri+"/v1/send-callback/qwerty", httpmock.NewBytesResponder(http.StatusOK, body))

	ctx := context.Background()

	resource := &WebhooksResource{ResourceAbstract: NewResourceAbstract(transport, cfg)}
	result, resp, err := resource.SendCallback(ctx, "qwerty")
	assert.NoError(t, err)
	assert.NotEmpty(t, resp)
	assert.NotEmpty(t, result)
	//result
	assert.True(t, result.IsSuccess())
	assert.Equal(t, 200, result.Code)
	assert.Equal(t, "success", result.Status)
	assert.Equal(t, "Callback Initiated Successfully", result.Message)
	assert.Equal(t, int64(613589), result.Data.Payload.ID)
	assert.Equal(t, float64(520000), result.Data.Payload.RequestAmount)
	assert.Equal(t, "XAF", result.Data.Payload.RequestCurrency)
	assert.Equal(t, 791.44, result.Data.Payload.AccountAmount)
	assert.Equal(t, "EUR", result.Data.Payload.AccountCurrency)
	assert.Equal(t, 38.7806, result.Data.Payload.TransactionFee)
	assert.Equal(t, "international_eur", result.Data.Payload.ProviderID)
	assert.Equal(t, "123456789", result.Data.Payload.MerchantReference)
	assert.Equal(t, "DUSUPAY5FNZCVUKZ8C0KZE", result.Data.Payload.InternalReference)
	assert.Equal(t, "COMPLETED", result.Data.Payload.TransactionStatus)
	assert.Equal(t, "collection", result.Data.Payload.TransactionType)
	assert.Equal(t, "Transaction Completed Successfully", result.Data.Payload.Message)
	//response
	defer resp.Body.Close()
	bodyRsp, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, body, bodyRsp)
}

func Test_Webhooks_WebhooksResource_SendCallbackJsonError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	cfg := BuildStubConfig()
	transport := BuildStubHttpTransport()

	body, _ := LoadStubResponseData("stubs/errors/401.json")
	httpmock.RegisterResponder(http.MethodGet, cfg.Uri+"/v1/send-callback/qwerty", httpmock.NewBytesResponder(http.StatusOK, body))

	ctx := context.Background()

	resource := &WebhooksResource{ResourceAbstract: NewResourceAbstract(transport, cfg)}
	result, resp, err := resource.SendCallback(ctx, "qwerty")
	assert.Error(t, err)
	assert.NotEmpty(t, resp)
	assert.NotEmpty(t, result)
	//result
	assert.False(t, result.IsSuccess())
	assert.Equal(t, 401, result.Code)
	assert.Equal(t, "error", result.Status)
	assert.Equal(t, "Unauthorized API access. Unknown Merchant", result.Message)
	assert.Empty(t, result.Data)
	//response
	defer resp.Body.Close()
	bodyRsp, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, body, bodyRsp)
	//error
	assert.Equal(t, "Unauthorized API access. Unknown Merchant", err.Error())
}

func Test_Webhooks_WebhooksResource_SendCallbackNonJsonError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	cfg := BuildStubConfig()
	transport := BuildStubHttpTransport()

	body, _ := LoadStubResponseData("stubs/errors/500.html")
	httpmock.RegisterResponder(http.MethodGet, cfg.Uri+"/v1/send-callback/qwerty", httpmock.NewBytesResponder(http.StatusOK, body))

	ctx := context.Background()

	resource := &WebhooksResource{ResourceAbstract: NewResourceAbstract(transport, cfg)}
	result, resp, err := resource.SendCallback(ctx, "qwerty")
	assert.Error(t, err)
	assert.NotEmpty(t, resp)
	assert.Empty(t, result)
	//response
	defer resp.Body.Close()
	bodyRsp, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, body, bodyRsp)
}
