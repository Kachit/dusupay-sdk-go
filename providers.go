package dusupay

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

//ProvidersFilter list of providers filter
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

//buildPath method
func (pf *ProvidersFilter) buildPath() string {
	return string(pf.TransactionType) + "/" + string(pf.Method) + "/" + string(pf.Country)
}

//ProvidersResponse struct
type ProvidersResponse struct {
	*ResponseBody
	Data *ProvidersResponseData `json:"data,omitempty"`
}

//ProvidersResponseData struct
type ProvidersResponseData []*ProvidersResponseDataItem

//UnmarshalJSON method
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

//ProvidersResponseDataItem struct
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

//ProvidersResource wrapper
type ProvidersResource struct {
	*ResourceAbstract
}

//GetList Get providers list (see https://docs.dusupay.com/appendix/payment-options/payment-providers)
func (r *ProvidersResource) GetList(ctx context.Context, filter *ProvidersFilter) (*ProvidersResponse, *http.Response, error) {
	err := filter.isValid()
	if err != nil {
		return nil, nil, fmt.Errorf("ProvidersResource.GetList error: %v", err)
	}
	rsp, err := r.ResourceAbstract.tr.Get(ctx, "v1/payment-options/"+filter.buildPath(), nil)
	if err != nil {
		return nil, nil, fmt.Errorf("ProvidersResource.GetList error: %v", err)
	}
	var result ProvidersResponse
	err = unmarshalResponse(rsp, &result)
	if err != nil {
		return nil, rsp, fmt.Errorf("ProvidersResource.GetList error: %v", err)
	}
	if !result.IsSuccess() {
		err = errors.New(result.Message)
	}
	return &result, rsp, err
}
