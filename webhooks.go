package dusupay

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

//CollectionWebhook struct
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

func (cw *CollectionWebhook) BuildPayloadString(url string) string {
	return fmt.Sprintf("%d:%s:%s:%s", cw.ID, cw.InternalReference, cw.TransactionStatus, url)
}

//PayoutWebhook struct
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

func (pw *PayoutWebhook) BuildPayloadString(url string) string {
	return fmt.Sprintf("%d:%s:%s:%s", pw.ID, pw.InternalReference, pw.TransactionStatus, url)
}

//RefundWebhook struct
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

func (rw *RefundWebhook) BuildPayloadString(url string) string {
	return fmt.Sprintf("%d:%s:%s:%s", rw.ID, rw.InternalReference, rw.TransactionStatus, url)
}

//WebhookResponse struct
type WebhookResponse struct {
	ResponseBody
	Data *WebhookResponseData `json:"data,omitempty"`
}

//WebhookResponseData struct
type WebhookResponseData struct {
	Payload *WebhookResponsePayload `json:"payload,omitempty"`
}

//WebhookResponsePayload struct
type WebhookResponsePayload struct {
	ID                int64   `json:"id"`
	RequestAmount     float64 `json:"request_amount"`
	RequestCurrency   string  `json:"request_currency"`
	AccountAmount     float64 `json:"account_amount"`
	AccountCurrency   string  `json:"account_currency"`
	TransactionFee    float64 `json:"transaction_fee"`
	ProviderID        string  `json:"provider_id"`
	MerchantReference string  `json:"merchant_reference"`
	InternalReference string  `json:"internal_reference"`
	TransactionStatus string  `json:"transaction_status"`
	TransactionType   string  `json:"transaction_type"`
	Message           string  `json:"message"`
}

//WebhooksResource wrapper
type WebhooksResource struct {
	ResourceAbstract
}

//SendCallback (see https://docs.dusupay.com/appendix/webhooks/webhook-trigger)
func (r *WebhooksResource) SendCallback(ctx context.Context, internalReference string) (*WebhookResponse, *http.Response, error) {
	rsp, err := r.ResourceAbstract.tr.Get(ctx, "v1/send-callback/"+internalReference, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("WebhooksResource.SendCallback error: %v", err)
	}
	var response WebhookResponse
	err = unmarshalResponse(rsp, &response)
	if err != nil {
		return nil, rsp, fmt.Errorf("WebhooksResource.SendCallback error: %v", err)
	}
	if !response.IsSuccess() {
		err = errors.New(response.Message)
	}
	return &response, rsp, err
}
