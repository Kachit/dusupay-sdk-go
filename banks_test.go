package dusupay

import (
	"context"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

func Test_Banks_BanksFilter_IsValidSuccess(t *testing.T) {
	filter := BanksFilter{Country: CountryCodeKenya, TransactionType: TransactionTypeCollection}
	assert.Nil(t, filter.isValid())
	assert.NoError(t, filter.isValid())
}

func Test_Banks_BanksFilter_IsValidEmptyCountryCode(t *testing.T) {
	filter := BanksFilter{TransactionType: "qwerty"}
	result := filter.isValid()
	assert.Error(t, result)
	assert.Equal(t, `parameter "country_code" is empty`, result.Error())
}

func Test_Banks_BanksFilter_IsValidEmptyMethod(t *testing.T) {
	filter := BanksFilter{Country: CountryCodeKenya}
	result := filter.isValid()
	assert.Error(t, result)
	assert.Equal(t, `parameter "transaction_type" is empty`, result.Error())
}

func Test_Banks_BanksFilter_BuildPath(t *testing.T) {
	filter := BanksFilter{Country: CountryCodeKenya, TransactionType: TransactionTypeCollection}
	result := filter.buildPath()
	assert.Equal(t, `collection/bank/ke`, result)
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
	assert.Equal(t, `ke/branches/qwerty`, result)
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
	httpmock.RegisterResponder(http.MethodGet, cfg.Uri+"/v1/payment-options/payout/bank/ke", httpmock.NewBytesResponder(http.StatusOK, body))

	ctx := context.Background()

	filter := &BanksFilter{Country: CountryCodeKenya, TransactionType: TransactionTypePayout}
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
	assert.Equal(t, "access_bank", (*result.Data)[0].BankCode)
	assert.Equal(t, "Access Bank", (*result.Data)[0].Name)
	assert.Equal(t, "NGN", (*result.Data)[0].TransactionCurrency)
	assert.Equal(t, float64(1000), (*result.Data)[0].MinAmount)
	assert.Equal(t, float64(380000), (*result.Data)[0].MaxAmount)
	assert.Equal(t, true, (*result.Data)[0].Available)
	assert.Empty(t, (*result.Data)[0].SandboxTestAccounts)
	//response
	defer resp.Body.Close()
	bodyRsp, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, body, bodyRsp)
}

func Test_Banks_BanksResource_GetListSuccessSandbox(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	cfg := BuildStubConfig()
	transport := BuildStubHttpTransport()

	body, _ := LoadStubResponseData("stubs/banks/list/success-sandbox.json")
	httpmock.RegisterResponder(http.MethodGet, cfg.Uri+"/v1/payment-options/payout/bank/ke", httpmock.NewBytesResponder(http.StatusOK, body))

	ctx := context.Background()

	filter := &BanksFilter{Country: CountryCodeKenya, TransactionType: TransactionTypePayout}
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
	assert.Equal(t, "access_bank", (*result.Data)[0].BankCode)
	assert.Equal(t, "Access Bank", (*result.Data)[0].Name)
	assert.Equal(t, "NGN", (*result.Data)[0].TransactionCurrency)
	assert.Equal(t, float64(1000), (*result.Data)[0].MinAmount)
	assert.Equal(t, float64(380000), (*result.Data)[0].MaxAmount)
	assert.Equal(t, true, (*result.Data)[0].Available)
	assert.Equal(t, "256777000456", (*result.Data)[0].SandboxTestAccounts.Failure)
	assert.Equal(t, "256777000123", (*result.Data)[0].SandboxTestAccounts.Success)
	//response
	defer resp.Body.Close()
	bodyRsp, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, body, bodyRsp)
}

func Test_Banks_BanksResource_GetListJsonError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	cfg := BuildStubConfig()
	transport := BuildStubHttpTransport()

	body, _ := LoadStubResponseData("stubs/errors/401.json")
	httpmock.RegisterResponder(http.MethodGet, cfg.Uri+"/v1/payment-options/payout/bank/ke", httpmock.NewBytesResponder(http.StatusOK, body))

	ctx := context.Background()

	filter := &BanksFilter{Country: CountryCodeKenya, TransactionType: TransactionTypePayout}
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

func Test_Banks_BanksResource_GetListNonJsonError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	cfg := BuildStubConfig()
	transport := BuildStubHttpTransport()

	body, _ := LoadStubResponseData("stubs/errors/500.html")
	httpmock.RegisterResponder(http.MethodGet, cfg.Uri+"/v1/payment-options/payout/bank/ke", httpmock.NewBytesResponder(http.StatusOK, body))

	ctx := context.Background()

	filter := &BanksFilter{Country: CountryCodeKenya, TransactionType: TransactionTypePayout}
	resource := &BanksResource{ResourceAbstract: NewResourceAbstract(transport, cfg)}
	result, resp, err := resource.GetList(ctx, filter)
	assert.Error(t, err)
	assert.NotEmpty(t, resp)
	assert.Empty(t, result)
	//response
	defer resp.Body.Close()
	bodyRsp, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, body, bodyRsp)
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
	httpmock.RegisterResponder(http.MethodGet, cfg.Uri+"/v1/bank/ke/branches/qwerty", httpmock.NewBytesResponder(http.StatusOK, body))

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

func Test_Banks_BanksResource_GetBranchesListJsonError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	cfg := BuildStubConfig()
	transport := BuildStubHttpTransport()

	body, _ := LoadStubResponseData("stubs/errors/401.json")
	httpmock.RegisterResponder(http.MethodGet, cfg.Uri+"/v1/bank/ke/branches/qwerty", httpmock.NewBytesResponder(http.StatusOK, body))

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

func Test_Banks_BanksResource_GetBranchesListNonJsonError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	cfg := BuildStubConfig()
	transport := BuildStubHttpTransport()

	body, _ := LoadStubResponseData("stubs/errors/500.html")
	httpmock.RegisterResponder(http.MethodGet, cfg.Uri+"/v1/bank/ke/branches/qwerty", httpmock.NewBytesResponder(http.StatusOK, body))

	ctx := context.Background()

	filter := &BanksBranchesFilter{Country: CountryCodeKenya, Bank: "qwerty"}
	resource := &BanksResource{ResourceAbstract: NewResourceAbstract(transport, cfg)}
	result, resp, err := resource.GetBranchesList(ctx, filter)
	assert.Error(t, err)
	assert.NotEmpty(t, resp)
	assert.Empty(t, result)
	//response
	defer resp.Body.Close()
	bodyRsp, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, body, bodyRsp)
}
