package dusupay

import "context"

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

func (pf *ProvidersFilter) isValid() error {
	return nil
}

func (r *ProvidersResource) GetList(ctx context.Context, filter ProvidersFilter) (*Response, error) {
	return nil, nil
}
