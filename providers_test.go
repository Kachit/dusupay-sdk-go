package dusupay

import (
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
