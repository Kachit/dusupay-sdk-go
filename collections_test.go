package dusupay

import (
	"context"
	"encoding/json"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
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
		Amount:            100,
		Method:            TransactionMethodBank,
		ProviderId:        "provider_id",
		MerchantReference: "merchant_reference",
		RedirectUrl:       "redirect_url",
		Narration:         "narration",
	}
	result := request.isValid()
	assert.Error(t, result)
	assert.Equal(t, `parameter "currency" is empty`, result.Error())
}

func Test_Collections_CollectionRequest_IsValidEmptyAmount(t *testing.T) {
	request := CollectionRequest{
		Currency:          CurrencyCodeKES,
		Method:            TransactionMethodBank,
		ProviderId:        "provider_id",
		MerchantReference: "merchant_reference",
		RedirectUrl:       "redirect_url",
		Narration:         "narration",
	}
	result := request.isValid()
	assert.Error(t, result)
	assert.Equal(t, `parameter "amount" is empty`, result.Error())
}

func Test_Collections_CollectionRequest_IsValidEmptyMethod(t *testing.T) {
	request := CollectionRequest{
		Currency:          CurrencyCodeKES,
		Amount:            100,
		ProviderId:        "provider_id",
		MerchantReference: "merchant_reference",
		RedirectUrl:       "redirect_url",
		Narration:         "narration",
	}
	result := request.isValid()
	assert.Error(t, result)
	assert.Equal(t, `parameter "method" is empty`, result.Error())
}

func Test_Collections_CollectionRequest_IsValidEmptyProviderId(t *testing.T) {
	request := CollectionRequest{
		Currency:          CurrencyCodeKES,
		Amount:            100,
		Method:            TransactionMethodBank,
		MerchantReference: "merchant_reference",
		RedirectUrl:       "redirect_url",
		Narration:         "narration",
	}
	result := request.isValid()
	assert.Error(t, result)
	assert.Equal(t, `parameter "provider_id" is empty`, result.Error())
}

func Test_Collections_CollectionRequest_IsValidEmptyMerchantReference(t *testing.T) {
	request := CollectionRequest{
		Currency:    CurrencyCodeKES,
		Amount:      100,
		Method:      TransactionMethodBank,
		ProviderId:  "provider_id",
		RedirectUrl: "redirect_url",
		Narration:   "narration",
	}
	result := request.isValid()
	assert.Error(t, result)
	assert.Equal(t, `parameter "merchant_reference" is empty`, result.Error())
}

func Test_Collections_CollectionRequest_IsValidEmptyNarration(t *testing.T) {
	request := CollectionRequest{
		Currency:          CurrencyCodeKES,
		Amount:            100,
		Method:            TransactionMethodBank,
		ProviderId:        "provider_id",
		MerchantReference: "merchant_reference",
		RedirectUrl:       "redirect_url",
	}
	result := request.isValid()
	assert.Error(t, result)
	assert.Equal(t, `parameter "narration" is empty`, result.Error())
}

func Test_Collections_CollectionRequest_IsValidEmptyRedirectUrlByDefault(t *testing.T) {
	request := CollectionRequest{
		Currency:          CurrencyCodeKES,
		Amount:            100,
		Method:            TransactionMethodBank,
		ProviderId:        "provider_id",
		MerchantReference: "merchant_reference",
		Narration:         "narration",
	}
	result := request.isValid()
	assert.Error(t, result)
	assert.Equal(t, `parameter "redirect_url" is empty`, result.Error())
}

func Test_Collections_CollectionRequest_IsValidEmptyRedirectUrlMobileMoney(t *testing.T) {
	request := CollectionRequest{
		Currency:          CurrencyCodeKES,
		Amount:            100,
		Method:            TransactionMethodMobileMoney,
		ProviderId:        "provider_id",
		MerchantReference: "merchant_reference",
		Narration:         "narration",
		AccountNumber:     "account_number",
	}
	result := request.isValid()
	assert.Nil(t, result)
	assert.NoError(t, result)
}

func Test_Collections_CollectionRequest_IsValidEmptyRedirectUrlMobileMoneyWithoutAccountNumber(t *testing.T) {
	request := CollectionRequest{
		Currency:          CurrencyCodeKES,
		Amount:            100,
		Method:            TransactionMethodMobileMoney,
		ProviderId:        "provider_id",
		MerchantReference: "merchant_reference",
		Narration:         "narration",
	}
	result := request.isValid()
	assert.Error(t, result)
	assert.Equal(t, `parameter "account_number" is empty`, result.Error())
}

func Test_Collections_CollectionResponse_UnmarshalSuccess(t *testing.T) {
	var response CollectionResponse
	body, _ := LoadStubResponseData("stubs/collections/create/success.json")
	err := json.Unmarshal(body, &response)
	assert.NoError(t, err)
	assert.Equal(t, 202, response.Code)
	assert.Equal(t, "accepted", response.Status)
	assert.Equal(t, "Transaction Initiated", response.Message)
	assert.Equal(t, int64(226), response.Data.ID)
	assert.Equal(t, 0.2, response.Data.RequestAmount)
	assert.Equal(t, "USD", response.Data.RequestCurrency)
	assert.Equal(t, 737.9934, response.Data.AccountAmount)
	assert.Equal(t, "UGX", response.Data.AccountCurrency)
	assert.Equal(t, 21.4018, response.Data.TransactionFee)
	assert.Equal(t, 716.5916, response.Data.TotalCredit)
	assert.Equal(t, "mtn_ug", response.Data.ProviderID)
	assert.Equal(t, "76859aae-f148-48c5-9901-2e474cf19b71", response.Data.MerchantReference)
	assert.Equal(t, "DUSUPAY405GZM1G5JXGA71IK", response.Data.InternalReference)
	assert.Equal(t, "PENDING", response.Data.TransactionStatus)
	assert.Equal(t, "collection", response.Data.TransactionType)
	assert.Equal(t, "Transaction Initiated", response.Data.Message)
	assert.Equal(t, false, response.Data.CustomerCharged)
	assert.Equal(t, "https://sandbox.dusupay.com/v1/complete-payment/DUSUPAY405GZM1G5JXGA71IK", response.Data.PaymentURL)
	assert.Equal(t, "Ensure that you have sufficient balance on your MTN Mobile Money account", response.Data.Instructions[0].Description)
	assert.Equal(t, "1", response.Data.Instructions[0].StepNo)
}

func Test_Collections_CollectionResponse_UnmarshalErrorUnauthorized(t *testing.T) {
	var response CollectionResponse
	body, _ := LoadStubResponseData("stubs/errors/401.json")
	err := json.Unmarshal(body, &response)
	assert.NoError(t, err)
	assert.Equal(t, 401, response.Code)
	assert.Equal(t, "error", response.Status)
	assert.Equal(t, "Unauthorized API access. Unknown Merchant", response.Message)
	assert.Empty(t, response.Data)
}

func Test_Collections_CollectionsResource_CreateSuccess(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	cfg := BuildStubConfig()
	transport := BuildStubHttpTransport()

	body, _ := LoadStubResponseData("stubs/collections/create/success.json")
	httpmock.RegisterResponder(http.MethodPost, cfg.Uri+"/v1/collections", httpmock.NewBytesResponder(http.StatusOK, body))

	ctx := context.Background()

	request := &CollectionRequest{
		Currency:          CurrencyCodeKES,
		Amount:            100,
		Method:            TransactionMethodBank,
		ProviderId:        "provider_id",
		MerchantReference: "merchant_reference",
		RedirectUrl:       "redirect_url",
		Narration:         "narration",
	}
	resource := &CollectionsResource{ResourceAbstract: NewResourceAbstract(transport, cfg)}
	result, resp, err := resource.Create(ctx, request)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp)
	assert.NotEmpty(t, result)
	//result
	assert.True(t, result.IsSuccess())
	assert.Equal(t, 202, result.Code)
	assert.Equal(t, "accepted", result.Status)
	assert.Equal(t, "Transaction Initiated", result.Message)
	assert.Equal(t, int64(226), result.Data.ID)
	assert.Equal(t, 0.2, result.Data.RequestAmount)
	assert.Equal(t, "USD", result.Data.RequestCurrency)
	assert.Equal(t, 737.9934, result.Data.AccountAmount)
	assert.Equal(t, "UGX", result.Data.AccountCurrency)
	assert.Equal(t, 21.4018, result.Data.TransactionFee)
	assert.Equal(t, 716.5916, result.Data.TotalCredit)
	assert.Equal(t, "mtn_ug", result.Data.ProviderID)
	assert.Equal(t, "76859aae-f148-48c5-9901-2e474cf19b71", result.Data.MerchantReference)
	assert.Equal(t, "DUSUPAY405GZM1G5JXGA71IK", result.Data.InternalReference)
	assert.Equal(t, "PENDING", result.Data.TransactionStatus)
	assert.Equal(t, "collection", result.Data.TransactionType)
	assert.Equal(t, "Transaction Initiated", result.Data.Message)
	assert.Equal(t, false, result.Data.CustomerCharged)
	assert.Equal(t, "https://sandbox.dusupay.com/v1/complete-payment/DUSUPAY405GZM1G5JXGA71IK", result.Data.PaymentURL)
	assert.Equal(t, "Ensure that you have sufficient balance on your MTN Mobile Money account", result.Data.Instructions[0].Description)
	assert.Equal(t, "1", result.Data.Instructions[0].StepNo)
	//response
	defer resp.Body.Close()
	bodyRsp, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, body, bodyRsp)
}

func Test_Collections_CollectionsResource_CreateJsonError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	cfg := BuildStubConfig()
	transport := BuildStubHttpTransport()

	body, _ := LoadStubResponseData("stubs/errors/401.json")
	httpmock.RegisterResponder(http.MethodPost, cfg.Uri+"/v1/collections", httpmock.NewBytesResponder(http.StatusOK, body))

	ctx := context.Background()

	request := &CollectionRequest{
		Currency:          CurrencyCodeKES,
		Amount:            100,
		Method:            TransactionMethodBank,
		ProviderId:        "provider_id",
		MerchantReference: "merchant_reference",
		RedirectUrl:       "redirect_url",
		Narration:         "narration",
	}
	resource := &CollectionsResource{ResourceAbstract: NewResourceAbstract(transport, cfg)}
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

func Test_Collections_CollectionsResource_CreateNonJsonError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	cfg := BuildStubConfig()
	transport := BuildStubHttpTransport()

	body, _ := LoadStubResponseData("stubs/errors/500.html")
	httpmock.RegisterResponder(http.MethodPost, cfg.Uri+"/v1/collections", httpmock.NewBytesResponder(http.StatusOK, body))

	ctx := context.Background()

	request := &CollectionRequest{
		Currency:          CurrencyCodeKES,
		Amount:            100,
		Method:            TransactionMethodBank,
		ProviderId:        "provider_id",
		MerchantReference: "merchant_reference",
		RedirectUrl:       "redirect_url",
		Narration:         "narration",
	}
	resource := &CollectionsResource{ResourceAbstract: NewResourceAbstract(transport, cfg)}
	result, resp, err := resource.Create(ctx, request)
	assert.Error(t, err)
	assert.NotEmpty(t, resp)
	assert.Empty(t, result)
	//response
	defer resp.Body.Close()
	bodyRsp, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, body, bodyRsp)
}

func Test_Collections_CollectionsResource_CreateInvalidRequest(t *testing.T) {
	config := BuildStubConfig()
	transport := NewHttpTransport(config, nil)
	ctx := context.Background()
	req := &CollectionRequest{}
	resource := &CollectionsResource{ResourceAbstract: NewResourceAbstract(transport, config)}
	result, rsp, err := resource.Create(ctx, req)
	assert.Nil(t, rsp)
	assert.Nil(t, result)
	assert.Error(t, err)
}
