# Dusupay API SDK GO (Unofficial)
[![Tests](https://github.com/Kachit/dusupay-sdk-go/actions/workflows/tests.yml/badge.svg?branch=master)](https://github.com/Kachit/dusupay-sdk-go/actions/workflows/tests.yml)
[![Codecov](https://codecov.io/gh/Kachit/dusupay-sdk-go/branch/master/graph/badge.svg)](https://codecov.io/gh/Kachit/dusupay-sdk-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/kachit/dusupay-sdk-go)](https://goreportcard.com/report/github.com/kachit/dusupay-sdk-go)
[![Version](https://img.shields.io/github/go-mod/go-version/Kachit/dusupay-sdk-go)](https://go.dev/doc/go1.24)
[![Release](https://img.shields.io/github/v/release/Kachit/dusupay-sdk-go.svg)](https://github.com/Kachit/dusupay-sdk-go/releases)
[![License](https://img.shields.io/github/license/mashape/apistatus.svg)](https://github.com/kachit/dusupay-sdk-go/blob/master/LICENSE)
[![GoDoc](https://pkg.go.dev/badge/github.com/kachit/dusupay-sdk-go/v2)](https://pkg.go.dev/github.com/kachit/dusupay-sdk-go/v2)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go#third-party-apis)

## Description
Unofficial Dusupay payment gateway API Client for Go

## API documentation
https://developer.dusupay.com/

## Installation
```shell
go get -u github.com/kachit/dusupay-sdk-go/v2
```

## Usage
```go
package main

import (
    "fmt"
    "context"
    dusupay "github.com/kachit/dusupay-sdk-go/v2"
)

func main(){
    cfg := dusupay.NewConfig("Your public key", "Your secret key")
    client, err := dusupay.NewClientFromConfig(cfg, nil)
    if err != nil {
        panic(err)
    }
}
```

### Get wallet balances
```go
ctx := context.Background()
// filter is optional, pass nil to get all wallets
balances, response, err := client.Merchants().GetBalances(ctx, &dusupay.BalancesFilter{Currency: dusupay.CurrencyCodeUGX})
if err != nil {
    panic(err)
}
fmt.Println(balances.Data.MerchantWallets[0].Currency)
fmt.Println(balances.Data.MerchantWallets[0].AvailableBalance)
```

### Verify transaction
```go
ctx := context.Background()
result, response, err := client.Merchants().Verify(ctx, "DUSUPAY5FNZCVUKZ8C0KZE")
if err != nil {
    panic(err)
}
fmt.Println(result.Data.InternalReference)
fmt.Println(result.Data.TransactionStatus)
```

### Get payment providers list
```go
ctx := context.Background()
filter := &dusupay.ProvidersFilter{Currency: dusupay.CurrencyCodeUGX, TransactionType: dusupay.TransactionTypeCollection}
providers, response, err := client.Providers().GetList(ctx, filter)
if err != nil {
    panic(err)
}
fmt.Println(providers.Data.PaymentProviders[0].ProviderCode)
fmt.Println(providers.Data.PaymentProviders[0].ProviderName)
```

### Get payout banks list
```go
ctx := context.Background()
filter := &dusupay.BanksFilter{ProviderCode: "bank_ng"}
banks, response, err := client.Banks().GetList(ctx, filter)
if err != nil {
    panic(err)
}
fmt.Println(banks.Data.PayoutBanks[0].BankName)
fmt.Println(banks.Data.PayoutBanks[0].BankCode)
```

### Create collection request
```go
ctx := context.Background()
request := &dusupay.CollectionRequest{
    Currency:          dusupay.CurrencyCodeUGX,
    Amount:            10000,
    TransactionMethod: dusupay.TransactionMethodMobileMoney,
    ProviderCode:      "airtel_ug",
    MerchantReference: "1234567891",
    Description:       "payment description",
    Msisdn:            "256752000123",
}
result, response, err := client.Collections().Create(ctx, request)
if err != nil {
    panic(err)
}
fmt.Println(result.Data.InternalReference)
fmt.Println(result.Data.MerchantReference)
fmt.Println(result.Data.PaymentURL)
```

### Create hosted page collection request
```go
request := &dusupay.CollectionRequest{
    Currency:          dusupay.CurrencyCodeUGX,
    Amount:            10000,
    TransactionMethod: dusupay.TransactionMethodMobileMoney,
    ProviderCode:      "mtn_ug",
    MerchantReference: "auto",
    Description:       "payment description",
    Msisdn:            "256777000001",
    MobileMoneyHPP:    true,
    RedirectURL:       "https://your-redirect-url",
}
result, response, err := client.Collections().Create(ctx, request)
// redirect the customer to result.Data.PaymentURL
```

### Create direct card (S2S) collection request
```go
encryptor, _ := dusupay.NewCardEncryptor(cfg.DusupayPublicKey())
cipher, err := encryptor.Encrypt(&dusupay.CardData{
    FirstName: "John",
    LastName:  "Doe",
    CardNo:    "4622000000005678",
    ExpMonth:  "06",
    ExpYear:   "22",
    CVV:       "123",
})
request := &dusupay.CollectionRequest{
    Currency:          dusupay.CurrencyCodeNGN,
    Amount:            1000,
    TransactionMethod: dusupay.TransactionMethodCard,
    ProviderCode:      "local_ngn",
    MerchantReference: "auto",
    Description:       "payment description",
    CardCipher:        cipher,
    RedirectURL:       "https://your-redirect-url",
}
result, response, err := client.Collections().Create(ctx, request)
```

### Authorize collection with OTP
```go
request := &dusupay.CollectionAuthorizeRequest{InternalReference: "DUSUPAYSTPDKZ6MLY5WPQ", OTP: "123456"}
result, response, err := client.Collections().Authorize(ctx, request)
```

### Create payout request
```go
ctx := context.Background()
request := &dusupay.PayoutRequest{
    Currency:          dusupay.CurrencyCodeUGX,
    Amount:            10000,
    TransactionMethod: dusupay.TransactionMethodMobileMoney,
    ProviderCode:      "airtel_ug",
    MerchantReference: "1234567892",
    Description:       "payout description",
    AccountNumber:     "256752000123",
    CustomerName:      "Foo Bar",
}
result, response, err := client.Payouts().Create(ctx, request)
if err != nil {
    panic(err)
}
fmt.Println(result.Data.InternalReference)
fmt.Println(result.Data.MerchantReference)
```

### Create bank payout request
```go
request := &dusupay.PayoutRequest{
    Currency:          dusupay.CurrencyCodeNGN,
    Amount:            4000,
    TransactionMethod: dusupay.TransactionMethodBank,
    ProviderCode:      "bank_ng",
    MerchantReference: "1234567893",
    Description:       "payout description",
    AccountNumber:     "2121562123",
    CustomerName:      "Foo Bar",
    ExtraParams:       &dusupay.PayoutExtraParams{BankCode: "access_bank_ng"},
}
result, response, err := client.Payouts().Create(ctx, request)
```

### Handle webhook callback and verify RSA signature
Callbacks are sent as `{"event": "...", "payload": {...}}`. Use `dusupay.PayoutCallback` for payout callbacks.
```go
var callback dusupay.CollectionCallback
_ = json.Unmarshal(requestBody, &callback)

validator, _ := dusupay.NewSignatureValidator(cfg.DusupayPublicKey())
err := validator.ValidateSignature(&callback, r.Header.Get("rsa-signature"))

if callback.Event == dusupay.CallbackEventTransactionProcessing &&
    callback.Payload.Authorization != nil &&
    callback.Payload.Authorization.Mode == dusupay.AuthorizationModeOTP {
    // collect OTP from the customer and call client.Collections().Authorize
}
```

### Verify HMAC webhook signature
```go
var callback dusupay.CollectionCallback
_ = json.Unmarshal(requestBody, &callback)

validator := dusupay.NewHMACSignatureValidator("your-signing-key")
err := validator.ValidateSignature(&callback, r.Header.Get("hmac-signature"))
```

### Verify redirect signature (hosted pages, cards)
```go
payload := dusupay.NewRedirectPayload(r.URL.Query())

validator := dusupay.NewHMACSignatureValidator("your-signing-key")
err := validator.ValidateSignature(payload, payload.HMACSignature)
```
