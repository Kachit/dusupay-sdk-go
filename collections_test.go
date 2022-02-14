package dusupay

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Collections_CollectionRequest_IsValidSuccess(t *testing.T) {
	request := CollectionRequest{
		Currency:          CurrencyCodeKES,
		Amount:            100,
		Method:            TransactionMethodBank,
		ProviderId:        "provider_id",
		MerchantReference: "merchant_reference",
		RedirectUrl:       "redirect_url",
		Narration:         "narration",
	}
	assert.Nil(t, request.isValid())
	assert.NoError(t, request.isValid())
}

func Test_Collections_CollectionRequest_IsValidEmptyCurrency(t *testing.T) {
	request := CollectionRequest{
		Amount:            100,
		Method:            TransactionMethodBank,
		ProviderId:        "provider_id",
		MerchantReference: "merchant_reference",
		RedirectUrl:       "redirect_url",
		Narration:         "narration",
	}
	result := request.isValid()
	assert.Error(t, result)
	assert.Equal(t, `parameter "currency" is empty`, result.Error())
}

func Test_Collections_CollectionRequest_IsValidEmptyAmount(t *testing.T) {
	request := CollectionRequest{
		Currency:          CurrencyCodeKES,
		Method:            TransactionMethodBank,
		ProviderId:        "provider_id",
		MerchantReference: "merchant_reference",
		RedirectUrl:       "redirect_url",
		Narration:         "narration",
	}
	result := request.isValid()
	assert.Error(t, result)
	assert.Equal(t, `parameter "amount" is empty`, result.Error())
}

func Test_Collections_CollectionRequest_IsValidEmptyMethod(t *testing.T) {
	request := CollectionRequest{
		Currency:          CurrencyCodeKES,
		Amount:            100,
		ProviderId:        "provider_id",
		MerchantReference: "merchant_reference",
		RedirectUrl:       "redirect_url",
		Narration:         "narration",
	}
	result := request.isValid()
	assert.Error(t, result)
	assert.Equal(t, `parameter "method" is empty`, result.Error())
}

func Test_Collections_CollectionRequest_IsValidEmptyProviderId(t *testing.T) {
	request := CollectionRequest{
		Currency:          CurrencyCodeKES,
		Amount:            100,
		Method:            TransactionMethodBank,
		MerchantReference: "merchant_reference",
		RedirectUrl:       "redirect_url",
		Narration:         "narration",
	}
	result := request.isValid()
	assert.Error(t, result)
	assert.Equal(t, `parameter "provider_id" is empty`, result.Error())
}

func Test_Collections_CollectionRequest_IsValidEmptyMerchantReference(t *testing.T) {
	request := CollectionRequest{
		Currency:    CurrencyCodeKES,
		Amount:      100,
		Method:      TransactionMethodBank,
		ProviderId:  "provider_id",
		RedirectUrl: "redirect_url",
		Narration:   "narration",
	}
	result := request.isValid()
	assert.Error(t, result)
	assert.Equal(t, `parameter "merchant_reference" is empty`, result.Error())
}

func Test_Collections_CollectionRequest_IsValidEmptyNarration(t *testing.T) {
	request := CollectionRequest{
		Currency:          CurrencyCodeKES,
		Amount:            100,
		Method:            TransactionMethodBank,
		ProviderId:        "provider_id",
		MerchantReference: "merchant_reference",
		RedirectUrl:       "redirect_url",
	}
	result := request.isValid()
	assert.Error(t, result)
	assert.Equal(t, `parameter "narration" is empty`, result.Error())
}

func Test_Collections_CollectionRequest_IsValidEmptyRedirectUrlByDefault(t *testing.T) {
	request := CollectionRequest{
		Currency:          CurrencyCodeKES,
		Amount:            100,
		Method:            TransactionMethodBank,
		ProviderId:        "provider_id",
		MerchantReference: "merchant_reference",
		Narration:         "narration",
	}
	result := request.isValid()
	assert.Error(t, result)
	assert.Equal(t, `parameter "redirect_url" is empty`, result.Error())
}

func Test_Collections_CollectionRequest_IsValidEmptyRedirectUrlMobileMoney(t *testing.T) {
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
	assert.Nil(t, result)
	assert.NoError(t, result)
}

func Test_Collections_CollectionRequest_IsValidEmptyRedirectUrlMobileMoneyWithoutAccountNumber(t *testing.T) {
	request := CollectionRequest{
		Currency:          CurrencyCodeKES,
		Amount:            100,
		Method:            TransactionMethodMobileMoney,
		ProviderId:        "provider_id",
		MerchantReference: "merchant_reference",
		Narration:         "narration",
	}
	result := request.isValid()
	assert.Error(t, result)
	assert.Equal(t, `parameter "account_number" is empty`, result.Error())
}
