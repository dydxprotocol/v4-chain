package simulation

// DONTCOVER

import (
	"github.com/dydxprotocol/v4-chain/protocol/app/module"
	"math"
	"math/big"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/sim_helpers"
	assettypes "github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/sending/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/sending/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

const (
	opWeightMsgCreateTransfer = "op_weight_msg_create_transfer"

	defaultWeightMsgCreateTransfer int = 100
)

var (
	typeMsgCreateTransfer = sdk.MsgTypeURL(&types.MsgCreateTransfer{})
)

// WeightedOperations returns all the operations from the module with their respective weights.
func WeightedOperations(
	appParams simtypes.AppParams,
	k keeper.Keeper,
	ak types.AccountKeeper,
	bk types.BankKeeper,
	sk types.SubaccountsKeeper,
) simulation.WeightedOperations {
	protoCdc := codec.NewProtoCodec(module.InterfaceRegistry)

	operations := make([]simtypes.WeightedOperation, 0)

	var weightMsgCreateTransfer int
	appParams.GetOrGenerate(opWeightMsgCreateTransfer, &weightMsgCreateTransfer, nil,
		func(_ *rand.Rand) {
			weightMsgCreateTransfer = defaultWeightMsgCreateTransfer
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreateTransfer,
		SimulateMsgCreateTransfer(k, ak, bk, sk, protoCdc),
	))

	return operations
}

// SimulateMsgCreateTransfer generates a random MsgCreateTransfer.
func SimulateMsgCreateTransfer(
	k keeper.Keeper,
	ak types.AccountKeeper,
	bk types.BankKeeper,
	sk types.SubaccountsKeeper,
	cdc *codec.ProtoCodec,
) simtypes.Operation {
	return func(
		r *rand.Rand,
		app *baseapp.BaseApp,
		ctx sdk.Context,
		accs []simtypes.Account,
		chainId string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		senderAccount, err := sk.GetRandomSubaccount(ctx, r)
		// Return a no-op message if we don't have any accounts
		if err != nil {
			return simtypes.NoOpMsg(
				types.ModuleName,
				typeMsgCreateTransfer,
				"Not enough subaccounts for transfer",
			), nil, nil
		}

		risk, err := sk.GetNetCollateralAndMarginRequirements(
			ctx,
			satypes.Update{
				SubaccountId: *senderAccount.GetId(),
			},
		)
		if err != nil {
			panic(err)
		}

		// Select a different subaccount as the recipient.
		recipientAccount, err := sk.GetRandomSubaccount(ctx, r)
		if err != nil || *recipientAccount.GetId() == *senderAccount.GetId() {
			// Return a no-op message if we don't have any accounts or choose the same account as the sender
			// which should be quite rare.
			return simtypes.NoOpMsg(
				types.ModuleName,
				typeMsgCreateTransfer,
				"Not enough subaccounts for transfer",
			), nil, nil
		}

		// Calculate the maximum amount that the receiver can receive without any integer overflow.
		bigAmountReceivable := new(big.Int).Sub(
			new(big.Int).SetUint64(math.MaxUint64),
			recipientAccount.GetUsdcPosition(),
		)

		// Calculate the maximum amount that can be sent without making the subaccount under-collateralized.
		bigAmountPayable := new(big.Int).Sub(risk.NC, risk.IMR)

		bigMaxAmountToSend := lib.BigMin(bigAmountPayable, bigAmountReceivable)
		if bigMaxAmountToSend.Sign() <= 0 {
			return simtypes.NoOpMsg(
				types.ModuleName,
				typeMsgCreateTransfer,
				"Sender does not have enough balance or receiver might overflow",
			), nil, nil
		}

		// Generate a random amount between [1, bigMaxAmountToSend].
		// Rand generates a number in [0, bigMaxAmountToSend), add 1 to make it in [1, bigMaxAmountToSend].
		amountToSend := new(big.Int).Add(
			new(big.Int).Rand(r, bigMaxAmountToSend),
			big.NewInt(1),
		)

		msg := &types.MsgCreateTransfer{
			Transfer: &types.Transfer{
				Sender:    *senderAccount.GetId(),
				Recipient: *recipientAccount.GetId(),
				AssetId:   assettypes.AssetUsdc.Id,
				Amount:    amountToSend.Uint64(),
			},
		}

		proposer, _ := simtypes.FindAccount(accs, senderAccount.GetId().MustGetAccAddress())
		opMsg, err := sim_helpers.GenerateAndDeliverTx(
			r,
			app,
			ctx,
			chainId,
			cdc,
			ak,
			bk,
			proposer,
			types.ModuleName,
			msg,
			typeMsgCreateTransfer,
			true, // fee does not apply when creating a transfer.
		)
		if err != nil {
			panic(err) // panic to halt/fail simulation.
		}

		return opMsg, nil, nil
	}
}
