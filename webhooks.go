package dusupay

type CollectionWebhook struct {
	ID                int64   `json:"id"`
	RequestAmount     float64 `json:"request_amount"`
	RequestCurrency   string  `json:"request_currency"`
	AccountAmount     float64 `json:"account_amount"`
	AccountCurrency   string  `json:"account_currency"`
	TransactionFee    float64 `json:"transaction_fee"`
	TotalCredit       float64 `json:"total_credit"`
	CustomerCharged   bool    `json:"customer_charged"`
	ProviderID        string  `json:"provider_id"`
	MerchantReference string  `json:"merchant_reference"`
	InternalReference string  `json:"internal_reference"`
	TransactionStatus string  `json:"transaction_status"`
	TransactionType   string  `json:"transaction_type"`
	Message           string  `json:"message"`
	AccountNumber     string  `json:"account_number"`
	AccountName       string  `json:"account_name"`
	InstitutionName   string  `json:"institution_name"`
}

type PayoutWebhook struct {
	ID                int64   `json:"id"`
	RequestAmount     float64 `json:"request_amount"`
	RequestCurrency   string  `json:"request_currency"`
	AccountAmount     float64 `json:"account_amount"`
	AccountCurrency   string  `json:"account_currency"`
	TransactionFee    float64 `json:"transaction_fee"`
	TotalDebit        float64 `json:"total_debit"`
	ProviderID        string  `json:"provider_id"`
	MerchantReference string  `json:"merchant_reference"`
	InternalReference string  `json:"internal_reference"`
	TransactionStatus string  `json:"transaction_status"`
	TransactionType   string  `json:"transaction_type"`
	Message           string  `json:"message"`
	AccountNumber     string  `json:"account_number"`
	AccountName       string  `json:"account_name"`
	InstitutionName   string  `json:"institution_name"`
}
type RefundWebhook struct {
	ID                  int64   `json:"id"`
	RefundAmount        float64 `json:"refund_amount"`
	RefundCurrency      string  `json:"refund_currency"`
	TransactionFee      float64 `json:"transaction_fee"`
	TotalDebit          float64 `json:"total_debit"`
	ProviderID          string  `json:"provider_id"`
	CollectionReference string  `json:"collection_reference"`
	InternalReference   string  `json:"internal_reference"`
	TransactionType     string  `json:"transaction_type"`
	TransactionStatus   string  `json:"transaction_status"`
	AccountNumber       string  `json:"account_number"`
	Message             string  `json:"message"`
}
