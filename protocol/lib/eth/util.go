package eth

import (
	"math/big"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/app/config"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/constants"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	bridgetypes "github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
	ethabi "github.com/ethereum/go-ethereum/accounts/abi"
	ethcoretypes "github.com/ethereum/go-ethereum/core/types"
)

var bridgeEventAbi *ethabi.ABI

const (
	MIN_ADDRESS_BYTES = 20
	MAX_ADDRESS_BYTES = 32
)

// getBridgeEventAbi returns the ABI (application binary interface) for the Bridge contract.
func getBridgeEventAbi() *ethabi.ABI {
	// Initialize the singleton if it does not exist.
	if bridgeEventAbi == nil {
		bridgeAbi, err := ethabi.JSON(strings.NewReader(constants.BridgeEventABI))
		if err != nil {
			panic(err)
		}
		bridgeEventAbi = &bridgeAbi
	}
	return bridgeEventAbi
}

// padOrTruncateAddress right-pads an address with zeros if it's shorter than `MIN_ADDRESS_BYTES` or
// takes the first `MAX_ADDRESS_BYTES` if it's longer than that.
func padOrTruncateAddress(address []byte) []byte {
	if len(address) > MAX_ADDRESS_BYTES {
		return address[:MAX_ADDRESS_BYTES]
	} else if len(address) < MIN_ADDRESS_BYTES {
		return append(address, make([]byte, MIN_ADDRESS_BYTES-len(address))...)
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
	bridgeEventData, err := getBridgeEventAbi().Unpack("Bridge", log.Data)
	if err != nil {
		panic(err)
	}
	amount := bridgeEventData[0].(*big.Int)
	address := padOrTruncateAddress(bridgeEventData[2].([]byte))

	// Unused daemon fields.
	// bridgeEventData[1] is the Ethereum address that sent the tokens
	// bridgeEventData[3] is the user-supplied memo

	return bridgetypes.BridgeEvent{
		Id:             id,
		Coin:           sdk.NewCoin(denom, sdk.NewIntFromBigInt(amount)),
		Address:        sdk.MustBech32ifyAddressBytes(config.Bech32PrefixAccAddr, address),
		EthBlockHeight: log.BlockNumber,
	}
}
