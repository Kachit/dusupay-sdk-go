package dusupay

import "encoding/json"

type TransactionStatusCode string

const TransactionStatusPending TransactionStatusCode = "PENDING"
const TransactionStatusFailed TransactionStatusCode = "FAILED"
const TransactionStatusCompleted TransactionStatusCode = "COMPLETED"
const TransactionStatusCancelled TransactionStatusCode = "CANCELLED"

type TransactionMethodCode string

const TransactionMethodMobileMoney TransactionMethodCode = "MOBILE_MONEY"
const TransactionMethodCard TransactionMethodCode = "CARD"
const TransactionMethodBank TransactionMethodCode = "BANK"
const TransactionMethodCrypto TransactionMethodCode = "CRYPTO"

type CountryCode string

const CountryCodeUganda CountryCode = "UG"
const CountryCodeKenya CountryCode = "KE"
const CountryCodeTanzania CountryCode = "TZ"
const CountryCodeRwanda CountryCode = "RW"
const CountryCodeBurundi CountryCode = "BI"
const CountryCodeGhana CountryCode = "GH"
const CountryCodeCameroon CountryCode = "CM"
const CountryCodeSouthAfrica CountryCode = "ZA"
const CountryCodeNigeria CountryCode = "NG"
const CountryCodeZambia CountryCode = "ZM"
const CountryCodeUSA CountryCode = "US"
const CountryCodeUnitedKingdom CountryCode = "GB"
const CountryCodeEurope CountryCode = "EU"

type CurrencyCode string

const CurrencyCodeUGX CurrencyCode = "UGX"
const CurrencyCodeKES CurrencyCode = "KES"
const CurrencyCodeTZS CurrencyCode = "TZS"
const CurrencyCodeRWF CurrencyCode = "RWF"
const CurrencyCodeBIF CurrencyCode = "BIF"
const CurrencyCodeGHS CurrencyCode = "GHS"
const CurrencyCodeXAF CurrencyCode = "XAF"
const CurrencyCodeZAR CurrencyCode = "ZAR"
const CurrencyCodeNGN CurrencyCode = "NGN"
const CurrencyCodeZMW CurrencyCode = "ZMW"
const CurrencyCodeUSD CurrencyCode = "USD"
const CurrencyCodeGBP CurrencyCode = "GBP"
const CurrencyCodeEUR CurrencyCode = "EUR"

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
