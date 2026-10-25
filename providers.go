package dusupay

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

// ProvidersFilter list of providers filter
type ProvidersFilter struct {
	Currency        CurrencyCode        `json:"currency"`
	TransactionType TransactionTypeCode `json:"transaction_type"`
	CountryCode     CountryCode         `json:"country_code,omitempty"`
}

func (pf *ProvidersFilter) isValid() error {
	if pf.Currency == "" {
		return fmt.Errorf(`parameter "currency" is empty`)
	}

	if pf.TransactionType == "" {
		return fmt.Errorf(`parameter "transaction_type" is empty`)
	}

	return nil
}

// ProvidersResponse struct
type ProvidersResponse struct {
	ResponseBody
	Data *ProvidersResponseData `json:"data,omitempty"`
}

// ProvidersResponseData struct
type ProvidersResponseData struct {
	PaymentProviders []*ProvidersResponseDataItem `json:"payment_providers"`
}

// ProvidersResponseDataItem struct
type ProvidersResponseDataItem struct {
	ProviderCode          string   `json:"provider_code"`
	ProviderName          string   `json:"provider_name"`
	TransactionMethod     string   `json:"transaction_method"`
	TransactionMethodName string   `json:"transaction_method_name"`
	TransactionType       string   `json:"transaction_type"`
	TransactionCurrency   string   `json:"transaction_currency"`
	MinAmount             float64  `json:"min_amount"`
	MaxAmount             float64  `json:"max_amount"`
	MaxDailyAmount        float64  `json:"max_daily_amount"`
	IsActive              bool     `json:"is_active"`
	OptionPrefixes        []string `json:"option_prefixes"`
}

// ProvidersResource wrapper
type ProvidersResource struct {
	ResourceAbstract
}

// GetList get providers list (see https://developer.dusupay.com/utility-functions/payment-options.md)
func (r *ProvidersResource) GetList(
	ctx context.Context, filter *ProvidersFilter,
) (*ProvidersResponse, *http.Response, error) {
	if err := filter.isValid(); err != nil {
		return nil, nil, fmt.Errorf("ProvidersResource.GetList error: %v", err)
	}

	query := map[string]interface{}{
		"currency":         filter.Currency,
		"transaction_type": filter.TransactionType,
	}
	if filter.CountryCode != "" {
		query["country_code"] = filter.CountryCode
	}

	rsp, err := r.tr.Get(ctx, "data/payment-providers", query)
	if err != nil {
		return nil, nil, fmt.Errorf("ProvidersResource.GetList error: %v", err)
	}

	var result ProvidersResponse
	if err = unmarshalResponse(rsp, &result); err != nil {
		return nil, rsp, fmt.Errorf("ProvidersResource.GetList error: %v", err)
	}

	if !result.IsSuccess() {
		err = errors.New(result.Message)
	}

	return &result, rsp, err
}
