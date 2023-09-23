package storage

import (
	"context"
	"io"

	"github.com/dipdup-io/celestia-indexer/pkg/types"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/lib/pq"
)

var Models = []any{
	&State{},
	&Constant{},
	&DenomMetadata{},
	&Balance{},
	&Address{},
	&Block{},
	&BlockStats{},
	&Tx{},
	&Message{},
	&Event{},
	&Namespace{},
	&NamespaceMessage{},
	&Signer{},
	&MsgAddress{},
	&Validator{},
}

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type Notificator interface {
	Notify(ctx context.Context, channel string, payload string) error
}

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type Listener interface {
	io.Closer

	Subscribe(ctx context.Context, channels ...string) error
	Listen() chan *pq.Notification
}

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type ListenerFactory interface {
	CreateListener() Listener
}

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type Transaction interface {
	sdk.Transaction

	SaveConstants(ctx context.Context, constants ...Constant) error
	SaveTransactions(ctx context.Context, txs ...Tx) error
	SaveNamespaces(ctx context.Context, namespaces ...*Namespace) (uint64, error)
	SaveAddresses(ctx context.Context, addresses ...*Address) (uint64, error)
	SaveBalances(ctx context.Context, balances ...Balance) error
	SaveMessages(ctx context.Context, msgs ...*Message) error
	SaveSigners(ctx context.Context, addresses ...Signer) error
	SaveMsgAddresses(ctx context.Context, addresses ...MsgAddress) error
	SaveNamespaceMessage(ctx context.Context, nsMsgs ...NamespaceMessage) error
	SaveValidators(ctx context.Context, validators ...*Validator) error
	SaveEvents(ctx context.Context, events ...Event) error
	LastBlock(ctx context.Context) (block Block, err error)
	State(ctx context.Context, name string) (state State, err error)
	Namespace(ctx context.Context, id uint64) (ns Namespace, err error)

	RollbackBlock(ctx context.Context, height types.Level) error
	RollbackBlockStats(ctx context.Context, height types.Level) (stats BlockStats, err error)
	RollbackAddresses(ctx context.Context, height types.Level) (address []Address, err error)
	RollbackTxs(ctx context.Context, height types.Level) (txs []Tx, err error)
	RollbackEvents(ctx context.Context, height types.Level) (events []Event, err error)
	RollbackMessages(ctx context.Context, height types.Level) (msgs []Message, err error)
	RollbackNamespaceMessages(ctx context.Context, height types.Level) (msgs []NamespaceMessage, err error)
	RollbackNamespaces(ctx context.Context, height types.Level) (ns []Namespace, err error)
	RollbackValidators(ctx context.Context, height types.Level) (err error)
	RollbackSigners(ctx context.Context, txIds []uint64) (err error)
	RollbackMessageAddresses(ctx context.Context, msgIds []uint64) (err error)
	DeleteBalances(ctx context.Context, ids []uint64) error
	LastAddressAction(ctx context.Context, address []byte) (uint64, error)
}

const (
	ChannelHead = "head"
	ChannelTx   = "tx"
)
