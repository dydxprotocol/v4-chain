package address

import (
	"github.com/cosmos/cosmos-sdk/types/bech32"
)

func ConvertAddressPrefix(oldAddress string, newPrefix string) (string, error) {
	_, bz, err := bech32.DecodeAndConvert(oldAddress)
	if err != nil {
		return "", err
	}

	newAddress, err := bech32.ConvertAndEncode(newPrefix, bz)
	if err != nil {
		return "", err
	}

	return newAddress, nil
}
