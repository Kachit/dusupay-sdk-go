package dusupay

import "context"

//Banks list filter
type BanksFilter struct {
	Method  TransactionMethodCode `json:"method"`
	Country CountryCode           `json:"country_code"`
}

func (bf *BanksFilter) isValid() error {
	return nil
}

//Banks branches list filter
type BanksBranchesFilter struct {
	Country CountryCode `json:"country_code"`
	Bank    string      `json:"bank_code"`
}

func (bbf *BanksBranchesFilter) isValid() error {
	return nil
}

//Payouts resource wrapper
type BanksResource struct {
	*ResourceAbstract
}

func (r *BanksResource) GetList(ctx context.Context, filter *BanksFilter) (*Response, error) {
	return nil, nil
}

func (r *BanksResource) GetBranchesList(ctx context.Context, filter *BanksBranchesFilter) (*Response, error) {
	return nil, nil
}
