package types

import (
	"encoding/json"
	"fmt"
	gometrics "github.com/armon/go-metrics"
	"github.com/cosmos/cosmos-sdk/telemetry"
	"github.com/dydxprotocol/v4/lib/metrics"
	"github.com/dydxprotocol/v4/x/prices/types"
	"sort"
	"strings"
	"sync"
)

// PricefeedMutableMarketConfigs is the interface that stores a single copy of all market state
// that can change dynamically and synchronizes access for running go routines within the daemon.
type PricefeedMutableMarketConfigs interface {
	AddExchangeConfigUpdater(updater ExchangeConfigUpdater)
	UpdateMarkets(marketParams []types.MarketParam) error
	GetExchangeMarketConfigCopy(
		id ExchangeId,
	) (
		mutableExchangeMarketConfig *MutableExchangeMarketConfig,
		err error,
	)
	GetMarketConfigCopies(
		markets []MarketId,
	) (
		mutableMarketConfigs []*MutableMarketConfig,
		err error,
	)
}

// Ensure the `PricefeedMutableMarketConfigsImpl` struct is implemented at compile time.
var _ PricefeedMutableMarketConfigs = &PricefeedMutableMarketConfigsImpl{}

// PricefeedMutableMarketConfigsImpl is the implementation of PricefeedMutableMarketConfigs.
type PricefeedMutableMarketConfigsImpl struct {
	sync.RWMutex

	// These maps are updated when the exchange market params are updated. Map reads
	// and updates are synchronized by the RWMutex.
	// mutableExchangeToConfigs contains the latest market configuration for each exchange.
	mutableExchangeToConfigs map[ExchangeId]*MutableExchangeMarketConfig

	// mutableMarketToConfigs contains the latest market configuration for each market, common
	// across all exchanges.
	mutableMarketToConfigs map[MarketId]*MutableMarketConfig

	// mutableExchangeConfigUpdaters contains a map of ExchangeConfigUpdaters for each exchange.
	// In reality, these ExchangeConfigUpdaters are pointers to price fetchers. It is initialized
	// as empty, and each price fetcher adds itself to the PricefeedMutableMarketConfigsImpl
	// on creation.
	// Whenever the mutable exchange market config for an exchange updates, the
	// PriceFeedMutableMarketConfigsImpl calls `UpdateMutableExchangeConfig` on the price fetcher.
	//
	// Once a key is populated, the value never changes. The map is populated when price fetchers are
	// created at application start. Individual price fetchers manage their own synchronization.
	mutableExchangeConfigUpdaters map[ExchangeId]ExchangeConfigUpdater
}

// NewPriceFeedMutableMarketConfigs creates a new PricefeedMutableMarketConfigsImpl with no markets assigned
// to any exchanges. Apply market settings by calling `UpdateMarkets`.
func NewPriceFeedMutableMarketConfigs(
	canonicalExchangeIds []ExchangeId,
) *PricefeedMutableMarketConfigsImpl {
	exchangeIdToMutableExchangeConfigUpdater := make(
		map[ExchangeId]ExchangeConfigUpdater,
		len(canonicalExchangeIds),
	)

	// Initialize the mutableExchangeConfigs for all exchanges with no markets.
	mutableExchangeToConfigs := make(map[ExchangeId]*MutableExchangeMarketConfig, len(canonicalExchangeIds))
	for _, exchangeId := range canonicalExchangeIds {
		mutableExchangeToConfigs[exchangeId] = &MutableExchangeMarketConfig{
			Id:             exchangeId,
			MarketToTicker: make(map[MarketId]string, 0),
		}
	}

	pfmmc := &PricefeedMutableMarketConfigsImpl{
		mutableExchangeToConfigs:      mutableExchangeToConfigs,
		mutableMarketToConfigs:        nil,
		mutableExchangeConfigUpdaters: exchangeIdToMutableExchangeConfigUpdater,
	}

	return pfmmc
}

// AddExchangeConfigUpdater adds a new exchange config updater to the pricefeed mutable market configs.
// This synchronized method is how a price fetcher reports to the PricefeedMutableMarketConfigs that it wants updates
// for a particular exchange.
// This method was added to the pricefeed mutable market configs because price fetchers are initialized
// with a pointer to the pricefeed mutable market configs, and the pricefeed mutable market configs also
// needs a reference to the price fetchers to be fully initialized - so it was decided to initialize the
// pricefeed mutable market configs first, and then add the price fetchers to it.
func (pfmmc *PricefeedMutableMarketConfigsImpl) AddExchangeConfigUpdater(
	updater ExchangeConfigUpdater,
) {
	pfmmc.Lock()
	defer pfmmc.Unlock()

	pfmmc.mutableExchangeConfigUpdaters[updater.GetExchangeId()] = updater
}

// ValidateAndTransformParams validates the market params and transforms them into the internal representation used
// by the PricefeedMutableMarketConfigsImpl. The method guarantees that the returned mutableExchangeConfigs will have
// an entry for all current exchange feeds. This method is exposed for testing.
func (pfmmc *PricefeedMutableMarketConfigsImpl) ValidateAndTransformParams(marketParams []types.MarketParam) (
	mutableExchangeConfigs map[ExchangeId]*MutableExchangeMarketConfig,
	mutableMarketConfigs map[MarketId]*MutableMarketConfig,
	err error,
) {
	if marketParams == nil {
		return nil, nil, fmt.Errorf("marketParams cannot be nil")
	}

	mutableMarketConfigs = make(map[MarketId]*MutableMarketConfig, len(marketParams))

	mutableExchangeConfigs = make(map[ExchangeId]*MutableExchangeMarketConfig, len(pfmmc.mutableExchangeToConfigs))
	// Initialize mutableExchangeConfigs with empty MutableExchangeMarketConfigs to make sure that each exchange
	// has an entry in the map. The set of exchanges is fixed and defined at compile time. We need
	// mutableExchangeMarketConfigs to be defined for all exchanges so that we can update the respective price fetchers.
	for exchangeId := range pfmmc.mutableExchangeToConfigs {
		mutableExchangeConfigs[exchangeId] = &MutableExchangeMarketConfig{
			Id:             exchangeId,
			MarketToTicker: map[MarketId]string{},
		}
	}

	// marketNameToId, exchangeNameToId used to validate and interpret ExchangeConfigJson values.
	marketNameToId := make(map[string]MarketId, len(pfmmc.mutableMarketToConfigs))
	for _, param := range marketParams {
		marketNameToId[param.Pair] = param.Id
	}

	exchangeNames := make([]ExchangeId, 0, len(pfmmc.mutableExchangeToConfigs))
	for exchangeName := range pfmmc.mutableExchangeToConfigs {
		exchangeNames = append(exchangeNames, exchangeName)
	}

	for i, marketParam := range marketParams {
		// Perform basic validation on the market params. Id and exponent may be zero, but pair should always
		// be populated.
		// Note: we do not call `Validate` on the marketParam, but only do daemon-specific validation on the param,
		// ignoring fields that aren't relevant to the daemon. We rely on the protocol itself to surface issues for
		// configuration values that don't make sense and don't relate to daemon operation. In our case, all we need to
		// do is make sure that the pair is not empty.
		if marketParam.Pair == "" {
			return nil, nil, fmt.Errorf("invalid market param %v: pair cannot be empty", i)
		}

		// Check for duplicate market params.
		if _, exists := mutableMarketConfigs[marketParam.Id]; exists {
			return nil, nil, fmt.Errorf("invalid market param %v: duplicate market id %v", i, marketParam.Id)
		}

		mutableMarketConfigs[marketParam.Id] = &MutableMarketConfig{
			Id:       marketParam.Id,
			Pair:     marketParam.Pair,
			Exponent: marketParam.Exponent,
		}

		var exchangeConfigJson ExchangeConfigJson
		err = json.Unmarshal([]byte(marketParam.ExchangeConfigJson), &exchangeConfigJson)
		if err != nil {
			wrappedErr := fmt.Errorf("invalid exchange config json for market param %v: %w", i, err)
			return nil, nil, wrappedErr
		}

		err = exchangeConfigJson.Validate(exchangeNames, marketNameToId)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid exchange config json for market param %v: %w", i, err)
		}

		for _, exchangeConfig := range exchangeConfigJson.Exchanges {
			exchangeId := exchangeConfig.ExchangeName
			mutableExchangeConfig, ok := mutableExchangeConfigs[exchangeId]
			if !ok {
				err := fmt.Errorf(
					"internal error: exchange '%v' not found in mutableExchangeConfigs for market %v",
					exchangeId,
					marketParam.Pair,
				)
				return nil, nil, err
			}

			mutableExchangeConfig.MarketToTicker[marketParam.Id] = exchangeConfig.Ticker
		}
	}
	return mutableExchangeConfigs, mutableMarketConfigs, nil
}

// UpdateMarkets parses the market params, validates them, and updates the pricefeed mutable market configs,
// broadcasting updates to the price fetchers when necessary.
// This method is synchronized.
// 1. Validate and parse market params into a new set of MutableExchangeMarketConfig and MutableMarketConfig maps.
// 2. As a sanity check, validate all new configs have a price fetcher.
// 3. Pre-compute updates to send to price fetchers.
// 4. Take the writer lock on the pricefeed mutable market configs.
// 5. Swap in new markets and exchange configs.
// 6. For each changed exchange config, send an update to the price fetcher.
func (pfmmc *PricefeedMutableMarketConfigsImpl) UpdateMarkets(marketParams []types.MarketParam) error {
	// Emit metrics periodically regardless of UpdateMarkets success/failure.
	defer pfmmc.emitMarketAndExchangeCountMetrics()

	// 1. Validate and parse market params into a new mapping of MutableExchangeMarketConfigs and MutableMarketConfigs.
	// maps.
	if marketParams == nil {
		return fmt.Errorf("UpdateMarkets: marketParams cannot be nil")
	}

	newMutableExchangeConfigs, newMutableMarketConfigs, err := pfmmc.ValidateAndTransformParams(marketParams)
	if err != nil {
		return fmt.Errorf("UpdateMarkets market param validation failed: %w", err)
	}

	// 2. As a sanity check, validate all new configs have a price fetcher.
	previousExchangeConfigs := pfmmc.mutableExchangeToConfigs
	for exchangeId := range newMutableExchangeConfigs {
		// Validate that a previous exchange config should always exist for each exchange.
		// An error here would be unexpected.
		if _, ok := previousExchangeConfigs[exchangeId]; !ok {
			return fmt.Errorf("internal error: exchange %v not found in previousExchangeConfigs", exchangeId)
		}

		// Validate a price fetcher exists for the exchange.
		// An error here would be unexpected.
		if _, ok := pfmmc.mutableExchangeConfigUpdaters[exchangeId]; !ok {
			return fmt.Errorf("internal error: price fetcher not found for exchange %v", exchangeId)
		}
	}

	// 3. Pre-compute updates to send to price fetchers.
	exchangeToUpdatedMarketConfigs := make(map[ExchangeId][]*MutableMarketConfig, len(newMutableExchangeConfigs))
	exchangeToUpdatedExchangeConfig := make(
		map[ExchangeId]*MutableExchangeMarketConfig,
		len(newMutableExchangeConfigs),
	)

	for exchangeId, mutableExchangeConfig := range newMutableExchangeConfigs {
		previousConfig := previousExchangeConfigs[exchangeId]

		// If the exchange config has changed, pre-compute the updates to send to the price fetcher, which will
		// be a copy of the updated exchange config, as well as a sorted list of copied market configs for each
		// market on the exchange.
		if !previousConfig.Equal(mutableExchangeConfig) {
			// Make a list of sorted copies of all market configurations for the exchange.
			marketConfigCopies := make([]*MutableMarketConfig, 0, len(mutableExchangeConfig.MarketToTicker))
			for marketId := range mutableExchangeConfig.MarketToTicker {
				marketConfigCopies = append(marketConfigCopies, newMutableMarketConfigs[marketId].Copy())
			}
			// Ensure markets are sorted to simplify testing.
			sort.Slice(marketConfigCopies, func(i, j int) bool {
				return marketConfigCopies[i].Id < marketConfigCopies[j].Id
			})

			exchangeToUpdatedMarketConfigs[exchangeId] = marketConfigCopies
			exchangeToUpdatedExchangeConfig[exchangeId] = mutableExchangeConfig.Copy()
		}
	}

	// 4. Take the writer lock.
	pfmmc.Lock()
	defer pfmmc.Unlock()

	// 5. Swap in new markets and exchange configs.
	pfmmc.mutableExchangeToConfigs = newMutableExchangeConfigs
	pfmmc.mutableMarketToConfigs = newMutableMarketConfigs

	// 6. For each changed exchange config, send an update to the price fetcher.
	// TODO(DEC-2020): use errors.Join once it's available.
	updateFetcherErrors := make([]string, 0, len(exchangeToUpdatedExchangeConfig))
	for exchangeId, mutableExchangeConfig := range exchangeToUpdatedExchangeConfig {
		priceFetcher := pfmmc.mutableExchangeConfigUpdaters[exchangeId]
		err := priceFetcher.UpdateMutableExchangeConfig(
			mutableExchangeConfig,
			exchangeToUpdatedMarketConfigs[exchangeId],
		)
		if err != nil {
			updateFetcherErrors = append(
				updateFetcherErrors,
				fmt.Errorf(
					"UpdateMarkets: failed to update price fetcher for exchange %v: %w",
					exchangeId,
					err,
				).Error(),
			)
		}
	}

	if len(updateFetcherErrors) > 0 {
		return fmt.Errorf("UpdateMarkets: failed to update some price fetcher : %v", strings.Join(updateFetcherErrors, ", "))
	}

	return nil
}

// GetExchangeMarketConfigCopy retrieves a copy of the current market-specific mutable configuration
// for all markets of an exchange, in order to maintain synchronization. Whenever a market is added
// or modified on an exchange, this data structure becomes stale.
func (pfmmc *PricefeedMutableMarketConfigsImpl) GetExchangeMarketConfigCopy(
	id ExchangeId,
) (
	mutableExchangeMarketConfig *MutableExchangeMarketConfig,
	err error,
) {
	pfmmc.RLock()
	defer func() { pfmmc.RUnlock() }()
	memc, ok := pfmmc.mutableExchangeToConfigs[id]
	if !ok {
		return nil, fmt.Errorf("mutableExchangeMarketConfig not found for exchange %v", id)
	}
	return memc.Copy(), nil
}

// GetMarketConfigCopies retrieves a copy of the current market-specific mutable configuration for
// the specified markets, in order to maintain synchronization. In the event of a market update,
// this data could become stale. MarketConfigs are returned in the same order as the input markets.
func (pfmmc *PricefeedMutableMarketConfigsImpl) GetMarketConfigCopies(
	markets []MarketId,
) (
	mutableMarketConfigs []*MutableMarketConfig,
	err error,
) {
	pfmmc.RLock()
	defer func() { pfmmc.RUnlock() }()

	mutableMarketConfigs = make([]*MutableMarketConfig, 0, len(markets))
	for _, market := range markets {
		config, ok := pfmmc.mutableMarketToConfigs[market]
		if !ok {
			return nil, fmt.Errorf("market %v not found in mutableMarketToConfigs", market)
		}
		mutableMarketConfigs = append(mutableMarketConfigs, config.Copy())
	}

	return mutableMarketConfigs, nil
}

// emitMarketAndExchangeCountMetrics emits metrics related to the number of configured markets and exchanges.
// This method is synchronized and invoked every time the pricefeed mutable market configs is updated.
// This method is not re-entrant and must be called via defer within other protected pricefeed mutable market
// config methods.
// Note: this method does not call into the metrics utility library shared by the daemon and the server because
// that library uses pricefeed constants, which has an import dependency on the types package here. The reason
// for this is that the MarketId, ExchangeId, and Exponent types are used to define the pricefeed constants
// for static config that is used within the types directory. The ultimate solution is to either remove this
// config or pass all of it
func (pfmmc *PricefeedMutableMarketConfigsImpl) emitMarketAndExchangeCountMetrics() {
	pfmmc.RLock()
	defer pfmmc.RUnlock()

	// Set configured market count.
	telemetry.SetGauge(
		float32(len(pfmmc.mutableMarketToConfigs)),
		metrics.PricefeedDaemon,
		metrics.ConfiguredMarketCount,
	)

	for exchangeId, exchangeConfig := range pfmmc.mutableExchangeToConfigs {
		// Report configured metric count with label.
		telemetry.SetGaugeWithLabels(
			[]string{metrics.PricefeedDaemon, metrics.ConfiguredMarketCountPerExchange},
			float32(len(exchangeConfig.MarketToTicker)),
			[]gometrics.Label{exchangeLabel(exchangeId)},
		)
	}

	// Set gauge for number of configured exchanges per market.
	// Compute exchanges per market.
	marketExchanges := make(map[MarketId][]ExchangeId, len(pfmmc.mutableMarketToConfigs))
	for marketId := range pfmmc.mutableMarketToConfigs {
		marketExchanges[marketId] = make([]ExchangeId, 0, len(pfmmc.mutableExchangeToConfigs))
	}
	for exchangeId, exchangeConfig := range pfmmc.mutableExchangeToConfigs {
		for marketId := range exchangeConfig.MarketToTicker {
			marketExchanges[marketId] = append(marketExchanges[marketId], exchangeId)
		}
	}
	// Report exchanges per market.
	for marketId, exchanges := range marketExchanges {
		// Report exchange count with market label.
		telemetry.SetGaugeWithLabels(
			[]string{metrics.PricefeedDaemon, metrics.ConfiguredExchangeCountPerMarket},
			float32(len(exchanges)),
			[]gometrics.Label{marketLabel(pfmmc.mutableMarketToConfigs[marketId])},
		)
	}
}

// exchangeLabel returns a metrics label for the specified exchange feed id. This logic is duplicated
// here to avoid import loops by introducing a dependency on the pricefeed constants.
func exchangeLabel(exchangeId ExchangeId) gometrics.Label {
	return metrics.GetLabelForStringValue(metrics.ExchangeId, exchangeId)
}

// marketLabel returns a metrics label for the specified market id. This logic is duplicated here to avoid
// import loops by introducing a dependency on the pricefeed constants.
func marketLabel(mutableMarketConfig *MutableMarketConfig) gometrics.Label {
	if mutableMarketConfig == nil {
		return metrics.GetLabelForStringValue(metrics.MarketId, "INVALID")
	}
	return metrics.GetLabelForStringValue(metrics.MarketId, mutableMarketConfig.Pair)
}
