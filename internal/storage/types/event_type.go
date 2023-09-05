package types

// EventType -
type EventType string

const (
	EventTypeUnknown            EventType = "unknown"
	EventTypeCoinReceived       EventType = "coin_received"
	EventTypeCoinbase           EventType = "coinbase"
	EventTypeCoinSpent          EventType = "coin_spent"
	EventTypeBurn               EventType = "burn"
	EventTypeMint               EventType = "mint"
	EventTypeMessage            EventType = "message"
	EventTypeProposerReward     EventType = "proposer_reward"
	EventTypeRewards            EventType = "rewards"
	EventTypeCommission         EventType = "commission"
	EventTypeLiveness           EventType = "liveness"
	EventTypeAttestationRequest EventType = "AttestationRequest"
	EventTypeTransfer           EventType = "transfer"
	EventTypePayForBlobs        EventType = "celestia.blob.v1.EventPayForBlobs"
	EventTypeRedelegate         EventType = "redelegate"
	EventTypeWithdrawRewards    EventType = "withdraw_rewards"
	EventTypeWaithdrawComission EventType = "withdraw_commission"
	EventTypeCreateValidator    EventType = "create_validator"
	EventTypeDelegate           EventType = "delegate"
	EventTypeEditValidator      EventType = "edit_validator"
	EventTypeUnbond             EventType = "unbond"
	EventTypeTx                 EventType = "tx"
	EventTypeUseFeegrant        EventType = "use_feegrant"
	EventTypeRevokeFeegrant     EventType = "revoke_feegrant"
	EventTypeSetFeegrant        EventType = "set_feegrant"
	EventTypeUpdateFeegrant     EventType = "update_feegrant"
)
