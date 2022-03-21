package dusupay

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
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
