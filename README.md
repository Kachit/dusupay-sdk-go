# Dusupay API SDK GO (Unofficial)
[![Build Status](https://travis-ci.org/Kachit/dusupay-sdk-go.svg?branch=master)](https://travis-ci.org/Kachit/dusupay-sdk-go)
[![codecov](https://codecov.io/gh/Kachit/dusupay-sdk-go/branch/master/graph/badge.svg)](https://codecov.io/gh/Kachit/dusupay-sdk-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/kachit/dusupay-sdk-go)](https://goreportcard.com/report/github.com/kachit/dusupay-sdk-go)
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
