package handler

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_isAddress(t *testing.T) {
	tests := []struct {
		name    string
		address string
		want    bool
	}{
		{
			name:    "test 1",
			address: "celestia12y6fchaufs4tmn8e8wlk3rtrrftpqp6vr228a7",
			want:    true,
		}, {
			name:    "test 2",
			address: "celestiavaloper1qycj0ymu9fqvwgyw4xz93p3n4a83jjk7sm2wzh",
			want:    true,
		}, {
			name:    "test 3",
			address: "invalid",
			want:    false,
		}, {
			name:    "test 4",
			address: "celestiavaloper111111111111111111111111111111111111111",
			want:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isAddress(tt.address)
			require.Equal(t, tt.want, got, tt.name)
		})
	}
}
