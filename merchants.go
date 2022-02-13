package dusupay

import "context"

//Merchants resource wrapper
type MerchantsResource struct {
	*ResourceAbstract
}

type BalanceResponse []*BalanceResponseItem

type BalanceResponseItem struct {
	Currency CurrencyCode `json:"currency"`
	Balance  float64      `json:"balance"`
}

func (r *MerchantsResource) GetBalance(ctx context.Context) (*Response, error) {
	return nil, nil
}
