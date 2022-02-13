package dusupay

import "context"

//Refunds resource wrapper
type RefundsResource struct {
	*ResourceAbstract
}

func (r *RefundsResource) create(ctx context.Context) (*Response, error) {
	return nil, nil
}
