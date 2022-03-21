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

func Test_Providers_ProvidersResponse_UnmarshalSuccessSandbox(t *testing.T) {
	var response ProvidersResponse
	body, _ := LoadStubResponseData("stubs/providers/payment-options/success-sandbox.json")
	err := json.Unmarshal(body, &response)
	assert.NoError(t, err)
	assert.Equal(t, 200, response.Code)
	assert.Equal(t, "success", response.Status)
	assert.Equal(t, "Request completed successfully.", response.Message)
	assert.Equal(t, "mtn_ug", (*response.Data)[0].ID)
	assert.Equal(t, "MTN Mobile Money", (*response.Data)[0].Name)
	assert.Equal(t, "UGX", (*response.Data)[0].TransactionCurrency)
	assert.Equal(t, float64(3000), (*response.Data)[0].MinAmount)
	assert.Equal(t, float64(5000000), (*response.Data)[0].MaxAmount)
	assert.Equal(t, true, (*response.Data)[0].Available)
	assert.Equal(t, "256777000456", (*response.Data)[0].SandboxTestAccounts.Failure)
	assert.Equal(t, "256777000123", (*response.Data)[0].SandboxTestAccounts.Success)
}

func Test_Providers_ProvidersResponse_UnmarshalSuccess(t *testing.T) {
	var response ProvidersResponse
	body, _ := LoadStubResponseData("stubs/providers/payment-options/success.json")
	err := json.Unmarshal(body, &response)
	assert.NoError(t, err)
	assert.Equal(t, 200, response.Code)
	assert.Equal(t, "success", response.Status)
	assert.Equal(t, "Request completed successfully.", response.Message)
	assert.Equal(t, "mtn_ug", (*response.Data)[0].ID)
	assert.Equal(t, "MTN Mobile Money", (*response.Data)[0].Name)
	assert.Equal(t, "UGX", (*response.Data)[0].TransactionCurrency)
	assert.Equal(t, float64(3000), (*response.Data)[0].MinAmount)
	assert.Equal(t, float64(5000000), (*response.Data)[0].MaxAmount)
	assert.Equal(t, true, (*response.Data)[0].Available)
	assert.Empty(t, (*response.Data)[0].SandboxTestAccounts)
}

func Test_Providers_ProvidersResponse_UnmarshalErrorUnauthorized(t *testing.T) {
	var response ProvidersResponse
	body, _ := LoadStubResponseData("stubs/errors/401.json")
	err := json.Unmarshal(body, &response)
	assert.NoError(t, err)
	assert.Equal(t, 401, response.Code)
	assert.Equal(t, "error", response.Status)
	assert.Equal(t, "Unauthorized API access. Unknown Merchant", response.Message)
	assert.Empty(t, response.Data)
}

func Test_Providers_ProvidersResource_GetListSuccess(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	cfg := BuildStubConfig()
	transport := BuildStubHttpTransport()

	body, _ := LoadStubResponseData("stubs/providers/payment-options/success.json")
	httpmock.RegisterResponder(http.MethodGet, cfg.Uri+"/v1/payment-options/COLLECTION/CARD/KE", httpmock.NewBytesResponder(http.StatusOK, body))

	ctx := context.Background()

	filter := &ProvidersFilter{Country: CountryCodeKenya, Method: TransactionMethodCard, TransactionType: TransactionTypeCollection}
	resource := &ProvidersResource{ResourceAbstract: NewResourceAbstract(transport, cfg)}
	result, resp, err := resource.GetList(ctx, filter)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp)
	assert.NotEmpty(t, result)
	//result
	assert.True(t, result.IsSuccess())
	assert.Equal(t, 200, result.Code)
	assert.Equal(t, "success", result.Status)
	assert.Equal(t, "Request completed successfully.", result.Message)
	assert.Equal(t, "mtn_ug", (*result.Data)[0].ID)
	assert.Equal(t, "MTN Mobile Money", (*result.Data)[0].Name)
	assert.Equal(t, "UGX", (*result.Data)[0].TransactionCurrency)
	assert.Equal(t, float64(3000), (*result.Data)[0].MinAmount)
	assert.Equal(t, float64(5000000), (*result.Data)[0].MaxAmount)
	assert.Equal(t, true, (*result.Data)[0].Available)
	assert.Empty(t, (*result.Data)[0].SandboxTestAccounts)
	//response
	defer resp.Body.Close()
	bodyRsp, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, body, bodyRsp)
}

func Test_Providers_ProvidersResource_GetListJsonError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	cfg := BuildStubConfig()
	transport := BuildStubHttpTransport()

	body, _ := LoadStubResponseData("stubs/errors/401.json")
	httpmock.RegisterResponder(http.MethodGet, cfg.Uri+"/v1/payment-options/COLLECTION/CARD/KE", httpmock.NewBytesResponder(http.StatusOK, body))

	ctx := context.Background()

	filter := &ProvidersFilter{Country: CountryCodeKenya, Method: TransactionMethodCard, TransactionType: TransactionTypeCollection}
	resource := &ProvidersResource{ResourceAbstract: NewResourceAbstract(transport, cfg)}
	result, resp, err := resource.GetList(ctx, filter)
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

func Test_Providers_ProvidersResource_GetListNonError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	cfg := BuildStubConfig()
	transport := BuildStubHttpTransport()

	body, _ := LoadStubResponseData("stubs/errors/500.html")
	httpmock.RegisterResponder(http.MethodGet, cfg.Uri+"/v1/payment-options/COLLECTION/CARD/KE", httpmock.NewBytesResponder(http.StatusOK, body))

	ctx := context.Background()

	filter := &ProvidersFilter{Country: CountryCodeKenya, Method: TransactionMethodCard, TransactionType: TransactionTypeCollection}
	resource := &ProvidersResource{ResourceAbstract: NewResourceAbstract(transport, cfg)}
	result, resp, err := resource.GetList(ctx, filter)
	assert.Error(t, err)
	assert.NotEmpty(t, resp)
	assert.Empty(t, result)
	//response
	defer resp.Body.Close()
	bodyRsp, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, body, bodyRsp)
}

func Test_Providers_ProvidersResource_GetListInvalidFilter(t *testing.T) {
	config := BuildStubConfig()
	transport := NewHttpTransport(config, nil)
	ctx := context.Background()
	filter := &ProvidersFilter{}
	resource := &ProvidersResource{ResourceAbstract: NewResourceAbstract(transport, config)}
	result, rsp, err := resource.GetList(ctx, filter)
	assert.Nil(t, rsp)
	assert.Nil(t, result)
	assert.Error(t, err)
}
