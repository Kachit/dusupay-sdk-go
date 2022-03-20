package dusupay

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

//PayoutRequest struct
type PayoutRequest struct {
	Currency          CurrencyCode          `json:"currency"`
	Amount            float64               `json:"amount"`
	Method            TransactionMethodCode `json:"method"`
	ProviderId        string                `json:"provider_id"`
	AccountNumber     string                `json:"account_number"`
	AccountName       string                `json:"account_name"`
	AccountEmail      string                `json:"account_email"`
	MerchantReference string                `json:"merchant_reference"`
	Narration         string                `json:"narration"`
	ExtraParams       struct {
		BankCode       string `json:"bank_code"`
		BankBranchCode string `json:"branch_code"`
	} `json:"extra_params"`
}

//Check is valid PayoutRequest parameters
func (pr *PayoutRequest) isValid() error {
	var err error
	if pr.Currency == "" {
		err = fmt.Errorf(`parameter "currency" is empty`)
	} else if pr.Amount == 0 {
		err = fmt.Errorf(`parameter "amount" is empty`)
	} else if pr.Method == "" {
		err = fmt.Errorf(`parameter "method" is empty`)
	} else if pr.ProviderId == "" {
		err = fmt.Errorf(`parameter "provider_id" is empty`)
	} else if pr.MerchantReference == "" {
		err = fmt.Errorf(`parameter "merchant_reference" is empty`)
	} else if pr.Narration == "" {
		err = fmt.Errorf(`parameter "narration" is empty`)
	} else if pr.AccountNumber == "" {
		err = fmt.Errorf(`parameter "account_number" is empty`)
	} else if pr.AccountName == "" {
		err = fmt.Errorf(`parameter "account_name" is empty`)
	}
	return err
}

//PayoutResponse struct
type PayoutResponse struct {
	*ResponseBody
	Data *PayoutResponseData `json:"data,omitempty"`
}

//PayoutResponseData struct
type PayoutResponseData struct {
	ID                int64   `json:"id"`
	RequestAmount     float64 `json:"request_amount"`
	RequestCurrency   string  `json:"request_currency"`
	AccountAmount     float64 `json:"account_amount"`
	AccountCurrency   string  `json:"account_currency"`
	TransactionFee    float64 `json:"transaction_fee"`
	TotalDebit        float64 `json:"total_debit"`
	ProviderID        string  `json:"provider_id"`
	MerchantReference string  `json:"merchant_reference"`
	InternalReference string  `json:"internal_reference"`
	TransactionStatus string  `json:"transaction_status"`
	TransactionType   string  `json:"transaction_type"`
	Message           string  `json:"message"`
}

//PayoutsResource wrapper
type PayoutsResource struct {
	*ResourceAbstract
}

//Create payout request (see https://docs.dusupay.com/sending-money/payouts/post-payout-request)
func (r *PayoutsResource) Create(ctx context.Context, req *PayoutRequest) (*PayoutResponse, *http.Response, error) {
	err := req.isValid()
	if err != nil {
		return nil, nil, fmt.Errorf("PayoutsResource.Create error: %v", err)
	}
	post, err := transformStructToMap(req)
	if err != nil {
		return nil, nil, fmt.Errorf("PayoutsResource.Create error: %v", err)
	}
	rsp, err := r.ResourceAbstract.tr.Post(ctx, "v1/payouts", post, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("PayoutsResource.Create error: %v", err)
	}
	var result PayoutResponse
	err = unmarshalResponse(rsp, &result)
	if err != nil {
		return nil, rsp, fmt.Errorf("PayoutsResource.Create error: %v", err)
	}
	if !result.IsSuccess() {
		err = errors.New(result.Message)
	}
	return &result, rsp, err
}
