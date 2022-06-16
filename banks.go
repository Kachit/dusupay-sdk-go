package dusupay

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

//BanksFilter (see https://docs.dusupay.com/sending-money/payouts/bank-codes)
type BanksFilter struct {
	TransactionType TransactionTypeCode `json:"transaction_type"`
	//ISO-2 country code from those supported
	Country CountryCode `json:"country_code"`
}

//isValid check is valid BanksFilter parameters
func (bf *BanksFilter) isValid() error {
	var err error
	if bf.Country == "" {
		err = fmt.Errorf(`parameter "country_code" is empty`)
	} else if bf.TransactionType == "" {
		err = fmt.Errorf(`parameter "transaction_type" is empty`)
	}
	return err
}

//buildPath
func (bf *BanksFilter) buildPath() string {
	return strings.ToLower(string(bf.TransactionType) + "/bank/" + string(bf.Country))
}

//BanksBranchesFilter branches list filter (see https://docs.dusupay.com/sending-money/payouts/bank-branches)
type BanksBranchesFilter struct {
	//ISO-2 country code from those supported
	Country CountryCode `json:"country_code"`
	Bank    string      `json:"bank_code"`
}

//isValid check is valid BanksBranchesFilter parameters
func (bbf *BanksBranchesFilter) isValid() error {
	var err error
	if bbf.Country == "" {
		err = fmt.Errorf(`parameter "country_code" is empty`)
	} else if bbf.Bank == "" {
		err = fmt.Errorf(`parameter "bank_code" is empty`)
	}
	return err
}

//buildPath
func (bbf *BanksBranchesFilter) buildPath() string {
	return strings.ToLower(string(bbf.Country) + "/branches/" + string(bbf.Bank))
}

//BanksResponse struct
type BanksResponse struct {
	ResponseBody
	Data *BanksResponseData `json:"data,omitempty"`
}

//BanksResponseData struct
type BanksResponseData []*BanksResponseDataItem

//UnmarshalJSON unmarshal json data
func (rsp *BanksResponseData) UnmarshalJSON(data []byte) error {
	if isEmptyObjectResponseData(data) {
		return nil
	}
	var arr []*BanksResponseDataItem
	err := json.Unmarshal(data, &arr)
	if err != nil {
		return err
	}
	*rsp = append(*rsp, arr...)
	return nil
}

//BanksResponseDataItem struct
type BanksResponseDataItem struct {
	Id                  string  `json:"id"`
	Name                string  `json:"name"`
	TransactionCurrency string  `json:"transaction_currency"`
	MinAmount           float64 `json:"min_amount"`
	MaxAmount           float64 `json:"max_amount"`
	BankCode            string  `json:"bank_code"`
	Available           bool    `json:"available"`
	SandboxTestAccounts struct {
		Success string `json:"success"`
		Failure string `json:"failure"`
	} `json:"sandbox_test_accounts,omitempty"`
}

//BanksBranchesResponse struct
type BanksBranchesResponse struct {
	ResponseBody
	Data *BanksBranchesResponseData `json:"data,omitempty"`
}

//BanksBranchesResponseData struct
type BanksBranchesResponseData []*BanksBranchesResponseDataItem

//UnmarshalJSON unmarshal json data
func (rsp *BanksBranchesResponseData) UnmarshalJSON(data []byte) error {
	if isEmptyObjectResponseData(data) {
		return nil
	}
	var arr []*BanksBranchesResponseDataItem
	err := json.Unmarshal(data, &arr)
	if err != nil {
		return err
	}
	*rsp = append(*rsp, arr...)
	return nil
}

//BanksBranchesResponseDataItem struct
type BanksBranchesResponseDataItem struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

//BanksResource wrapper
type BanksResource struct {
	ResourceAbstract
}

//GetList get banks list (see https://docs.dusupay.com/sending-money/payouts/bank-codes)
func (r *BanksResource) GetList(ctx context.Context, filter *BanksFilter) (*BanksResponse, *http.Response, error) {
	err := filter.isValid()
	if err != nil {
		return nil, nil, fmt.Errorf("BanksResource.GetList error: %v", err)
	}
	rsp, err := r.ResourceAbstract.tr.Get(ctx, "v1/payment-options/"+filter.buildPath(), nil)
	if err != nil {
		return nil, nil, fmt.Errorf("BanksResource.GetList error: %v", err)
	}
	var result BanksResponse
	err = unmarshalResponse(rsp, &result)
	if err != nil {
		return nil, rsp, fmt.Errorf("BanksResource.GetList error: %v", err)
	}
	if !result.IsSuccess() {
		err = errors.New(result.Message)
	}
	return &result, rsp, err
}

//GetBranchesList get banks branches list (see https://docs.dusupay.com/sending-money/payouts/bank-branches)
func (r *BanksResource) GetBranchesList(ctx context.Context, filter *BanksBranchesFilter) (*BanksBranchesResponse, *http.Response, error) {
	err := filter.isValid()
	if err != nil {
		return nil, nil, fmt.Errorf("BanksResource.GetBranchesList error: %v", err)
	}
	rsp, err := r.ResourceAbstract.tr.Get(ctx, "v1/bank/"+filter.buildPath(), nil)
	if err != nil {
		return nil, nil, fmt.Errorf("BanksResource.GetBranchesList error: %v", err)
	}
	var result BanksBranchesResponse
	err = unmarshalResponse(rsp, &result)
	if err != nil {
		return nil, rsp, fmt.Errorf("BanksResource.GetBranchesList error: %v", err)
	}
	if !result.IsSuccess() {
		err = errors.New(result.Message)
	}
	return &result, rsp, err
}
