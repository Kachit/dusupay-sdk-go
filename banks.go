// Package dusupay provides the Dusupay payment API v2 SDK.
package dusupay

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

// BanksFilter (see https://developer.dusupay.com/utility-functions/payout-bank-codes.md)
type BanksFilter struct {
	ProviderCode string `json:"provider_code"`
}

func (bf *BanksFilter) isValid() error {
	if bf.ProviderCode == "" {
		return fmt.Errorf(`parameter "provider_code" is empty`)
	}

	return nil
}

// BanksResponse struct
type BanksResponse struct {
	ResponseBody
	Data *BanksResponseData `json:"data,omitempty"`
}

// BanksResponseData struct
type BanksResponseData struct {
	PayoutBanks []*BanksResponseDataItem `json:"payout_banks"`
}

// BanksResponseDataItem struct
type BanksResponseDataItem struct {
	BankName string `json:"bank_name"`
	BankCode string `json:"bank_code"`
	IsActive bool   `json:"is_active"`
}

// BanksResource wrapper
type BanksResource struct {
	ResourceAbstract
}

// GetList get payout banks list (see https://developer.dusupay.com/utility-functions/payout-bank-codes.md)
func (r *BanksResource) GetList(ctx context.Context, filter *BanksFilter) (*BanksResponse, *http.Response, error) {
	if err := filter.isValid(); err != nil {
		return nil, nil, fmt.Errorf("BanksResource.GetList error: %v", err)
	}

	query := map[string]interface{}{
		"provider_code": filter.ProviderCode,
	}

	rsp, err := r.tr.Get(ctx, "data/payout-bank-codes", query)
	if err != nil {
		return nil, nil, fmt.Errorf("BanksResource.GetList error: %v", err)
	}

	var result BanksResponse
	if err = unmarshalResponse(rsp, &result); err != nil {
		return nil, rsp, fmt.Errorf("BanksResource.GetList error: %v", err)
	}

	if !result.IsSuccess() {
		err = errors.New(result.Message)
	}

	return &result, rsp, err
}
