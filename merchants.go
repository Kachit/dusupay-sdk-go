package dusupay

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

// BalancesFilter optional balances filter
type BalancesFilter struct {
	Currency CurrencyCode `json:"currency,omitempty"`
}

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
	ID                  int64   `json:"id"`
	MerchantReference   string  `json:"merchant_reference"`
	InternalReference   string  `json:"internal_reference"`
	TransactionType     string  `json:"transaction_type"`
	RequestCurrency     string  `json:"request_currency"`
	RequestAmount       float64 `json:"request_amount"`
	TransactionCurrency string  `json:"transaction_currency"`
	TransactionAmount   float64 `json:"transaction_amount"`
	TransactionCharge   float64 `json:"transaction_charge"`
	TransactionAccount  string  `json:"transaction_account"`
	ChargeCustomer      bool    `json:"charge_customer"`
	TotalCredit         float64 `json:"total_credit,omitempty"`
	TotalDebit          float64 `json:"total_debit,omitempty"`
	ProviderCode        string  `json:"provider_code"`
	CustomerName        string  `json:"customer_name"`
	InstitutionName     string  `json:"institution_name,omitempty"`
	TransactionStatus   string  `json:"transaction_status"`
	StatusMessage       string  `json:"status_message"`
}

// MerchantsResource wrapper
type MerchantsResource struct {
	ResourceAbstract
}

// GetBalances get balances list, filter is optional.
// See https://developer.dusupay.com/utility-functions/balance-inquiry.md
func (r *MerchantsResource) GetBalances(
	ctx context.Context, filter *BalancesFilter,
) (*BalancesResponse, *http.Response, error) {
	var query map[string]interface{}
	if filter != nil && filter.Currency != "" {
		query = map[string]interface{}{"currency": filter.Currency}
	}

	rsp, err := r.tr.Get(ctx, "data/wallet-balances", query)
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
