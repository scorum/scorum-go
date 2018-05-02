# Deposits

The example represent a deposits system and consists of two main parts:

* Incoming transfers monitor
  * Pull history operations each 10 seconds
  * Take only operations created in blocks with number <= LastIrreversibleBlockNum
  * Increase deposits balances with the transferred amounts

* Payout
  * Choose a random deposit
  * Make a transfer transaction
  * Update deposits balances with the transferred amounts

![alt monitor transfers](https://github.com/scorum/scorum-go/blob/master/examples/deposits/diagrams/monitor.png "Monitor transfers")
![alt broadcast transfers](https://github.com/scorum/scorum-go/blob/master/examples/deposits/diagrams/broadcast.png "Broadcast transfers")