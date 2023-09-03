Returns string value by passed table and function.

### Availiable tables
* `block`
* `tx`
* `message`
* `event`


### Availiable functions
* `sum`
* `min`
* `max`
* `avg`
* `count`


`Column` query parameter is required for functions `sum`, `min`, `max` and `avg` and should not pass for `count`.


###  Availiable columns and functions for tables:

#### Block
* `height`         -- min max
* `time`           -- min max
* `tx_count`       -- min max sum avg
* `events_count`   -- min max sum avg
* `blobs_size`     -- min max sum avg
* `fee`            -- min max sum avg

#### Tx
* `height`         -- min max
* `time`           -- min max
* `gas_wanted`     -- min max sum avg
* `gas_used`       -- min max sum avg
* `timeout_height` -- min max avg
* `events_count`   -- min max sum avg
* `messages_count` -- min max sum avg
* `fee`            -- min max sum avg

#### Event
* `height`         -- min max
* `time`           -- min max

#### Message
* `height`         -- min max
* `time`           -- min max