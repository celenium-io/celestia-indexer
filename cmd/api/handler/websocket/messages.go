package websocket

// methods
const (
	MethodSubscribe   = "subscribe"
	MethodUnsubscribe = "unsubscribe"
)

// channels
const (
	ChannelHead = "head"
	ChannelTx   = "tx"
)

type Message struct {
	Method string `json:"method" validate:"required,oneof=subscribe,unsubscribe"`
	Body   []byte `json:"body"   validate:"required"`
}

type Subscribe struct {
	Channel string `json:"channel" validate:"required,oneof=head"`
	Filters []byte `json:"filters" validate:"required"`
}

type Unsubscribe struct {
	Channel string `json:"channel" validate:"required,oneof=head"`
}

type TransactionFilters struct {
	Status []string `json:"status,omitempty"`
}
