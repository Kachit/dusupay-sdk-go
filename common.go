package dusupay

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

const CountryCodeUGX CurrencyCode = "UGX"
const CountryCodeKES CurrencyCode = "KES"
const CountryCodeTZS CurrencyCode = "TZS"
const CountryCodeRWF CurrencyCode = "RWF"
const CountryCodeBIF CurrencyCode = "BIF"
const CountryCodeGHS CurrencyCode = "GHS"
const CountryCodeXAF CurrencyCode = "XAF"
const CountryCodeZAR CurrencyCode = "ZAR"
const CountryCodeNGN CurrencyCode = "NGN"
const CountryCodeZMW CurrencyCode = "ZMW"
const CountryCodeUSD CurrencyCode = "USD"
const CountryCodeGBP CurrencyCode = "GBP"
const CountryCodeEUR CurrencyCode = "EUR"
