package dusupay

import "context"

type PayoutRequest struct {
	Currency CurrencyCode          `json:"currency"`
	Amount   float64               `json:"amount"`
	Method   TransactionMethodCode `json:"method"`
}

//Check is valid PayoutRequest parameters
func (pr *PayoutRequest) isValid() error {
	return nil
}

//Payouts resource wrapper
type PayoutsResource struct {
	*ResourceAbstract
}

func (r *PayoutsResource) Create(ctx context.Context, req *PayoutRequest) (*Response, error) {
	return nil, nil
}
