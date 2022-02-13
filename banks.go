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

//Banks resource wrapper
type BanksResource struct {
	*ResourceAbstract
}

//Get banks list
//see https://docs.dusupay.com/sending-money/payouts/bank-codes
func (r *BanksResource) GetList(ctx context.Context, filter *BanksFilter) (*Response, error) {
	err := filter.isValid()
	if err != nil {
		return nil, err
	}
	return nil, nil
}

//Get banks list
//see https://docs.dusupay.com/sending-money/payouts/bank-branches
func (r *BanksResource) GetBranchesList(ctx context.Context, filter *BanksBranchesFilter) (*Response, error) {
	err := filter.isValid()
	if err != nil {
		return nil, err
	}
	return nil, nil
}
