package keeper

import (
	"context"
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

	lastTenEvents := k.sDAIEventManager.GetLastTensDAIEventsUnordered()

	for _, event := range lastTenEvents {

		if event.EthereumBlockNumber == "" || event.ConversionRate == "" {
			continue
		}

		blockNumber, err := ConvertStringToBigInt(event.EthereumBlockNumber)
		if err != nil {
			return nil, err
		}

		if blockNumber.Cmp(bigEthereumBlockNumber) == 0 {

			conversionRate, err := ConvertStringToBigInt(event.ConversionRate)
			if err != nil {
				return nil, err
			}

			if bigConversionRate.Cmp(conversionRate) == 0 {

				// TODO [YBCP-20]: Handle initializations better
				currentRate, initialized := k.GetSDAIPrice(ctx)

				if initialized && conversionRate.Cmp(currentRate) <= 0 {
					return nil, errorsmod.Wrap(
						types.ErrInvalidSDAIConversionRate,
						"The suggested sDAI conversion rate must be greater than the curret one",
					)
				}

				if !initialized {
					k.SetAssetYieldIndex(ctx, new(big.Rat).SetInt64(0))
				}

				k.SetSDAIPrice(ctx, conversionRate)

				k.MintNewTDaiAndSetNewYieldIndex(ctx)

				return &types.MsgUpdateSDAIConversionRateResponse{}, nil
			}
		}
	}

	return nil, errorsmod.Wrap(
		types.ErrInvalidSDAIConversionRate,
		"The suggested sDAI conversion rate is not valid",
	)
}

func ConvertStringToBigInt(str string) (*big.Int, error) {

	bigint, ok := new(big.Int).SetString(str, 10)
	if !ok {
		return nil, errorsmod.Wrap(
			types.ErrUnableToDecodeBigInt,
			"Unable to convert the sDAI conversion rate to a big int",
		)
	}

	return bigint, nil
}
