// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package decoder

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCleanString(t *testing.T) {
	for _, tt := range []struct {
		name string
		in   string
		want string
	}{
		{name: "empty", in: "", want: ""},
		{name: "no nul", in: "memo", want: "memo"},
		{name: "nul inside", in: "me\x00mo", want: "memo"},
		{name: "only nul", in: "\x00\x00", want: ""},
		{name: "unicode", in: "валидатор\x00", want: "валидатор"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, CleanString(tt.in))
		})
	}
}
