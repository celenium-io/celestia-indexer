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

Now 2 channels are supported:

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
