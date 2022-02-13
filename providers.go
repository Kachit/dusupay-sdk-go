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

func (r *ProvidersResource) GetList(ctx context.Context, filter ProvidersFilter) (*Response, error) {
	err := filter.isValid()
	if err != nil {
		return nil, err
	}
	return nil, nil
}
