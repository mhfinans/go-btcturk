# go-btcturk

golang client for btcturk api

<div align="center">

[![Go Report Card](https://goreportcard.com/badge/github.com/aliereno/go-btcturk)](https://goreportcard.com/report/github.com/aliereno/go-btcturk)
[![HitCount](http://hits.dwyl.com/aliereno/go-btcturk.svg)](http://hits.dwyl.com/aliereno/go-btcturk)

</div>

#### This repo has been forked from [vural/go-btcturk](https://github.com/vural/go-btcturk) and has been improved for Btcturk's new api. So if you like this repo please give a star actual [repo](https://github.com/vural/go-btcturk).

## documentation

the documentation is available on [godoc](http://godoc.org/github.com/aliereno/go-btcturk/btcturk)

## install

```sh
go get -u github.com/aliereno/go-btcturk/btcturk
```

## usage
```go
package main

import (
	"github.com/aliereno/go-btcturk/btcturk"
)

func main() {
    api := btcturk.NewBTCTurkClient()
    t, err := api.Ticker()
    if err != nil {
        print(err)
        os.Exit(1)
    }
    
    for _, v := range t {
        println(v.Ask)
    }

    // if you don't plan to call authenticated api methods. SetAuthKey not required.
    api.SetAuthKey("publicKey", "privateKey")
}

```

## Passing params

[Endpoint Params](https://github.com/aliereno/go-btcturk/blob/master/btcturk/params.go)
```go
package main

import (
	"github.com/aliereno/btcturk/btcturk"
)

func main() {
    api := btcturk.NewBTCTurkClient()
    api.SetAuthKey("publicKey", "privateKey")

    api.Quantity(0.001).
        Price(50000).
        StopPrice(0).
        OrderMethod("limit").
        PairSymbol(btcturk.BTCTRY).
        Buy()
}

```

## Notes
 - you can get your private/public key peer from your account

**[BTCTurk API documentation](https://docs.btcturk.com)**
