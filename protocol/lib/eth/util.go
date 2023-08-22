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

// BridgeLogToEvent converts an Ethereum log from Bridge contract to a BridgeEvent.
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
	address := bridgeEventData[2].([]byte)

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
