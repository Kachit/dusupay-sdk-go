package dusupay

import (
	"context"
	"fmt"
)

type BalancesResponse struct {
	*ResponseBody
	Data []*BalancesResponseDataItem `json:"data,omitempty"`
}

type BalancesResponseDataItem struct {
	Currency CurrencyCode `json:"currency"`
	Balance  float64      `json:"balance"`
}

//Merchants resource wrapper
type MerchantsResource struct {
	*ResourceAbstract
}

//get balances list
func (r *MerchantsResource) GetBalances(ctx context.Context) (*Response, error) {
	rsp, err := r.ResourceAbstract.get(ctx, "v1/merchants/balance", nil)
	if err != nil {
		return nil, fmt.Errorf("MerchantsResource@GetBalances error: %v", err)
	}
	return rsp, err
}
