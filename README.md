# Dusupay API SDK GO (Unofficial)
[![Build Status](https://travis-ci.org/Kachit/dusupay-sdk-go.svg?branch=master)](https://travis-ci.org/Kachit/dusupay-sdk-go)
[![codecov](https://codecov.io/gh/Kachit/dusupay-sdk-go/branch/master/graph/badge.svg)](https://codecov.io/gh/Kachit/dusupay-sdk-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/kachit/dusupay-sdk-go)](https://goreportcard.com/report/github.com/kachit/dusupay-sdk-go)
[![Release](https://img.shields.io/github/v/release/Kachit/dusupay-sdk-go.svg)](https://github.com/Kachit/dusupay-sdk-go/releases)
[![License](https://img.shields.io/github/license/mashape/apistatus.svg)](https://github.com/kachit/dusupay-sdk-go/blob/master/LICENSE)

## Description
Unofficial Dusupay payment gateway API Client for Go

## API documentation
https://docs.dusupay.com/

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
	cfg := dusupay.NewConfig("Your public key", "Your secret key")
    client, err := dusupay.NewClientFromConfig(cfg, nil)
    if err != nil {
        fmt.Printf("config parameter error " + err.Error())
    }

    ctx := context.Background()
    response, err := client.Merchants().GetBalances(ctx)   
}
```
