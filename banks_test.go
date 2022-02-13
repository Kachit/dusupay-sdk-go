package dusupay

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Banks_BanksFilter_IsValidSuccess(t *testing.T) {
	filter := BanksFilter{Country: CountryCodeKenya, Method: "qwerty"}
	assert.Nil(t, filter.isValid())
	assert.NoError(t, filter.isValid())
}

func Test_Banks_BanksFilter_IsValidEmptyCountryCode(t *testing.T) {
	filter := BanksFilter{Method: "qwerty"}
	result := filter.isValid()
	assert.Error(t, result)
	assert.Equal(t, `parameter "country_code" is empty`, result.Error())
}

func Test_Banks_BanksFilter_IsValidEmptyMethod(t *testing.T) {
	filter := BanksFilter{Country: CountryCodeKenya}
	result := filter.isValid()
	assert.Error(t, result)
	assert.Equal(t, `parameter "method" is empty`, result.Error())
}

func Test_Banks_BanksBranchesFilter_IsValidSuccess(t *testing.T) {
	filter := BanksBranchesFilter{Country: CountryCodeKenya, Bank: "qwerty"}
	assert.Nil(t, filter.isValid())
	assert.NoError(t, filter.isValid())
}

func Test_Banks_BanksBranchesFilter_IsValidEmptyCountryCode(t *testing.T) {
	filter := BanksBranchesFilter{Bank: "qwerty"}
	result := filter.isValid()
	assert.Error(t, result)
	assert.Equal(t, `parameter "country_code" is empty`, result.Error())
}

func Test_Banks_BanksBranchesFilter_IsValidEmptyBankCode(t *testing.T) {
	filter := BanksBranchesFilter{Country: CountryCodeKenya}
	result := filter.isValid()
	assert.Error(t, result)
	assert.Equal(t, `parameter "bank_code" is empty`, result.Error())
}
