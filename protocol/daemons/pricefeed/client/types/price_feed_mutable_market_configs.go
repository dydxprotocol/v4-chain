package types

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/cosmos/cosmos-sdk/telemetry"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	"github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	gometrics "github.com/hashicorp/go-metrics"
)

const (
	expectedUpdatersPerExchange = 2
)

// PricefeedMutableMarketConfigs stores a single copy of all market state that can change dynamically and synchronizes
// access for running go routines within the daemon.
type PricefeedMutableMarketConfigs interface {
	AddPriceFetcher(updater ExchangeConfigUpdater)
	AddPriceEncoder(updater ExchangeConfigUpdater)
	UpdateMarkets(marketParams []types.MarketParam) (marketParamErrors map[MarketId]error, err error)
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

// UpdatersForExchange contains named references to all ExchangeConfigUpdaters for a single exchange.
type UpdatersForExchange struct {
	PriceFetcher ExchangeConfigUpdater
	PriceEncoder ExchangeConfigUpdater
}

// UpdateParameters contains the parameters to send to an ExchangeConfigUpdater when the exchange config changes.
type UpdateParameters struct {
	ExchangeConfig *MutableExchangeMarketConfig
	MarketConfigs  []*MutableMarketConfig
}

func (ufe *UpdatersForExchange) Validate() error {
	if ufe.PriceFetcher == nil {
		return fmt.Errorf("price fetcher cannot be nil")
	}
	if ufe.PriceEncoder == nil {
		return fmt.Errorf("price encoder cannot be nil")
	}
	return nil
}

// PricefeedMutableMarketConfigsImpl implements PricefeedMutableMarketConfigs.
type PricefeedMutableMarketConfigsImpl struct {
	sync.Mutex

	// mutableExchangeToConfigs contains the latest market configuration for each exchange.
	// These maps are updated when the exchange market params are updated. Map reads
	// and updates are synchronized by the Mutex.
	mutableExchangeToConfigs map[ExchangeId]*MutableExchangeMarketConfig

	// mutableMarketToConfigs contains the latest market configuration for each market, common
	// across all exchanges.
	mutableMarketToConfigs map[MarketId]*MutableMarketConfig

	// mutableExchangeConfigUpdaters contains a map of `ExchangeConfigUpdater`s for each exchange.
	// In reality, these ExchangeConfigUpdaters are pointers to price fetchers and price encoders.
	// It is initialized with no updaters, and each price fetcher and encoder subscribes to updates
	// from the PricefeedMutableMarketConfigsImpl on creation.
	//
	// Whenever the mutable exchange market config for an exchange updates, the
	// PriceFeedMutableMarketConfigsImpl calls `UpdateMutableExchangeConfig` on each subscribed fetcher
	// and encoder. The encoder is called first so that it has all necessary market config to support
	// price updates coming from the fetcher in the event of adding a new market.
	//
	// Once a key is populated, the value never changes. The map is populated when updaters are
	// created at application start. Individual updaters manage their own synchronization.
	mutableExchangeConfigUpdaters map[ExchangeId]UpdatersForExchange

	// updatersInitialized tracks whether all expected updaters have been added to the pricefeed mutable market configs.
	// This is used to ensure that all updaters are subscribed before the pricefeed mutable market configs processes or
	// emits updates. The pfmmc only emit updates to exchange config updaters when the config changes, so any missing
	// subscribers would never receive an update until the next change.
	updatersInitialized sync.WaitGroup
}

// NewPriceFeedMutableMarketConfigs creates a new PricefeedMutableMarketConfigsImpl with no markets assigned
// to any exchanges. Apply market settings by calling `UpdateMarkets`.
func NewPriceFeedMutableMarketConfigs(
	canonicalExchangeIds []ExchangeId,
) *PricefeedMutableMarketConfigsImpl {
	exchangeIdToMutableExchangeConfigUpdater := make(
		map[ExchangeId]UpdatersForExchange,
		len(canonicalExchangeIds),
	)

	// Initialize the mutableExchangeConfigs for all exchanges with no markets.
	mutableExchangeToConfigs := make(map[ExchangeId]*MutableExchangeMarketConfig, len(canonicalExchangeIds))
	for _, exchangeId := range canonicalExchangeIds {
		mutableExchangeToConfigs[exchangeId] = &MutableExchangeMarketConfig{
			Id:                   exchangeId,
			MarketToMarketConfig: make(map[MarketId]MarketConfig, 0),
		}
	}

	pfmmc := &PricefeedMutableMarketConfigsImpl{
		mutableExchangeToConfigs:      mutableExchangeToConfigs,
		mutableMarketToConfigs:        nil,
		mutableExchangeConfigUpdaters: exchangeIdToMutableExchangeConfigUpdater,
	}

	// Add the expected number of registered updaters to the wait group.
	pfmmc.updatersInitialized.Add(expectedUpdatersPerExchange * len(canonicalExchangeIds))

	return pfmmc
}

// AddPriceFetcher adds a new price fetcher to the pricefeed mutable market configs. This method is synchronized.
func (pfmmc *PricefeedMutableMarketConfigsImpl) AddPriceFetcher(
	priceFetcher ExchangeConfigUpdater,
) {
	pfmmc.addExchangeConfigUpdater(priceFetcher, true)
}

// AddPriceEncoder adds a new price encoder to the pricefeed mutable market configs. This method is synchronized.
func (pfmmc *PricefeedMutableMarketConfigsImpl) AddPriceEncoder(
	priceEncoder ExchangeConfigUpdater,
) {
	pfmmc.addExchangeConfigUpdater(priceEncoder, false)
}

// AddExchangeConfigUpdater adds a new exchange config updater to the pricefeed mutable market configs.
// This synchronized method is how a price fetcher or encoder subscribes itself in PricefeedMutableMarketConfigs
// for exchange configuration updates.
//
// This method was added to the pricefeed mutable market configs because fetchers and encoders are initialized
// with a pointer to the pricefeed mutable market configs, and the pricefeed mutable market configs also
// needs a reference to the updater to be fully initialized - so it was decided to initialize the
// pricefeed mutable market configs first, and then have updaters add themselves.
func (pfmmc *PricefeedMutableMarketConfigsImpl) addExchangeConfigUpdater(
	updater ExchangeConfigUpdater,
	isPriceFetcher bool,
) {
	pfmmc.Lock()
	defer pfmmc.Unlock()

	updatersForExchange, exists := pfmmc.mutableExchangeConfigUpdaters[updater.GetExchangeId()]
	if !exists {
		updatersForExchange = UpdatersForExchange{}
	}
	if isPriceFetcher {
		// Enforce that each updater can be added only once.
		if updatersForExchange.PriceFetcher != nil {
			panic(fmt.Sprintf("internal error: price fetcher already exists for exchange %v", updater.GetExchangeId()))
		}
		updatersForExchange.PriceFetcher = updater
	} else {
		// Enforce that each updater can be added only once.
		if updatersForExchange.PriceEncoder != nil {
			panic(fmt.Sprintf("internal error: price encoder already exists for exchange %v", updater.GetExchangeId()))
		}
		updatersForExchange.PriceEncoder = updater
	}

	pfmmc.mutableExchangeConfigUpdaters[updater.GetExchangeId()] = updatersForExchange
	pfmmc.updatersInitialized.Done()
}

// ValidateAndTransformParams validates the market params and transforms them into the internal representation used
// by the PricefeedMutableMarketConfigsImpl. The method guarantees that the returned mutableExchangeConfigs will have
// an entry for all current exchange feeds. This method is exposed for testing.
// MarketParams are validated and applied independently. If any market param is invalid, the method will populate
// marketParamErrors with the error and continue processing the remaining market params. If the entire validation fails,
// the method will return an error.
func (pfmmc *PricefeedMutableMarketConfigsImpl) ValidateAndTransformParams(marketParams []types.MarketParam) (
	mutableExchangeConfigs map[ExchangeId]*MutableExchangeMarketConfig,
	mutableMarketConfigs map[MarketId]*MutableMarketConfig,
	marketParamErrors map[MarketId]error,
	err error,
) {
	// Track individual errors for each market param that fails to apply.
	marketParamErrors = make(map[MarketId]error, len(marketParams))

	if marketParams == nil {
		return nil, nil, nil, fmt.Errorf("marketParams cannot be nil")
	}

	mutableMarketConfigs = make(map[MarketId]*MutableMarketConfig, len(marketParams))

	mutableExchangeConfigs = make(map[ExchangeId]*MutableExchangeMarketConfig, len(pfmmc.mutableExchangeToConfigs))
	// Initialize mutableExchangeConfigs with empty MutableExchangeMarketConfigs to make sure that each exchange
	// has an entry in the map. The set of exchanges is fixed and defined at compile time. We need
	// mutableExchangeMarketConfigs to be defined for all exchanges so that we can update the respective updaters.
	for exchangeId := range pfmmc.mutableExchangeToConfigs {
		mutableExchangeConfigs[exchangeId] = &MutableExchangeMarketConfig{
			Id:                   exchangeId,
			MarketToMarketConfig: map[MarketId]MarketConfig{},
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

	for _, marketParam := range marketParams {
		// Perform validation on the market params.
		if err := marketParam.Validate(); err != nil {
			marketParamErrors[marketParam.Id] = fmt.Errorf("invalid market param %v: %w", marketParam.Id, err)
			continue
		}

		// Check for duplicate market params.
		if _, exists := mutableMarketConfigs[marketParam.Id]; exists {
			// In this case, return an error, because we do not know which market param is the correct one.
			return nil,
				nil,
				nil,
				fmt.Errorf("invalid market params: duplicate market id %v", marketParam.Id)
		}

		var exchangeConfigJson ExchangeConfigJson
		err = json.Unmarshal([]byte(marketParam.ExchangeConfigJson), &exchangeConfigJson)
		if err != nil {
			wrappedErr := fmt.Errorf("invalid exchange config json for market param %v: %w", marketParam.Id, err)
			marketParamErrors[marketParam.Id] = wrappedErr
			continue
		}

		err = exchangeConfigJson.Validate(exchangeNames, marketNameToId)
		if err != nil {
			marketParamErrors[marketParam.Id] = fmt.Errorf(
				"invalid exchange config json for market param %v: %w",
				marketParam.Id,
				err,
			)
			continue
		}

		// Errors in the following loop are unexpected because we have already validated the exchange config json.
		// In this case, we return an error.
		for _, exchangeConfig := range exchangeConfigJson.Exchanges {
			exchangeId := exchangeConfig.ExchangeName
			mutableExchangeConfig, ok := mutableExchangeConfigs[exchangeId]
			if !ok {
				err := fmt.Errorf(
					"unexpected internal error: exchange '%v' not found in mutableExchangeConfigs for market %v",
					exchangeId,
					marketParam.Pair,
				)
				return nil, nil, nil, err
			}
			marketConfig := MarketConfig{
				Ticker: exchangeConfig.Ticker,
				Invert: exchangeConfig.Invert,
			}

			// Populate the adjustByMarket only if it is specified in the config.
			if exchangeConfig.AdjustByMarket != "" {
				adjustByMarketId, ok := marketNameToId[exchangeConfig.AdjustByMarket]
				if !ok {
					return nil, nil, nil, fmt.Errorf(
						"unexpected internal error: invalid exchange config json for exchange '%v' "+
							"on market param %v: adjustByMarket '%v' not found",
						exchangeConfig.ExchangeName,
						marketParam.Id,
						exchangeConfig.AdjustByMarket,
					)
				}
				marketConfig.AdjustByMarket = new(MarketId)
				*marketConfig.AdjustByMarket = adjustByMarketId
			}
			mutableExchangeConfig.MarketToMarketConfig[marketParam.Id] = marketConfig
		}

		// If we've reached this point, the market param is valid. Add it to the mutable market configs.
		mutableMarketConfigs[marketParam.Id] = &MutableMarketConfig{
			Id:           marketParam.Id,
			Pair:         marketParam.Pair,
			Exponent:     marketParam.Exponent,
			MinExchanges: marketParam.MinExchanges,
		}
	}
	return mutableExchangeConfigs, mutableMarketConfigs, marketParamErrors, nil
}

// getUpdateParametersForExchange returns a copy of the exchange config and all relevant market configs for the
// exchange. This parameter list can be used to update the price fetcher or the price encoder whenever an exchange
// config changes.
func getUpdateParametersForExchange(
	mutableExchangeConfig *MutableExchangeMarketConfig,
	mutableMarketConfigs map[MarketId]*MutableMarketConfig,
) (
	updateParameters UpdateParameters,
) {
	updateParameters.ExchangeConfig = mutableExchangeConfig.Copy()

	// Make a list of sorted copies of all market configurations for the exchange.
	marketConfigCopies := make([]*MutableMarketConfig, 0, len(mutableExchangeConfig.MarketToMarketConfig))

	// Detect which markets are needed. Due to adjustment markets, we need to deduplicate markets.
	marketsOnExchange := make(map[MarketId]struct{})
	for marketId, config := range mutableExchangeConfig.MarketToMarketConfig {
		marketsOnExchange[marketId] = struct{}{}
		if config.AdjustByMarket != nil {
			marketsOnExchange[*config.AdjustByMarket] = struct{}{}
		}
	}

	// Copy the market configs for each market on the exchange.
	for marketId := range marketsOnExchange {
		marketConfigCopies = append(marketConfigCopies, mutableMarketConfigs[marketId].Copy())
	}

	// Ensure markets are sorted by id in order to make behavior deterministic for testing.
	sort.Slice(marketConfigCopies, func(i, j int) bool {
		return marketConfigCopies[i].Id < marketConfigCopies[j].Id
	})

	updateParameters.MarketConfigs = marketConfigCopies

	return updateParameters
}

// UpdateMarkets parses the market params, validates them, and updates the pricefeed mutable market configs,
// broadcasting updates to the price fetchers and encoders when necessary.
// This method is synchronized.
// 1. Validate and parse market params into a new set of MutableExchangeMarketConfig and MutableMarketConfig maps.
// 2. As a sanity check, validate all new configs have 2 entries: a price fetcher and encoder.
// 3. Pre-compute updates to send to updaters.
// 4. Take the writer lock on the pricefeed mutable market configs.
// 5. Swap in new markets and exchange configs.
// 6. For each changed exchange config, send an update to each updater.
// UpdateMarkets applies market settings independently. If any market param is invalid, the method will populate
// marketParamErrors with the error and continue processing the remaining market params. If the entire validation fails,
// the method will return an error.
func (pfmmc *PricefeedMutableMarketConfigsImpl) UpdateMarkets(marketParams []types.MarketParam) (
	marketParamErrors map[MarketId]error,
	err error,
) {
	// Wait for all updaters to be added. This should happen quickly after the daemon starts.
	pfmmc.updatersInitialized.Wait()

	// Emit metrics periodically regardless of UpdateMarkets success/failure.
	defer pfmmc.emitMarketAndExchangeCountMetrics()

	// 1. Validate and parse market params into a new mapping of MutableExchangeMarketConfigs and MutableMarketConfigs.
	// maps.
	if marketParams == nil {
		return nil, fmt.Errorf("UpdateMarkets: marketParams cannot be nil")
	}

	newMutableExchangeConfigs,
		newMutableMarketConfigs,
		marketParamErrors,
		err := pfmmc.ValidateAndTransformParams(marketParams)
	if err != nil {
		return nil, fmt.Errorf("UpdateMarkets market param validation failed: %w", err)
	}

	// 2. As a sanity check, validate all new configs have a set of updaters.
	previousExchangeConfigs := pfmmc.mutableExchangeToConfigs
	for exchangeId := range newMutableExchangeConfigs {
		// Validate that a previous exchange config should always exist for each exchange.
		// An error here would be unexpected.
		if _, ok := previousExchangeConfigs[exchangeId]; !ok {
			return nil, fmt.Errorf("internal error: exchange %v not found in previousExchangeConfigs", exchangeId)
		}

		// Validate we have an encoder and fetcher subscribed for updates to each exchange.
		// An error here would be unexpected.
		if _, ok := pfmmc.mutableExchangeConfigUpdaters[exchangeId]; !ok {
			return nil, fmt.Errorf("internal error: price fetcher not found for exchange %v", exchangeId)
		}

		exchangeUpdaters := pfmmc.mutableExchangeConfigUpdaters[exchangeId]
		if err := exchangeUpdaters.Validate(); err != nil {
			return nil, fmt.Errorf("internal error for exchange %v: %w", exchangeId, err)
		}
	}

	// 3. Pre-compute updates to send to updaters.
	updaterToUpdateParameters := make(
		map[ExchangeConfigUpdater]UpdateParameters,
		len(newMutableExchangeConfigs)*expectedUpdatersPerExchange,
	)

	for exchangeId, updaters := range pfmmc.mutableExchangeConfigUpdaters {
		mutableExchangeConfig := newMutableExchangeConfigs[exchangeId]
		previousConfig := previousExchangeConfigs[exchangeId]

		// If the exchange config has changed, pre-compute the updates to send to the price fetcher and encoder, which
		// will be a copy of the updated exchange config, as well as a sorted list of copied market configs for each
		// market on the exchange.
		if !previousConfig.Equal(mutableExchangeConfig) {
			updaterToUpdateParameters[updaters.PriceFetcher] = getUpdateParametersForExchange(
				mutableExchangeConfig,
				newMutableMarketConfigs,
			)
			updaterToUpdateParameters[updaters.PriceEncoder] = getUpdateParametersForExchange(
				mutableExchangeConfig,
				newMutableMarketConfigs,
			)
		}
	}

	// 4. Take the writer lock.
	pfmmc.Lock()
	defer pfmmc.Unlock()

	// 5. Swap in new markets and exchange configs.
	pfmmc.mutableExchangeToConfigs = newMutableExchangeConfigs
	pfmmc.mutableMarketToConfigs = newMutableMarketConfigs

	// 6. For each changed exchange config, send an update to the associated price fetcher and encoder.
	// Update the encoder before the fetcher so that the encoder has all necessary market config to support
	// price updates coming from the fetcher in the event of adding a new market.
	// TODO(DEC-2020): use errors.Join once it's available.
	updaterErrors := make([]string, 0, len(updaterToUpdateParameters))

	for exchangeId, updaters := range pfmmc.mutableExchangeConfigUpdaters {
		// Update the encoder first.
		if updateParams, ok := updaterToUpdateParameters[updaters.PriceEncoder]; ok {
			err = updaters.PriceEncoder.UpdateMutableExchangeConfig(updateParams.ExchangeConfig, updateParams.MarketConfigs)
			if err != nil {
				updaterErrors = append(
					updaterErrors,
					fmt.Errorf(
						"UpdateMarkets: failed to update price encoder for exchange %v: %w",
						exchangeId,
						err,
					).Error(),
				)
			}
		}

		// Update the fetcher second.
		if updateParams, ok := updaterToUpdateParameters[updaters.PriceFetcher]; ok {
			err = updaters.PriceFetcher.UpdateMutableExchangeConfig(updateParams.ExchangeConfig, updateParams.MarketConfigs)
			if err != nil {
				updaterErrors = append(
					updaterErrors,
					fmt.Errorf(
						"UpdateMarkets: failed to update price fetcher for exchange %v: %w",
						exchangeId,
						err,
					).Error(),
				)
			}
		}
	}

	if len(updaterErrors) > 0 {
		err = fmt.Errorf("UpdateMarkets: failed to update some fetchers or encoders: %v", strings.Join(updaterErrors, ", "))
	}

	return marketParamErrors, err
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
	pfmmc.Lock()
	defer pfmmc.Unlock()
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
	pfmmc.Lock()
	defer pfmmc.Unlock()

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
	pfmmc.Lock()
	defer pfmmc.Unlock()

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
			float32(len(exchangeConfig.MarketToMarketConfig)),
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
		for marketId := range exchangeConfig.MarketToMarketConfig {
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
