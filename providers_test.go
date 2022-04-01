package dusupay

import (
	"context"
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
	assert.Equal(t, `collection/card/ke`, result)
}

func Test_Providers_ProvidersResource_GetListSuccess(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	cfg := BuildStubConfig()
	transport := BuildStubHttpTransport()

	body, _ := LoadStubResponseData("stubs/providers/payment-options/success.json")
	httpmock.RegisterResponder(http.MethodGet, cfg.Uri+"/v1/payment-options/collection/card/ke", httpmock.NewBytesResponder(http.StatusOK, body))

	ctx := context.Background()

	filter := &ProvidersFilter{Country: CountryCodeKenya, Method: TransactionMethodCard, TransactionType: TransactionTypeCollection}
	resource := &ProvidersResource{ResourceAbstract: NewResourceAbstract(transport, cfg)}
	result, resp, err := resource.GetList(ctx, filter)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp)
	assert.NotEmpty(t, result)
	//result
	assert.True(t, result.IsSuccess())
	assert.Equal(t, http.StatusOK, result.Code)
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

func Test_Providers_ProvidersResource_GetListSuccessSandbox(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	cfg := BuildStubConfig()
	transport := BuildStubHttpTransport()

	body, _ := LoadStubResponseData("stubs/providers/payment-options/success-sandbox.json")
	httpmock.RegisterResponder(http.MethodGet, cfg.Uri+"/v1/payment-options/collection/card/ke", httpmock.NewBytesResponder(http.StatusOK, body))

	ctx := context.Background()

	filter := &ProvidersFilter{Country: CountryCodeKenya, Method: TransactionMethodCard, TransactionType: TransactionTypeCollection}
	resource := &ProvidersResource{ResourceAbstract: NewResourceAbstract(transport, cfg)}
	result, resp, err := resource.GetList(ctx, filter)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp)
	assert.NotEmpty(t, result)
	//result
	assert.True(t, result.IsSuccess())
	assert.Equal(t, http.StatusOK, result.Code)
	assert.Equal(t, "success", result.Status)
	assert.Equal(t, "Request completed successfully.", result.Message)
	assert.Equal(t, "mtn_ug", (*result.Data)[0].ID)
	assert.Equal(t, "MTN Mobile Money", (*result.Data)[0].Name)
	assert.Equal(t, "UGX", (*result.Data)[0].TransactionCurrency)
	assert.Equal(t, float64(3000), (*result.Data)[0].MinAmount)
	assert.Equal(t, float64(5000000), (*result.Data)[0].MaxAmount)
	assert.Equal(t, true, (*result.Data)[0].Available)
	assert.Equal(t, "256777000456", (*result.Data)[0].SandboxTestAccounts.Failure)
	assert.Equal(t, "256777000123", (*result.Data)[0].SandboxTestAccounts.Success)
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
	httpmock.RegisterResponder(http.MethodGet, cfg.Uri+"/v1/payment-options/collection/card/ke", httpmock.NewBytesResponder(http.StatusOK, body))

	ctx := context.Background()

	filter := &ProvidersFilter{Country: CountryCodeKenya, Method: TransactionMethodCard, TransactionType: TransactionTypeCollection}
	resource := &ProvidersResource{ResourceAbstract: NewResourceAbstract(transport, cfg)}
	result, resp, err := resource.GetList(ctx, filter)
	assert.Error(t, err)
	assert.NotEmpty(t, resp)
	assert.NotEmpty(t, result)
	//result
	assert.False(t, result.IsSuccess())
	assert.Equal(t, http.StatusUnauthorized, result.Code)
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
	httpmock.RegisterResponder(http.MethodGet, cfg.Uri+"/v1/payment-options/collection/card/ke", httpmock.NewBytesResponder(http.StatusOK, body))

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
