package dusupay

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

//RefundRequest struct
type RefundRequest struct {
	Amount            float64 `json:"amount"`
	InternalReference string  `json:"internal_reference"`
}

//Check is valid PayoutRequest parameters
func (rr *RefundRequest) isValid() error {
	var err error
	if rr.InternalReference == "" {
		err = fmt.Errorf(`parameter "internal_reference" is empty`)
	}
	return err
}

//RefundResponse struct
type RefundResponse struct {
	ResponseBody
	Data *RefundResponseData `json:"data,omitempty"`
}

//RefundResponseData struct
type RefundResponseData struct {
	ID                  int64   `json:"id"`
	RefundAmount        float64 `json:"refund_amount"`
	RefundCurrency      string  `json:"refund_currency"`
	TransactionFee      float64 `json:"transaction_fee"`
	TotalDebit          float64 `json:"total_debit"`
	ProviderID          string  `json:"provider_id"`
	MerchantReference   string  `json:"merchant_reference"`
	CollectionReference string  `json:"collection_reference"`
	InternalReference   string  `json:"internal_reference"`
	TransactionType     string  `json:"transaction_type"`
	TransactionStatus   string  `json:"transaction_status"`
	AccountNumber       string  `json:"account_number"`
	Message             string  `json:"message"`
}

//RefundsResource wrapper
type RefundsResource struct {
	ResourceAbstract
}

//Create refund request (see https://docs.dusupay.com/appendix/refunds)
func (r *RefundsResource) Create(ctx context.Context, req *RefundRequest) (*RefundResponse, *http.Response, error) {
	err := req.isValid()
	if err != nil {
		return nil, nil, fmt.Errorf("RefundsResource.Create error: %v", err)
	}
	post, err := transformStructToMap(req)
	if err != nil {
		return nil, nil, fmt.Errorf("RefundsResource.Create error: %v", err)
	}
	rsp, err := r.ResourceAbstract.tr.Post(ctx, "v1/refund", post, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("RefundsResource.Create error: %v", err)
	}
	var result RefundResponse
	err = unmarshalResponse(rsp, &result)
	if err != nil {
		return nil, rsp, fmt.Errorf("RefundsResource.Create error: %v", err)
	}
	if !result.IsSuccess() {
		err = errors.New(result.Message)
	}
	return &result, rsp, err
}
