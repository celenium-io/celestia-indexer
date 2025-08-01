// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package types

import (
	"encoding/base64"
	"encoding/hex"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddress_Decimal(t *testing.T) {
	tests := []struct {
		name    string
		a       Address
		want    string
		wantErr bool
	}{
		{
			name: "test 1",
			a:    Address("celestia1vsvx8n7f8dh5udesqqhgrjutyun7zqrgehdq2l"),
			want: "571442989789389314489155229272138443602249449576",
		}, {
			name: "test 2",
			a:    Address("celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6"),
			want: "847144225136750676593624093723244210182685458328",
		}, {
			name: "test 3",
			a:    Address("celestiavaloper1fg9l3xvfuu9wxremv2229966zawysg4r40gw5x"),
			want: "422732281331565053903020826810589057922481267363",
		}, {
			name: "test 4",
			a:    Address("celestia1ws4hfsl8hlylt38ptk5cn9ura20slu2fnkre76"),
			want: "663212033147379327340199490680921163866055897417",
		}, {
			name: "test 5",
			a:    Address("celestiavaloper12c6cwd0kqlg48sdhjnn9f0z82g0c82fmrl7j9y"),
			want: "492166924306969235254183868234764112087912524091",
		}, {
			name: "test 6",
			a:    Address("celestia1vysgwc9mykfz5249g9thjlffx6nha0kkwsvs37"),
			want: "554497521690734140602662053179950865913302925014",
		}, {
			name: "test 7",
			a:    Address("celestiavaloper1vysgwc9mykfz5249g9thjlffx6nha0kkt0wf8c"),
			want: "554497521690734140602662053179950865913302925014",
		}, {
			name: "test 8",
			a:    Address("celestia170qq26qenw420ufd5py0r59kpg3tj2m7gl5cja"),
			want: "1391566971373145669833734855247902575571779529598",
		}, {
			name: "test 9",
			a:    Address("celestiavaloper170qq26qenw420ufd5py0r59kpg3tj2m7dqkpym"),
			want: "1391566971373145669833734855247902575571779529598",
		}, {
			name: "test 10",
			a:    Address("celestia1zefjxuq43xmjq9x4hhw23wkvvz6st5uhv40tys"),
			want: "127451845879864064811426441084571049438795387799",
		}, {
			name: "test 11",
			a:    Address("celestia18r6ujzzkg6ku9sr39nxy4847q4qea5kg4a8pxv"),
			want: "325184678146186827908099091810224693115038388936",
		}, {
			name: "test 12",
			a:    Address("celestia1vnflc6322f8z7cpl28r7un5dxhmjxghc20aydq"),
			want: "575626522741514541054636655673170065879262765816",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.a.Decimal()
			require.Equal(t, err != nil, tt.wantErr)
			if err == nil {
				require.Equal(t, tt.want, got.String())
			}
		})
	}
}

func TestAddress_Decode(t *testing.T) {
	tests := []struct {
		name    string
		a       Address
		want    string
		want1   []byte
		wantErr bool
	}{
		{
			name:    "celestia1mm8yykm46ec3t0dgwls70g0jvtm055wk9ayal8",
			a:       Address("celestia1mm8yykm46ec3t0dgwls70g0jvtm055wk9ayal8"),
			want:    AddressPrefixCelestia,
			want1:   []byte{0xde, 0xce, 0x42, 0x5b, 0x75, 0xd6, 0x71, 0x15, 0xbd, 0xa8, 0x77, 0xe1, 0xe7, 0xa1, 0xf2, 0x62, 0xf6, 0xfa, 0x51, 0xd6},
			wantErr: false,
		}, {
			name:    "celestia1jc92qdnty48pafummfr8ava2tjtuhfdw774w60",
			a:       Address("celestia1jc92qdnty48pafummfr8ava2tjtuhfdw774w60"),
			want:    AddressPrefixCelestia,
			want1:   []byte{0x96, 0xa, 0xa0, 0x36, 0x6b, 0x25, 0x4e, 0x1e, 0xa7, 0x9b, 0xda, 0x46, 0x7e, 0xb3, 0xaa, 0x5c, 0x97, 0xcb, 0xa5, 0xae},
			wantErr: false,
		}, {
			name:    "celestia1vsvx8n7f8dh5udesqqhgrjutyun7zqrgehdq2l",
			a:       Address("celestia1vsvx8n7f8dh5udesqqhgrjutyun7zqrgehdq2l"),
			want:    AddressPrefixCelestia,
			want1:   []byte{0x64, 0x18, 0x63, 0xcf, 0xc9, 0x3b, 0x6f, 0x4e, 0x37, 0x30, 0x0, 0x2e, 0x81, 0xcb, 0x8b, 0x27, 0x27, 0xe1, 0x0, 0x68},
			wantErr: false,
		}, {
			name:    "celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6",
			a:       Address("celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6"),
			want:    AddressPrefixCelestia,
			want1:   []byte{0x94, 0x63, 0x42, 0xc7, 0x73, 0x2f, 0x6, 0x22, 0x37, 0xb6, 0x73, 0xf5, 0x26, 0xf2, 0x67, 0x4c, 0x7d, 0xb, 0xff, 0x98},
			wantErr: false,
		}, {
			name:    "celestiavaloper1fg9l3xvfuu9wxremv2229966zawysg4r40gw5x",
			a:       Address("celestiavaloper1fg9l3xvfuu9wxremv2229966zawysg4r40gw5x"),
			want:    AddressPrefixValoper,
			want1:   []byte{0x4a, 0xb, 0xf8, 0x99, 0x89, 0xe7, 0xa, 0xe3, 0xf, 0x3b, 0x62, 0x94, 0xa2, 0x97, 0x5a, 0x17, 0x5c, 0x48, 0x22, 0xa3},
			wantErr: false,
		}, {
			name:    "celestia1ws4hfsl8hlylt38ptk5cn9ura20slu2fnkre76",
			a:       Address("celestia1ws4hfsl8hlylt38ptk5cn9ura20slu2fnkre76"),
			want:    AddressPrefixCelestia,
			want1:   []byte{0x74, 0x2b, 0x74, 0xc3, 0xe7, 0xbf, 0xc9, 0xf5, 0xc4, 0xe1, 0x5d, 0xa9, 0x89, 0x97, 0x83, 0xea, 0x9f, 0xf, 0xf1, 0x49},
			wantErr: false,
		}, {
			name:    "celestiavaloper12c6cwd0kqlg48sdhjnn9f0z82g0c82fmrl7j9y",
			a:       Address("celestiavaloper12c6cwd0kqlg48sdhjnn9f0z82g0c82fmrl7j9y"),
			want:    AddressPrefixValoper,
			want1:   []byte{0x56, 0x35, 0x87, 0x35, 0xf6, 0x7, 0xd1, 0x53, 0xc1, 0xb7, 0x94, 0xe6, 0x54, 0xbc, 0x47, 0x52, 0x1f, 0x83, 0xa9, 0x3b},
			wantErr: false,
		}, {
			name:    "celestia1vysgwc9mykfz5249g9thjlffx6nha0kkwsvs37",
			a:       Address("celestia1vysgwc9mykfz5249g9thjlffx6nha0kkwsvs37"),
			want:    AddressPrefixCelestia,
			want1:   []byte{0x61, 0x20, 0x87, 0x60, 0xbb, 0x25, 0x92, 0x2a, 0x2a, 0xa5, 0x41, 0x57, 0x79, 0x7d, 0x29, 0x36, 0xa7, 0x7e, 0xbe, 0xd6},
			wantErr: false,
		}, {
			name:    "celestiavaloper1vysgwc9mykfz5249g9thjlffx6nha0kkt0wf8c",
			a:       Address("celestiavaloper1vysgwc9mykfz5249g9thjlffx6nha0kkt0wf8c"),
			want:    AddressPrefixValoper,
			want1:   []byte{0x61, 0x20, 0x87, 0x60, 0xbb, 0x25, 0x92, 0x2a, 0x2a, 0xa5, 0x41, 0x57, 0x79, 0x7d, 0x29, 0x36, 0xa7, 0x7e, 0xbe, 0xd6},
			wantErr: false,
		}, {
			name:    "celestiavaloper170qq26qenw420ufd5py0r59kpg3tj2m7dqkpym",
			a:       Address("celestiavaloper170qq26qenw420ufd5py0r59kpg3tj2m7dqkpym"),
			want:    AddressPrefixValoper,
			want1:   []byte{0xf3, 0xc0, 0x5, 0x68, 0x19, 0x9b, 0xaa, 0xa7, 0xf1, 0x2d, 0xa0, 0x48, 0xf1, 0xd0, 0xb6, 0xa, 0x22, 0xb9, 0x2b, 0x7e},
			wantErr: false,
		}, {
			name:    "celestia1zefjxuq43xmjq9x4hhw23wkvvz6st5uhv40tys",
			a:       Address("celestia1zefjxuq43xmjq9x4hhw23wkvvz6st5uhv40tys"),
			want:    AddressPrefixCelestia,
			want1:   []byte{0x16, 0x53, 0x23, 0x70, 0x15, 0x89, 0xb7, 0x20, 0x14, 0xd5, 0xbd, 0xdc, 0xa8, 0xba, 0xcc, 0x60, 0xb5, 0x5, 0xd3, 0x97},
			wantErr: false,
		}, {
			name:    "celestia18r6ujzzkg6ku9sr39nxy4847q4qea5kg4a8pxv",
			a:       Address("celestia18r6ujzzkg6ku9sr39nxy4847q4qea5kg4a8pxv"),
			want:    AddressPrefixCelestia,
			want1:   []byte{0x38, 0xf5, 0xc9, 0x8, 0x56, 0x46, 0xad, 0xc2, 0xc0, 0x71, 0x2c, 0xcc, 0x4a, 0x9e, 0xbe, 0x5, 0x41, 0x9e, 0xd2, 0xc8},
			wantErr: false,
		}, {
			name:    "celestia1vnflc6322f8z7cpl28r7un5dxhmjxghc20aydq",
			a:       Address("celestia1vnflc6322f8z7cpl28r7un5dxhmjxghc20aydq"),
			want:    AddressPrefixCelestia,
			want1:   []byte{0x64, 0xd3, 0xfc, 0x6a, 0x2a, 0x52, 0x4e, 0x2f, 0x60, 0x3f, 0x51, 0xc7, 0xee, 0x4e, 0x8d, 0x35, 0xf7, 0x23, 0x22, 0xf8},
			wantErr: false,
		}, {
			name:    "celestia1xzsdn65hyljcmenlxyjmdmvghhd0w4ut27k3fx56jp2p69eh6srs8p3rss",
			a:       Address("celestia1xzsdn65hyljcmenlxyjmdmvghhd0w4ut27k3fx56jp2p69eh6srs8p3rss"),
			want:    AddressPrefixCelestia,
			want1:   []byte{48, 160, 217, 234, 151, 39, 229, 141, 230, 127, 49, 37, 182, 237, 136, 189, 218, 247, 87, 139, 87, 173, 20, 154, 154, 144, 84, 29, 23, 55, 212, 7},
			wantErr: false,
		}, {
			name:    "celestia1vdjkcetnw35kzvtxd3chzwrjdc6hw7tt0qmkkctnwea8g6rnwfehsdrwwpkk5epkdf6nvdmswdcsaw2mhz",
			a:       Address("celestia1vdjkcetnw35kzvtxd3chzwrjdc6hw7tt0qmkkctnwea8g6rnwfehsdrwwpkk5epkdf6nvdmswdcsaw2mhz"),
			want:    AddressPrefixCelestia,
			want1:   []byte{99, 101, 108, 101, 115, 116, 105, 97, 49, 102, 108, 113, 113, 56, 114, 110, 53, 119, 121, 107, 120, 55, 107, 97, 115, 118, 122, 116, 104, 115, 114, 115, 120, 52, 110, 112, 109, 106, 100, 54, 106, 117, 54, 55, 112, 115, 113},
			wantErr: false,
		}, {
			name:    "celestia13djgs6acn677nlygmsnc7m5qst2s0llaa8wdzmjj04zn3785elzqulvmvt",
			a:       Address("celestia13djgs6acn677nlygmsnc7m5qst2s0llaa8wdzmjj04zn3785elzqulvmvt"),
			want:    AddressPrefixCelestia,
			want1:   []byte{139, 100, 136, 107, 184, 158, 189, 233, 252, 136, 220, 39, 143, 110, 128, 130, 213, 7, 255, 253, 233, 220, 209, 110, 82, 125, 69, 56, 248, 244, 207, 196},
			wantErr: false,
		}, {
			name:    "celestia1lumaemzetgx7fa8gcw3waagevs6rqhejfwxdvszzxlneselshwxqm5rmxn",
			a:       Address("celestia1lumaemzetgx7fa8gcw3waagevs6rqhejfwxdvszzxlneselshwxqm5rmxn"),
			want:    AddressPrefixCelestia,
			want1:   []byte{255, 55, 220, 236, 89, 90, 13, 228, 244, 232, 195, 162, 238, 245, 25, 100, 52, 48, 95, 50, 75, 140, 214, 64, 66, 55, 231, 152, 103, 240, 187, 140},
			wantErr: false,
		}, {
			name:    "celestia1vdjkcetnw35kzvf50f3njuphwpshgu3nxfmxuum90fn85cenvehrjerww9shy6rvd4mxx7t4w3nqtn9hsr",
			a:       Address("celestia1vdjkcetnw35kzvf50f3njuphwpshgu3nxfmxuum90fn85cenvehrjerww9shy6rvd4mxx7t4w3nqtn9hsr"),
			want:    AddressPrefixCelestia,
			want1:   []byte{99, 101, 108, 101, 115, 116, 105, 97, 49, 52, 122, 99, 57, 112, 55, 112, 97, 116, 114, 51, 50, 118, 110, 115, 101, 122, 102, 122, 99, 51, 102, 110, 57, 100, 110, 113, 97, 114, 104, 108, 109, 118, 99, 121, 117, 116, 102},
			wantErr: false,
		}, {
			name:    "celestia15xvs73zpnfkrkdvm0k5wrdfz2kjfhd5vz9nnujyvagpq9sqjh3jqkp2tp3",
			a:       Address("celestia15xvs73zpnfkrkdvm0k5wrdfz2kjfhd5vz9nnujyvagpq9sqjh3jqkp2tp3"),
			want:    AddressPrefixCelestia,
			want1:   []byte{0xa1, 0x99, 0xf, 0x44, 0x41, 0x9a, 0x6c, 0x3b, 0x35, 0x9b, 0x7d, 0xa8, 0xe1, 0xb5, 0x22, 0x55, 0xa4, 0x9b, 0xb6, 0x8c, 0x11, 0x67, 0x3e, 0x48, 0x8c, 0xea, 0x2, 0x2, 0xc0, 0x12, 0xbc, 0x64},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := tt.a.Decode()
			require.Equal(t, err != nil, tt.wantErr)
			if err == nil {
				require.Equal(t, tt.want, got)
				require.Equal(t, tt.want1, got1)
			}
		})
	}
}

func TestGetConsAddressBytesFromPubKey(t *testing.T) {
	tests := []struct {
		name string
		data string
		want string
	}{
		{
			name: "test 1",
			data: "BAYGw7aWpwoWogn2t5JZdEXvu3g6CEjgpXf4aq7JjmA=",
			want: "5BE36E57293ECCF5BE945A114F21C379B6C5F496",
		}, {
			name: "test 2",
			data: "hykgLmpkKVGwfjJRxafw9ymRmckZ8b8bQ7LS6rrWi58=",
			want: "D74EC29E6E4597943942E7E97B1F519A0615E3B4",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := base64.StdEncoding.DecodeString(tt.data)
			require.NoError(t, err)

			got := GetConsAddressBytesFromPubKey(b)
			require.Equal(t, tt.want, strings.ToUpper(hex.EncodeToString(got)))
		})
	}
}
