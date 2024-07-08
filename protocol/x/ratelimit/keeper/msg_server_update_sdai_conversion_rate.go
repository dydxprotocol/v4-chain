package keeper

import (
	"context"
	"fmt"
	"math/big"

	errorsmod "cosmossdk.io/errors"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/lib"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/types"
)

// UpdateSDAIConversionRate updates the sDAI conversion rate.
func (k msgServer) UpdateSDAIConversionRate(
	goCtx context.Context,
	msg *types.MsgUpdateSDAIConversionRate,
) (*types.MsgUpdateSDAIConversionRateResponse, error) {
	ctx := lib.UnwrapSDKContext(goCtx, types.ModuleName)

	bigConversionRate, err := ConvertStringToBigInt(msg.ConversionRate)
	if err != nil {
		return nil, err
	}

	bigEthereumBlockNumber, err := ConvertStringToBigInt(msg.EthereumBlockNumber)
	if err != nil {
		return nil, err
	}

	lastTenEvents := k.sDAIEventManager.GetLastTensDAIEvents()

	for _, event := range lastTenEvents {

		// todo what if not all the array is full

		blockNumber, err := ConvertStringToBigInt(event.EthereumBlockNumber)
		if err != nil {
			return nil, err
		}

		if blockNumber == bigEthereumBlockNumber {

			conversionRate, err := ConvertStringToBigInt(event.ConversionRate)
			if err != nil {
				return nil, err
			}

			if bigConversionRate == conversionRate {

				currentRate, ok := k.GetSDAIPrice(ctx)
				if !ok {
					return nil, errorsmod.Wrap(
						types.ErrSDAIConversionRateNotInitisialised,
						fmt.Sprintf(
							"The suggested sDAI conversion rate is not valid",
						),
					)
				}

				if conversionRate.Cmp(currentRate) <= 0 {
					return nil, errorsmod.Wrap(
						types.ErrInvalidSDAIConversionRate,
						fmt.Sprintf(
							"The suggested sDAI conversion rate must be greater than the curret one",
						),
					)
				}

				k.SetSDAIPrice(ctx, conversionRate)
				return &types.MsgUpdateSDAIConversionRateResponse{}, nil
			}
		}
	}

	return nil, errorsmod.Wrap(
		types.ErrInvalidSDAIConversionRate,
		fmt.Sprintf(
			"The suggested sDAI conversion rate is not valid",
		),
	)
}

func ConvertStringToBigInt(str string) (*big.Int, error) {

	bigint, ok := new(big.Int).SetString(str, 10)
	if !ok {
		return nil, errorsmod.Wrap(
			types.ErrUnableToDecodeBigInt,
			fmt.Sprintf(
				"Unable to convert the sDAI conversion rate to a big int",
			),
		)
	}

	return bigint, nil
}
