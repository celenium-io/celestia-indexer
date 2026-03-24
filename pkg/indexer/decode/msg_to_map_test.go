// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package decode

import (
	"testing"

	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/stretchr/testify/require"
)

func Test_snakeToPascal(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"from_address", "FromAddress"},
		{"to_address", "ToAddress"},
		{"amount", "Amount"},
		{"denom", "Denom"},
		{"client_id", "ClientId"},
		{"packet_sequence", "PacketSequence"},
		{"source_port", "SourcePort"},
		{"", ""},
		{"already", "Already"},
		{"AlreadyPascal", "AlreadyPascal"},
		{"a_b_c", "ABC"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			require.Equal(t, tt.want, transformKey(tt.input))
		})
	}
}

func transformKey(s string) string {
	// reuse the same logic via convertMapKeys with a single key
	result := convertMapKeys(map[string]any{s: nil})
	for k := range result {
		return k
	}
	return ""
}

func Test_convertMapKeys_flat(t *testing.T) {
	input := map[string]any{
		"from_address": "celestia1abc",
		"to_address":   "celestia1def",
		"amount":       "1000",
	}
	got := convertMapKeys(input)
	require.Equal(t, map[string]any{
		"FromAddress": "celestia1abc",
		"ToAddress":   "celestia1def",
		"Amount":      "1000",
	}, got)
}

func Test_convertMapKeys_nested_map(t *testing.T) {
	input := map[string]any{
		"outer_key": map[string]any{
			"inner_key": "value",
		},
	}
	got := convertMapKeys(input)
	require.Equal(t, map[string]any{
		"OuterKey": map[string]any{
			"InnerKey": "value",
		},
	}, got)
}

func Test_convertMapKeys_nested_array_of_maps(t *testing.T) {
	input := map[string]any{
		"coin_list": []any{
			map[string]any{"denom": "utia", "amount": "500"},
			map[string]any{"denom": "uatom", "amount": "200"},
		},
	}
	got := convertMapKeys(input)
	require.Equal(t, map[string]any{
		"CoinList": []any{
			map[string]any{"Denom": "utia", "Amount": "500"},
			map[string]any{"Denom": "uatom", "Amount": "200"},
		},
	}, got)
}

func Test_convertMapKeys_nested_array_of_arrays(t *testing.T) {
	input := map[string]any{
		"matrix": []any{
			[]any{
				map[string]any{"inner_key": "v1"},
			},
		},
	}
	got := convertMapKeys(input)
	require.Equal(t, map[string]any{
		"Matrix": []any{
			[]any{
				map[string]any{"InnerKey": "v1"},
			},
		},
	}, got)
}

func Test_convertMapKeys_primitives(t *testing.T) {
	input := map[string]any{
		"str_val":  "hello",
		"num_val":  float64(42),
		"bool_val": true,
		"nil_val":  nil,
	}
	got := convertMapKeys(input)
	require.Equal(t, map[string]any{
		"StrVal":  "hello",
		"NumVal":  float64(42),
		"BoolVal": true,
		"NilVal":  nil,
	}, got)
}

func Test_keysToPascalCase_returns_PackedBytes(t *testing.T) {
	input := map[string]any{
		"from_address": "celestia1abc",
		"nested_obj": map[string]any{
			"inner_key": "val",
		},
	}
	got := keysToPascalCase(input)

	require.IsType(t, storageTypes.PackedBytes{}, got)
	require.Equal(t, "celestia1abc", got["FromAddress"])

	// nested must be plain map[string]any, not PackedBytes
	nested, ok := got["NestedObj"].(map[string]any)
	require.True(t, ok, "nested map must be map[string]any, not PackedBytes")
	require.Equal(t, "val", nested["InnerKey"])
}

func Test_transformValue_passthrough(t *testing.T) {
	require.Equal(t, "str", transformValue("str"))
	require.Equal(t, float64(1), transformValue(float64(1)))
	require.Equal(t, true, transformValue(true))
	require.Nil(t, transformValue(nil))
}

func Test_transformValue_map(t *testing.T) {
	v := transformValue(map[string]any{"snake_key": "val"})
	m, ok := v.(map[string]any)
	require.True(t, ok)
	require.Equal(t, "val", m["SnakeKey"])
}

func Test_transformValue_slice(t *testing.T) {
	v := transformValue([]any{
		map[string]any{"from_address": "addr1"},
		"plain_string",
		float64(99),
	})
	s, ok := v.([]any)
	require.True(t, ok)

	m, ok := s[0].(map[string]any)
	require.True(t, ok)
	require.Equal(t, "addr1", m["FromAddress"])

	require.Equal(t, "plain_string", s[1])
	require.Equal(t, float64(99), s[2])
}
