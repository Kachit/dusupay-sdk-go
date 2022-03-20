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

func Test_Merchants_BalancesResponse_UnmarshalSuccess(t *testing.T) {
	var response BalancesResponse
	body, _ := LoadStubResponseData("stubs/merchants/balance/success.json")
	err := json.Unmarshal(body, &response)
	assert.NoError(t, err)
	assert.True(t, response.IsSuccess())
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
	assert.False(t, response.IsSuccess())
	assert.Equal(t, 401, response.Code)
	assert.Equal(t, "error", response.Status)
	assert.Equal(t, "Unauthorized API access. Unknown Merchant", response.Message)
	assert.Empty(t, response.Data)
}

func Test_Merchants_MerchantsResource_GetBalancesSuccess(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	cfg := BuildStubConfig()
	transport := BuildStubHttpTransport()

	body, _ := LoadStubResponseData("stubs/merchants/balance/success.json")
	httpmock.RegisterResponder(http.MethodGet, cfg.Uri+"/v1/merchants/balance", httpmock.NewBytesResponder(http.StatusOK, body))

	ctx := context.Background()

	resource := &MerchantsResource{ResourceAbstract: NewResourceAbstract(transport, cfg)}
	result, resp, err := resource.GetBalances(ctx)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp)
	assert.NotEmpty(t, result)
	//result
	assert.True(t, result.IsSuccess())
	assert.Equal(t, 200, result.Code)
	assert.Equal(t, "success", result.Status)
	assert.Equal(t, "Request completed successfully.", result.Message)
	assert.Equal(t, "UGX", (*result.Data)[0].Currency)
	assert.Equal(t, 5475.816, (*result.Data)[0].Balance)
	assert.Equal(t, "USD", (*result.Data)[1].Currency)
	assert.Equal(t, float64(12), (*result.Data)[1].Balance)
	//response
	defer resp.Body.Close()
	bodyRsp, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, body, bodyRsp)
}

func Test_Merchants_MerchantsResource_GetBalancesJsonError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	cfg := BuildStubConfig()
	transport := BuildStubHttpTransport()

	body, _ := LoadStubResponseData("stubs/errors/401.json")
	httpmock.RegisterResponder(http.MethodGet, cfg.Uri+"/v1/merchants/balance", httpmock.NewBytesResponder(http.StatusOK, body))

	ctx := context.Background()

	resource := &MerchantsResource{ResourceAbstract: NewResourceAbstract(transport, cfg)}
	result, resp, err := resource.GetBalances(ctx)
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

func Test_Merchants_MerchantsResource_GetBalancesNonJsonError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	cfg := BuildStubConfig()
	transport := BuildStubHttpTransport()

	body, _ := LoadStubResponseData("stubs/errors/500.html")
	httpmock.RegisterResponder(http.MethodGet, cfg.Uri+"/v1/merchants/balance", httpmock.NewBytesResponder(http.StatusOK, body))

	ctx := context.Background()

	resource := &MerchantsResource{ResourceAbstract: NewResourceAbstract(transport, cfg)}
	result, resp, err := resource.GetBalances(ctx)
	assert.Error(t, err)
	assert.NotEmpty(t, resp)
	assert.Empty(t, result)
	//response
	defer resp.Body.Close()
	bodyRsp, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, body, bodyRsp)
}
