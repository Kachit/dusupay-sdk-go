package dusupay

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Merchants_BalancesResponse_UnmarshalSuccess(t *testing.T) {
	var response BalancesResponse
	body, _ := LoadStubResponseData("stubs/merchants/balance/success.json")
	err := json.Unmarshal(body, &response)
	assert.NoError(t, err)
}

func Test_Merchants_BalancesResponse_UnmarshalUnauthorized(t *testing.T) {
	var response BalancesResponse
	body, _ := LoadStubResponseData("stubs/errors/401.json")
	err := json.Unmarshal(body, &response)
	assert.Error(t, err)
}
