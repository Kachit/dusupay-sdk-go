# Dusupay API SDK GO (Unofficial)
[![Build Status](https://app.travis-ci.com/Kachit/dusupay-sdk-go.svg?branch=master)](https://app.travis-ci.com/github/Kachit/dusupay-sdk-go)
[![Codecov](https://codecov.io/gh/Kachit/dusupay-sdk-go/branch/master/graph/badge.svg)](https://codecov.io/gh/Kachit/dusupay-sdk-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/kachit/dusupay-sdk-go)](https://goreportcard.com/report/github.com/kachit/dusupay-sdk-go)
[![Version](https://img.shields.io/github/go-mod/go-version/Kachit/dusupay-sdk-go)](https://go.dev/doc/go1.14)
[![Release](https://img.shields.io/github/v/release/Kachit/dusupay-sdk-go.svg)](https://github.com/Kachit/dusupay-sdk-go/releases)
[![License](https://img.shields.io/github/license/mashape/apistatus.svg)](https://github.com/kachit/dusupay-sdk-go/blob/master/LICENSE)
[![GoDoc](https://pkg.go.dev/badge/github.com/kachit/dusupay-sdk-go)](https://pkg.go.dev/github.com/kachit/dusupay-sdk-go)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go#third-party-apis)

## Description
Unofficial Dusupay payment gateway API Client for Go

## API documentation
https://docs.dusupay.com

## Installation
```shell
go get -u github.com/kachit/dusupay-sdk-go
```

## Usage
```go
package main

import (
    "fmt"
    "context"
    dusupay "github.com/kachit/dusupay-sdk-go"
)

func main(){
    // Create a client instance
    cfg := dusupay.NewConfig("Your public key", "Your secret key")
    client, err := dusupay.NewClientFromConfig(cfg, nil)
    if err != nil {
        fmt.Printf("config parameter error " + err.Error())
        panic(err)
    }
}
```
### Get balances list
```go
ctx := context.Background()
balances, response, err := client.Merchants().GetBalances(ctx)

if err != nil {
    fmt.Printf("Wrong API request " + err.Error())
    panic(err)
}

//Dump raw response
fmt.Println(response)

//Dump result
fmt.Println(balances.Status)
fmt.Println(balances.Code)
fmt.Println(balances.Message)
fmt.Println((*balances.Data)[0].Currency)
fmt.Println((*balances.Data)[0].Balance)
```

### Get banks list
```go
ctx := context.Background()
filter := &dusupay.BanksFilter{TransactionType: dusupay.TransactionTypePayout, Country: dusupay.CountryCodeGhana}
banks, response, err := client.Banks().GetList(ctx, filter)

if err != nil {
    fmt.Printf("Wrong API request " + err.Error())
    panic(err)
}

//Dump raw response
fmt.Println(response)

//Dump result
fmt.Println(banks.Status)
fmt.Println(banks.Code)
fmt.Println(banks.Message)
fmt.Println((*banks.Data)[0].Id)
fmt.Println((*banks.Data)[0].Name)
fmt.Println((*banks.Data)[0].BankCode)
```

### Get banks branches list
```go
ctx := context.Background()
filter := &dusupay.BanksBranchesFilter{Country: dusupay.CountryCodeGhana, Bank: "BankCode"}
branches, response, err := client.Banks().GetBranchesList(ctx, filter)

if err != nil {
    fmt.Printf("Wrong API request " + err.Error())
    panic(err)
}

//Dump raw response
fmt.Println(response)

//Dump result
fmt.Println(branches.Status)
fmt.Println(branches.Code)
fmt.Println(branches.Message)
fmt.Println((*branches.Data)[0].Name)
fmt.Println((*branches.Data)[0].Code)
```

### Get providers list
```go
ctx := context.Background()
filter := &dusupay.ProvidersFilter{Country: dusupay.CountryCodeUganda, Method: dusupay.TransactionMethodMobileMoney, TransactionType: dusupay.TransactionTypeCollection}
providers, response, err := client.Providers().GetList(ctx, filter)

if err != nil {
    fmt.Printf("Wrong API request " + err.Error())
    panic(err)
}

//Dump raw response
fmt.Println(response)

//Dump result
fmt.Println(providers.Status)
fmt.Println(providers.Code)
fmt.Println(providers.Message)
fmt.Println((*providers.Data)[0].ID)
fmt.Println((*providers.Data)[0].Name)
```

### Create collection request
```go
ctx := context.Background()
request := &dusupay.CollectionRequest{
    Currency:          dusupay.CurrencyCodeUGX,
    Amount:            10000,
    Method:            dusupay.TransactionMethodMobileMoney,
    ProviderId:        "airtel_ug",
    MerchantReference: "1234567891",
    RedirectUrl:       "http://foo.bar",
    Narration:         "narration",
    AccountNumber:         "256752000123",
    MobileMoneyHpp:         true,
}
result, response, err := client.Collections().Create(ctx, request)

if err != nil {
    fmt.Printf("Wrong API request " + err.Error())
    panic(err)
}

//Dump raw response
fmt.Println(response)

//Dump result
fmt.Println(result.Status)
fmt.Println(result.Code)
fmt.Println(result.Message)
fmt.Println((*result.Data).ID)
fmt.Println((*result.Data).PaymentURL)
```

### Create payout request
```go
ctx := context.Background()
request := &dusupay.PayoutRequest{
    Currency:          dusupay.CurrencyCodeUGX,
    Amount:            10000,
    Method:            dusupay.TransactionMethodMobileMoney,
    ProviderId:        "airtel_ug",
    MerchantReference: "1234567892",
    Narration:         "narration",
    AccountNumber:         "256752000123",
    AccountName:         "Foo Bar",
}
result, response, err := client.Payouts().Create(ctx, request)

if err != nil {
    fmt.Printf("Wrong API request " + err.Error())
    panic(err)
}

//Dump raw response
fmt.Println(response)

//Dump result
fmt.Println(result.Status)
fmt.Println(result.Code)
fmt.Println(result.Message)
fmt.Println((*result.Data).ID)
```

### Create refund request
```go
ctx := context.Background()
request := &dusupay.RefundRequest{
    Amount:            100,
    InternalReference:            "DUSUPAY5FNZCVUKZ8C0KZE",
}
result, response, err := client.Refunds().Create(ctx, request)

if err != nil {
    fmt.Printf("Wrong API request " + err.Error())
    panic(err)
}

//Dump raw response
fmt.Println(response)

//Dump result
fmt.Println(result.Status)
fmt.Println(result.Code)
fmt.Println(result.Message)
fmt.Println((*result.Data).ID)
```

### Verify webhook signature
```go
requestPayload := `
{
    "id": 226,
    "request_amount": 10,
    "request_currency": "USD",
    "account_amount": 737.9934,
    "account_currency": "UGX",
    "transaction_fee": 21.4018,
    "total_credit": 716.5916,
    "customer_charged": false,
    "provider_id": "mtn_ug",
    "merchant_reference": "76859aae-f148-48c5-9901-2e474cf19b71",
    "internal_reference": "DUSUPAY405GZM1G5JXGA71IK",
    "transaction_status": "COMPLETED",
    "transaction_type": "collection",
    "message": "Transaction Completed Successfully"
}
`
requestUri := "https://www.sample-url.com/callback"
signature := "value from 'dusupay-signature' http header"

var webhook dusupay.CollectionWebhook
_ = json.Unmarshal(requestPayload, &webhook)

rawBytes, _ := ioutil.ReadFile("path/to/dusupay-public-key.pem")

validator, _ := dusupay.NewSignatureValidator(rawBytes)
err := validator.ValidateSignature(webhook, requestUri, signature)
```
