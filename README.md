# scorum/scorum-go
[![Go Report Card](https://goreportcard.com/badge/github.com/scorum/scorum-go)](https://goreportcard.com/report/github.com/scorum/scorum-go)
[![GoDoc](https://godoc.org/github.com/scorum/scorum-go?status.svg)](https://godoc.org/github.com/scorum/scorum-go)
[![Build Status](https://travis-ci.org/scorum/scorum-go.svg?branch=master)](https://travis-ci.org/scorum/scorum-go)

Golang RPC client library for [Scorum](https://scorumcoins.com). Both http and websocket transports are supported.
The websocket one allows to set callbacks.

## Usage

```go
import "github.com/scorum/scorum-go"
```

## Example
```go
import (
    scorum "github.com/scorum/scorum-go"
    "github.com/scorum/scorum-go/rpc"
)

const testNet = "https://testnet.scorum.com"

// create client
transport := rpc.NewHTTPTransport(testNet)
client := NewClient(transport)

// get last 100 transactions of the particular account
history, _ := client.AccountHistory.GetAccountHistory("acc1", -1, 100)

for seq, trx := range history {
    for _, op := range trx.Operations {
        switch body := op.(type) {
        case *types.TransferOperation:
            // do something with the incoming transaction
        }
    }
}

```

