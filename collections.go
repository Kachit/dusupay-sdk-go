package dusupay

import (
	"context"
	"fmt"
)

//Collections resource wrapper
type CollectionsResource struct {
	*ResourceAbstract
}

func (r *CollectionsResource) create(ctx context.Context, req *CollectionRequest) (*Response, error) {
	return nil, nil
}

type CollectionRequest struct {
	Currency          CurrencyCode          `json:"currency"`
	Amount            float64               `json:"amount"`
	Method            TransactionMethodCode `json:"method"`
	ProviderId        string                `json:"provider_id"`
	AccountNumber     string                `json:"account_number"`
	MerchantReference string                `json:"merchant_reference"`
	Narration         string                `json:"narration"`
	MobileMoneyHpp    bool                  `json:"mobile_money_hpp"`
	RedirectUrl       string                `json:"redirect_url"`
	AccountName       string                `json:"account_name"`
	AccountEmail      string                `json:"account_email"`
}

//Check is valid CollectionRequest parameters
func (cr *CollectionRequest) isValid() error {
	var err error
	if cr.Currency == "" {
		err = fmt.Errorf(`parameter "currency" is empty`)
	} else if cr.Amount == 0 {
		err = fmt.Errorf(`parameter "amount" is empty`)
	} else if cr.Method == "" {
		err = fmt.Errorf(`parameter "method" is empty`)
	} else if cr.ProviderId == "" {
		err = fmt.Errorf(`parameter "provider_id" is empty`)
	} else if cr.MerchantReference == "" {
		err = fmt.Errorf(`parameter "merchant_reference" is empty`)
	} else if cr.Narration == "" {
		err = fmt.Errorf(`parameter "narration" is empty`)
	} else if cr.RedirectUrl == "" && cr.Method != TransactionMethodMobileMoney {
		err = fmt.Errorf(`parameter "redirect_url" is empty`)
	} else if cr.AccountNumber == "" && cr.Method == TransactionMethodMobileMoney {
		err = fmt.Errorf(`parameter "account_number" is empty`)
	}
	return err
}

type CollectionResponse struct {
}
