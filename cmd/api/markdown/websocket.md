## Documentation for websocket API

### Notification

The structure of notification is following in all channels:

```json
{
    "channel": "channel_name",
    "body": "<object or array>"  // depends on channel
}
```

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

Now 3 channels are supported:

* `head` - receive information about indexer state. Channel does not have any filters. Subscribe message should looks like:

```json
{
    "method": "subscribe",
    "body": {
        "channel": "head"
    }
}
```

Notification body of `responses.State` type will be sent to the channel.

* `blocks` - receive information about new blocks. Channel does not have any filters. Subscribe message should looks like:

```json
{
    "method": "subscribe",
    "body": {
        "channel": "blocks"
    }
}
```

Notification body of `responses.Block` type will be sent to the channel.

* `gas_price` - receive information about current gas price. Channel does not have any filters. Subscribe message should looks like:

```json
{
    "method": "subscribe",
    "body": {
        "channel": "gas_price"
    }
}
```

Notification body of `responses.GasPrice` type will be sent to the channel.


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

### Errors

If an incoming message can not be handled (malformed payload, unknown method or unknown channel), the server replies with a notification in the `error` channel. The connection stays open. Internal error details are never exposed; only a stable error code and a generic message are sent.

```json
{
    "channel": "error",
    "body": {
        "code": 2,
        "message": "unknown method"
    }
}
```

Supported error codes:

| Code | Message           | Description                                                          |
|------|-------------------|---------------------------------------------------------------------|
| 1    | `invalid message` | The message could not be parsed (malformed JSON or invalid payload). |
| 2    | `unknown method`  | The `method` field is not `subscribe` or `unsubscribe`.             |
| 3    | `unknown channel` | The requested channel is not one of `head`, `blocks`, `gas_price`.  |
