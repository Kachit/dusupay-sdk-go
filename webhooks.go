package dusupay

import (
	"fmt"
	"net/url"
	"strconv"
)

// CallbackPayloadInterface implemented by all callback payload types
type CallbackPayloadInterface interface {
	BuildPayloadString() string
}

// buildPayloadString builds the RSA/HMAC verification payload string
// Payload format: event:merchant_reference:internal_reference:transaction_type:transaction_status
func buildPayloadString(event CallbackEventCode, merchantRef, internalRef, txType, txStatus string) string {
	return fmt.Sprintf("%s:%s:%s:%s:%s", event, merchantRef, internalRef, txType, txStatus)
}

// CollectionCallback webhook callback request body for collection transactions.
// See https://developer.dusupay.com/funds-collection/getting-started.md
type CollectionCallback struct {
	Event   CallbackEventCode          `json:"event"`
	Payload *CollectionCallbackPayload `json:"payload"`
}

// BuildPayloadString builds the RSA/HMAC verification payload string
func (cc *CollectionCallback) BuildPayloadString() string {
	p := cc.Payload
	if p == nil {
		p = &CollectionCallbackPayload{}
	}

	return buildPayloadString(cc.Event, p.MerchantReference, p.InternalReference, p.TransactionType, p.TransactionStatus)
}

// CollectionCallbackPayload collection transaction details sent in the callback
type CollectionCallbackPayload struct {
	ID                  int64                  `json:"id"`
	MerchantReference   string                 `json:"merchant_reference"`
	InternalReference   string                 `json:"internal_reference"`
	TransactionType     string                 `json:"transaction_type"`
	RequestCurrency     string                 `json:"request_currency"`
	RequestAmount       float64                `json:"request_amount"`
	TransactionCurrency string                 `json:"transaction_currency"`
	TransactionAmount   float64                `json:"transaction_amount"`
	TransactionCharge   float64                `json:"transaction_charge"`
	TransactionAccount  string                 `json:"transaction_account"`
	ChargeCustomer      bool                   `json:"charge_customer"`
	TotalCredit         float64                `json:"total_credit"`
	ProviderCode        string                 `json:"provider_code"`
	CustomerName        string                 `json:"customer_name"`
	InstitutionName     string                 `json:"institution_name,omitempty"`
	TransactionStatus   string                 `json:"transaction_status"`
	StatusMessage       string                 `json:"status_message"`
	PaymentInstructions string                 `json:"payment_instructions,omitempty"`
	Authorization       *CallbackAuthorization `json:"authorization,omitempty"`
	SubAccountDetails   *CallbackSubAccount    `json:"sub_account_details,omitempty"`
}

// CallbackAuthorization authorization requirements sent with transaction.processing callbacks
type CallbackAuthorization struct {
	Mode               AuthorizationModeCode             `json:"mode"`
	RequiredParameters []*CallbackAuthorizationParameter `json:"required_parameters"`
}

// CallbackAuthorizationParameter parameter the merchant must collect and submit (e.g. otp)
type CallbackAuthorizationParameter struct {
	ParameterName string `json:"parameter_name"`
}

// CallbackSubAccount virtual bank account details for bank transfer collections
type CallbackSubAccount struct {
	AccountName string  `json:"account_name"`
	AccountNo   string  `json:"account_no"`
	BankName    string  `json:"bank_name"`
	Amount      float64 `json:"amount"`
}

// PayoutCallback webhook callback request body for payout transactions.
// See https://developer.dusupay.com/payouts-disbursements/getting-started.md
type PayoutCallback struct {
	Event   CallbackEventCode      `json:"event"`
	Payload *PayoutCallbackPayload `json:"payload"`
}

// BuildPayloadString builds the RSA/HMAC verification payload string
func (pc *PayoutCallback) BuildPayloadString() string {
	p := pc.Payload
	if p == nil {
		p = &PayoutCallbackPayload{}
	}

	return buildPayloadString(pc.Event, p.MerchantReference, p.InternalReference, p.TransactionType, p.TransactionStatus)
}

// PayoutCallbackPayload payout transaction details sent in the callback
type PayoutCallbackPayload struct {
	ID                  int64   `json:"id"`
	MerchantReference   string  `json:"merchant_reference"`
	InternalReference   string  `json:"internal_reference"`
	TransactionType     string  `json:"transaction_type"`
	RequestCurrency     string  `json:"request_currency"`
	RequestAmount       float64 `json:"request_amount"`
	TransactionCurrency string  `json:"transaction_currency"`
	TransactionAmount   float64 `json:"transaction_amount"`
	TransactionCharge   float64 `json:"transaction_charge"`
	TransactionAccount  string  `json:"transaction_account"`
	ChargeCustomer      bool    `json:"charge_customer"`
	TotalDebit          float64 `json:"total_debit"`
	ProviderCode        string  `json:"provider_code"`
	CustomerName        string  `json:"customer_name"`
	TransactionStatus   string  `json:"transaction_status"`
	StatusMessage       string  `json:"status_message"`
}

// RedirectPayload data passed as query parameters to the merchant redirect_url (hosted pages and card payments).
// See https://developer.dusupay.com/funds-collection/card-payments/hosted-payment-page.md
type RedirectPayload struct {
	ID                int64
	Event             CallbackEventCode
	MerchantReference string
	InternalReference string
	TransactionType   string
	TransactionStatus string
	StatusMessage     string
	RSASignature      string
	HMACSignature     string
}

// NewRedirectPayload parses redirect query parameters
func NewRedirectPayload(query url.Values) *RedirectPayload {
	id, _ := strconv.ParseInt(query.Get("id"), 10, 64)

	return &RedirectPayload{
		ID:                id,
		Event:             CallbackEventCode(query.Get("event")),
		MerchantReference: query.Get("merchant_reference"),
		InternalReference: query.Get("internal_reference"),
		TransactionType:   query.Get("transaction_type"),
		TransactionStatus: query.Get("transaction_status"),
		StatusMessage:     query.Get("status_message"),
		RSASignature:      query.Get("rsa_signature"),
		HMACSignature:     query.Get("hmac_signature"),
	}
}

// BuildPayloadString builds the RSA/HMAC verification payload string
func (rp *RedirectPayload) BuildPayloadString() string {
	return buildPayloadString(
		rp.Event, rp.MerchantReference, rp.InternalReference, rp.TransactionType, rp.TransactionStatus,
	)
}
