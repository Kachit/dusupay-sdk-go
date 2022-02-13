package dusupay

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Common_TransformStructToMapSuccess(t *testing.T) {
	req := &CollectionRequest{}
	req.Amount = 100
	req.ProviderId = "foo"
	req.Currency = CurrencyCodeEUR
	req.Method = TransactionMethodCard
	result, err := transformStructToMap(req)
	assert.NoError(t, err)
	assert.Equal(t, req.Amount, result["amount"])
	assert.Equal(t, req.ProviderId, result["provider_id"])
	assert.Equal(t, string(req.Method), result["method"])
	assert.Equal(t, string(req.Method), result["method"])
}

func Test_Common_TransformStructToMapError(t *testing.T) {
	result, err := transformStructToMap("foo")
	assert.Error(t, err)
	assert.Nil(t, result)
}
