package dusupay

import (
	"context"
	"fmt"
)

//Banks list filter
//see https://docs.dusupay.com/sending-money/payouts/bank-codes
type BanksFilter struct {
	Method TransactionMethodCode `json:"method"`
	//ISO-2 country code from those supported
	Country CountryCode `json:"country_code"`
}

//Check is valid BanksFilter parameters
func (bf *BanksFilter) isValid() error {
	var err error
	if bf.Country == "" {
		err = fmt.Errorf(`parameter "country_code" is empty`)
	} else if bf.Method == "" {
		err = fmt.Errorf(`parameter "method" is empty`)
	}
	return err
}

func (bf *BanksFilter) buildPath() string {
	return string(bf.Method) + "/bank/" + string(bf.Country)
}

//Banks branches list filter
//see https://docs.dusupay.com/sending-money/payouts/bank-branches
type BanksBranchesFilter struct {
	//ISO-2 country code from those supported
	Country CountryCode `json:"country_code"`
	Bank    string      `json:"bank_code"`
}

//Check is valid BanksBranchesFilter parameters
func (bbf *BanksBranchesFilter) isValid() error {
	var err error
	if bbf.Country == "" {
		err = fmt.Errorf(`parameter "country_code" is empty`)
	} else if bbf.Bank == "" {
		err = fmt.Errorf(`parameter "bank_code" is empty`)
	}
	return err
}

func (bbf *BanksBranchesFilter) buildPath() string {
	return "bank/" + string(bbf.Country) + "/branches/" + string(bbf.Bank)
}

//BanksBranchesResponse struct
type BanksResponse struct {
	*ResponseBody
	Data []*BanksResponseDataItem `json:"data,omitempty"`
}

type BanksResponseDataItem struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

//BanksBranchesResponse struct
type BanksBranchesResponse struct {
	*ResponseBody
	Data []*BanksBranchesResponseDataItem `json:"data,omitempty"`
}

type BanksBranchesResponseDataItem struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

//Banks resource wrapper
type BanksResource struct {
	*ResourceAbstract
}

//get banks list
//see https://docs.dusupay.com/sending-money/payouts/bank-codes
func (r *BanksResource) GetList(ctx context.Context, filter *BanksFilter) (*Response, error) {
	err := filter.isValid()
	if err != nil {
		return nil, err
	}
	rsp, err := r.ResourceAbstract.get(ctx, "v1/payment-options/"+filter.buildPath(), nil)
	if err != nil {
		return nil, fmt.Errorf("BanksResource@GetList error: %v", err)
	}
	return rsp, err
}

//get banks branches list
//see https://docs.dusupay.com/sending-money/payouts/bank-branches
func (r *BanksResource) GetBranchesList(ctx context.Context, filter *BanksBranchesFilter) (*Response, error) {
	err := filter.isValid()
	if err != nil {
		return nil, err
	}
	rsp, err := r.ResourceAbstract.get(ctx, "v1/bank/"+filter.buildPath(), nil)
	if err != nil {
		return nil, fmt.Errorf("BanksResource@GetBranchesList error: %v", err)
	}
	return rsp, err
}
