package dusupay

import (
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Refunds_RefundsRequest_IsValidSuccess(t *testing.T) {
	request := RefundRequest{
		Amount:            100,
		InternalReference: "internal_reference",
	}
	assert.Nil(t, request.isValid())
	assert.NoError(t, request.isValid())
}

func Test_Refunds_RefundsRequest_IsValidEmptyInternalReference(t *testing.T) {
	request := RefundRequest{
		Amount: 100,
	}
	result := request.isValid()
	assert.Error(t, result)
	assert.Equal(t, `parameter "internal_reference" is empty`, result.Error())
}

func Test_Refunds_RefundResponse_UnmarshalSuccess(t *testing.T) {
	var response RefundResponse
	body, _ := LoadStubResponseData("stubs/refunds/create/success.json")
	err := json.Unmarshal(body, &response)
	assert.NoError(t, err)
	assert.Equal(t, 202, response.Code)
	assert.Equal(t, "accepted", response.Status)
	assert.Equal(t, "Refund Initiated Successfully", response.Message)
	assert.Equal(t, int64(65205), response.Data.ID)
	assert.Equal(t, float64(1054), response.Data.RefundAmount)
	assert.Equal(t, "UGX", response.Data.RefundCurrency)
	assert.Equal(t, float64(0), response.Data.TransactionFee)
	assert.Equal(t, float64(1054), response.Data.TotalDebit)
	assert.Equal(t, "international_ugx", response.Data.ProviderID)
	assert.Equal(t, "hAkEROAdhIsHrEnB", response.Data.MerchantReference)
	assert.Equal(t, "DUSUPAYXYXYXYXYXYXYXYXYX", response.Data.CollectionReference)
	assert.Equal(t, "RFD-DUSUPAYXYXYXYXYXYXYXYXYX-3486003", response.Data.InternalReference)
	assert.Equal(t, "PENDING", response.Data.TransactionStatus)
	assert.Equal(t, "refund", response.Data.TransactionType)
	assert.Equal(t, "4860610032773134", response.Data.AccountNumber)
	assert.Equal(t, "Request Initiated", response.Data.Message)
}

func Test_Refunds_RefundResponse_UnmarshalErrorUnauthorized(t *testing.T) {
	var response RefundResponse
	body, _ := LoadStubResponseData("stubs/errors/401.json")
	err := json.Unmarshal(body, &response)
	assert.NoError(t, err)
	assert.Equal(t, 401, response.Code)
	assert.Equal(t, "error", response.Status)
	assert.Equal(t, "Unauthorized API access. Unknown Merchant", response.Message)
	assert.Empty(t, response.Data)
}

func Test_Refunds_RefundsResource_CreateInvalidRequest(t *testing.T) {
	config := BuildStubConfig()
	transport := NewHttpTransport(config, nil)
	ctx := context.Background()
	req := &RefundRequest{}
	resource := &RefundsResource{ResourceAbstract: NewResourceAbstract(transport, config)}
	rsp, err := resource.Create(ctx, req)
	assert.Nil(t, rsp)
	assert.Error(t, err)
}
