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


### Distribution

> IMPORTANT: `valoper` address does not have balance. Only delegator has balance.

* `proposer_reward`

Example of event:

```json
{
    "type": "proposer_reward",
    "attributes": [
        {
            "key": "amount",
            "value": null // Amount nullable
        },
        {
            "key": "validator",
            "value": "celestiavaloper1jwzamm3ltkzce7ey5tn7uadt8uxg6k89a9tj94"
        }
    ]
}
```

Building balance updates:

```json
// Map event to balance updates
[{
    "address": "delegator_address_by_validator_value",
    "amount": "amount_value"
}]
```

* `rewards`

Example of event:

```json
{
    "type": "rewards",
    "attributes": [
        {
            "key": "amount",
            "value": "9.272443826506892545utia"
        },
        {
            "key": "validator",
            "value": "celestiavaloper189vk0yl8ce5wfh6h36hmgvjlmwrz2sgl5q5zp6"
        }
    ]
}
```

Building balance updates:

```json
// Map event to balance updates
[{
    "address": "delegator_address_by_validator_value",
    "amount": "amount_value"
}]
```

* `commision`

Example of event:

```json
{
    "type": "commission",
    "attributes": [
        {
            "key": "amount",
            "value": "0.927244382650689254utia"
        },
        {
            "key": "validator",
            "value": "celestiavaloper189vk0yl8ce5wfh6h36hmgvjlmwrz2sgl5q5zp6"
        }
    ]
}
```

Building balance updates:

```json
// Map event to balance updates
[{
    "address": "delegator_address_by_validator_value",
    "amount": "amount_value"
}]
```

### Minting

* `coinbase`

Example of event:

```json
 {
    "type": "coinbase",
    "attributes": [
        {
            "key": "minter",
            "value": "celestia1m3h30wlvsf8llruxtpukdvsy0km2kum8emkgad"
        },
        {
            "key": "amount",
            "value": "30862303utia"
        }
    ]
}
```

Building balance updates:

```json
// Map event to balance updates
[{
    "address": "minter_value",
    "amount": "amount_value"
}]

// Example
[{
    "address": "celestia1m3h30wlvsf8llruxtpukdvsy0km2kum8emkgad",
    "amount": "30862303utia"
}]
```