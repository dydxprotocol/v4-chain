package module

import (
	"cosmossdk.io/x/tx/signing"
	addresscodec "github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	gogoproto "github.com/cosmos/gogoproto/proto"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var InterfaceRegistry types.InterfaceRegistry

// TODO(CORE-846): Consider having app injected messages return an error instead of empty signers list.
func noSigners(_ proto.Message) ([][]byte, error) {
	return [][]byte{}, nil
}

func init() {
	var err error

	if InterfaceRegistry, err = NewInterfaceRegistry(
		sdk.GetConfig().GetBech32AccountAddrPrefix(),
		sdk.GetConfig().GetBech32ValidatorAddrPrefix(),
	); err != nil {
		panic(err)
	}
}

func NewInterfaceRegistry(addrPrefix string, valAddrPrefix string) (types.InterfaceRegistry, error) {
	return types.NewInterfaceRegistryWithOptions(types.InterfaceRegistryOptions{
		ProtoFiles: gogoproto.HybridResolver,
		SigningOptions: signing.Options{
			AddressCodec:          addresscodec.NewBech32Codec(addrPrefix),
			ValidatorAddressCodec: addresscodec.NewBech32Codec(valAddrPrefix),
			// TODO(CORE-840): cosmos.msg.v1.signer annotation doesn't supported nested messages beyond a depth of 1
			// which requires any message where the authority is nested further to implement its own accessor. Once
			// https://github.com/cosmos/cosmos-sdk/issues/18722 is fixed, replace this with the cosmos.msg.v1.signing
			// annotation on the protos.
			CustomGetSigners: map[protoreflect.FullName]signing.GetSignersFunc{
				// App injected messages have no signers.
				"dydxprotocol.bridge.MsgAcknowledgeBridges":  noSigners,
				"dydxprotocol.clob.MsgProposedOperations":    noSigners,
				"dydxprotocol.perpetuals.MsgAddPremiumVotes": noSigners,
				"dydxprotocol.prices.MsgUpdateMarketPrices":  noSigners,
			},
		},
	})
}
