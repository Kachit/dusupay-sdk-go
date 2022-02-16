package dusupay

import (
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Payouts_PayoutsRequest_IsValidSuccess(t *testing.T) {
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
	assert.Nil(t, request.isValid())
	assert.NoError(t, request.isValid())
}

func Test_Payouts_PayoutsRequest_IsValidEmptyCurrency(t *testing.T) {
	request := PayoutRequest{
		Amount:            100,
		Method:            TransactionMethodBank,
		ProviderId:        "provider_id",
		MerchantReference: "merchant_reference",
		Narration:         "narration",
	}
	result := request.isValid()
	assert.Error(t, result)
	assert.Equal(t, `parameter "currency" is empty`, result.Error())
}

func Test_Payouts_PayoutsRequest_IsValidEmptyAmount(t *testing.T) {
	request := PayoutRequest{
		Currency:          CurrencyCodeKES,
		Method:            TransactionMethodBank,
		ProviderId:        "provider_id",
		MerchantReference: "merchant_reference",
		Narration:         "narration",
	}
	result := request.isValid()
	assert.Error(t, result)
	assert.Equal(t, `parameter "amount" is empty`, result.Error())
}

func Test_Payouts_PayoutsRequest_IsValidEmptyMethod(t *testing.T) {
	request := PayoutRequest{
		Currency:          CurrencyCodeKES,
		Amount:            100,
		ProviderId:        "provider_id",
		MerchantReference: "merchant_reference",
		Narration:         "narration",
	}
	result := request.isValid()
	assert.Error(t, result)
	assert.Equal(t, `parameter "method" is empty`, result.Error())
}

func Test_Payouts_PayoutsRequest_IsValidEmptyProviderId(t *testing.T) {
	request := PayoutRequest{
		Currency:          CurrencyCodeKES,
		Amount:            100,
		Method:            TransactionMethodBank,
		MerchantReference: "merchant_reference",
		Narration:         "narration",
	}
	result := request.isValid()
	assert.Error(t, result)
	assert.Equal(t, `parameter "provider_id" is empty`, result.Error())
}

func Test_Payouts_PayoutsRequest_IsValidEmptyMerchantReference(t *testing.T) {
	request := PayoutRequest{
		Currency:   CurrencyCodeKES,
		Amount:     100,
		Method:     TransactionMethodBank,
		ProviderId: "provider_id",
		Narration:  "narration",
	}
	result := request.isValid()
	assert.Error(t, result)
	assert.Equal(t, `parameter "merchant_reference" is empty`, result.Error())
}

func Test_Payouts_PayoutsRequest_IsValidEmptyNarration(t *testing.T) {
	request := PayoutRequest{
		Currency:          CurrencyCodeKES,
		Amount:            100,
		Method:            TransactionMethodBank,
		ProviderId:        "provider_id",
		MerchantReference: "merchant_reference",
	}
	result := request.isValid()
	assert.Error(t, result)
	assert.Equal(t, `parameter "narration" is empty`, result.Error())
}

func Test_Payouts_PayoutsRequest_IsValidEmptyAccountNumber(t *testing.T) {
	request := PayoutRequest{
		Currency:          CurrencyCodeKES,
		Amount:            100,
		Method:            TransactionMethodBank,
		ProviderId:        "provider_id",
		MerchantReference: "merchant_reference",
		Narration:         "narration",
	}
	result := request.isValid()
	assert.Error(t, result)
	assert.Equal(t, `parameter "account_number" is empty`, result.Error())
}

func Test_Payouts_PayoutsRequest_IsValidEmptyAccountName(t *testing.T) {
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
	assert.Error(t, result)
	assert.Equal(t, `parameter "account_name" is empty`, result.Error())
}

func Test_Payouts_PayoutResponse_UnmarshalSuccess(t *testing.T) {
	var response PayoutResponse
	body, _ := LoadStubResponseData("stubs/payouts/create/success.json")
	err := json.Unmarshal(body, &response)
	assert.NoError(t, err)
	assert.Equal(t, 202, response.Code)
	assert.Equal(t, "accepted", response.Status)
	assert.Equal(t, "Transaction Initiated", response.Message)
	assert.Equal(t, int64(124468), response.Data.ID)
	assert.Equal(t, float64(700), response.Data.RequestAmount)
	assert.Equal(t, "UGX", response.Data.RequestCurrency)
	assert.Equal(t, float64(700), response.Data.AccountAmount)
	assert.Equal(t, "UGX", response.Data.AccountCurrency)
	assert.Equal(t, float64(1500), response.Data.TransactionFee)
	assert.Equal(t, float64(2200), response.Data.TotalDebit)
	assert.Equal(t, "mtn_ug", response.Data.ProviderID)
	assert.Equal(t, "payout-1005", response.Data.MerchantReference)
	assert.Equal(t, "DUSUPAY405GZMDVTKASJL8UQ", response.Data.InternalReference)
	assert.Equal(t, "PENDING", response.Data.TransactionStatus)
	assert.Equal(t, "payout", response.Data.TransactionType)
	assert.Equal(t, "Transaction Initiated", response.Data.Message)
}

func Test_Payouts_PayoutResponse_UnmarshalErrorUnauthorized(t *testing.T) {
	var response RefundResponse
	body, _ := LoadStubResponseData("stubs/errors/401.json")
	err := json.Unmarshal(body, &response)
	assert.NoError(t, err)
	assert.Equal(t, 401, response.Code)
	assert.Equal(t, "error", response.Status)
	assert.Equal(t, "Unauthorized API access. Unknown Merchant", response.Message)
	assert.Empty(t, response.Data)
}

func Test_Payouts_PayoutsResource_CreateInvalidRequest(t *testing.T) {
	config := BuildStubConfig()
	transport := NewHttpTransport(config, nil)
	ctx := context.Background()
	req := &PayoutRequest{}
	resource := &PayoutsResource{ResourceAbstract: NewResourceAbstract(transport, config)}
	rsp, err := resource.Create(ctx, req)
	assert.Nil(t, rsp)
	assert.Error(t, err)
}
