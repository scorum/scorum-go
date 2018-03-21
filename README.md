# scorum/scorum-go
[![GoDoc](https://godoc.org/github.com/scorum/scorum-go?status.svg)](https://godoc.org/github.com/scorum/scorum-go)

Golang RPC client library for [Scorum](https://scorumcoins.com).

## Usage

```go
import "github.com/scorum/scorum-go"
```

## Example
```go
import scorum "github.com/scorum/scorum-go"

const testNet = "http://blockchain.scorum.com:8003"

client := scorum.NewClient(caller.NewHttpCaller(testNet))

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

