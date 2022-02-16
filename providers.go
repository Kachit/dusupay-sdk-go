package dusupay

import (
	"context"
	"encoding/json"
	"fmt"
)

//Providers list filter
type ProvidersFilter struct {
	TransactionType TransactionTypeCode   `json:"transaction_type"`
	Method          TransactionMethodCode `json:"method"`
	Country         CountryCode           `json:"country"`
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

type ProvidersResponse struct {
	*ResponseBody
	Data *ProvidersResponseData `json:"data,omitempty"`
}

type ProvidersResponseData []*ProvidersResponseDataItem

//UnmarshalJSON
func (rsp *ProvidersResponseData) UnmarshalJSON(data []byte) error {
	if isEmptyObjectResponseData(data) {
		return nil
	}
	var arr []*ProvidersResponseDataItem
	err := json.Unmarshal(data, &arr)
	if err != nil {
		return err
	}
	*rsp = append(*rsp, arr...)
	return nil
}

type ProvidersResponseDataItem struct {
	ID                  string  `json:"id"`
	Name                string  `json:"name"`
	TransactionCurrency string  `json:"transaction_currency"`
	MinAmount           float64 `json:"min_amount"`
	MaxAmount           float64 `json:"max_amount"`
	Available           bool    `json:"available"`
	SandboxTestAccounts struct {
		Success string `json:"success"`
		Failure string `json:"failure"`
	} `json:"sandbox_test_accounts,omitempty"`
}

//Providers resource wrapper
type ProvidersResource struct {
	*ResourceAbstract
}

//Get providers list (see https://docs.dusupay.com/appendix/payment-options/payment-providers)
func (r *ProvidersResource) GetList(ctx context.Context, filter *ProvidersFilter) (*Response, error) {
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
