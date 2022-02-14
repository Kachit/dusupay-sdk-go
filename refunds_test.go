package dusupay

import (
	"context"
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
