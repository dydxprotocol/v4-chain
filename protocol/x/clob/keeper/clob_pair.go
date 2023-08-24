package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

// CreatePerpetualClobPair creates a new perpetual CLOB pair in the store.
// Additionally, it creates an order book matching the ID of the newly created CLOB pair.
//
// An error will occur if any of the fields fail validation (see validateClobPair for details),
// or if the `perpetualId` cannot be found.
// In the event of an error, the store will not be updated nor will a matching order book be created.
//
// Returns the newly created CLOB pair and an error if one occurs.
func (k Keeper) CreatePerpetualClobPair(
	ctx sdk.Context,
	perpetualId uint32,
	stepSizeBaseQuantums satypes.BaseQuantums,
	quantumConversionExponent int32,
	subticksPerTick uint32,
	status types.ClobPair_Status,
) (types.ClobPair, error) {
	nextId := k.GetNumClobPairs(ctx)

	clobPair := types.ClobPair{
		Metadata: &types.ClobPair_PerpetualClobMetadata{
			PerpetualClobMetadata: &types.PerpetualClobMetadata{
				PerpetualId: perpetualId,
			},
		},
		Id:                        nextId,
		StepBaseQuantums:          stepSizeBaseQuantums.ToUint64(),
		QuantumConversionExponent: quantumConversionExponent,
		SubticksPerTick:           subticksPerTick,
		Status:                    status,
	}
	if err := k.validateClobPair(ctx, &clobPair); err != nil {
		return clobPair, err
	}

	k.setClobPair(ctx, clobPair)
	k.setNumClobPairs(ctx, nextId+1)

	return clobPair, nil
}

// validateClobPair validates a CLOB pair's fields are suitable for CLOB pair creation.
//
// - Metadata:
//   - Must be a perpetual CLOB pair with a perpetualId matching a perpetual in the store.
//
// - Status:
//   - Must be a supported status.
//
// - StepBaseQuantums:
//   - Must be greater than zero.
//
// - SubticksPerTick:
//   - Must be greater than zero.
func (k Keeper) validateClobPair(ctx sdk.Context, clobPair *types.ClobPair) error {
	if isSupported := types.IsSupportedClobPairStatus(clobPair.Status); !isSupported {
		return sdkerrors.Wrapf(
			types.ErrInvalidClobPairParameter,
			"CLOB pair (%+v) has unsupported status %+v",
			clobPair,
			clobPair.Status,
		)
	}

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

// mustGetClobPair fetches a ClobPair from state given its id.
// This function panics if the ClobPair is not found.
func (k Keeper) mustGetClobPair(
	ctx sdk.Context,
	clobPairId types.ClobPairId,
) types.ClobPair {
	clobPair, found := k.GetClobPair(ctx, clobPairId)
	if !found {
		panic(
			fmt.Sprintf(
				"mustGetClobPair: ClobPair with id %+v not found",
				clobPairId,
			),
		)
	}
	return clobPair
}

// SetClobPairStatus fetches a ClobPair by id and sets its
// Status property equal to the provided ClobPair_Status. This function returns
// an error if the proposed status transition is not supported.
func (k Keeper) SetClobPairStatus(
	ctx sdk.Context,
	clobPairId types.ClobPairId,
	clobPairStatus types.ClobPair_Status,
) error {
	clobPair := k.mustGetClobPair(ctx, clobPairId)

	if isSupported := types.IsSupportedClobPairStatusTransition(clobPair.Status, clobPairStatus); !isSupported {
		return sdkerrors.Wrapf(
			types.ErrInvalidClobPairStatusTransition,
			"Cannot transition from status %+v to status %+v",
			clobPair.Status,
			clobPairStatus,
		)
	}

	clobPair.Status = clobPairStatus
	if err := k.validateClobPair(ctx, &clobPair); err != nil {
		return err
	}

	k.setClobPair(ctx, clobPair)

	return nil
}
