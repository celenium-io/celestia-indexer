// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package stats

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
)

type testType struct {
	bun.BaseModel `bun:"test"`

	Field1 uint64 `bun:"field1"     stats:"func:avg min,filterable"`
	Field2 uint64 `bun:",notnull"   stats:"func:max"`
	Field3 uint64 `stats:"func:sum"`
	Field4 uint64 `bun:"field4"`
	Field5 uint64 `bun:"field5"     stats:"-"`
	Field6 uint64 `bun:"-"          stats:"func:min"`
}

type testFailedType struct {
	bun.BaseModel `bun:"test"`

	Field1 uint64 `bun:"field1"     stats:"func:avg min,filterable"`
	Field2 uint64 `bun:",notnull"   stats:"func:invalid"`
	Field3 uint64 `stats:"func:sum"`
	Field4 uint64 `bun:"field4"`
}

func TestInit(t *testing.T) {
	tests := []struct {
		name    string
		models  []any
		wantErr bool
		want    map[string]Table
	}{
		{
			name: "test 1",
			models: []any{
				&testType{},
			},
			wantErr: false,
			want: map[string]Table{
				"test": {
					Columns: map[string]Column{
						"field1": {
							Functions: map[string]struct{}{
								"avg": {},
								"min": {},
							},
							Filterable: true,
						},
						"field2": {
							Functions: map[string]struct{}{"max": {}},
						},
					},
				},
			},
		}, {
			name: "test 2",
			models: []any{
				&testFailedType{},
			},
			wantErr: true,
			want:    make(map[string]Table),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Tables = make(map[string]Table)
			err := Init(tt.models...)
			require.Equal(t, tt.wantErr, err != nil)
			require.Equal(t, tt.want, Tables)
		})
	}
}
