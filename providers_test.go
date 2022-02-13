package dusupay

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Providers_ProvidersFilter_IsValidSuccess(t *testing.T) {
	filter := ProvidersFilter{Country: CountryCodeKenya, Method: TransactionMethodCard, TransactionType: TransactionTypeCollection}
	assert.Nil(t, filter.isValid())
	assert.NoError(t, filter.isValid())
}

func Test_Providers_ProvidersFilter_IsValidEmptyCountryCode(t *testing.T) {
	filter := ProvidersFilter{Method: TransactionMethodCard}
	result := filter.isValid()
	assert.Error(t, result)
	assert.Equal(t, `parameter "country_code" is empty`, result.Error())
}

func Test_Providers_ProvidersFilter_IsValidEmptyMethod(t *testing.T) {
	filter := ProvidersFilter{Country: CountryCodeKenya}
	result := filter.isValid()
	assert.Error(t, result)
	assert.Equal(t, `parameter "method" is empty`, result.Error())
}

func Test_Providers_ProvidersFilter_IsValidEmptyTransactionType(t *testing.T) {
	filter := ProvidersFilter{Country: CountryCodeKenya, Method: TransactionMethodCard}
	result := filter.isValid()
	assert.Error(t, result)
	assert.Equal(t, `parameter "transaction_type" is empty`, result.Error())
}

func Test_Providers_ProvidersFilter_BuildPath(t *testing.T) {
	filter := ProvidersFilter{Country: CountryCodeKenya, Method: TransactionMethodCard, TransactionType: TransactionTypeCollection}
	result := filter.buildPath()
	assert.Equal(t, `COLLECTION/CARD/KE`, result)
}

func Test_Providers_ProvidersResponse_UnmarshalSandbox(t *testing.T) {
	var response ProvidersResponse
	body, _ := LoadStubResponseData("stubs/providers/payment-options/success-sandbox.json")
	err := json.Unmarshal(body, &response)
	assert.NoError(t, err)
}

func Test_Providers_ProvidersResponse_Unmarshal(t *testing.T) {
	var response ProvidersResponse
	body, _ := LoadStubResponseData("stubs/providers/payment-options/success.json")
	err := json.Unmarshal(body, &response)
	assert.NoError(t, err)
}
