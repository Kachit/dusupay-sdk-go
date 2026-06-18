# Dusupay API SDK GO (Unofficial)
[![Go Report Card](https://goreportcard.com/badge/github.com/kachit/dusupay-sdk-go)](https://goreportcard.com/report/github.com/kachit/dusupay-sdk-go)
[![Release](https://img.shields.io/github/v/release/Kachit/dusupay-sdk-go.svg)](https://github.com/Kachit/dusupay-sdk-go/releases)
[![License](https://img.shields.io/github/license/mashape/apistatus.svg)](https://github.com/kachit/dusupay-sdk-go/blob/master/LICENSE)
[![GoDoc](https://pkg.go.dev/badge/github.com/kachit/dusupay-sdk-go/v2)](https://pkg.go.dev/github.com/kachit/dusupay-sdk-go/v2)

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
balances, response, err := client.Merchants().GetBalances(ctx)
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

### Verify RSA webhook signature
```go
var payload dusupay.CollectionCallbackPayload
_ = json.Unmarshal(requestBody, &payload)

rawBytes, _ := os.ReadFile("path/to/dusupay-public-key.pem")
validator, _ := dusupay.NewSignatureValidator(rawBytes)
err := validator.ValidateSignature(&payload, signatureHeader)
```

### Verify HMAC webhook signature
```go
var payload dusupay.CollectionCallbackPayload
_ = json.Unmarshal(requestBody, &payload)

validator := dusupay.NewHMACSignatureValidator("your-signing-key")
err := validator.ValidateSignature(&payload, signatureHeader)
```
