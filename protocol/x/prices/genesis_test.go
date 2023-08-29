package prices_test

import (
	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"

	indexerevents "github.com/dydxprotocol/v4-chain/protocol/indexer/events"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	sdktest "github.com/dydxprotocol/v4-chain/protocol/testutil/sdk"
	"github.com/dydxprotocol/v4-chain/protocol/x/prices"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := constants.Prices_DefaultGenesisState

	ctx, k, _, _, _, _ := keepertest.PricesKeepers(t)
	prices.InitGenesis(ctx, *k, genesisState)
	got := prices.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	require.ElementsMatch(t, genesisState.MarketParams, got.MarketParams)
}

func TestGenesisEmitsPriceUpdates(t *testing.T) {
	p := &mocks.PricesKeeper{}
	ctx := sdktest.NewContextWithBlockHeightAndTime(0, time.Now())

	updateEvent0 := indexerevents.NewMarketPriceUpdateEvent(0, constants.FiveBillion)
	updateEvent1 := indexerevents.NewMarketPriceUpdateEvent(1, constants.ThreeBillion)

	p.On("InitializeForGenesis", ctx).Return().Once()
	p.On("CreateMarket", ctx, mock.Anything, mock.Anything).Return(types.MarketParam{}, nil).Times(2)

	// Unclear why the "maybe" is needed here since the result is clearly used by the following calls.
	p.On("GenerateMarketPriceUpdateEvents", constants.Prices_DefaultGenesisState.MarketPrices).Return(
		[]*indexerevents.MarketEventV1{
			updateEvent0,
			updateEvent1,
		}).Maybe()

	indexerEventManager := &mocks.IndexerEventManager{}
	indexerEventManager.On("AddTxnEvent", ctx, indexerevents.SubtypeMarket, indexer_manager.GetB64EncodedEventMessage(updateEvent0)).Return().Once()
	indexerEventManager.On("AddTxnEvent", ctx, indexerevents.SubtypeMarket, indexer_manager.GetB64EncodedEventMessage(updateEvent1)).Return().Once()
	p.On("GetIndexerEventManager").Return(indexerEventManager).Times(2)

	prices.InitGenesis(ctx, p, constants.Prices_DefaultGenesisState)

	indexerEventManager.AssertExpectations(t)
	p.AssertExpectations(t)
}
