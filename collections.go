package dusupay

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

// CollectionRequest struct
type CollectionRequest struct {
	Currency               CurrencyCode          `json:"currency"`
	Amount                 float64               `json:"amount"`
	TransactionMethod      TransactionMethodCode `json:"transaction_method"`
	ProviderCode           string                `json:"provider_code"`
	MerchantReference      string                `json:"merchant_reference"`
	Description            string                `json:"description"`
	Msisdn                 string                `json:"msisdn,omitempty"`
	CustomerName           string                `json:"customer_name,omitempty"`
	CustomerEmail          string                `json:"customer_email,omitempty"`
	ChargeCustomer         bool                  `json:"charge_customer,omitempty"`
	AllowFinalStatusChange bool                  `json:"allow_final_status_change,omitempty"`
}

func (cr *CollectionRequest) isValid() error {
	if cr.Currency == "" {
		return fmt.Errorf(`parameter "currency" is empty`)
	}

	if cr.Amount == 0 {
		return fmt.Errorf(`parameter "amount" is empty`)
	}

	if cr.TransactionMethod == "" {
		return fmt.Errorf(`parameter "transaction_method" is empty`)
	}

	if cr.ProviderCode == "" {
		return fmt.Errorf(`parameter "provider_code" is empty`)
	}

	if cr.MerchantReference == "" {
		return fmt.Errorf(`parameter "merchant_reference" is empty`)
	}

	if cr.Description == "" {
		return fmt.Errorf(`parameter "description" is empty`)
	}

	if cr.Msisdn == "" && cr.TransactionMethod == TransactionMethodMobileMoney {
		return fmt.Errorf(`parameter "msisdn" is empty`)
	}

	return nil
}

// CollectionResponse struct
type CollectionResponse struct {
	ResponseBody
	Data *CollectionResponseData `json:"data,omitempty"`
}

// CollectionResponseData struct
type CollectionResponseData struct {
	InternalReference string `json:"internal_reference"`
	MerchantReference string `json:"merchant_reference"`
	PaymentURL        string `json:"payment_url,omitempty"`
}

// CollectionAuthorizeRequest struct
type CollectionAuthorizeRequest struct {
	InternalReference string `json:"internal_reference"`
	OTP               string `json:"otp"`
}

// CollectionAuthorizeResponse struct
type CollectionAuthorizeResponse struct {
	ResponseBody
	Data *CollectionAuthorizeResponseData `json:"data,omitempty"`
}

// CollectionAuthorizeResponseData struct
type CollectionAuthorizeResponseData struct {
	InternalReference string `json:"internal_reference"`
	MerchantReference string `json:"merchant_reference"`
}

// CollectionsResource wrapper
type CollectionsResource struct {
	ResourceAbstract
}

// Create collection request (see https://developer.dusupay.com/funds-collection/getting-started.md)
func (r *CollectionsResource) Create(
	ctx context.Context, req *CollectionRequest,
) (*CollectionResponse, *http.Response, error) {
	if err := req.isValid(); err != nil {
		return nil, nil, fmt.Errorf("CollectionsResource.Create error: %v", err)
	}

	post, err := transformStructToMap(req)
	if err != nil {
		return nil, nil, fmt.Errorf("CollectionsResource.Create error: %v", err)
	}

	rsp, err := r.tr.Post(ctx, "collections/initialize", post, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("CollectionsResource.Create error: %v", err)
	}

	var result CollectionResponse
	if err = unmarshalResponse(rsp, &result); err != nil {
		return nil, rsp, fmt.Errorf("CollectionsResource.Create error: %v", err)
	}

	if !result.IsSuccess() {
		err = errors.New(result.Message)
	}

	return &result, rsp, err
}

// Authorize collection OTP.
// See https://developer.dusupay.com/funds-collection/mobile-money-collection/mobile-money-direct-charge.md
func (r *CollectionsResource) Authorize(
	ctx context.Context, req *CollectionAuthorizeRequest,
) (*CollectionAuthorizeResponse, *http.Response, error) {
	if req.InternalReference == "" {
		return nil, nil, fmt.Errorf(`CollectionsResource.Authorize error: parameter "internal_reference" is empty`)
	}

	if req.OTP == "" {
		return nil, nil, fmt.Errorf(`CollectionsResource.Authorize error: parameter "otp" is empty`)
	}

	post, err := transformStructToMap(req)
	if err != nil {
		return nil, nil, fmt.Errorf("CollectionsResource.Authorize error: %v", err)
	}

	rsp, err := r.tr.Post(ctx, "collections/authorize", post, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("CollectionsResource.Authorize error: %v", err)
	}

	var result CollectionAuthorizeResponse
	if err = unmarshalResponse(rsp, &result); err != nil {
		return nil, rsp, fmt.Errorf("CollectionsResource.Authorize error: %v", err)
	}

	if !result.IsSuccess() {
		err = errors.New(result.Message)
	}

	return &result, rsp, err
}
