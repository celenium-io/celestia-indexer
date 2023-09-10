package types

import (
	"encoding/json"
	"time"

	"github.com/dipdup-io/celestia-indexer/pkg/types"
	"github.com/tendermint/tendermint/libs/bytes"
)

type Genesis struct {
	GenesisTime     time.Time             `json:"genesis_time"`
	ChainID         string                `json:"chain_id"`
	InitialHeight   int64                 `json:"initial_height,string"`
	ConsensusParams types.ConsensusParams `json:"consensus_params"`
	AppHash         bytes.HexBytes        `json:"app_hash"`
	AppState        AppState              `json:"app_state"`
}

type AuthParams struct {
	MaxMemoCharacters      string `json:"max_memo_characters"`
	TxSigLimit             string `json:"tx_sig_limit"`
	TxSizeCostPerByte      string `json:"tx_size_cost_per_byte"`
	SigVerifyCostEd25519   string `json:"sig_verify_cost_ed25519"`
	SigVerifyCostSecp256K1 string `json:"sig_verify_cost_secp256k1"`
}

type BaseAccount struct {
	Address       string      `json:"address"`
	PubKey        interface{} `json:"pub_key"`
	AccountNumber string      `json:"account_number"`
	Sequence      string      `json:"sequence"`
}

type BaseVestingAccount struct {
	BaseAccount      BaseAccount `json:"base_account"`
	OriginalVesting  []Coins     `json:"original_vesting"`
	DelegatedFree    []Coins     `json:"delegated_free"`
	DelegatedVesting []Coins     `json:"delegated_vesting"`
	EndTime          string      `json:"end_time"`
}

type Accounts struct {
	Type               string             `json:"@type"`
	Address            string             `json:"address,omitempty"`
	PubKey             interface{}        `json:"pub_key,omitempty"`
	AccountNumber      string             `json:"account_number,omitempty"`
	Sequence           string             `json:"sequence,omitempty"`
	BaseAccount        BaseAccount        `json:"base_account,omitempty"`
	BaseVestingAccount BaseVestingAccount `json:"base_vesting_account,omitempty"`
	Name               string             `json:"name,omitempty"`
	Permissions        []interface{}      `json:"permissions,omitempty"`
}

type Auth struct {
	Params   AuthParams `json:"params"`
	Accounts []Accounts `json:"accounts"`
}

type Authz struct {
	Authorization []interface{} `json:"authorization"`
}

type BankParams struct {
	SendEnabled        []interface{} `json:"send_enabled"`
	DefaultSendEnabled bool          `json:"default_send_enabled"`
}

type Coins struct {
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}

type Balances struct {
	Address string  `json:"address"`
	Coins   []Coins `json:"coins"`
}

type Supply struct {
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}

type DenomUnits struct {
	Denom    string   `json:"denom"`
	Exponent int      `json:"exponent"`
	Aliases  []string `json:"aliases"`
}

type DenomMetadata struct {
	Description string       `json:"description"`
	DenomUnits  []DenomUnits `json:"denom_units"`
	Base        string       `json:"base"`
	Display     string       `json:"display"`
	Name        string       `json:"name"`
	Symbol      string       `json:"symbol"`
	URI         string       `json:"uri"`
	URIHash     string       `json:"uri_hash"`
}

type Bank struct {
	Params        BankParams      `json:"params"`
	Balances      []Balances      `json:"balances"`
	Supply        []Supply        `json:"supply"`
	DenomMetadata []DenomMetadata `json:"denom_metadata"`
}

type BlobParams struct {
	GasPerBlobByte   int    `json:"gas_per_blob_byte"`
	GovMaxSquareSize string `json:"gov_max_square_size"`
}

type BlobState struct {
	Params BlobParams `json:"params"`
}

type Capability struct {
	Index  string        `json:"index"`
	Owners []interface{} `json:"owners"`
}

type Crisis struct {
	ConstantFee Coins `json:"constant_fee"`
}

type DistributionParams struct {
	CommunityTax        string `json:"community_tax"`
	BaseProposerReward  string `json:"base_proposer_reward"`
	BonusProposerReward string `json:"bonus_proposer_reward"`
	WithdrawAddrEnabled bool   `json:"withdraw_addr_enabled"`
}

type FeePool struct {
	CommunityPool []interface{} `json:"community_pool"`
}

type Distribution struct {
	Params                          DistributionParams `json:"params"`
	FeePool                         FeePool            `json:"fee_pool"`
	DelegatorWithdrawInfos          []interface{}      `json:"delegator_withdraw_infos"`
	PreviousProposer                string             `json:"previous_proposer"`
	OutstandingRewards              []interface{}      `json:"outstanding_rewards"`
	ValidatorAccumulatedCommissions []interface{}      `json:"validator_accumulated_commissions"`
	ValidatorHistoricalRewards      []interface{}      `json:"validator_historical_rewards"`
	ValidatorCurrentRewards         []interface{}      `json:"validator_current_rewards"`
	DelegatorStartingInfos          []interface{}      `json:"delegator_starting_infos"`
	ValidatorSlashEvents            []interface{}      `json:"validator_slash_events"`
}

type Evidence struct {
	Evidence []interface{} `json:"evidence"`
}

type Feegrant struct {
	Allowances []interface{} `json:"allowances"`
}

type Description struct {
	Moniker         string `json:"moniker"`
	Identity        string `json:"identity"`
	Website         string `json:"website"`
	SecurityContact string `json:"security_contact"`
	Details         string `json:"details"`
}

type Commission struct {
	Rate          string `json:"rate"`
	MaxRate       string `json:"max_rate"`
	MaxChangeRate string `json:"max_change_rate"`
}

type Pubkey struct {
	Type string `json:"@type"`
	Key  string `json:"key"`
}

type Messages struct {
	Type              string      `json:"@type"`
	Description       Description `json:"description"`
	Commission        Commission  `json:"commission"`
	MinSelfDelegation string      `json:"min_self_delegation"`
	DelegatorAddress  string      `json:"delegator_address"`
	ValidatorAddress  string      `json:"validator_address"`
	Pubkey            Pubkey      `json:"pubkey"`
	Value             Coins       `json:"value"`
	EvmAddress        string      `json:"evm_address"`
}

type Body struct {
	Messages                    []Messages    `json:"messages"`
	Memo                        string        `json:"memo"`
	TimeoutHeight               uint64        `json:"timeout_height,string"`
	ExtensionOptions            []interface{} `json:"extension_options"`
	NonCriticalExtensionOptions []interface{} `json:"non_critical_extension_options"`
}

type Single struct {
	Mode string `json:"mode"`
}

type ModeInfo struct {
	Single Single `json:"single"`
}

type SignerInfos struct {
	PublicKey Pubkey   `json:"public_key"`
	ModeInfo  ModeInfo `json:"mode_info"`
	Sequence  string   `json:"sequence"`
}

type Fee struct {
	Amount   []interface{} `json:"amount"`
	GasLimit string        `json:"gas_limit"`
	Payer    string        `json:"payer"`
	Granter  string        `json:"granter"`
}

type AuthInfo struct {
	SignerInfos []SignerInfos `json:"signer_infos"`
	Fee         Fee           `json:"fee"`
	Tip         interface{}   `json:"tip"`
}

type GenTxs struct {
	Body       Body     `json:"body"`
	AuthInfo   AuthInfo `json:"auth_info"`
	Signatures []string `json:"signatures"`
}

type Genutil struct {
	GenTxs []json.RawMessage `json:"gen_txs"`
}

type DepositParams struct {
	MinDeposit       []Coins `json:"min_deposit"`
	MaxDepositPeriod string  `json:"max_deposit_period"`
}

type VotingParams struct {
	VotingPeriod string `json:"voting_period"`
}

type TallyParams struct {
	Quorum        string `json:"quorum"`
	Threshold     string `json:"threshold"`
	VetoThreshold string `json:"veto_threshold"`
}

type Gov struct {
	StartingProposalID string        `json:"starting_proposal_id"`
	Deposits           []interface{} `json:"deposits"`
	Votes              []interface{} `json:"votes"`
	Proposals          []interface{} `json:"proposals"`
	DepositParams      DepositParams `json:"deposit_params"`
	VotingParams       VotingParams  `json:"voting_params"`
	TallyParams        TallyParams   `json:"tally_params"`
}

type ClientGenesisParams struct {
	AllowedClients []string `json:"allowed_clients"`
}

type ClientGenesis struct {
	Clients            []interface{}       `json:"clients"`
	ClientsConsensus   []interface{}       `json:"clients_consensus"`
	ClientsMetadata    []interface{}       `json:"clients_metadata"`
	Params             ClientGenesisParams `json:"params"`
	CreateLocalhost    bool                `json:"create_localhost"`
	NextClientSequence string              `json:"next_client_sequence"`
}

type ConnectionGenesisParams struct {
	MaxExpectedTimePerBlock string `json:"max_expected_time_per_block"`
}

type ConnectionGenesis struct {
	Connections            []interface{}           `json:"connections"`
	ClientConnectionPaths  []interface{}           `json:"client_connection_paths"`
	NextConnectionSequence string                  `json:"next_connection_sequence"`
	Params                 ConnectionGenesisParams `json:"params"`
}

type ChannelGenesis struct {
	Channels            []interface{} `json:"channels"`
	Acknowledgements    []interface{} `json:"acknowledgements"`
	Commitments         []interface{} `json:"commitments"`
	Receipts            []interface{} `json:"receipts"`
	SendSequences       []interface{} `json:"send_sequences"`
	RecvSequences       []interface{} `json:"recv_sequences"`
	AckSequences        []interface{} `json:"ack_sequences"`
	NextChannelSequence string        `json:"next_channel_sequence"`
}

type Ibc struct {
	ClientGenesis     ClientGenesis     `json:"client_genesis"`
	ConnectionGenesis ConnectionGenesis `json:"connection_genesis"`
	ChannelGenesis    ChannelGenesis    `json:"channel_genesis"`
}

type Minter struct {
	InflationRate     string      `json:"inflation_rate"`
	AnnualProvisions  string      `json:"annual_provisions"`
	PreviousBlockTime interface{} `json:"previous_block_time"`
	BondDenom         string      `json:"bond_denom"`
}

type Mint struct {
	Minter Minter `json:"minter"`
}

type QgbParams struct {
	DataCommitmentWindow string `json:"data_commitment_window"`
}

type Qgb struct {
	Params QgbParams `json:"params"`
}

type SlashingParams struct {
	SignedBlocksWindow      string `json:"signed_blocks_window"`
	MinSignedPerWindow      string `json:"min_signed_per_window"`
	DowntimeJailDuration    string `json:"downtime_jail_duration"`
	SlashFractionDoubleSign string `json:"slash_fraction_double_sign"`
	SlashFractionDowntime   string `json:"slash_fraction_downtime"`
}

type Slashing struct {
	Params       SlashingParams `json:"params"`
	SigningInfos []interface{}  `json:"signing_infos"`
	MissedBlocks []interface{}  `json:"missed_blocks"`
}

type StakingParams struct {
	UnbondingTime     string `json:"unbonding_time"`
	MaxValidators     int    `json:"max_validators"`
	MaxEntries        int    `json:"max_entries"`
	HistoricalEntries int    `json:"historical_entries"`
	BondDenom         string `json:"bond_denom"`
	MinCommissionRate string `json:"min_commission_rate"`
}

type Staking struct {
	Params               StakingParams `json:"params"`
	LastTotalPower       string        `json:"last_total_power"`
	LastValidatorPowers  []interface{} `json:"last_validator_powers"`
	Validators           []interface{} `json:"validators"`
	Delegations          []interface{} `json:"delegations"`
	UnbondingDelegations []interface{} `json:"unbonding_delegations"`
	Redelegations        []interface{} `json:"redelegations"`
	Exported             bool          `json:"exported"`
}

type TransferParams struct {
	SendEnabled    bool `json:"send_enabled"`
	ReceiveEnabled bool `json:"receive_enabled"`
}

type Transfer struct {
	PortID      string         `json:"port_id"`
	DenomTraces []interface{}  `json:"denom_traces"`
	Params      TransferParams `json:"params"`
}

type Vesting struct {
}

type AppState struct {
	Auth         Auth         `json:"auth"`
	Authz        Authz        `json:"authz"`
	Bank         Bank         `json:"bank"`
	Blob         Blob         `json:"blob"`
	Capability   Capability   `json:"capability"`
	Crisis       Crisis       `json:"crisis"`
	Distribution Distribution `json:"distribution"`
	Evidence     Evidence     `json:"evidence"`
	Feegrant     Feegrant     `json:"feegrant"`
	Genutil      Genutil      `json:"genutil"`
	Gov          Gov          `json:"gov"`
	Ibc          Ibc          `json:"ibc"`
	Mint         Mint         `json:"mint"`
	Params       interface{}  `json:"params"`
	Qgb          Qgb          `json:"qgb"`
	Slashing     Slashing     `json:"slashing"`
	Staking      Staking      `json:"staking"`
	Transfer     Transfer     `json:"transfer"`
	Vesting      Vesting      `json:"vesting"`
}
