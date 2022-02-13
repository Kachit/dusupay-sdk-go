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
