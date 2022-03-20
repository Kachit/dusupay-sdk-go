package dusupay

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

//CollectionRequest struct
type CollectionRequest struct {
	Currency          CurrencyCode          `json:"currency"`
	Amount            float64               `json:"amount"`
	Method            TransactionMethodCode `json:"method"`
	ProviderId        string                `json:"provider_id"`
	AccountNumber     string                `json:"account_number"`
	MerchantReference string                `json:"merchant_reference"`
	Narration         string                `json:"narration"`
	MobileMoneyHpp    bool                  `json:"mobile_money_hpp"`
	RedirectUrl       string                `json:"redirect_url"`
	AccountName       string                `json:"account_name"`
	AccountEmail      string                `json:"account_email"`
	Voucher           string                `json:"voucher"`
}

//Check is valid CollectionRequest parameters
func (cr *CollectionRequest) isValid() error {
	var err error
	if cr.Currency == "" {
		err = fmt.Errorf(`parameter "currency" is empty`)
	} else if cr.Amount == 0 {
		err = fmt.Errorf(`parameter "amount" is empty`)
	} else if cr.Method == "" {
		err = fmt.Errorf(`parameter "method" is empty`)
	} else if cr.ProviderId == "" {
		err = fmt.Errorf(`parameter "provider_id" is empty`)
	} else if cr.MerchantReference == "" {
		err = fmt.Errorf(`parameter "merchant_reference" is empty`)
	} else if cr.Narration == "" {
		err = fmt.Errorf(`parameter "narration" is empty`)
	} else if cr.RedirectUrl == "" && cr.Method != TransactionMethodMobileMoney {
		err = fmt.Errorf(`parameter "redirect_url" is empty`)
	} else if cr.AccountNumber == "" && cr.Method == TransactionMethodMobileMoney {
		err = fmt.Errorf(`parameter "account_number" is empty`)
	}
	return err
}

//CollectionResponse struct
type CollectionResponse struct {
	*ResponseBody
	Data *CollectionResponseData `json:"data,omitempty"`
}

//CollectionResponseData struct
type CollectionResponseData struct {
	ID                int64   `json:"id"`
	RequestAmount     float64 `json:"request_amount"`
	RequestCurrency   string  `json:"request_currency"`
	AccountAmount     float64 `json:"account_amount"`
	AccountCurrency   string  `json:"account_currency"`
	TransactionFee    float64 `json:"transaction_fee"`
	TotalCredit       float64 `json:"total_credit"`
	ProviderID        string  `json:"provider_id"`
	MerchantReference string  `json:"merchant_reference"`
	InternalReference string  `json:"internal_reference"`
	TransactionStatus string  `json:"transaction_status"`
	TransactionType   string  `json:"transaction_type"`
	Message           string  `json:"message"`
	CustomerCharged   bool    `json:"customer_charged"`
	PaymentURL        string  `json:"payment_url"`
	Instructions      []struct {
		StepNo      string `json:"step_no"`
		Description string `json:"description"`
	} `json:"instructions"`
}

//CollectionsResource wrapper
type CollectionsResource struct {
	*ResourceAbstract
}

//Create collection request (see https://docs.dusupay.com/receiving-money/collections/post-collection-request)
func (r *CollectionsResource) Create(ctx context.Context, req *CollectionRequest) (*CollectionResponse, *http.Response, error) {
	err := req.isValid()
	if err != nil {
		return nil, nil, fmt.Errorf("CollectionsResource.Create error: %v", err)
	}
	post, err := transformStructToMap(req)
	if err != nil {
		return nil, nil, fmt.Errorf("CollectionsResource.Create error: %v", err)
	}
	rsp, err := r.ResourceAbstract.tr.Post(ctx, "v1/collections", post, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("CollectionsResource.Create error: %v", err)
	}
	var result CollectionResponse
	err = unmarshalResponse(rsp, &result)
	if err != nil {
		return nil, rsp, fmt.Errorf("CollectionsResource.Create error: %v", err)
	}
	if !result.IsSuccess() {
		err = errors.New(result.Message)
	}
	return &result, rsp, err
}
