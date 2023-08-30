## Documentation for websocket API

### Subscribe

To receive updates from websocket API send `subscribe` request to server.

```json
{
    "method": "subscribe",
    "body": {
        "channel": "<CHANNEL_NAME>",
        "filters": {
            // pass channel filters
        }
    }
}
```

Now 2 channels are supported:

* `head` - receive information about new block. Channel does not have any filters. Subscribe message should looks like:

```json
{
    "method": "subscribe",
    "body": {
        "channel": "head"
    }
}
```

In that channel messages of `responses.Block` type will be sent.

* `tx` - receive information about new transactions. The channel has filters for target receiving information. Now 2 filters are supported:

```json
{
    "method": "subscribe",
    "body": {
        "channel": "tx",
        "filters": {
            "status": [  // array of transaction status. Can be emtpy.
                "success", 
                "failed"
            ],
            "msg_type": [  // array of containing message types status. Can be emtpy.
                "WithdrawValidatorCommission",
                "WithdrawDelegatorReward",
                "EditValidator",
                "BeginRedelegate",
                "CreateValidator",
                "Delegate",
                "Undelegate",
                "Unjail",
                "Send",
                "CreateVestingAccount",
                "CreatePeriodicVestingAccount",
                "PayForBlobs"
            ]
        }
    }
}
```

If all filers are empty subscription to all transaction will be created.

In that channel messages of `responses.Tx` type will be sent.


### Unsubscribe

To unsubscribe send `unsubscribe` message containing one of channel name describing above.


```json
{
    "method": "unsubscribe",
    "body": {
        "channel": "<CHANNEL_NAME>",
    }
}
```
