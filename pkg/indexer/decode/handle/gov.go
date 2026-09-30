// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package handle

import (
	stdJSON "encoding/json"
	"fmt"
	"maps"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	consensusv1 "cosmossdk.io/api/cosmos/consensus/v1"
	json "github.com/bytedance/sonic"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	blobTypes "github.com/celestiaorg/celestia-app/v10/x/blob/types"
	"github.com/cosmos/cosmos-sdk/codec"
	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	distributionTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	v1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	cosmosGovTypesV1Beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	paramsV1Beta "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	ibcTypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	"github.com/pkg/errors"
	"github.com/stoewer/go-strcase"
)

func moduleFromTypeURL(typeURL string) string {
	splitted := strings.Split(typeURL, ".")
	if len(splitted) < 2 {
		return typeURL
	}
	return splitted[1]
}

// ProtoJSON encodes google.protobuf.Duration as seconds with the "s" suffix
var durationRe = regexp.MustCompile(`^-?\d+(\.\d+)?s$`)

// UseNumber keeps big integers exact instead of float64
var paramsJSON = json.Config{UseNumber: true}.Froze()

func flatten(module, prefix string, value map[string]any, changes *[]paramsV1Beta.ParamChange) error {
	// sorted keys: changes are stored in proposal.changes and must be stable
	for _, key := range slices.Sorted(maps.Keys(value)) {
		name := key
		if prefix != "" {
			name = prefix + "_" + key
		}

		switch typed := value[key].(type) {
		case nil:
			continue

		case map[string]any:
			if err := flatten(module, name, typed, changes); err != nil {
				return err
			}

		case []any:
			if amount, ok := firstCoinAmount(typed); ok {
				appendChange(changes, module, name, amount)
				continue
			}
			if len(typed) == 0 {
				continue
			}
			raw, err := paramsJSON.MarshalToString(typed)
			if err != nil {
				return errors.Wrap(err, name)
			}
			appendChange(changes, module, name, raw)

		case string:
			if durationRe.MatchString(typed) {
				d, err := time.ParseDuration(typed)
				if err != nil {
					return errors.Wrap(err, name)
				}
				appendChange(changes, module, name, strconv.FormatInt(d.Nanoseconds(), 10))
				continue
			}
			appendChange(changes, module, name, typed)

		case stdJSON.Number: // alias of json.Number
			appendChange(changes, module, name, typed.String())

		case bool:
			appendChange(changes, module, name, strconv.FormatBool(typed))

		default:
			return errors.Errorf("unexpected param type %T: %s", typed, name)
		}
	}
	return nil
}

// firstCoinAmount matches sdk.Coins: the constant keeps the amount of the first coin and its denom
func firstCoinAmount(items []any) (string, bool) {
	if len(items) == 0 {
		return "", false
	}
	coin, ok := items[0].(map[string]any)
	if !ok || len(coin) != 2 {
		return "", false
	}
	amount, ok := coin["amount"]
	if !ok {
		return "", false
	}
	amountString, ok := amount.(string)
	if !ok {
		return "", false
	}
	denom, ok := coin["denom"]
	if !ok {
		return "", false
	}
	denomString, ok := denom.(string)
	if !ok {
		return "", false
	}
	return fmt.Sprintf("%s%s", amountString, denomString), true
}

func appendChange(changes *[]paramsV1Beta.ParamChange, module, key, value string) {
	if changes == nil {
		return
	}
	*changes = append(*changes, paramsV1Beta.NewParamChange(module, key, value))
}

// MsgSubmitProposalV1
func MsgSubmitProposalV1(ctx *context.Context, codec codec.Codec, status storageTypes.Status, msgId uint64, msg *v1.MsgSubmitProposal) (storageTypes.MsgType, []any, *storage.Proposal, error) {
	msgType := storageTypes.MsgSubmitProposal
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeProposer, address: msg.Proposer},
	}, ctx.Block.Height, msgId)
	if err != nil {
		return msgType, nil, nil, err
	}

	if status != storageTypes.StatusSuccess {
		return msgType, nil, nil, nil
	}

	prpsl := &storage.Proposal{
		Height: ctx.Block.Height,
		Proposer: &storage.Address{
			Address: msg.Proposer,
		},
		CreatedAt:   ctx.Block.Time,
		Status:      storageTypes.ProposalStatusInactive,
		Type:        storageTypes.ProposalTypeText,
		Title:       msg.Title,
		Description: msg.Summary,
		Metadata:    msg.Metadata,
		Expedited:   &msg.Expedited,
	}

	if prpsl.Title == "" {
		prpsl.Title = "Proposal with messages"
	}

	changes := make([]paramsV1Beta.ParamChange, 0)
	recoveries := make([]map[string]any, 0)
	var sb strings.Builder
	if _, err := sb.WriteString("Proposal contains messages:\r\n"); err != nil {
		return msgType, nil, nil, errors.Wrap(err, "building proposal description from messages")
	}
	for i := range msg.Messages {
		if _, err := fmt.Fprintf(&sb, "%d. %s\r\n", i+1, msg.Messages[i].TypeUrl); err != nil {
			return msgType, nil, nil, errors.Wrap(err, "building proposal description from messages")
		}

		switch msg.Messages[i].TypeUrl {
		case "/celestia.blob.v1.MsgUpdateBlobParams":
			var params blobTypes.MsgUpdateBlobParams
			if err := codec.Unmarshal(msg.Messages[i].Value, &params); err != nil {
				return msgType, nil, nil, errors.Wrap(err, "unmarshalling proposal with cosmos.gov.v1.MsgUpdateParams")
			}
			prpsl.Type = storageTypes.ProposalTypeParamChanged

			gasBlobPerByte := strconv.FormatUint(uint64(params.Params.GetGasPerBlobByte()), 10)
			changes = append(changes, paramsV1Beta.NewParamChange(storageTypes.ModuleNameBlob.String(), "gas_per_blob_byte", gasBlobPerByte))

			maxSquareSize := strconv.FormatUint(params.Params.GetGovMaxSquareSize(), 10)
			changes = append(changes, paramsV1Beta.NewParamChange(storageTypes.ModuleNameBlob.String(), "gov_max_square_size", maxSquareSize))

		case "/ibc.core.client.v1.MsgRecoverClient":
			var recoverMsg ibcTypes.MsgRecoverClient
			if err := codec.Unmarshal(msg.Messages[i].Value, &recoverMsg); err != nil {
				return msgType, nil, nil, errors.Wrap(err, "unmarshalling proposal with ibc.core.client.v1.MsgRecoverClient")
			}

			prpsl.Type = storageTypes.ProposalTypeClientUpdate
			// the storage module reads them on execution to find the substitute
			recoveries = append(recoveries, map[string]any{
				"SubjectClientId":    recoverMsg.SubjectClientId,
				"SubstituteClientId": recoverMsg.SubstituteClientId,
			})

		case "/cosmos.consensus.v1.MsgUpdateParams":
			var params consensusv1.MsgUpdateParams
			if err := codec.Unmarshal(msg.Messages[i].Value, &params); err != nil {
				return msgType, nil, nil, errors.Wrap(err, "unmarshalling proposal with cosmos.gov.v1.MsgUpdateParams")
			}

			prpsl.Type = storageTypes.ProposalTypeParamChanged

			if block := params.GetBlock(); block != nil {
				maxBytes := strconv.FormatInt(block.GetMaxBytes(), 10)
				changes = append(changes, paramsV1Beta.NewParamChange(storageTypes.ModuleNameConsensus.String(), "block_max_bytes", maxBytes))

				maxGas := strconv.FormatInt(block.GetMaxGas(), 10)
				changes = append(changes, paramsV1Beta.NewParamChange(storageTypes.ModuleNameConsensus.String(), "block_max_gas", maxGas))
			}

			if evidence := params.GetEvidence(); evidence != nil {
				maxAgeNumBlocks := strconv.FormatInt(evidence.GetMaxAgeNumBlocks(), 10)
				changes = append(changes, paramsV1Beta.NewParamChange(storageTypes.ModuleNameConsensus.String(), "evidence_max_age_num_blocks", maxAgeNumBlocks))

				maxBytes := strconv.FormatInt(evidence.GetMaxBytes(), 10)
				changes = append(changes, paramsV1Beta.NewParamChange(storageTypes.ModuleNameConsensus.String(), "evidence_max_bytes", maxBytes))

				if age := evidence.GetMaxAgeDuration(); age != nil {
					value := strconv.FormatInt(age.AsDuration().Nanoseconds(), 10)
					changes = append(changes, paramsV1Beta.NewParamChange(storageTypes.ModuleNameConsensus.String(), "evidence_max_age_duration", value))
				}
			}
		default:
			var sdkMsg cosmosTypes.Msg
			if err := codec.UnpackAny(msg.Messages[i], &sdkMsg); err != nil {
				continue // skip unknown message
			}
			module, err := storageTypes.ParseModuleName(moduleFromTypeURL(msg.Messages[i].TypeUrl))
			if err != nil {
				continue // not a module we keep constants for
			}
			raw, err := codec.MarshalJSON(sdkMsg)
			if err != nil {
				return msgType, nil, prpsl, err
			}
			var body struct {
				Params map[string]any `json:"params"`
			}
			if err := paramsJSON.Unmarshal(raw, &body); err != nil {
				return msgType, nil, prpsl, err
			}
			if len(body.Params) > 0 {
				if err := flatten(module.String(), "", body.Params, &changes); err != nil {
					return msgType, nil, prpsl, err
				}
				prpsl.Type = storageTypes.ProposalTypeParamChanged
			}
		}
	}
	if prpsl.Description == "" {
		prpsl.Description = sb.String()
	}
	switch {
	case len(changes) > 0:
		prpsl.Changes, err = json.Marshal(changes)
		if err != nil {
			return msgType, nil, nil, errors.Wrap(err, "marshalling changes proposal v1")
		}
	case len(recoveries) > 0:
		prpsl.Changes, err = json.Marshal(recoveries)
		if err != nil {
			return msgType, nil, nil, errors.Wrap(err, "marshalling client recoveries proposal v1")
		}
	}
	return msgType, nil, prpsl, nil
}

// MsgSubmitProposalV1Beta
func MsgSubmitProposalV1Beta(ctx *context.Context, codec codec.Codec, status storageTypes.Status, msgId uint64, msg *cosmosGovTypesV1Beta1.MsgSubmitProposal) (storageTypes.MsgType, any, *storage.Proposal, error) {
	msgType := storageTypes.MsgSubmitProposal
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeProposer, address: msg.Proposer},
	}, ctx.Block.Height, msgId)
	if err != nil {
		return msgType, nil, nil, err
	}
	if status != storageTypes.StatusSuccess {
		return msgType, nil, nil, nil
	}

	prpsl := &storage.Proposal{
		Height: ctx.Block.Height,
		Proposer: &storage.Address{
			Address: msg.Proposer,
		},
		CreatedAt: ctx.Block.Time,
		Status:    storageTypes.ProposalStatusInactive,
	}

	switch msg.Content.TypeUrl {
	case "/cosmos.gov.v1beta1.TextProposal":
		var proposal cosmosGovTypesV1Beta1.TextProposal
		if err := proposal.Unmarshal(msg.Content.Value); err != nil {
			return msgType, nil, nil, errors.Wrap(err, "unmarshalling text proposal for submit proposal content")
		}
		prpsl.Title = proposal.Title
		prpsl.Description = proposal.Description
		prpsl.Type = storageTypes.ProposalTypeText
		return msgType, proposal, prpsl, nil
	case "/cosmos.params.v1beta1.ParameterChangeProposal":
		var proposal paramsV1Beta.ParameterChangeProposal
		if err := proposal.Unmarshal(msg.Content.Value); err != nil {
			return msgType, nil, nil, errors.Wrap(err, "unmarshalling param change proposal for submit proposal content")
		}
		prpsl.Title = proposal.Title
		prpsl.Description = proposal.Description
		prpsl.Type = storageTypes.ProposalTypeParamChanged
		changes := make([]paramsV1Beta.ParamChange, 0, len(proposal.Changes))
		for i := range proposal.Changes {
			moduleName, err := storageTypes.ParseModuleName(proposal.Changes[i].GetSubspace())
			if err != nil {
				return msgType, nil, nil, errors.Wrapf(err, "parsing module name in proposal changes: %s", proposal.Changes[i].GetSubspace())
			}
			key := proposal.Changes[i].GetKey()
			value := proposal.Changes[i].GetValue()

			switch moduleName {
			case storageTypes.ModuleNameConsensus, storageTypes.ModuleNameBaseapp:

				switch key {
				case "BlockParams":
					if err := parseParamsToConstants(storageTypes.ModuleNameConsensus, "block_", value, &changes); err != nil {
						return msgType, nil, nil, errors.Wrap(err, "parse block params")
					}
				case "EvidenceParams":
					if err := parseParamsToConstants(storageTypes.ModuleNameConsensus, "evidence_", value, &changes); err != nil {
						return msgType, nil, nil, errors.Wrap(err, "parse evidence params")
					}
				case "ValidatorParams":
					if err := parseParamsToConstants(storageTypes.ModuleNameConsensus, "validator_", value, &changes); err != nil {
						return msgType, nil, nil, errors.Wrap(err, "parse validator params")
					}
				}

			case storageTypes.ModuleNameGov:
				if key == "votingparams" {
					if err := parseParamsToConstants(moduleName, "", value, &changes); err != nil {
						return msgType, nil, nil, errors.Wrap(err, "parse voting params")
					}
				}

			case storageTypes.ModuleNameBlob:
				val := value
				if key == "GovMaxSquareSize" {
					val, err = strconv.Unquote(value)
					if err != nil {
						return msgType, nil, nil, errors.Wrap(err, value)
					}
				}
				change := paramsV1Beta.NewParamChange(
					proposal.Changes[i].GetSubspace(),
					strcase.SnakeCase(key),
					val,
				)
				changes = append(changes, change)

			default:
				change := paramsV1Beta.NewParamChange(
					proposal.Changes[i].GetSubspace(),
					strcase.SnakeCase(key),
					value,
				)
				changes = append(changes, change)

			}
		}
		data, err := paramsJSON.Marshal(changes)
		if err != nil {
			return msgType, nil, nil, errors.Wrap(err, "marshalling changes proposal for submit proposal content")
		}
		prpsl.Changes = data

		return msgType, proposal, prpsl, nil
	case "/ibc.core.client.v1.ClientUpdateProposal":
		var proposal ibcTypes.ClientUpdateProposal //nolint:staticcheck
		if err := proposal.Unmarshal(msg.Content.Value); err != nil {
			return msgType, nil, nil, errors.Wrap(err, "unmarshalling client update proposal for submit proposal content")
		}
		prpsl.Title = proposal.Title
		prpsl.Description = proposal.Description
		prpsl.Type = storageTypes.ProposalTypeClientUpdate
		// a list, like v1 MsgRecoverClient proposals
		prpsl.Changes, err = json.Marshal([]map[string]any{{
			"SubjectClientId":    proposal.SubjectClientId,
			"SubstituteClientId": proposal.SubstituteClientId,
		}})
		if err != nil {
			return msgType, nil, nil, errors.Wrap(err, "marshalling changes proposal for submit proposal content")
		}
		return msgType, proposal, prpsl, nil

	case "/cosmos.distribution.v1beta1.CommunityPoolSpendProposal":
		var proposal distributionTypes.CommunityPoolSpendProposal //nolint:staticcheck
		if err := proposal.Unmarshal(msg.Content.Value); err != nil {
			return msgType, nil, nil, errors.Wrap(err, "unmarshalling community pool spend proposal for submit proposal content")
		}
		prpsl.Title = proposal.Title
		prpsl.Description = proposal.Description
		prpsl.Type = storageTypes.ProposalTypeCommunityPoolSpend
		prpsl.Changes, err = json.Marshal(map[string]any{
			"Recipient": proposal.Recipient,
			"Amount":    proposal.Amount,
		})
		if err != nil {
			return msgType, nil, nil, errors.Wrap(err, "marshalling changes proposal for submit proposal content")
		}
		return msgType, proposal, prpsl, nil

	default:
		return msgType, nil, nil, errors.Errorf("unknown content type in submit proposal: %s", msg.Content.TypeUrl)
	}
}

// MsgExecLegacyContent is used to wrap the legacy content field into a message.
// This ensures backwards compatibility with v1beta1.MsgSubmitProposal.
func MsgExecLegacyContent(ctx *context.Context, msgId uint64, m *v1.MsgExecLegacyContent) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgExecLegacyContent
	err := createAddresses(
		ctx,
		addressesData{
			{t: storageTypes.MsgAddressTypeAuthority, address: m.Authority},
		}, ctx.Block.Height, msgId)
	return msgType, err
}

// MsgVote defines a message to cast a vote.
func MsgVote(ctx *context.Context, msgId uint64, voterAddress string) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgVote
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeVoter, address: voterAddress},
	}, ctx.Block.Height, msgId)
	return msgType, err
}

// MsgVoteWeighted defines a message to cast a vote.
func MsgVoteWeighted(ctx *context.Context, msgId uint64, voterAddress string) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgVoteWeighted
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeVoter, address: voterAddress},
	}, ctx.Block.Height, msgId)
	return msgType, err
}

// MsgDeposit defines a message to submit a deposit to an existing proposal.
func MsgDeposit(ctx *context.Context, msgId uint64, depositorAddress string) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgDeposit
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeDepositor, address: depositorAddress},
	}, ctx.Block.Height, msgId)
	return msgType, err
}

func MsgUpdateParamsGov(ctx *context.Context, msgId uint64, m *v1.MsgUpdateParams) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgUpdateParams
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeAuthority, address: m.Authority},
	}, ctx.Block.Height, msgId)
	return msgType, err
}

func parseParamsToConstants(moduleName storageTypes.ModuleName, keyPrefix, value string, changes *[]paramsV1Beta.ParamChange) error {
	var params map[string]string
	if err := json.Unmarshal([]byte(value), &params); err != nil {
		return errors.Wrap(err, "unmarshal params")
	}
	for _, key := range slices.Sorted(maps.Keys(params)) {
		change := paramsV1Beta.NewParamChange(
			moduleName.String(),
			keyPrefix+key,
			params[key],
		)
		*changes = append(*changes, change)
	}
	return nil
}
