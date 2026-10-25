package dusupay

import "encoding/json"

// TransactionTypeCode type
type TransactionTypeCode string

// TransactionTypeCollection const
const TransactionTypeCollection TransactionTypeCode = "COLLECTION"

// TransactionTypePayout const
const TransactionTypePayout TransactionTypeCode = "PAYOUT"

// TransactionTypeRefund const
const TransactionTypeRefund TransactionTypeCode = "REFUND"

// TransactionStatusCode type
type TransactionStatusCode string

// TransactionStatusPending const
const TransactionStatusPending TransactionStatusCode = "PENDING"

// TransactionStatusProcessing const
const TransactionStatusProcessing TransactionStatusCode = "PROCESSING"

// TransactionStatusFailed const
const TransactionStatusFailed TransactionStatusCode = "FAILED"

// TransactionStatusCompleted const
const TransactionStatusCompleted TransactionStatusCode = "COMPLETED"

// TransactionStatusCancelled const
const TransactionStatusCancelled TransactionStatusCode = "CANCELLED"

// CallbackEventCode type
// See https://developer.dusupay.com/utility-functions/handling-notifications-callbacks/callback-events.md
type CallbackEventCode string

// CallbackEventRequestFailed const
const CallbackEventRequestFailed CallbackEventCode = "request.failed"

// CallbackEventTransactionCharges const
const CallbackEventTransactionCharges CallbackEventCode = "transaction.charges"

// CallbackEventTransactionProcessing const
const CallbackEventTransactionProcessing CallbackEventCode = "transaction.processing"

// CallbackEventTransactionFailed const
const CallbackEventTransactionFailed CallbackEventCode = "transaction.failed"

// CallbackEventTransactionCompleted const
const CallbackEventTransactionCompleted CallbackEventCode = "transaction.completed"

// CallbackEventInvoiceGenerated const
const CallbackEventInvoiceGenerated CallbackEventCode = "invoice.generated"

// CallbackEventInvoicePaid const
const CallbackEventInvoicePaid CallbackEventCode = "invoice.paid"

// AuthorizationModeCode type
type AuthorizationModeCode string

// AuthorizationModeNone const
const AuthorizationModeNone AuthorizationModeCode = "NONE"

// AuthorizationModeOTP const
const AuthorizationModeOTP AuthorizationModeCode = "OTP"

// WalletTypeCode type
type WalletTypeCode string

// WalletTypeMerchantCollection const
const WalletTypeMerchantCollection WalletTypeCode = "MERCHANT_COLLECTION"

// WalletTypeMerchantPayout const
const WalletTypeMerchantPayout WalletTypeCode = "MERCHANT_PAYOUT"

// LanguageCode type (Accept-Language header)
type LanguageCode string

// LanguageEnglish const
const LanguageEnglish LanguageCode = "en"

// LanguageFrench const
const LanguageFrench LanguageCode = "fr"

// LanguageSwahili const
const LanguageSwahili LanguageCode = "sw"

// TransactionMethodCode type
type TransactionMethodCode string

// TransactionMethodMobileMoney type
const TransactionMethodMobileMoney TransactionMethodCode = "MOBILE_MONEY"

// TransactionMethodCard type
const TransactionMethodCard TransactionMethodCode = "CARD"

// TransactionMethodBank type
const TransactionMethodBank TransactionMethodCode = "BANK"

// TransactionMethodCrypto type
const TransactionMethodCrypto TransactionMethodCode = "CRYPTO"

// CountryCode type
type CountryCode string

// CountryCodeUganda type
const CountryCodeUganda CountryCode = "UG"

// CountryCodeKenya type
const CountryCodeKenya CountryCode = "KE"

// CountryCodeTanzania type
const CountryCodeTanzania CountryCode = "TZ"

// CountryCodeRwanda type
const CountryCodeRwanda CountryCode = "RW"

// CountryCodeBurundi type
const CountryCodeBurundi CountryCode = "BI"

// CountryCodeGhana type
const CountryCodeGhana CountryCode = "GH"

// CountryCodeCameroon type
const CountryCodeCameroon CountryCode = "CM"

// CountryCodeSouthAfrica type
const CountryCodeSouthAfrica CountryCode = "ZA"

// CountryCodeNigeria type
const CountryCodeNigeria CountryCode = "NG"

// CountryCodeZambia type
const CountryCodeZambia CountryCode = "ZM"

// CountryCodeUSA type
const CountryCodeUSA CountryCode = "US"

// CountryCodeUnitedKingdom type
const CountryCodeUnitedKingdom CountryCode = "GB"

// CountryCodeEurope type
const CountryCodeEurope CountryCode = "EU"

// CurrencyCode type
type CurrencyCode string

// CurrencyCodeUGX type
const CurrencyCodeUGX CurrencyCode = "UGX"

// CurrencyCodeKES type
const CurrencyCodeKES CurrencyCode = "KES"

// CurrencyCodeTZS type
const CurrencyCodeTZS CurrencyCode = "TZS"

// CurrencyCodeRWF type
const CurrencyCodeRWF CurrencyCode = "RWF"

// CurrencyCodeBIF type
const CurrencyCodeBIF CurrencyCode = "BIF"

// CurrencyCodeGHS type
const CurrencyCodeGHS CurrencyCode = "GHS"

// CurrencyCodeXAF type
const CurrencyCodeXAF CurrencyCode = "XAF"

// CurrencyCodeZAR type
const CurrencyCodeZAR CurrencyCode = "ZAR"

// CurrencyCodeNGN type
const CurrencyCodeNGN CurrencyCode = "NGN"

// CurrencyCodeZMW type
const CurrencyCodeZMW CurrencyCode = "ZMW"

// CurrencyCodeUSD type
const CurrencyCodeUSD CurrencyCode = "USD"

// CurrencyCodeGBP type
const CurrencyCodeGBP CurrencyCode = "GBP"

// CurrencyCodeEUR type
const CurrencyCodeEUR CurrencyCode = "EUR"

// transformStructToMap method
func transformStructToMap(st interface{}) (map[string]interface{}, error) {
	bytes, err := json.Marshal(st)
	if err != nil {
		return nil, err
	}

	jsonMap := make(map[string]interface{})

	err = json.Unmarshal(bytes, &jsonMap)
	if err != nil {
		return nil, err
	}

	return jsonMap, nil
}
