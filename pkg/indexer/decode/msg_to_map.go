// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package decode

import (
	"github.com/bytedance/sonic"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/ettle/strcase"
)

var jsonApi = sonic.ConfigFastest

// MsgToMap converts a cosmos Msg to PackedBytes using the codec JSON marshaler.
// Exported for use in tests.
func MsgToMap(msg cosmosTypes.Msg) (storageTypes.PackedBytes, error) {
	return msgToMap(msg)
}

func msgToMap(msg cosmosTypes.Msg) (storageTypes.PackedBytes, error) {
	b, err := cfg.Codec.MarshalJSON(msg)
	if err != nil {
		// Fallback for messages containing unresolvable Any types (e.g. empty TypeUrl).
		b, err = jsonApi.Marshal(msg)
		if err != nil {
			return nil, err
		}
	}
	var raw map[string]any
	if err := jsonApi.Unmarshal(b, &raw); err != nil {
		return nil, err
	}
	return keysToPascalCase(raw), nil
}

// keysToPascalCase converts top-level map keys from snake_case to PascalCase.
// Nested maps are converted via convertMapKeys which returns plain map[string]any
// so that type assertions .(map[string]any) work throughout the codebase.
func keysToPascalCase(m map[string]any) storageTypes.PackedBytes {
	result := make(storageTypes.PackedBytes, len(m))
	for k, v := range m {
		result[strcase.ToPascal(k)] = transformValue(v)
	}
	return result
}

func convertMapKeys(m map[string]any) map[string]any {
	result := make(map[string]any, len(m))
	for k, v := range m {
		result[strcase.ToPascal(k)] = transformValue(v)
	}
	return result
}

func transformValue(v any) any {
	switch val := v.(type) {
	case map[string]any:
		return convertMapKeys(val)
	case []any:
		for i, item := range val {
			val[i] = transformValue(item)
		}
	}
	return v
}
