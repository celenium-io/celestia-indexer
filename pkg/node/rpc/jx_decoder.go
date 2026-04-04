// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package rpc

// jx_decoder.go — streaming JSON decoder for CometBFT block/results responses.
//
// Uses go-faster/jx for allocation-efficient token-based parsing:
//   - ObjBytes: zero-copy key access (no string(key) alloc)
//   - StrBytes: zero-copy string view into jx's read buffer
//   - knownEventStrings: intern table for repetitive event type / attr key strings

import (
	"encoding/base64"
	"encoding/hex"
	stdjson "encoding/json"
	"strconv"
	"time"
	"unsafe"

	jxpkg "github.com/go-faster/jx"
	"github.com/pkg/errors"

	"github.com/celenium-io/celestia-indexer/internal/pool"
	nodeTypes "github.com/celenium-io/celestia-indexer/pkg/node/types"
	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"
	cmtTypes "github.com/cometbft/cometbft/types"
)

// base64BufPool reuses scratch buffers for base64-encoded transaction strings.
// StrAppend(buf) writes the JSON string content directly into the provided
// buffer (never calls strSlow), so large blob transactions no longer cause
// per-tx heap allocations.
var base64BufPool = pool.New(func() []byte { return make([]byte, 0, 4*1024*1024) })

const keyHeight = "height"

// knownEventStrings contains all well-known event type and attribute key strings
// in Celestia. Populated once at init; the backing string constants live in
// read-only data and are never GC-ed.
var knownEventStrings = func() map[string]string {
	ss := []string{
		// event types — sourced from internal/storage/types/event_type.go
		// base Cosmos / staking
		"coin_received", "coinbase", "coin_spent", "burn", "mint",
		"message", "proposer_reward", "rewards", "commission",
		"liveness", "transfer",
		"redelegate", "AttestationRequest",
		"withdraw_rewards", "withdraw_commission", "set_withdraw_address",
		"create_validator", "delegate", "edit_validator", "unbond", "tx",
		"complete_redelegation", "complete_unbonding",
		// feegrant
		"use_feegrant", "revoke_feegrant", "set_feegrant", "update_feegrant",
		// slash / gov
		"slash",
		"proposal_vote", "proposal_deposit", "submit_proposal",
		"active_proposal", "inactive_proposal",
		// authz
		"cosmos.authz.v1beta1.EventGrant",
		"cosmos.authz.v1beta1.EventRevoke",
		"cosmos.authz.v1.EventRevoke",
		"cancel_unbonding_delegation",
		// IBC core
		"send_packet", "ibc_transfer",
		"fungible_token_packet", "acknowledge_packet",
		"create_client", "update_client", "update_client_proposal",
		"connection_open_try", "connection_open_init",
		"connection_open_confirm", "connection_open_ack",
		"channel_open_try", "channel_open_init",
		"channel_open_confirm", "channel_open_ack", "channel_close_confirm",
		"recv_packet", "write_acknowledgement",
		"timeout", "timeout_packet",
		"ics27_packet", "ibccallbackerror-ics27_packet",
		// Celestia
		"celestia.blob.v1.EventPayForBlobs",
		"celestia.forwarding.v1.EventTokenForwarded",
		"celestia.forwarding.v1.EventForwardingComplete",
		"celestia.zkism.v1.EventCreateInterchainSecurityModule",
		"celestia.zkism.v1.EventUpdateInterchainSecurityModule",
		"celestia.zkism.v1.EventSubmitMessages",
		// Hyperlane
		"hyperlane.core.v1.EventDispatch",
		"hyperlane.core.v1.EventProcess",
		"hyperlane.core.v1.EventCreateMailbox",
		"hyperlane.core.v1.EventSetMailbox",
		"hyperlane.warp.v1.EventCreateSyntheticToken",
		"hyperlane.warp.v1.EventCreateCollateralToken",
		"hyperlane.warp.v1.EventSetToken",
		"hyperlane.warp.v1.EventEnrollRemoteRouter",
		"hyperlane.warp.v1.EventUnrollRemoteRouter",
		"hyperlane.warp.v1.EventSendRemoteTransfer",
		"hyperlane.warp.v1.EventReceiveRemoteTransfer",
		"hyperlane.core.post_dispatch.v1.EventCreateMerkleTreeHook",
		"hyperlane.core.post_dispatch.v1.EventInsertedIntoTree",
		"hyperlane.core.post_dispatch.v1.EventGasPayment",
		"hyperlane.core.post_dispatch.v1.EventCreateNoopHook",
		"hyperlane.core.post_dispatch.v1.EventCreateIgp",
		"hyperlane.core.post_dispatch.v1.EventSetIgp",
		"hyperlane.core.post_dispatch.v1.EventSetDestinationGasConfig",
		"hyperlane.core.post_dispatch.v1.EventClaimIgp",
		"hyperlane.core.interchain_security.v1.EventCreateNoopIsm",
		"hyperlane.core.interchain_security.v1.EventSetRoutingIsmDomain",
		"hyperlane.core.interchain_security.v1.EventSetRoutingIsm",
		"hyperlane.core.interchain_security.v1.EventCreateRoutingIsm",
		// signal
		"signal_version",
		// attribute keys (plain, Cosmos SDK ≥ v0.47)
		"amount", "msg_index", "validator", "mode", "sender",
		"receiver", "delegator", "spender", "recipient",
		"action", "module", "granter", "grantee",
		"proposal_id", "option", "deposit", "voter",
		"packet_data", "packet_src_channel", "packet_dst_channel",
		"packet_src_port", "packet_dst_port", "packet_sequence",
		"packet_timeout_height", "packet_timeout_timestamp",
		"packet_connection", "connection_id", "channel_id",
		"port_id", "counterparty_channel_id", "counterparty_port_id",
		"acc_seq", "fee", "tip", "success", "error",
		"denom", "new_shares", "completion_time",
		"withdraw_address", "validator_address",
		// new plain keys found in 9232135 / other blocks
		"fee_payer", "signature", "signer",
		"blob_sizes", "namespaces",
		"address", "missed_blocks", "height",
		"authz_msg_index", "source_validator", "destination_validator",
		"minter", "inflation_rate", "annual_provisions",
		"packet_data_hex", "packet_channel_ordering",
		"packet_ack", "packet_ack_hex", "host_channel_id",
		"client_id", "client_type", "consensus_height", "consensus_heights",
		"header", "memo",
		// governance
		"voting_period_start", "proposal_messages",
		// IBC connection/client
		"counterparty_client_id", "counterparty_connection_id",
		// Celestia forwarding
		"forward_addr", "tokens_forwarded", "tokens_failed",
		// Hyperlane
		"message_id", "token_id", "igp_id", "ism_id", "mailbox_id", "origin_mailbox_id",
		"new_owner", "renounce_ownership",
		"default_hook", "default_ism", "required_hook",
		"dest_domain", "dest_recipient", "local_domain", "remote_domain", "origin_mailbox",
		"gas_amount", "gas_overhead", "gas_price", "token_exchange_rate",
		"merkle_tree_address", "origin_denom",
		// zkism
		"state_root", "groth16_vkey", "state_transition_vkey", "state_membership_vkey",
	}
	m := make(map[string]string, len(ss))
	for _, s := range ss {
		m[s] = s
	}
	return m
}()

// jxInternStr reads a string from the decoder without allocating for known strings.
// An ephemeral string header is created over jx's buffer for the map lookup only —
// it must not be stored or passed to another goroutine.
func jxInternStr(d *jxpkg.Decoder) (string, error) {
	raw, err := d.StrBytes()
	if err != nil {
		return "", err
	}
	tmp := unsafe.String(unsafe.SliceData(raw), len(raw))
	if known, ok := knownEventStrings[tmp]; ok {
		return known, nil
	}
	return string(raw), nil
}

// jxHex decodes a JSON hex string using StrBytes to avoid an intermediate
// string allocation.
func jxHex(d *jxpkg.Decoder) ([]byte, error) {
	raw, err := d.StrBytes()
	if err != nil || len(raw) == 0 {
		return nil, err
	}
	out := make([]byte, hex.DecodedLen(len(raw)))
	_, err = hex.Decode(out, raw)
	return out, err
}

func jxInt64(d *jxpkg.Decoder) (int64, error) {
	raw, err := d.StrBytes()
	if err != nil {
		return 0, err
	}

	return strconv.ParseInt(unsafe.String(unsafe.SliceData(raw), len(raw)), 10, 64)
}

func jxUint64(d *jxpkg.Decoder) (uint64, error) {
	raw, err := d.StrBytes()
	if err != nil {
		return 0, err
	}

	return strconv.ParseUint(unsafe.String(unsafe.SliceData(raw), len(raw)), 10, 64)
}

func jxDuration(d *jxpkg.Decoder) (time.Duration, error) {
	v, err := jxInt64(d)
	if err != nil {
		return 0, nil
	}
	return time.Duration(v), nil
}

func jxTime(d *jxpkg.Decoder) (time.Time, error) {
	raw, err := d.StrBytes()
	if err != nil {
		return time.Time{}, err
	}
	return time.Parse(time.RFC3339Nano, unsafe.String(unsafe.SliceData(raw), len(raw)))
}

func jxEventAttribute(d *jxpkg.Decoder) (pkgTypes.EventAttribute, error) {
	var a pkgTypes.EventAttribute
	return a, d.ObjBytes(func(d *jxpkg.Decoder, key []byte) error {
		switch string(key) {
		case "key":
			s, err := jxInternStr(d)
			if err != nil {
				return err
			}
			a.Key = s
		case "value":
			s, err := d.StrBytes()
			if err != nil {
				return err
			}
			a.Value = string(s)
		default:
			return d.Skip()
		}
		return nil
	})
}

func jxEvent(d *jxpkg.Decoder) (pkgTypes.Event, error) {
	var ev pkgTypes.Event
	return ev, d.ObjBytes(func(d *jxpkg.Decoder, key []byte) error {
		switch string(key) {
		case "type":
			s, err := jxInternStr(d)
			if err != nil {
				return err
			}
			ev.Type = s
		case "attributes":
			return d.Arr(func(d *jxpkg.Decoder) error {
				attr, err := jxEventAttribute(d)
				if err != nil {
					return err
				}
				ev.Attributes = append(ev.Attributes, attr)
				return nil
			})
		default:
			return d.Skip()
		}
		return nil
	})
}

func jxTxResult(d *jxpkg.Decoder) (pkgTypes.ResponseDeliverTx, error) {
	var tx pkgTypes.ResponseDeliverTx
	return tx, d.ObjBytes(func(d *jxpkg.Decoder, key []byte) error {
		switch string(key) {
		case "code":
			v, err := d.UInt32()
			if err != nil {
				return err
			}
			tx.Code = v
		case "log":
			raw, err := d.Raw()
			if err != nil {
				return err
			}
			tx.Log = stdjson.RawMessage(raw)
		case "gas_wanted":
			raw, err := d.StrBytes()
			if err != nil {
				return err
			}
			v, err := strconv.ParseInt(unsafe.String(unsafe.SliceData(raw), len(raw)), 10, 64)
			if err != nil {
				return err
			}
			tx.GasWanted = v
		case "gas_used":
			raw, err := d.StrBytes()
			if err != nil {
				return err
			}
			v, err := strconv.ParseInt(unsafe.String(unsafe.SliceData(raw), len(raw)), 10, 64)
			if err != nil {
				return err
			}
			tx.GasUsed = v
		case "events":
			return d.Arr(func(d *jxpkg.Decoder) error {
				ev, err := jxEvent(d)
				if err != nil {
					return err
				}
				tx.Events = append(tx.Events, ev)
				return nil
			})
		case "codespace":
			s, err := d.Str()
			if err != nil {
				return err
			}
			tx.Codespace = s
		default:
			return d.Skip()
		}
		return nil
	})
}

// jxResultBlockResults fully decodes the block_results payload.
// consensus_param_updates and validator_updates are skipped.
func jxResultBlockResults(d *jxpkg.Decoder) (pkgTypes.ResultBlockResults, error) {
	var r pkgTypes.ResultBlockResults
	return r, d.ObjBytes(func(d *jxpkg.Decoder, key []byte) error {
		switch string(key) {
		case keyHeight:
			raw, err := d.StrBytes()
			if err != nil {
				return err
			}
			h, err := strconv.ParseInt(unsafe.String(unsafe.SliceData(raw), len(raw)), 10, 64)
			if err != nil {
				return err
			}
			r.Height = pkgTypes.Level(h)
		case "txs_results":
			if d.Next() == jxpkg.Null {
				return d.Null()
			}
			return d.Arr(func(d *jxpkg.Decoder) error {
				tx, err := jxTxResult(d)
				if err != nil {
					return err
				}
				r.TxsResults = append(r.TxsResults, tx)
				return nil
			})
		case "finalize_block_events":
			return d.Arr(func(d *jxpkg.Decoder) error {
				ev, err := jxEvent(d)
				if err != nil {
					return err
				}
				r.FinalizeBlockEvents = append(r.FinalizeBlockEvents, ev)
				return nil
			})
		case "consensus_param_updates":
			if d.Next() == jxpkg.Null {
				return d.Null()
			}
			params, err := jxConsensusParams(d)
			if err != nil {
				return errors.Wrap(err, "consensus params update")
			}
			r.ConsensusParamUpdates = &params
		default:
			return d.Skip()
		}
		return nil
	})
}

func jxConsensusParams(d *jxpkg.Decoder) (params pkgTypes.ConsensusParams, err error) {
	err = d.ObjBytes(func(d *jxpkg.Decoder, key []byte) error {
		switch string(key) {
		case "block":
			if d.Next() == jxpkg.Null {
				return d.Null()
			}
			bp, err := jxConsensusParamsBlock(d)
			if err != nil {
				return errors.Wrap(err, "block params")
			}
			params.Block = &bp
		case "evidence":
			if d.Next() == jxpkg.Null {
				return d.Null()
			}
			ep, err := jxConsensusParamsEvidence(d)
			if err != nil {
				return errors.Wrap(err, "block params")
			}
			params.Evidence = &ep
		case "validator":
			if d.Next() == jxpkg.Null {
				return d.Null()
			}
			vp, err := jxConsensusParamsValidator(d)
			if err != nil {
				return errors.Wrap(err, "validator params")
			}
			params.Validator = &vp
		default:
			return d.Skip()
		}
		return nil
	})
	return
}

func jxConsensusParamsBlock(d *jxpkg.Decoder) (params pkgTypes.BlockParams, err error) {
	err = d.ObjBytes(func(d *jxpkg.Decoder, key []byte) error {
		switch string(key) {
		case "max_bytes":
			v, err := jxInt64(d)
			if err != nil {
				return errors.Wrap(err, "max_bytes")
			}
			params.MaxBytes = v
			return nil
		case "max_gas":
			v, err := jxInt64(d)
			if err != nil {
				return errors.Wrap(err, "max_gas")
			}
			params.MaxGas = v
			return nil
		default:
			return d.Skip()
		}
	})
	return
}

func jxConsensusParamsEvidence(d *jxpkg.Decoder) (params pkgTypes.EvidenceParams, err error) {
	err = d.ObjBytes(func(d *jxpkg.Decoder, key []byte) error {
		switch string(key) {
		case "max_age_num_blocks":
			v, err := jxInt64(d)
			if err != nil {
				return errors.Wrap(err, "max_age_num_blocks")
			}
			params.MaxAgeNumBlocks = v
		case "max_age_duration":
			v, err := jxDuration(d)
			if err != nil {
				return errors.Wrap(err, "max_age_duration")
			}
			params.MaxAgeDuration = v
		case "max_bytes":
			v, err := jxInt64(d)
			if err != nil {
				return errors.Wrap(err, "max_bytes")
			}
			params.MaxBytes = v
		default:
			return d.Skip()
		}

		return nil
	})
	return
}

func jxConsensusParamsValidator(d *jxpkg.Decoder) (params pkgTypes.ValidatorParams, err error) {
	err = d.ObjBytes(func(d *jxpkg.Decoder, key []byte) error {
		switch string(key) {
		case "pub_key_types":
			return d.Arr(func(d *jxpkg.Decoder) error {
				raw, err := d.StrBytes()
				if err != nil {
					return err
				}
				params.PubKeyTypes = append(params.PubKeyTypes, string(raw))
				return nil
			})
		default:
			return d.Skip()
		}
	})
	return
}

func jxHeader(d *jxpkg.Decoder, h *pkgTypes.Header) error {
	return d.ObjBytes(func(d *jxpkg.Decoder, key []byte) error {
		switch string(key) {
		case "chain_id":
			raw, err := d.StrBytes()
			if err != nil {
				return err
			}
			tmp := unsafe.String(unsafe.SliceData(raw), len(raw))
			if known, ok := knownEventStrings[tmp]; ok {
				h.ChainID = known
			} else {
				h.ChainID = string(raw)
			}
		case keyHeight:
			raw, err := d.StrBytes()
			if err != nil {
				return err
			}
			v, err := strconv.ParseInt(unsafe.String(unsafe.SliceData(raw), len(raw)), 10, 64)
			if err != nil {
				return err
			}
			h.Height = v
		case "time":
			t, err := jxTime(d)
			if err != nil {
				return err
			}
			h.Time = t
		case "last_block_id":
			return d.ObjBytes(func(d *jxpkg.Decoder, key []byte) error {
				if string(key) == "hash" {
					hx, err := jxHex(d)
					if err != nil {
						return err
					}
					h.LastBlockID.Hash = hx
					return nil
				}
				return d.Skip()
			})
		case "last_commit_hash":
			hx, err := jxHex(d)
			if err != nil {
				return err
			}
			h.LastCommitHash = hx
		case "data_hash":
			hx, err := jxHex(d)
			if err != nil {
				return err
			}
			h.DataHash = hx
		case "validators_hash":
			hx, err := jxHex(d)
			if err != nil {
				return err
			}
			h.ValidatorsHash = hx
		case "next_validators_hash":
			hx, err := jxHex(d)
			if err != nil {
				return err
			}
			h.NextValidatorsHash = hx
		case "consensus_hash":
			hx, err := jxHex(d)
			if err != nil {
				return err
			}
			h.ConsensusHash = hx
		case "app_hash":
			hx, err := jxHex(d)
			if err != nil {
				return err
			}
			h.AppHash = hx
		case "last_results_hash":
			hx, err := jxHex(d)
			if err != nil {
				return err
			}
			h.LastResultsHash = hx
		case "evidence_hash":
			hx, err := jxHex(d)
			if err != nil {
				return err
			}
			h.EvidenceHash = hx
		case "proposer_address":
			hx, err := jxHex(d)
			if err != nil {
				return err
			}
			h.ProposerAddress = hx
		case "version":
			err := d.ObjBytes(func(d *jxpkg.Decoder, key []byte) error {
				switch string(key) {
				case "app":
					version, err := jxUint64(d)
					if err != nil {
						return errors.Wrap(err, "app")
					}
					h.Version.App = version
					return nil
				default:
					return d.Skip()
				}
			})
			if err != nil {
				return err
			}
		default:
			return d.Skip()
		}
		return nil
	})
}

// jxData decodes the block's data section
func jxData(d *jxpkg.Decoder, data *pkgTypes.Data) error {
	return d.ObjBytes(func(d *jxpkg.Decoder, key []byte) error {
		switch string(key) {
		case "txs":
			return d.Arr(func(d *jxpkg.Decoder) error {
				buf := base64BufPool.Get()
				buf, err := d.StrAppend(buf[:0])
				if err != nil {
					base64BufPool.Put(buf)
					return err
				}
				decoded, err := base64.StdEncoding.AppendDecode(nil, buf)
				base64BufPool.Put(buf)
				if err != nil {
					return errors.Wrap(err, "data base64 decode")
				}
				data.Txs = append(data.Txs, decoded)
				return nil
			})
		case "square_size":
			v, err := jxUint64(d)
			if err != nil {
				return err
			}
			data.SquareSize = v
		default:
			return d.Skip()
		}
		return nil
	})
}

func jxCommitSig(d *jxpkg.Decoder) (pkgTypes.CommitSig, error) {
	var s pkgTypes.CommitSig
	return s, d.ObjBytes(func(d *jxpkg.Decoder, key []byte) error {
		switch string(key) {
		case "block_id_flag":
			v, err := d.Int()
			if err != nil {
				return err
			}
			s.BlockIDFlag = cmtTypes.BlockIDFlag(v)
		case "validator_address":
			hx, err := jxHex(d)
			if err != nil {
				return err
			}
			s.ValidatorAddress = hx
		case "timestamp":
			t, err := jxTime(d)
			if err != nil {
				return err
			}
			s.Timestamp = t
		default:
			return d.Skip()
		}
		return nil
	})
}

func jxCommit(d *jxpkg.Decoder) (pkgTypes.Commit, error) {
	var c pkgTypes.Commit
	return c, d.ObjBytes(func(d *jxpkg.Decoder, key []byte) error {
		switch string(key) {
		case keyHeight:
			v, err := jxInt64(d)
			if err != nil {
				return err
			}
			c.Height = v
		case "signatures":
			return d.Arr(func(d *jxpkg.Decoder) error {
				sig, err := jxCommitSig(d)
				if err != nil {
					return err
				}
				c.Signatures = append(c.Signatures, sig)
				return nil
			})
		default:
			return d.Skip()
		}
		return nil
	})
}

// jxResultBlock decodes the block payload.
func jxResultBlock(d *jxpkg.Decoder) (pkgTypes.ResultBlock, error) {
	var rb pkgTypes.ResultBlock
	return rb, d.ObjBytes(func(d *jxpkg.Decoder, key []byte) error {
		switch string(key) {
		case "block_id":
			return d.ObjBytes(func(d *jxpkg.Decoder, key []byte) error {
				if string(key) == "hash" {
					h, err := jxHex(d)
					if err != nil {
						return err
					}
					rb.BlockID.Hash = h
					return nil
				}
				return d.Skip()
			})
		case "block":
			if rb.Block == nil {
				rb.Block = new(pkgTypes.Block)
			}
			return d.ObjBytes(func(d *jxpkg.Decoder, key []byte) error {
				switch string(key) {
				case "header":
					return jxHeader(d, &rb.Block.Header)
				case "data":
					return jxData(d, &rb.Block.Data)
				case "last_commit":
					lc, err := jxCommit(d)
					if err != nil {
						return err
					}
					rb.Block.LastCommit = &lc
				default:
					return d.Skip()
				}
				return nil
			})
		default:
			return d.Skip()
		}
	})
}

// jxResponse decodes a single JSON-RPC response envelope, calling fn with the
// decoder positioned at the "result" value. Returns ErrRequest immediately
// when an "error" field is present and non-null.
func jxResponse(d *jxpkg.Decoder, fn func(*jxpkg.Decoder) error) error {
	return d.ObjBytes(func(d *jxpkg.Decoder, key []byte) error {
		switch string(key) {
		case "result":
			return fn(d)
		case "error":
			if d.Next() == jxpkg.Null {
				return d.Null()
			}
			var msg string
			if err := d.ObjBytes(func(d *jxpkg.Decoder, key []byte) error {
				if string(key) == "message" {
					var err error
					msg, err = d.Str()
					return err
				}
				return d.Skip()
			}); err != nil {
				return err
			}
			return errors.Wrapf(nodeTypes.ErrRequest, "request error: %s", msg)
		default:
			return d.Skip()
		}
	})
}

// jxStatusMinimal decodes the "sync_info.latest_block_height" field from the
// /status response, skipping all other fields.
func jxStatusMinimal(d *jxpkg.Decoder) (pkgTypes.Level, error) {
	var level pkgTypes.Level
	return level, d.ObjBytes(func(d *jxpkg.Decoder, key []byte) error {
		if string(key) != "sync_info" {
			return d.Skip()
		}
		return d.ObjBytes(func(d *jxpkg.Decoder, key []byte) error {
			if string(key) != "latest_block_height" {
				return d.Skip()
			}
			v, err := jxInt64(d)
			if err != nil {
				return err
			}
			level = pkgTypes.Level(v)
			return nil
		})
	})
}

// jxGenesisChunk decodes a GenesisChunk from the genesis_chunked response.
// The "data" field is base64-decoded from the JSON string value.
func jxGenesisChunk(d *jxpkg.Decoder, gc *GenesisChunk) error {
	return d.ObjBytes(func(d *jxpkg.Decoder, key []byte) error {
		switch string(key) {
		case "chunk":
			v, err := jxInt64(d)
			if err != nil {
				return err
			}
			gc.Chunk = v
		case "total":
			v, err := jxInt64(d)
			if err != nil {
				return err
			}
			gc.Total = v
		case "data":
			raw, err := d.StrBytes()
			if err != nil {
				return err
			}
			decoded, err := base64.StdEncoding.AppendDecode(nil, raw)
			if err != nil {
				return errors.Wrap(err, "genesis chunk base64 decode")
			}
			gc.Data = decoded
		default:
			return d.Skip()
		}
		return nil
	})
}

// jxBatchResponse decodes a JSON-RPC batch response array, calling fn for each
// pair of block + block_results responses. Handles JSON-RPC error objects.
func jxBatchResponse(d *jxpkg.Decoder, fn func(pkgTypes.BlockData) error) error {
	var current pkgTypes.BlockData
	idx := 0
	return d.Arr(func(d *jxpkg.Decoder) error {
		var rpcErr string
		return d.ObjBytes(func(d *jxpkg.Decoder, key []byte) error {
			switch string(key) {
			case "result":
				if idx%2 == 0 {
					rb, err := jxResultBlock(d)
					if err != nil {
						return err
					}
					current.ResultBlock = rb
				} else {
					rbr, err := jxResultBlockResults(d)
					if err != nil {
						return err
					}
					current.ResultBlockResults = rbr
					if err := fn(current); err != nil {
						return err
					}
					current = pkgTypes.BlockData{}
				}
				idx++
			case "error":
				if d.Next() == jxpkg.Null {
					return d.Null()
				}
				if err := d.ObjBytes(func(d *jxpkg.Decoder, key []byte) error {
					if string(key) == "message" {
						var err error
						rpcErr, err = d.Str()
						return err
					}
					return d.Skip()
				}); err != nil {
					return err
				}
			default:
				return d.Skip()
			}
			if rpcErr != "" {
				return errors.Wrapf(nodeTypes.ErrRequest, "request error: %s", rpcErr)
			}
			return nil
		})
	})
}
