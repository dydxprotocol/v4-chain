package eth

import (
	"math/big"
	"strings"
	"sync"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/app/config"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/constants"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	bridgetypes "github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
	ethabi "github.com/ethereum/go-ethereum/accounts/abi"
	ethcoretypes "github.com/ethereum/go-ethereum/core/types"
)

const (
	MinAddrLen = 20
	MaxAddrLen = 32
)

// bridgeEventAbi is the ABI (application binary interface) for the Bridge contract.
// It is initialized at most once.
var bridgeEventAbi = sync.OnceValue[*ethabi.ABI](
	func() *ethabi.ABI {
		bAbi, err := ethabi.JSON(strings.NewReader(constants.BridgeEventABI))
		if err != nil {
			panic(err)
		}
		return &bAbi
	},
)

// GetBridgeEventAbi returns the ABI (application binary interface) for the Bridge contract.
func GetBridgeEventAbi() *ethabi.ABI {
	return bridgeEventAbi()
}

// PadOrTruncateAddress right-pads an address with zeros if it's shorter than `MinAddrLen` or
// takes the first `MaxAddrLen` if it's longer than that.
func PadOrTruncateAddress(address []byte) []byte {
	if len(address) > MaxAddrLen {
		return address[:MaxAddrLen]
	} else if len(address) < MinAddrLen {
		return append(address, make([]byte, MinAddrLen-len(address))...)
	}
	return address
}

/*
BridgeLogToEvent converts an Ethereum log from Bridge contract to a BridgeEvent.
Note: The format of a dYdX address is [prefix][separator][address][checksum], where `prefix` is `dydx`,
`separator` is `1`, `address` is the actual address portion, and `checksum` occupies last 6 characters.
An address in Ethereum logs is in hexadecimal format and in Cosmos bech32 format. For example, a
20-byte address in hexadecimal format is 20*8=160 bits, which is 160/5=32 bech32 characters.
*/
func BridgeLogToEvent(
	log ethcoretypes.Log,
	denom string,
) bridgetypes.BridgeEvent {
	// Unpack the topics.
	id := lib.MustConvertIntegerToUint32(log.Topics[1].Big().Uint64())

	// Unpack the data.
	bridgeEventData, err := GetBridgeEventAbi().Unpack("Bridge", log.Data)
	if err != nil {
		panic(err)
	}
	amount := bridgeEventData[0].(*big.Int)
	address := PadOrTruncateAddress(bridgeEventData[2].([]byte))

	// Unused daemon fields.
	// bridgeEventData[1] is the Ethereum address that sent the tokens
	// bridgeEventData[3] is the user-supplied memo

	return bridgetypes.BridgeEvent{
		Id:             id,
		Coin:           sdk.NewCoin(denom, sdkmath.NewIntFromBigInt(amount)),
		Address:        sdk.MustBech32ifyAddressBytes(config.Bech32PrefixAccAddr, address),
		EthBlockHeight: log.BlockNumber,
	}
}
