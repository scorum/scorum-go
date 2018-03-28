# Payments processing

The example represent a payment system and consists of two main parts:

* Incoming transfers monitor
  * Pulls history operations each 10 seconds
  * Sort the operations and filters only transfers
  * Increase deposits balances with the transferred amounts
* Payout
  * Choose a random deposit
  * Make a transfer transaction
  * Updates deposits balances with the transferred amounts
