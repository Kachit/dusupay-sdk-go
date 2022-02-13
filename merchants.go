package dusupay

import (
	"context"
	"fmt"
)

//Merchants resource wrapper
type MerchantsResource struct {
	*ResourceAbstract
}

type BalancesResponse struct {
	*ResponseBody
	Data []*BalancesResponseDataItem `json:"data"`
}

type BalancesResponseDataItem struct {
	Currency string  `json:"currency"`
	Balance  float64 `json:"balance"`
}

//Get balances list
func (r *MerchantsResource) GetBalances(ctx context.Context) (*Response, error) {
	query := make(map[string]interface{})
	query["api_key"] = r.ResourceAbstract.cfg.PublicKey
	rsp, err := r.ResourceAbstract.Get(ctx, "v1/merchants/balance", query)
	if err != nil {
		return nil, fmt.Errorf("MerchantsResource@GetBalances error: %v", err)
	}
	return rsp, err
}
