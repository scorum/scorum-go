# Payments processing

The example represent a payment system and consists of two main parts:

* Incoming transfers monitor
⋅⋅1. Pulls history operations each 10 seconds
⋅⋅2. Sort the operations and filters only transfers
⋅⋅3. Increase deposits balances with the transferred amounts
* Payout
⋅⋅1. Choose a random deposit
⋅⋅2. Make a transfer transaction
⋅⋅3. Updates deposits balances with the transferred amounts