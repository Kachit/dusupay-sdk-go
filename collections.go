package dusupay

type CollectionsResource struct {
	*ResourceAbstract
}

type CollectionRequest struct {
	ApiKey            string                `json:"api_key"`
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

type CollectionResponse struct {
}
