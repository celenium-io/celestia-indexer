# Balance computing

The text describes algorithm how our indexer computes tia balances of accounts. Updates of account balances are computed from module events which will be described below.

## Events and balance updates

To compute balance updates we need to know message and events connected with it. 

### General case

* `transfer`

> IMPORTANT NOTE: an event type can contain different attributes. Therefore, you should only create balance updates if all required attributes are found 

Example of event:

```json
{
    "type": "transfer",
    "attributes": [
        {
            "key": "recipient",
            "value": "celestia17xpfvakm2amg962yls6f84z3kell8c5lpnjs3s"
        },
        {
            "key": "sender",
            "value": "celestia1vsvx8n7f8dh5udesqqhgrjutyun7zqrgehdq2l"
        },
        {
            "key": "amount",
            "value": "25000utia"
        }
    ]
}
```

Building balance updates:

```json
// Map event to balance updates
[{
    "address": "recipient_value",
    "amount": "amount_value"
},{
    "address": "sender_value",
    "amount": "-amount_value"
}]

// Example
[{
    "address": "celestia17xpfvakm2amg962yls6f84z3kell8c5lpnjs3s",
    "amount": "25000"
},{
    "address": "celestia1vsvx8n7f8dh5udesqqhgrjutyun7zqrgehdq2l",
    "amount": "-25000"
}]
```

> IMPORTANT NOTE: you have to trim suffix `utia` from amount
