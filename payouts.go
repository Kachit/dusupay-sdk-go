package dusupay

import "context"

//Payouts resource wrapper
type PayoutsResource struct {
	*ResourceAbstract
}

type PayoutRequest struct {
	ApiKey   string                `json:"api_key"`
	Currency CurrencyCode          `json:"currency"`
	Amount   float64               `json:"amount"`
	Method   TransactionMethodCode `json:"method"`
}

func (pr *PayoutRequest) isValid() error {
	return nil
}

func (r *PayoutsResource) create(ctx context.Context, req *PayoutRequest) (*Response, error) {
	return nil, nil
}
