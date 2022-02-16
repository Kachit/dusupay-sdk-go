package dusupay

import (
	"context"
	"encoding/json"
	"fmt"
)

type BalancesResponse struct {
	*ResponseBody
	Data *BalancesResponseData `json:"data,omitempty"`
}

type BalancesResponseData []*BalancesResponseDataItem

//UnmarshalJSON
func (brd *BalancesResponseData) UnmarshalJSON(data []byte) error {
	if isEmptyObjectResponseData(data) {
		return nil
	}
	var arr []*BalancesResponseDataItem
	err := json.Unmarshal(data, &arr)
	if err != nil {
		return err
	}
	*brd = append(*brd, arr...)
	return nil
}

type BalancesResponseDataItem struct {
	Currency string  `json:"currency"`
	Balance  float64 `json:"balance"`
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
