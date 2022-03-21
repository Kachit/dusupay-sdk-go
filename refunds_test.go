package dusupay

import (
	"context"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
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

func Test_Refunds_RefundsResource_CreateSuccess(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	cfg := BuildStubConfig()
	transport := BuildStubHttpTransport()

	body, _ := LoadStubResponseData("stubs/refunds/create/success.json")
	httpmock.RegisterResponder(http.MethodPost, cfg.Uri+"/v1/refund", httpmock.NewBytesResponder(http.StatusOK, body))

	ctx := context.Background()

	request := &RefundRequest{
		Amount:            100,
		InternalReference: "internal_reference",
	}
	resource := &RefundsResource{ResourceAbstract: NewResourceAbstract(transport, cfg)}
	result, resp, err := resource.Create(ctx, request)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp)
	assert.NotEmpty(t, result)
	//result
	assert.True(t, result.IsSuccess())
	assert.Equal(t, 202, result.Code)
	assert.Equal(t, "accepted", result.Status)
	assert.Equal(t, "Refund Initiated Successfully", result.Message)
	assert.Equal(t, int64(65205), result.Data.ID)
	assert.Equal(t, float64(1054), result.Data.RefundAmount)
	assert.Equal(t, "UGX", result.Data.RefundCurrency)
	assert.Equal(t, float64(0), result.Data.TransactionFee)
	assert.Equal(t, float64(1054), result.Data.TotalDebit)
	assert.Equal(t, "international_ugx", result.Data.ProviderID)
	assert.Equal(t, "hAkEROAdhIsHrEnB", result.Data.MerchantReference)
	assert.Equal(t, "DUSUPAYXYXYXYXYXYXYXYXYX", result.Data.CollectionReference)
	assert.Equal(t, "RFD-DUSUPAYXYXYXYXYXYXYXYXYX-3486003", result.Data.InternalReference)
	assert.Equal(t, "PENDING", result.Data.TransactionStatus)
	assert.Equal(t, "refund", result.Data.TransactionType)
	assert.Equal(t, "4860610032773134", result.Data.AccountNumber)
	assert.Equal(t, "Request Initiated", result.Data.Message)
	//response
	defer resp.Body.Close()
	bodyRsp, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, body, bodyRsp)
}

func Test_Refunds_RefundsResource_CreateJsonError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	cfg := BuildStubConfig()
	transport := BuildStubHttpTransport()

	body, _ := LoadStubResponseData("stubs/errors/401.json")
	httpmock.RegisterResponder(http.MethodPost, cfg.Uri+"/v1/refund", httpmock.NewBytesResponder(http.StatusOK, body))

	ctx := context.Background()

	request := &RefundRequest{
		Amount:            100,
		InternalReference: "internal_reference",
	}
	resource := &RefundsResource{ResourceAbstract: NewResourceAbstract(transport, cfg)}
	result, resp, err := resource.Create(ctx, request)
	assert.Error(t, err)
	assert.NotEmpty(t, resp)
	assert.NotEmpty(t, result)
	//result
	assert.False(t, result.IsSuccess())
	assert.Equal(t, 401, result.Code)
	assert.Equal(t, "error", result.Status)
	assert.Equal(t, "Unauthorized API access. Unknown Merchant", result.Message)
	assert.Empty(t, result.Data)
	//response
	defer resp.Body.Close()
	bodyRsp, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, body, bodyRsp)
	//error
	assert.Equal(t, "Unauthorized API access. Unknown Merchant", err.Error())
}

func Test_Refunds_RefundsResource_CreateNonJsonError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	cfg := BuildStubConfig()
	transport := BuildStubHttpTransport()

	body, _ := LoadStubResponseData("stubs/errors/500.html")
	httpmock.RegisterResponder(http.MethodPost, cfg.Uri+"/v1/refund", httpmock.NewBytesResponder(http.StatusOK, body))

	ctx := context.Background()

	request := &RefundRequest{
		Amount:            100,
		InternalReference: "internal_reference",
	}
	resource := &RefundsResource{ResourceAbstract: NewResourceAbstract(transport, cfg)}
	result, resp, err := resource.Create(ctx, request)
	assert.Error(t, err)
	assert.NotEmpty(t, resp)
	assert.Empty(t, result)
	//response
	defer resp.Body.Close()
	bodyRsp, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, body, bodyRsp)
}

func Test_Refunds_RefundsResource_CreateInvalidRequest(t *testing.T) {
	config := BuildStubConfig()
	transport := NewHttpTransport(config, nil)
	ctx := context.Background()
	req := &RefundRequest{}
	resource := &RefundsResource{ResourceAbstract: NewResourceAbstract(transport, config)}
	result, rsp, err := resource.Create(ctx, req)
	assert.Nil(t, result)
	assert.Nil(t, rsp)
	assert.Error(t, err)
}
