package dusupay

import "context"

//Payouts resource wrapper
type PayoutsResource struct {
	*ResourceAbstract
}

type PayoutRequest struct {
	Currency CurrencyCode          `json:"currency"`
	Amount   float64               `json:"amount"`
	Method   TransactionMethodCode `json:"method"`
}

//Check is valid PayoutRequest parameters
func (pr *PayoutRequest) isValid() error {
	return nil
}

func (r *PayoutsResource) Create(ctx context.Context, req *PayoutRequest) (*Response, error) {
	return nil, nil
}
