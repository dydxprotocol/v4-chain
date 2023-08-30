package constants

import (
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/app/config"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

func init() {
	// This package does not contain the `app/config` package in its import chain, and therefore needs to call
	// SetAddressPrefixes() explicitly in order to set the `dydx` address prefixes.
	config.SetAddressPrefixes()
}

var (
	Alice_Num0 = satypes.SubaccountId{
		Owner:  types.MustBech32ifyAddressBytes(config.Bech32PrefixAccAddr, AliceAccAddress),
		Number: 0,
	}
	Alice_Num1 = satypes.SubaccountId{
		Owner:  types.MustBech32ifyAddressBytes(config.Bech32PrefixAccAddr, AliceAccAddress),
		Number: 1,
	}
	Bob_Num0 = satypes.SubaccountId{
		Owner:  types.MustBech32ifyAddressBytes(config.Bech32PrefixAccAddr, BobAccAddress),
		Number: 0,
	}
	Bob_Num1 = satypes.SubaccountId{
		Owner:  types.MustBech32ifyAddressBytes(config.Bech32PrefixAccAddr, BobAccAddress),
		Number: 1,
	}
	Bob_Num2 = satypes.SubaccountId{
		Owner:  types.MustBech32ifyAddressBytes(config.Bech32PrefixAccAddr, BobAccAddress),
		Number: 2,
	}
	Carl_Num0 = satypes.SubaccountId{
		Owner:  types.MustBech32ifyAddressBytes(config.Bech32PrefixAccAddr, CarlAccAddress),
		Number: 0,
	}
	Carl_Num1 = satypes.SubaccountId{
		Owner:  types.MustBech32ifyAddressBytes(config.Bech32PrefixAccAddr, CarlAccAddress),
		Number: 1,
	}
	Dave_Num0 = satypes.SubaccountId{
		Owner:  types.MustBech32ifyAddressBytes(config.Bech32PrefixAccAddr, DaveAccAddress),
		Number: 0,
	}
	Dave_Num1 = satypes.SubaccountId{
		Owner:  types.MustBech32ifyAddressBytes(config.Bech32PrefixAccAddr, DaveAccAddress),
		Number: 1,
	}
)
