package responses

import "github.com/btcsuite/btcutil/bech32"

func DecodeAddress(address string) ([]byte, error) {
	_, data, err := bech32.Decode(address)
	return data, err
}

func EncodeAddress(hash []byte) (string, error) {
	return bech32.Encode("celestia", hash)
}
