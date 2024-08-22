package voteweighted

import (
	"fmt"
	"math/big"
	"sort"

	"cosmossdk.io/log"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ethos-works/ethos/ethos-chain/x/ccv/consumer/types"
)

type CCValidatorStore interface {
	GetAllCCValidator(ctx sdk.Context) []types.CrossChainValidator
	GetCCValidator(ctx sdk.Context, addr []byte) (types.CrossChainValidator, bool)
}

type AggregatorPricePair struct {
	SpotPrice *big.Int
	PnlPrice  *big.Int
}

type AggregateFn func(ctx sdk.Context, vePrices map[string]map[string]AggregatorPricePair) (map[string]AggregatorPricePair, error)

// DefaultPowerThreshold defines the total voting power % that must be
// submitted in order for a currency pair to be considered for the
// final oracle price. We provide a default supermajority threshold
// of 2/3+.
var DefaultPowerThreshold = math.LegacyNewDecWithPrec(667, 3)

type (
	// VoteWeightPriceInfo tracks the stake weight(s) + price(s) for a given currency pair.
	PriceInfo struct {
		SpotPrices  []PricePerValidator
		PnlPrices   []PricePerValidator
		TotalWeight math.Int
	}

	// VoteWeightPrice defines a price update that includes the stake weight of the validator.
	PricePerValidator struct {
		VoteWeight int64
		Price      *big.Int
	}
)

func Median(
	logger log.Logger,
	validatorStore CCValidatorStore,
	threshold math.LegacyDec,
) AggregateFn {
	return func(
		ctx sdk.Context,
		vePricesPerValidator map[string]map[string]AggregatorPricePair,
	) (map[string]AggregatorPricePair, error) {
		priceInfo := make(map[string]PriceInfo)

		for validatorAddr, validatorPrices := range vePricesPerValidator {
			validatorPower, err := getValidatorPowerByAddress(ctx, validatorStore, validatorAddr)
			if err != nil {
				logger.Info(
					"failed to get validator power, skipping",
					"validator_address", validatorAddr,
					"err", err,
				)
				continue
			}

			for market, pricePair := range validatorPrices {
				if pricePair.SpotPrice == nil && pricePair.PnlPrice == nil {
					logger.Info(
						"both spot and pnl prices are nil, skipping",
						"validator_address", validatorAddr,
						"currency_pair", market,
					)
					continue
				}

				if _, ok := priceInfo[market]; !ok {
					priceInfo[market] = PriceInfo{
						SpotPrices:  make([]PricePerValidator, 0),
						PnlPrices:   make([]PricePerValidator, 0),
						TotalWeight: math.ZeroInt(),
					}
				}

				pInfo := priceInfo[market]

				if pricePair.SpotPrice != nil {
					pInfo.SpotPrices = append(pInfo.SpotPrices, PricePerValidator{
						VoteWeight: validatorPower,
						Price:      pricePair.SpotPrice,
					})

					pInfo.PnlPrices = append(pInfo.PnlPrices, PricePerValidator{
						VoteWeight: validatorPower,
						Price:      pricePair.PnlPrice,
					})

					pInfo.TotalWeight = pInfo.TotalWeight.Add(math.NewInt(validatorPower))

				}

				priceInfo[market] = pInfo
			}
		}

		finalPrices := make(map[string]AggregatorPricePair)

		totalPower := GetTotalPower(ctx, validatorStore)

		for pair, info := range priceInfo {
			// The total voting power % that submitted a price update for the given currency pair must be
			// greater than the threshold to be included in the final oracle price.
			percentSubmitted := math.LegacyNewDecFromInt(info.TotalWeight).Quo(math.LegacyNewDecFromInt(totalPower))

			if percentSubmitted.GTE(threshold) {
				finalPrices[pair] = AggregatorPricePair{
					SpotPrice: ComputeMedian(info.SpotPrices, info.TotalWeight),
					PnlPrice:  ComputeMedian(info.PnlPrices, info.TotalWeight),
				}

				logger.Info(
					"computed stake-weighted median prices for currency pair",
					"currency_pair", pair,
					"percent_submitted", percentSubmitted.String(),
					"threshold", threshold.String(),
					"final_spot_price", finalPrices[pair].SpotPrice.String(),
					"final_pnl_price", finalPrices[pair].PnlPrice.String(),
					"num_validators", len(info.SpotPrices),
				)
			} else {
				logger.Info(
					"not enough voting power to compute stake-weighted median prices for currency pair",
					"currency_pair", pair,
					"threshold", threshold.String(),
					"percent_submitted", percentSubmitted.String(),
					"num_validators", len(info.SpotPrices),
				)
			}
		}
		return finalPrices, nil
	}
}

func ComputeMedian(prices []PricePerValidator, totalWeight math.Int) *big.Int {
	// Sort the prices by price.
	sort.SliceStable(prices, func(i, j int) bool {
		switch prices[i].Price.Cmp(prices[j].Price) {
		case -1:
			return true
		case 1:
			return false
		default:
			return true
		}
	})

	// Compute the median weight.
	middle := totalWeight.QuoRaw(2)

	// Iterate through the prices and compute the median price.
	sum := math.ZeroInt()
	for index, price := range prices {
		sum = sum.Add(math.NewInt(price.VoteWeight))

		if sum.GTE(middle) {
			return price.Price
		}

		// If we reached the end of the list, return the last price.
		if index == len(prices)-1 {
			return price.Price
		}
	}

	return nil
}

func getValidatorPowerByAddress(
	ctx sdk.Context,
	validatorStore CCValidatorStore,
	validatorAddr string,
) (int64, error) {
	addr, err := sdk.ConsAddressFromBech32(validatorAddr)
	if err != nil {
		return 0, err
	}

	validator, found := validatorStore.GetCCValidator(ctx, addr.Bytes())
	if !found {
		return 0, fmt.Errorf("validator not found")
	}

	validatorPower := validator.GetPower()
	return validatorPower, nil
}
