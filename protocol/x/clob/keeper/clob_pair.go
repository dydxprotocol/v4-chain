package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

// CreatePerpetualClobPair creates a new perpetual CLOB pair in the store.
// It stores the perpetual clob pair based off of the clob pair id, and will overwrite
// any previous clobPairs stored with the same Id.
// Additionally, it creates an order book matching the ID of the newly created CLOB pair.
//
// An error will occur if any of the fields fail validation (see validateClobPair for details),
// or if the `perpetualId` cannot be found.
// In the event of an error, the store will not be updated nor will a matching order book be created.
//
// Returns the newly created CLOB pair and an error if one occurs.
func (k Keeper) CreatePerpetualClobPair(
	ctx sdk.Context,
	clobPairId uint32,
	perpetualId uint32,
	stepSizeBaseQuantums satypes.BaseQuantums,
	quantumConversionExponent int32,
	subticksPerTick uint32,
	status types.ClobPair_Status,
	makerFeePpm uint32,
	takerFeePpm uint32,
) (types.ClobPair, error) {
	clobPair := types.ClobPair{
		Metadata: &types.ClobPair_PerpetualClobMetadata{
			PerpetualClobMetadata: &types.PerpetualClobMetadata{
				PerpetualId: perpetualId,
			},
		},
		Id:                        clobPairId,
		StepBaseQuantums:          stepSizeBaseQuantums.ToUint64(),
		QuantumConversionExponent: quantumConversionExponent,
		SubticksPerTick:           subticksPerTick,
		Status:                    status,
		MakerFeePpm:               makerFeePpm,
		TakerFeePpm:               takerFeePpm,
	}
	if err := k.validateClobPair(ctx, &clobPair); err != nil {
		return clobPair, err
	}

	k.setClobPair(ctx, clobPair)

	return clobPair, nil
}

// validateClobPair validates a CLOB pair's fields are suitable for CLOB pair creation.
//
// - MakerFeePpm:
//   - Must be <= MaxFeePpm.
//   - Must be <= TakerFeePpm.
//
// - Metadata:
//   - Must be a perpetual CLOB pair with a perpetualId matching a perpetual in the store.
//
// - Status:
//   - Must be a status other than ClobPair_STATUS_UNSPECIFIED.
//
// - StepBaseQuantums:
//   - Must be greater than zero.
//
// - SubticksPerTick:
//   - Must be greater than zero.
//
// - TakerFeePpm:
//   - Must be <= MaxFeePpm.
func (k Keeper) validateClobPair(ctx sdk.Context, clobPair *types.ClobPair) error {
	// TODO(DEC-1535): update this validation when we implement "spot"/"asset" clob pairs.
	switch clobPair.Metadata.(type) {
	case *types.ClobPair_PerpetualClobMetadata:
		perpetualId, err := clobPair.GetPerpetualId()
		if err != nil {
			return sdkerrors.Wrapf(
				err,
				"CLOB pair (%+v) has invalid perpetual.",
				clobPair,
			)
		}
		// Validate the perpetual referenced by the CLOB pair exists.
		if _, err := k.perpetualsKeeper.GetPerpetual(ctx, perpetualId); err != nil {
			return sdkerrors.Wrapf(
				err,
				"CLOB pair (%+v) has invalid perpetual.",
				clobPair,
			)
		}
	default:
		return sdkerrors.Wrapf(
			types.ErrInvalidClobPairParameter,
			// TODO(DEC-1535): update this error message when we implement "spot"/"asset" clob pairs.
			"CLOB pair (%+v) is not a perpetual CLOB.",
			clobPair,
		)
	}

	if clobPair.StepBaseQuantums <= 0 {
		return sdkerrors.Wrapf(
			types.ErrInvalidClobPairParameter,
			"invalid ClobPair parameter: StepBaseQuantums must be > 0. Got %v",
			clobPair.StepBaseQuantums,
		)
	}

	// Since a subtick will be calculated as (1 tick/SubticksPerTick), the denominator cannot be 0
	// and negative numbers do not make sense.
	if clobPair.SubticksPerTick <= 0 {
		return sdkerrors.Wrapf(
			types.ErrInvalidClobPairParameter,
			"invalid ClobPair parameter: SubticksPerTick must be > 0. Got %v",
			clobPair.SubticksPerTick,
		)
	}

	if clobPair.Status == types.ClobPair_STATUS_UNSPECIFIED {
		return sdkerrors.Wrap(types.ErrInvalidClobPairParameter, "invalid ClobPair parameter: Status must be specified.")
	}

	if clobPair.MakerFeePpm > types.MaxFeePpm || clobPair.TakerFeePpm > types.MaxFeePpm {
		return sdkerrors.Wrapf(
			types.ErrInvalidClobPairParameter,
			"invalid ClobPair parameter: MakerFeePpm (%v) and TakerFeePpm (%v) must both be <= MaxFeePpm (%v).",
			clobPair.MakerFeePpm,
			clobPair.TakerFeePpm,
			types.MaxFeePpm,
		)
	}

	// TODO(DEC-1549): Revisit the decision to allow MakerFeePpm == TakerFeePpm.
	// DEC-133 was initially slated to set TakerFeePpm > MakerFeePpm, but NoFee constants used in many tests prevented it.
	if clobPair.MakerFeePpm > clobPair.TakerFeePpm {
		return sdkerrors.Wrapf(
			types.ErrInvalidClobPairParameter,
			"invalid ClobPair parameter: MakerFeePpm (%v) must be <= TakerFeePpm (%v).",
			clobPair.MakerFeePpm,
			clobPair.TakerFeePpm,
		)
	}

	return nil
}

// setClobPair sets a specific `ClobPair` in the store from its index, and additionally creates a new orderbook
// to store all memclob orders.
func (k Keeper) setClobPair(ctx sdk.Context, clobPair types.ClobPair) {
	b := k.cdc.MustMarshal(&clobPair)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ClobPairKeyPrefix))
	// Write the `ClobPair` to state.
	store.Set(types.ClobPairKey(clobPair.GetClobPairId()), b)

	// Create the corresponding orderbook in the memclob.
	k.MemClob.CreateOrderbook(ctx, clobPair)
}

// InitMemClobOrderbooks initializes the memclob with `ClobPair`s from state.
// This is called during app initialization in `app.go`, before any ABCI calls are received.
func (k Keeper) InitMemClobOrderbooks(ctx sdk.Context) {
	clobPairs := k.GetAllClobPair(ctx)
	for _, clobPair := range clobPairs {
		// Create the corresponding orderbook in the memclob.
		k.MemClob.CreateOrderbook(
			ctx,
			clobPair,
		)
	}
}

// Sets the total count of CLOB pairs in the store to `num`.
func (k Keeper) setNumClobPairs(ctx sdk.Context, num uint32) {
	// Get necessary stores.
	store := ctx.KVStore(k.storeKey)

	// Set `numClobPairs`.
	store.Set(types.KeyPrefix(types.NumClobPairsKey), lib.Uint32ToBytes(num))
}

// Returns the total count of CLOB pairs, read from the store.
func (k Keeper) GetNumClobPairs(
	ctx sdk.Context,
) uint32 {
	store := ctx.KVStore(k.storeKey)
	numClobPairBytes := store.Get(types.KeyPrefix(types.NumClobPairsKey))
	return lib.BytesToUint32(numClobPairBytes)
}

// GetClobPair returns a clobPair from its index
func (k Keeper) GetClobPair(
	ctx sdk.Context,
	id types.ClobPairId,

) (val types.ClobPair, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ClobPairKeyPrefix))

	b := store.Get(types.ClobPairKey(
		id,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveClobPair removes a clobPair from the store
func (k Keeper) RemoveClobPair(
	ctx sdk.Context,
	id types.ClobPairId,

) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ClobPairKeyPrefix))
	store.Delete(types.ClobPairKey(
		id,
	))
}

// GetAllClobPair returns all clobPair
func (k Keeper) GetAllClobPair(ctx sdk.Context) (list []types.ClobPair) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ClobPairKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.ClobPair
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
