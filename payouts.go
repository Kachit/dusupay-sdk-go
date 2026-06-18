package dusupay

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

// PayoutRequest struct
type PayoutRequest struct {
	Currency          CurrencyCode          `json:"currency"`
	Amount            float64               `json:"amount"`
	TransactionMethod TransactionMethodCode `json:"transaction_method"`
	ProviderCode      string                `json:"provider_code"`
	AccountNumber     string                `json:"account_number"`
	CustomerName      string                `json:"customer_name"`
	MerchantReference string                `json:"merchant_reference"`
	Description       string                `json:"description"`
	ExtraParams       *PayoutExtraParams    `json:"extra_params,omitempty"`
}

// PayoutExtraParams struct
type PayoutExtraParams struct {
	BankCode string `json:"bank_code"`
}

func (pr *PayoutRequest) isValid() error {
	if pr.Currency == "" {
		return fmt.Errorf(`parameter "currency" is empty`)
	}
	if pr.Amount == 0 {
		return fmt.Errorf(`parameter "amount" is empty`)
	}
	if pr.TransactionMethod == "" {
		return fmt.Errorf(`parameter "transaction_method" is empty`)
	}
	if pr.ProviderCode == "" {
		return fmt.Errorf(`parameter "provider_code" is empty`)
	}
	if pr.MerchantReference == "" {
		return fmt.Errorf(`parameter "merchant_reference" is empty`)
	}
	if pr.Description == "" {
		return fmt.Errorf(`parameter "description" is empty`)
	}
	if pr.AccountNumber == "" {
		return fmt.Errorf(`parameter "account_number" is empty`)
	}
	if pr.CustomerName == "" {
		return fmt.Errorf(`parameter "customer_name" is empty`)
	}
	return nil
}

// PayoutResponse struct
type PayoutResponse struct {
	ResponseBody
	Data *PayoutResponseData `json:"data,omitempty"`
}

// PayoutResponseData struct
type PayoutResponseData struct {
	InternalReference string `json:"internal_reference"`
	MerchantReference string `json:"merchant_reference"`
}

// PayoutsResource wrapper
type PayoutsResource struct {
	ResourceAbstract
}

// Create payout request (see https://developer.dusupay.com/payouts-disbursements/mobile-money-payouts.md)
func (r *PayoutsResource) Create(ctx context.Context, req *PayoutRequest) (*PayoutResponse, *http.Response, error) {
	if err := req.isValid(); err != nil {
		return nil, nil, fmt.Errorf("PayoutsResource.Create error: %v", err)
	}
	post, err := transformStructToMap(req)
	if err != nil {
		return nil, nil, fmt.Errorf("PayoutsResource.Create error: %v", err)
	}
	rsp, err := r.tr.Post(ctx, "payout/send-funds", post, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("PayoutsResource.Create error: %v", err)
	}
	var result PayoutResponse
	if err = unmarshalResponse(rsp, &result); err != nil {
		return nil, rsp, fmt.Errorf("PayoutsResource.Create error: %v", err)
	}
	if !result.IsSuccess() {
		err = errors.New(result.Message)
	}
	return &result, rsp, err
}
