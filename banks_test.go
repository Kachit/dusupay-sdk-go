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

func Test_Banks_BanksFilter_IsValidSuccess(t *testing.T) {
	filter := BanksFilter{Country: CountryCodeKenya, Method: TransactionMethodCard}
	assert.Nil(t, filter.isValid())
	assert.NoError(t, filter.isValid())
}

func Test_Banks_BanksFilter_IsValidEmptyCountryCode(t *testing.T) {
	filter := BanksFilter{Method: "qwerty"}
	result := filter.isValid()
	assert.Error(t, result)
	assert.Equal(t, `parameter "country_code" is empty`, result.Error())
}

func Test_Banks_BanksFilter_IsValidEmptyMethod(t *testing.T) {
	filter := BanksFilter{Country: CountryCodeKenya}
	result := filter.isValid()
	assert.Error(t, result)
	assert.Equal(t, `parameter "method" is empty`, result.Error())
}

func Test_Banks_BanksFilter_BuildPath(t *testing.T) {
	filter := BanksFilter{Country: CountryCodeKenya, Method: TransactionMethodCard}
	result := filter.buildPath()
	assert.Equal(t, `CARD/bank/KE`, result)
}

func Test_Banks_BanksBranchesFilter_IsValidSuccess(t *testing.T) {
	filter := BanksBranchesFilter{Country: CountryCodeKenya, Bank: "qwerty"}
	assert.Nil(t, filter.isValid())
	assert.NoError(t, filter.isValid())
}

func Test_Banks_BanksBranchesFilter_IsValidEmptyCountryCode(t *testing.T) {
	filter := BanksBranchesFilter{Bank: "qwerty"}
	result := filter.isValid()
	assert.Error(t, result)
	assert.Equal(t, `parameter "country_code" is empty`, result.Error())
}

func Test_Banks_BanksBranchesFilter_IsValidEmptyBankCode(t *testing.T) {
	filter := BanksBranchesFilter{Country: CountryCodeKenya}
	result := filter.isValid()
	assert.Error(t, result)
	assert.Equal(t, `parameter "bank_code" is empty`, result.Error())
}

func Test_Banks_BanksBranchesFilter_BuildPath(t *testing.T) {
	filter := BanksBranchesFilter{Country: CountryCodeKenya, Bank: "qwerty"}
	result := filter.buildPath()
	assert.Equal(t, `KE/branches/qwerty`, result)
}

func Test_Banks_BanksResource_GetListInvalidFilter(t *testing.T) {
	config := BuildStubConfig()
	transport := NewHttpTransport(config, nil)
	ctx := context.Background()
	filter := &BanksFilter{}
	resource := &BanksResource{ResourceAbstract: NewResourceAbstract(transport, config)}
	result, rsp, err := resource.GetList(ctx, filter)
	assert.Nil(t, result)
	assert.Nil(t, rsp)
	assert.Error(t, err)
}

func Test_Banks_BanksResource_GetListSuccess(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	cfg := BuildStubConfig()
	transport := BuildStubHttpTransport()

	body, _ := LoadStubResponseData("stubs/banks/list/success.json")
	httpmock.RegisterResponder(http.MethodGet, cfg.Uri+"/v1/payment-options/CARD/bank/KE", httpmock.NewBytesResponder(http.StatusOK, body))

	ctx := context.Background()

	filter := &BanksFilter{Country: CountryCodeKenya, Method: TransactionMethodCard}
	resource := &BanksResource{ResourceAbstract: NewResourceAbstract(transport, cfg)}
	result, resp, err := resource.GetList(ctx, filter)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp)
	assert.NotEmpty(t, result)
	//result
	assert.True(t, result.IsSuccess())
	assert.Equal(t, 200, result.Code)
	assert.Equal(t, "success", result.Status)
	assert.Equal(t, "Request completed successfully.", result.Message)
	assert.Equal(t, "GH030243", (*result.Data)[0].Code)
	assert.Equal(t, "BARCLAYS BANK(GH) LTD-NKAWKAW", (*result.Data)[0].Name)
	//response
	defer resp.Body.Close()
	bodyRsp, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, body, bodyRsp)
}

func Test_Banks_BanksResource_GetListError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	cfg := BuildStubConfig()
	transport := BuildStubHttpTransport()

	body, _ := LoadStubResponseData("stubs/errors/401.json")
	httpmock.RegisterResponder(http.MethodGet, cfg.Uri+"/v1/payment-options/CARD/bank/KE", httpmock.NewBytesResponder(http.StatusOK, body))

	ctx := context.Background()

	filter := &BanksFilter{Country: CountryCodeKenya, Method: TransactionMethodCard}
	resource := &BanksResource{ResourceAbstract: NewResourceAbstract(transport, cfg)}
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

func Test_Banks_BanksResource_GetBranchesListInvalidFilter(t *testing.T) {
	config := BuildStubConfig()
	transport := NewHttpTransport(config, nil)
	ctx := context.Background()
	filter := &BanksBranchesFilter{}
	resource := &BanksResource{ResourceAbstract: NewResourceAbstract(transport, config)}
	result, rsp, err := resource.GetBranchesList(ctx, filter)
	assert.Nil(t, result)
	assert.Nil(t, rsp)
	assert.Error(t, err)
}

func Test_Banks_BanksResource_GetBranchesListSuccess(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	cfg := BuildStubConfig()
	transport := BuildStubHttpTransport()

	body, _ := LoadStubResponseData("stubs/banks/branches/success.json")
	httpmock.RegisterResponder(http.MethodGet, cfg.Uri+"/v1/bank/KE/branches/qwerty", httpmock.NewBytesResponder(http.StatusOK, body))

	ctx := context.Background()

	filter := &BanksBranchesFilter{Country: CountryCodeKenya, Bank: "qwerty"}
	resource := &BanksResource{ResourceAbstract: NewResourceAbstract(transport, cfg)}
	result, resp, err := resource.GetBranchesList(ctx, filter)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp)
	assert.NotEmpty(t, result)
	//result
	assert.True(t, result.IsSuccess())
	assert.Equal(t, 200, result.Code)
	assert.Equal(t, "success", result.Status)
	assert.Equal(t, "Request completed successfully.", result.Message)
	assert.Equal(t, "GH030243", (*result.Data)[0].Code)
	assert.Equal(t, "BARCLAYS BANK(GH) LTD-NKAWKAW", (*result.Data)[0].Name)
	//response
	defer resp.Body.Close()
	bodyRsp, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, body, bodyRsp)
}

func Test_Banks_BanksResource_GetBranchesListError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	cfg := BuildStubConfig()
	transport := BuildStubHttpTransport()

	body, _ := LoadStubResponseData("stubs/errors/401.json")
	httpmock.RegisterResponder(http.MethodGet, cfg.Uri+"/v1/bank/KE/branches/qwerty", httpmock.NewBytesResponder(http.StatusOK, body))

	ctx := context.Background()

	filter := &BanksBranchesFilter{Country: CountryCodeKenya, Bank: "qwerty"}
	resource := &BanksResource{ResourceAbstract: NewResourceAbstract(transport, cfg)}
	result, resp, err := resource.GetBranchesList(ctx, filter)
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

func Test_Banks_BanksResponse_UnmarshalSuccess(t *testing.T) {
	var response BanksResponse
	body, _ := LoadStubResponseData("stubs/banks/list/success.json")
	err := json.Unmarshal(body, &response)
	assert.NoError(t, err)
	assert.Equal(t, 200, response.Code)
	assert.Equal(t, "success", response.Status)
	assert.Equal(t, "Request completed successfully.", response.Message)
	assert.Equal(t, "GH030243", (*response.Data)[0].Code)
	assert.Equal(t, "BARCLAYS BANK(GH) LTD-NKAWKAW", (*response.Data)[0].Name)
}

func Test_Banks_BanksResponse_UnmarshalErrorUnauthorized(t *testing.T) {
	var response BanksResponse
	body, _ := LoadStubResponseData("stubs/errors/401.json")
	err := json.Unmarshal(body, &response)
	assert.NoError(t, err)
	assert.Equal(t, 401, response.Code)
	assert.Equal(t, "error", response.Status)
	assert.Equal(t, "Unauthorized API access. Unknown Merchant", response.Message)
	assert.Empty(t, response.Data)
}

func Test_Banks_BanksBranchesResponse_UnmarshalSuccess(t *testing.T) {
	var response BanksBranchesResponse
	body, _ := LoadStubResponseData("stubs/banks/branches/success.json")
	err := json.Unmarshal(body, &response)
	assert.NoError(t, err)
	assert.Equal(t, 200, response.Code)
	assert.Equal(t, "success", response.Status)
	assert.Equal(t, "Request completed successfully.", response.Message)
	assert.Equal(t, "GH030243", (*response.Data)[0].Code)
	assert.Equal(t, "BARCLAYS BANK(GH) LTD-NKAWKAW", (*response.Data)[0].Name)
}

func Test_Banks_BanksBranchesResponse_UnmarshalErrorUnauthorized(t *testing.T) {
	var response BanksBranchesResponse
	body, _ := LoadStubResponseData("stubs/errors/401.json")
	err := json.Unmarshal(body, &response)
	assert.NoError(t, err)
	assert.Equal(t, 401, response.Code)
	assert.Equal(t, "error", response.Status)
	assert.Equal(t, "Unauthorized API access. Unknown Merchant", response.Message)
	assert.Empty(t, response.Data)
}
