package module

import (
	"fmt"

	"cosmossdk.io/x/tx/signing"
	addresscodec "github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	gogoproto "github.com/cosmos/gogoproto/proto"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var InterfaceRegistry types.InterfaceRegistry

func getLegacyMsgSignerFn(path []string) func(msg proto.Message) ([][]byte, error) {
	if len(path) == 0 {
		panic("path is expected to contain at least one value.")
	}

	return func(msg proto.Message) ([][]byte, error) {
		m := msg.ProtoReflect()
		for _, p := range path[:len(path)-1] {
			fieldDesc := m.Descriptor().Fields().ByName(protoreflect.Name(p))
			if fieldDesc.Kind() != protoreflect.MessageKind {
				return nil, fmt.Errorf("Expected for field %s to be Message type in path %+v for msg %+v.", p, path, msg)
			}
			v := m.Get(fieldDesc)
			if !v.IsValid() {
				return nil, fmt.Errorf("Expected for field %s to be populated in path %+v for msg %+v.", p, path, msg)
			}
			m = v.Message()
		}

		fieldDesc := m.Descriptor().Fields().ByName(protoreflect.Name(path[len(path)-1]))
		if fieldDesc.Kind() != protoreflect.StringKind {
			return nil, fmt.Errorf(
				"Expected for final field %s to be String type in path %+v for msg %+v.",
				path[len(path)-1],
				path,
				msg,
			)
		}
		signer, err := sdk.AccAddressFromBech32(m.Get(fieldDesc).String())
		if err != nil {
			return nil, err
		}
		return [][]byte{signer}, nil
	}
}

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
				"dydxprotocol.clob.MsgBatchCancel": getLegacyMsgSignerFn(
					[]string{"subaccount_id", "owner"},
				),
				"dydxprotocol.clob.MsgCancelOrder": getLegacyMsgSignerFn(
					[]string{"order_id", "subaccount_id", "owner"},
				),
				"dydxprotocol.clob.MsgPlaceOrder": getLegacyMsgSignerFn(
					[]string{"order", "order_id", "subaccount_id", "owner"},
				),
				"dydxprotocol.sending.MsgCreateTransfer": getLegacyMsgSignerFn(
					[]string{"transfer", "sender", "owner"},
				),
				"dydxprotocol.sending.MsgWithdrawFromSubaccount": getLegacyMsgSignerFn(
					[]string{"sender", "owner"},
				),
				"dydxprotocol.vault.MsgDepositToMegavault": getLegacyMsgSignerFn(
					[]string{"subaccount_id", "owner"},
				),
				"dydxprotocol.vault.MsgWithdrawFromMegavault": getLegacyMsgSignerFn(
					[]string{"subaccount_id", "owner"},
				),
				"dydxprotocol.clob.MsgUpdateLeverage": getLegacyMsgSignerFn(
					[]string{"subaccount_id", "owner"},
				),
				"dydxprotocol.listing.MsgCreateMarketPermissionless": getLegacyMsgSignerFn(
					[]string{"subaccount_id", "owner"},
				),

				// App injected messages have no signers.
				"dydxprotocol.bridge.MsgAcknowledgeBridges":  noSigners,
				"dydxprotocol.clob.MsgProposedOperations":    noSigners,
				"dydxprotocol.perpetuals.MsgAddPremiumVotes": noSigners,
				"dydxprotocol.prices.MsgUpdateMarketPrices":  noSigners,
			},
		},
	})
}
