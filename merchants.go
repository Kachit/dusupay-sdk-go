package dusupay

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

// BalancesResponse struct
type BalancesResponse struct {
	ResponseBody
	Data *BalancesResponseData `json:"data,omitempty"`
}

// BalancesResponseData struct
type BalancesResponseData struct {
	MerchantWallets []*WalletBalance `json:"merchant_wallets"`
}

// WalletBalance struct
type WalletBalance struct {
	AvailableBalance float64 `json:"available_balance"`
	ActualBalance    float64 `json:"actual_balance"`
	Currency         string  `json:"currency"`
	UnclearedBalance float64 `json:"uncleared_balance"`
	ReserveBalance   float64 `json:"reserve_balance"`
	WalletType       string  `json:"wallet_type"`
}

// VerifyResponse struct
type VerifyResponse struct {
	ResponseBody
	Data *VerifyResponseData `json:"data,omitempty"`
}

// VerifyResponseData struct
type VerifyResponseData struct {
	ID                int64   `json:"id"`
	MerchantReference string  `json:"merchant_reference"`
	InternalReference string  `json:"internal_reference"`
	TransactionStatus string  `json:"transaction_status"`
	TransactionAmount float64 `json:"transaction_amount"`
	TransactionCharge float64 `json:"transaction_charge"`
	TotalCredit       float64 `json:"total_credit"`
	ProviderCode      string  `json:"provider_code"`
	CustomerName      string  `json:"customer_name"`
	StatusMessage     string  `json:"status_message"`
}

// MerchantsResource wrapper
type MerchantsResource struct {
	ResourceAbstract
}

// GetBalances get balances list (see https://developer.dusupay.com/utility-functions/balance-inquiry.md)
func (r *MerchantsResource) GetBalances(ctx context.Context) (*BalancesResponse, *http.Response, error) {
	rsp, err := r.tr.Get(ctx, "data/wallet-balances", nil)
	if err != nil {
		return nil, nil, fmt.Errorf("MerchantsResource.GetBalances error: %v", err)
	}

	var balances BalancesResponse
	if err = unmarshalResponse(rsp, &balances); err != nil {
		return nil, rsp, fmt.Errorf("MerchantsResource.GetBalances error: %v", err)
	}

	if !balances.IsSuccess() {
		err = errors.New(balances.Message)
	}

	return &balances, rsp, err
}

// Verify transaction status by merchant reference.
// See https://developer.dusupay.com/utility-functions/transaction-status-verification.md
func (r *MerchantsResource) Verify(
	ctx context.Context, merchantReference string,
) (*VerifyResponse, *http.Response, error) {
	if merchantReference == "" {
		return nil, nil, fmt.Errorf(`MerchantsResource.Verify error: parameter "merchant_reference" is empty`)
	}

	rsp, err := r.tr.Get(ctx, "data/transaction/verify/"+merchantReference, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("MerchantsResource.Verify error: %v", err)
	}

	var result VerifyResponse
	if err = unmarshalResponse(rsp, &result); err != nil {
		return nil, rsp, fmt.Errorf("MerchantsResource.Verify error: %v", err)
	}

	if !result.IsSuccess() {
		err = errors.New(result.Message)
	}

	return &result, rsp, err
}
