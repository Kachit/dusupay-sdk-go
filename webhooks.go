package dusupay

import "fmt"

// CallbackPayloadInterface implemented by all callback payload types
type CallbackPayloadInterface interface {
	BuildPayloadString() string
}

// CollectionCallbackPayload webhook callback for collection transactions
type CollectionCallbackPayload struct {
	Event              string  `json:"event"`
	ID                 int64   `json:"id"`
	MerchantReference  string  `json:"merchant_reference"`
	InternalReference  string  `json:"internal_reference"`
	TransactionType    string  `json:"transaction_type"`
	TransactionStatus  string  `json:"transaction_status"`
	RequestAmount      float64 `json:"request_amount"`
	RequestCurrency    string  `json:"request_currency"`
	TransactionCharge  float64 `json:"transaction_charge"`
	TotalCredit        float64 `json:"total_credit"`
	TransactionAccount string  `json:"transaction_account"`
	ChargeCustomer     bool    `json:"charge_customer"`
	ProviderCode       string  `json:"provider_code"`
	CustomerName       string  `json:"customer_name"`
	StatusMessage      string  `json:"status_message"`
}

// BuildPayloadString builds the RSA/HMAC verification payload string
func (cp *CollectionCallbackPayload) BuildPayloadString() string {
	return fmt.Sprintf("%s:%s:%s:%s:%s",
		cp.Event, cp.MerchantReference, cp.InternalReference,
		cp.TransactionType, cp.TransactionStatus)
}

// PayoutCallbackPayload webhook callback for payout transactions
type PayoutCallbackPayload struct {
	Event             string  `json:"event"`
	ID                int64   `json:"id"`
	MerchantReference string  `json:"merchant_reference"`
	InternalReference string  `json:"internal_reference"`
	TransactionType   string  `json:"transaction_type"`
	TransactionStatus string  `json:"transaction_status"`
	RequestAmount     float64 `json:"request_amount"`
	RequestCurrency   string  `json:"request_currency"`
	TransactionCharge float64 `json:"transaction_charge"`
	TotalDebit        float64 `json:"total_debit"`
	ProviderCode      string  `json:"provider_code"`
	CustomerName      string  `json:"customer_name"`
	AccountNumber     string  `json:"account_number"`
	StatusMessage     string  `json:"status_message"`
}

// BuildPayloadString builds the RSA/HMAC verification payload string
func (pp *PayoutCallbackPayload) BuildPayloadString() string {
	return fmt.Sprintf("%s:%s:%s:%s:%s",
		pp.Event, pp.MerchantReference, pp.InternalReference,
		pp.TransactionType, pp.TransactionStatus)
}
