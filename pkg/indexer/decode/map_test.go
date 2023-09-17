package decode

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecimalFromMap(t *testing.T) {
	tests := []struct {
		name string
		m    map[string]any
		key  string
		want string
	}{
		{
			name: "test 1",
			m: map[string]any{
				"amount": "123utia",
			},
			key:  "amount",
			want: "123",
		}, {
			name: "test 2",
			m: map[string]any{
				"amount": "123utia",
			},
			key:  "invalid",
			want: "0",
		}, {
			name: "test 3",
			m: map[string]any{
				"amount": "123uta",
			},
			key:  "amount",
			want: "123",
		}, {
			name: "test 4",
			m: map[string]any{
				"amount": 123,
			},
			key:  "amount",
			want: "0",
		}, {
			name: "test 5",
			m: map[string]any{
				"amount": "123test",
			},
			key:  "amount",
			want: "123",
		}, {
			name: "test 6",
			m: map[string]any{
				"amount": "1-23test",
			},
			key:  "amount",
			want: "0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DecimalFromMap(tt.m, tt.key)
			require.Equal(t, tt.want, got.String())
		})
	}
}
