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
	assert.Equal(t, 200, response.Code)
	assert.Equal(t, "success", response.Status)
	assert.Equal(t, "Request completed successfully.", response.Message)
	assert.Equal(t, "UGX", (*response.Data)[0].Currency)
	assert.Equal(t, 5475.816, (*response.Data)[0].Balance)
	assert.Equal(t, "USD", (*response.Data)[1].Currency)
	assert.Equal(t, float64(12), (*response.Data)[1].Balance)
}

func Test_Merchants_BalancesResponse_UnmarshalErrorUnauthorized(t *testing.T) {
	var response BalancesResponse
	body, _ := LoadStubResponseData("stubs/errors/401.json")
	err := json.Unmarshal(body, &response)
	assert.NoError(t, err)
	assert.Equal(t, 401, response.Code)
	assert.Equal(t, "error", response.Status)
	assert.Equal(t, "Unauthorized API access. Unknown Merchant", response.Message)
	assert.Empty(t, response.Data)
}
