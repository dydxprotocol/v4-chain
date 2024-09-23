package types

import (
	"context"
	"math/big"
	"time"

	perptypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/perpetuals/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types" //nolint:staticcheck
	ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"
)

// BankKeeper defines the expected bank keeper used for simulations.
type BankKeeper interface {
	GetSupply(ctx context.Context, denom string) sdk.Coin
	SendCoinsFromAccountToModule(
		ctx context.Context,
		senderAddr sdk.AccAddress,
		recipientModule string,
		amt sdk.Coins,
	) error
	SendCoinsFromModuleToAccount(ctx context.Context,
		senderModule string,
		recipientAddr sdk.AccAddress,
		amt sdk.Coins,
	) error
	SendCoins(ctx context.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) error
	MintCoins(ctx context.Context, moduleName string, amt sdk.Coins) error
	BurnCoins(ctx context.Context, moduleName string, amt sdk.Coins) error
	GetBalance(ctx context.Context, addr sdk.AccAddress, denom string) sdk.Coin
}

type AssetsKeeper interface {
	ConvertCoinToAsset(ctx sdk.Context, assetId uint32, coin sdk.Coin) (quantums *big.Int, convertedDenom *big.Int, err error)
	ConvertAssetToCoin(ctx sdk.Context, assetId uint32, quantums *big.Int) (convertedQuantums *big.Int, coin sdk.Coin, err error)
}

type BlockTimeKeeper interface {
	GetTimeSinceLastBlock(ctx sdk.Context) time.Duration
}

type PerpetualsKeeper interface {
	UpdateYieldIndexToNewMint(ctx sdk.Context, totalTDaiPreMint *big.Int, totalTDaiMinted *big.Int) error
	GetAllPerpetuals(ctx sdk.Context) (list []perptypes.Perpetual)
	GetInsuranceFundModuleAddress(ctx sdk.Context, perpetualId uint32) (sdk.AccAddress, error)
}

// ICS4Wrapper defines the expected ICS4Wrapper for middleware
type ICS4Wrapper interface {
	WriteAcknowledgement(
		ctx sdk.Context,
		chanCap *capabilitytypes.Capability,
		packet ibcexported.PacketI,
		acknowledgement ibcexported.Acknowledgement,
	) error
	SendPacket(
		ctx sdk.Context,
		chanCap *capabilitytypes.Capability,
		sourcePort string,
		sourceChannel string,
		timeoutHeight clienttypes.Height,
		timeoutTimestamp uint64,
		data []byte,
	) (sequence uint64, err error)
	GetAppVersion(ctx sdk.Context, portID string, channelID string) (string, bool)
}
