package responses

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_formatFoat64(t *testing.T) {
	tests := []struct {
		name  string
		value float64
		want  string
	}{
		{
			name:  "test 1",
			value: .0123456,
			want:  "0.0123456",
		}, {
			name:  "test 2",
			value: 1234567.0123456,
			want:  "1234567.0123456",
		}, {
			name:  "test 3",
			value: 1234567.0123456789,
			want:  "1234567.0123456789",
		}, {
			name:  "test 4",
			value: -1234567.0123456789,
			want:  "-1234567.0123456789",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatFoat64(tt.value)
			require.Equal(t, tt.want, got)
		})
	}
}
