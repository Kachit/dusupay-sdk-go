package dusupay

import (
	"context"
	"fmt"
)

//Providers resource wrapper
type ProvidersResource struct {
	*ResourceAbstract
}

//Providers list filter
type ProvidersFilter struct {
	TransactionType TransactionTypeCode   `json:"transaction_type"`
	Method          TransactionMethodCode `json:"method"`
	Country         CountryCode           `json:"country"`
}

type ProvidersResponse struct {
	*ResponseBody
	Data []*ProvidersResponseDataItem `json:"data,omitempty"`
}

type ProvidersResponseDataItem struct {
	ID                  string `json:"id"`
	Name                string `json:"name"`
	TransactionCurrency string `json:"transaction_currency"`
	MinAmount           int    `json:"min_amount"`
	MaxAmount           int    `json:"max_amount"`
	Available           bool   `json:"available"`
	SandboxTestAccounts struct {
		Success string `json:"success"`
		Failure string `json:"failure"`
	} `json:"sandbox_test_accounts"`
}

//Check is valid ProvidersFilter parameters
func (pf *ProvidersFilter) isValid() error {
	var err error
	if pf.Country == "" {
		err = fmt.Errorf(`parameter "country_code" is empty`)
	} else if pf.Method == "" {
		err = fmt.Errorf(`parameter "method" is empty`)
	} else if pf.TransactionType == "" {
		err = fmt.Errorf(`parameter "transaction_type" is empty`)
	}
	return err
}

func (pf *ProvidersFilter) buildPath() string {
	return string(pf.TransactionType) + "/" + string(pf.Method) + "/" + string(pf.Country)
}

func (r *ProvidersResource) GetList(ctx context.Context, filter ProvidersFilter) (*Response, error) {
	err := filter.isValid()
	if err != nil {
		return nil, fmt.Errorf("ProvidersResource@GetList error: %v", err)
	}
	rsp, err := r.ResourceAbstract.get(ctx, "v1/payment-options/"+filter.buildPath(), nil)
	if err != nil {
		return nil, fmt.Errorf("ProvidersResource@GetList error: %v", err)
	}
	return rsp, err
}
