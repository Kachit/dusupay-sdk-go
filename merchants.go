package dusupay

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

//BalancesResponse struct
type BalancesResponse struct {
	*ResponseBody
	Data *BalancesResponseData `json:"data,omitempty"`
}

//BalancesResponseData struct
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

//BalancesResponseDataItem struct
type BalancesResponseDataItem struct {
	Currency string  `json:"currency"`
	Balance  float64 `json:"balance"`
}

//MerchantsResource wrapper
type MerchantsResource struct {
	*ResourceAbstract
}

//GetBalances get balances list (see https://docs.dusupay.com/appendix/account-balance)
func (r *MerchantsResource) GetBalances(ctx context.Context) (*BalancesResponse, *http.Response, error) {
	rsp, err := r.ResourceAbstract.tr.Get(ctx, "v1/merchants/balance", nil)
	if err != nil {
		return nil, nil, fmt.Errorf("MerchantsResource.GetBalances error: %v", err)
	}
	var balances BalancesResponse
	err = unmarshalResponse(rsp, &balances)
	if err != nil {
		return nil, rsp, fmt.Errorf("MerchantsResource.GetBalances error: %v", err)
	}
	if !balances.IsSuccess() {
		err = errors.New(balances.Message)
	}
	return &balances, rsp, err
}
