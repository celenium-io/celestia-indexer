package parser

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_getDecimalFromMap(t *testing.T) {
	tests := []struct {
		name string
		m    map[string]any
		key  string
		want string
	}{
		{
			name: "test 1",
			m:    map[string]any{},
			key:  "any_key",
			want: "0",
		}, {
			name: "test 2",
			m: map[string]any{
				"any_key": "qwertyui",
			},
			key:  "any_key",
			want: "0",
		}, {
			name: "test 3",
			m: map[string]any{
				"any_key": 123123,
			},
			key:  "any_key",
			want: "0",
		}, {
			name: "test 4",
			m: map[string]any{
				"any_key": "123123",
			},
			key:  "any_key",
			want: "123123",
		}, {
			name: "test 5",
			m: map[string]any{
				"any_key": "123123utia",
			},
			key:  "any_key",
			want: "123123",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getDecimalFromMap(tt.m, tt.key)
			require.Equal(t, tt.want, got.String())
		})
	}
}
