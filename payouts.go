package dusupay

type PayoutsResource struct {
	*ResourceAbstract
}

type PayoutRequest struct {
	ApiKey   string                `json:"api_key"`
	Currency CurrencyCode          `json:"currency"`
	Amount   float64               `json:"amount"`
	Method   TransactionMethodCode `json:"method"`
}
