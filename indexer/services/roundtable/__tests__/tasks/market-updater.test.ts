import {
  dbHelpers,
  FillTable,
  OraclePriceCreateObject,
  OrderTable,
  PerpetualMarketColumns,
  PerpetualMarketFromDatabase,
  PerpetualMarketTable,
  PerpetualMarketUpdateObject,
  PerpetualPositionTable,
  testConstants,
  OraclePriceTable,
  PriceMap,
  BlockTable,
  LiquidityTiersFromDatabase,
  LiquidityTiersTable,
  LiquidityTiersMap,
  LiquidityTiersColumns,
} from '@dydxprotocol-indexer/postgres';
import _ from 'lodash';

import { getUpdatedMarkets } from '../../src/helpers/websocket';
import marketUpdaterTask, { getPriceChange } from '../../src/tasks/market-updater';
import { expectMarketWebsocketMessage } from '../helpers/websocket-helpers';
import { producer } from '@dydxprotocol-indexer/kafka';
import { wrapBackgroundTask } from '@dydxprotocol-indexer/base';

import { synchronizeWrapBackgroundTask } from '@dydxprotocol-indexer/dev';
import { NextFundingCache, redis } from '@dydxprotocol-indexer/redis';
import { redisClient } from '../../src/helpers/redis';
import Big from 'big.js';
import { DateTime } from 'luxon';
import * as SubaccountTable from '@dydxprotocol-indexer/postgres/build/src/stores/subaccount-table';
import {
  defaultAsset,
  defaultAsset2,
  defaultAsset3,
  defaultBlock,
  defaultBlock2,
  defaultLiquidityTier,
  defaultLiquidityTier2,
  defaultMarket,
  defaultMarket2,
  defaultMarket3,
  defaultPerpetualMarket,
  defaultPerpetualMarket2,
  defaultPerpetualMarket3,
  defaultSubaccount,
  defaultSubaccount2,
  defaultTendermintEvent,
  defaultTendermintEvent2,
  defaultTendermintEvent3,
  defaultTendermintEvent4,
  defaultWallet,
} from '@dydxprotocol-indexer/postgres/build/__tests__/helpers/constants';
import * as MarketTable from '@dydxprotocol-indexer/postgres/build/src/stores/market-table';
import * as TendermintEventTable from '@dydxprotocol-indexer/postgres/build/src/stores/tendermint-event-table';
import * as AssetTable from '@dydxprotocol-indexer/postgres/build/src/stores/asset-table';
import * as WalletTable from '@dydxprotocol-indexer/postgres/build/src/stores/wallet-table';

jest.mock('@dydxprotocol-indexer/base', () => ({
  ...jest.requireActual('@dydxprotocol-indexer/base'),
  wrapBackgroundTask: jest.fn(),
}));

async function seedData() {
  await Promise.all([
    SubaccountTable.create(defaultSubaccount),
    SubaccountTable.create(defaultSubaccount2),
  ]);
  await Promise.all([
    MarketTable.create(defaultMarket),
    MarketTable.create(defaultMarket2),
    MarketTable.create(defaultMarket3),
  ]);
  await Promise.all([
    LiquidityTiersTable.create(defaultLiquidityTier),
    LiquidityTiersTable.create(defaultLiquidityTier2),
  ]);
  await Promise.all([
    PerpetualMarketTable.create(defaultPerpetualMarket),
    PerpetualMarketTable.create(defaultPerpetualMarket2),
    PerpetualMarketTable.create(defaultPerpetualMarket3),
  ]);
  await Promise.all([
    BlockTable.create(defaultBlock),
    BlockTable.create(defaultBlock2),
  ]);
  await Promise.all([
    TendermintEventTable.create(defaultTendermintEvent),
    TendermintEventTable.create(defaultTendermintEvent2),
    TendermintEventTable.create(defaultTendermintEvent3),
    TendermintEventTable.create(defaultTendermintEvent4),
  ]);
  await Promise.all([
    AssetTable.create(defaultAsset),
    AssetTable.create(defaultAsset2),
    AssetTable.create(defaultAsset3),
  ]);
  await Promise.all([
    WalletTable.create(defaultWallet),
  ]);
}

describe('market-updater', () => {

  const perpMarketUpdate1: PerpetualMarketUpdateObject = {
    id: testConstants.defaultPerpetualPosition.perpetualId,
    trades24H: 1,
    volume24H: testConstants.defaultFill.quoteAmount,
    openInterest: testConstants.defaultPerpetualPosition.size,
    nextFundingRate: '0.005',
  };
  const perpMarketUpdate2: PerpetualMarketUpdateObject = {
    id: testConstants.defaultPerpetualMarket2.id,
    trades24H: 0,
    volume24H: '0',
    openInterest: '0',
  };
  const perpMarketUpdate3: PerpetualMarketUpdateObject = {
    id: testConstants.defaultPerpetualMarket3.id,
    trades24H: 0,
    volume24H: '0',
    openInterest: '0',
  };

  beforeAll(async () => {
    await dbHelpers.migrate();
    await dbHelpers.clearData();
  });

  beforeEach(async () => {
    await seedData();
  });

  afterAll(async () => {
    await dbHelpers.teardown();
    jest.resetAllMocks();
  });

  afterEach(async () => {
    await dbHelpers.clearData();
    await redis.deleteAllAsync(redisClient);
    jest.clearAllMocks();
  });

  it('succeeds with no fills, positions or funding samples', async () => {
    const producerSendSpy: jest.SpyInstance = jest.spyOn(producer, 'send');
    const perpetualMarkets: PerpetualMarketFromDatabase[] = await PerpetualMarketTable.findAll(
      {},
      [],
      {},
    );

    const perpetualMarketMap: _.Dictionary<PerpetualMarketFromDatabase> = _.keyBy(
      perpetualMarkets,
      PerpetualMarketColumns.id,
    );
    const liquidityTiers:
    LiquidityTiersFromDatabase[] = await LiquidityTiersTable.findAll({}, []);

    const liquidityTiersMap: LiquidityTiersMap = _.keyBy(
      liquidityTiers,
      LiquidityTiersColumns.id,
    );
    await marketUpdaterTask();

    const newPerpetualMarketMap:
    _.Dictionary<PerpetualMarketFromDatabase> = _.chain(perpetualMarkets)
      .keyBy('id')
      .mapValues((perpetualMarket) => ({
        ...perpetualMarketMap[perpetualMarket.id],
        id: perpetualMarket.id,
        trades24H: 0,
        volume24H: '0',
        openInterest: '0',
      }))
      .value();

    const contents: string = JSON.stringify(
      getUpdatedMarkets(perpetualMarketMap, newPerpetualMarketMap, liquidityTiersMap),
    );
    await expectMarketWebsocketMessage(producerSendSpy, contents);
  });

  it('getPriceChange', () => {
    const latestPrices: PriceMap = {
      [testConstants.defaultOraclePrice.marketId]: '2',
      [testConstants.defaultOraclePrice2.marketId]: '3',
    };
    const previousPrices: PriceMap = {
      [testConstants.defaultOraclePrice.marketId]: '1',
    };
    expect(
      getPriceChange(testConstants.defaultOraclePrice.marketId, latestPrices, previousPrices),
    ).toEqual('1');
    expect(
      getPriceChange(testConstants.defaultOraclePrice2.marketId, latestPrices, previousPrices),
    ).toEqual(undefined);
  });

  it('getPriceChange with prices < 1e-6', () => {
    const latestPrices: PriceMap = {
      [testConstants.defaultOraclePrice.marketId]: '0.00000008',
      [testConstants.defaultOraclePrice2.marketId]: '0.00000009',
    };
    const previousPrices: PriceMap = {
      [testConstants.defaultOraclePrice.marketId]: '0.00000007',
    };
    expect(
      getPriceChange(testConstants.defaultOraclePrice.marketId, latestPrices, previousPrices),
    ).toEqual('0.00000001');
    expect(
      getPriceChange(testConstants.defaultOraclePrice2.marketId, latestPrices, previousPrices),
    ).toEqual(undefined);
  });

  it('succeeds with 24h price change', async () => {
    const producerSendSpy: jest.SpyInstance = jest.spyOn(producer, 'send');
    synchronizeWrapBackgroundTask(wrapBackgroundTask);

    const now: string = DateTime.local().toISO();
    const lessThan24HAgo: string = DateTime.local().minus({ hour: 23 }).toISO();
    const moreThan24HAgo: string = DateTime.local().minus({ hour: 24, minute: 5 }).toISO();

    const blockHeights = ['3', '4', '6', '7'];

    const blockPromises = blockHeights.map((height) => BlockTable.create({
      ...testConstants.defaultBlock,
      blockHeight: height,
    }),
    );

    await Promise.all(blockPromises);

    const oraclePrice3: OraclePriceCreateObject = {
      ...testConstants.defaultOraclePrice,
      price: '3',
      effectiveAtHeight: '3',
      effectiveAt: lessThan24HAgo,
    };
    const oraclePrice4: OraclePriceCreateObject = {
      ...testConstants.defaultOraclePrice,
      price: '4',
      effectiveAtHeight: '4',
      effectiveAt: moreThan24HAgo,
    };
    const oraclePrice6: OraclePriceCreateObject = {
      ...testConstants.defaultOraclePrice,
      price: '6',
      effectiveAtHeight: '6',
      effectiveAt: now,
    };
    await Promise.all([
      OraclePriceTable.create(oraclePrice3),
      OraclePriceTable.create(oraclePrice4),
      OraclePriceTable.create(oraclePrice6),
    ]);

    const perpetualMarkets: PerpetualMarketFromDatabase[] = await PerpetualMarketTable.findAll(
      {},
      [],
      {},
    );

    const perpetualMarketMap: _.Dictionary<PerpetualMarketFromDatabase> = _.keyBy(
      perpetualMarkets,
      PerpetualMarketColumns.id,
    );
    const liquidityTiers:
    LiquidityTiersFromDatabase[] = await LiquidityTiersTable.findAll({}, []);

    const liquidityTiersMap: LiquidityTiersMap = _.keyBy(
      liquidityTiers,
      LiquidityTiersColumns.id,
    );

    await marketUpdaterTask();

    const newPerpetualMarketMap: _.Dictionary<PerpetualMarketFromDatabase> = {};
    newPerpetualMarketMap[testConstants.defaultPerpetualPosition.perpetualId] = {
      ...perpetualMarketMap[testConstants.defaultPerpetualPosition.perpetualId],
      trades24H: 0,
      volume24H: '0',
      openInterest: '0',
      priceChange24H: '2',
    };
    newPerpetualMarketMap[testConstants.defaultPerpetualMarket2.id] = {
      ...perpetualMarketMap[testConstants.defaultPerpetualMarket2.id],
      trades24H: 0,
      volume24H: '0',
      openInterest: '0',
    };
    newPerpetualMarketMap[testConstants.defaultPerpetualMarket3.id] = {
      ...perpetualMarketMap[testConstants.defaultPerpetualMarket3.id],
      trades24H: 0,
      volume24H: '0',
      openInterest: '0',
    };

    const contents: string = JSON.stringify(
      getUpdatedMarkets(perpetualMarketMap, newPerpetualMarketMap, liquidityTiersMap),
    );
    await expectMarketWebsocketMessage(producerSendSpy, contents);
  });

  it('succeeds with one position, one fill and one funding sample', async () => {
    const producerSendSpy: jest.SpyInstance = jest.spyOn(producer, 'send');
    synchronizeWrapBackgroundTask(wrapBackgroundTask);
    await OrderTable.create(testConstants.defaultOrder);
    await Promise.all([
      FillTable.create(testConstants.defaultFill),
      PerpetualPositionTable.create(testConstants.defaultPerpetualPosition),
    ]);
    await NextFundingCache.addFundingSample(
      testConstants.defaultPerpetualMarket.ticker,
      new Big(perpMarketUpdate1.nextFundingRate!),
      redisClient,
    );

    const perpetualMarkets: PerpetualMarketFromDatabase[] = await PerpetualMarketTable.findAll(
      {},
      [],
    );

    const perpetualMarketMap: _.Dictionary<PerpetualMarketFromDatabase> = _.keyBy(
      perpetualMarkets,
      PerpetualMarketColumns.id,
    );

    const liquidityTiers:
    LiquidityTiersFromDatabase[] = await LiquidityTiersTable.findAll({}, []);

    const liquidityTiersMap: LiquidityTiersMap = _.keyBy(
      liquidityTiers,
      LiquidityTiersColumns.id,
    );
    await marketUpdaterTask();

    const newPerpetualMarketMap: _.Dictionary<PerpetualMarketFromDatabase> = {};
    newPerpetualMarketMap[testConstants.defaultPerpetualPosition.perpetualId] = {
      ...perpetualMarketMap[testConstants.defaultPerpetualPosition.perpetualId],
      ...perpMarketUpdate1,
    };
    newPerpetualMarketMap[testConstants.defaultPerpetualMarket2.id] = {
      ...perpetualMarketMap[testConstants.defaultPerpetualMarket2.id],
      ...perpMarketUpdate2,
    };
    newPerpetualMarketMap[testConstants.defaultPerpetualMarket3.id] = {
      ...perpetualMarketMap[testConstants.defaultPerpetualMarket3.id],
      ...perpMarketUpdate3,
    };

    const contents: string = JSON.stringify(
      getUpdatedMarkets(perpetualMarketMap, newPerpetualMarketMap, liquidityTiersMap),
    );
    await expectMarketWebsocketMessage(producerSendSpy, contents);
  });

  it('no message sent if no update, and no funding samples', async () => {
    const producerSendSpy: jest.SpyInstance = jest.spyOn(producer, 'send');
    synchronizeWrapBackgroundTask(wrapBackgroundTask);
    await OrderTable.create(testConstants.defaultOrder);
    await Promise.all([
      FillTable.create(testConstants.defaultFill),
      PerpetualPositionTable.create(testConstants.defaultPerpetualPosition),
    ]);

    await Promise.all([
      PerpetualMarketTable.update(perpMarketUpdate1),
      PerpetualMarketTable.update({
        ...perpMarketUpdate2,
        // set initial `nextFundingRate` to non-zero
        // we still don't expect update message if no funding samples.
        nextFundingRate: '0.000225',
      }),
      PerpetualMarketTable.update(perpMarketUpdate3),
    ]);

    await marketUpdaterTask();
    expect(producerSendSpy).toHaveBeenCalledTimes(0);
  });

  it('(non-zero default funding) no message sent if no update, and no funding samples', async () => {
    const producerSendSpy: jest.SpyInstance = jest.spyOn(producer, 'send');
    synchronizeWrapBackgroundTask(wrapBackgroundTask);
    await OrderTable.create(testConstants.defaultOrder);
    await Promise.all([
      FillTable.create(testConstants.defaultFill),
      PerpetualPositionTable.create(testConstants.defaultPerpetualPosition),
    ]);

    await Promise.all([
      PerpetualMarketTable.update({
        ...perpMarketUpdate1,
        defaultFundingRate1H: '0.0000125',
      }),
      PerpetualMarketTable.update({
        ...perpMarketUpdate2,
        defaultFundingRate1H: '0.0000125',
      }),
      PerpetualMarketTable.update(perpMarketUpdate3),
    ]);

    await marketUpdaterTask();
    expect(producerSendSpy).toHaveBeenCalledTimes(0);
  });

  it('(non-zero default funding) message send with funding sample of 0', async () => {
    const producerSendSpy: jest.SpyInstance = jest.spyOn(producer, 'send');
    synchronizeWrapBackgroundTask(wrapBackgroundTask);
    await OrderTable.create(testConstants.defaultOrder);
    await Promise.all([
      FillTable.create(testConstants.defaultFill),
      PerpetualPositionTable.create(testConstants.defaultPerpetualPosition),
    ]);

    await Promise.all([
      NextFundingCache.addFundingSample(
        testConstants.defaultPerpetualMarket.ticker,
        new Big(0), // `0` sample for BTC-USD
        redisClient,
      ),
      NextFundingCache.addFundingSample(
        testConstants.defaultPerpetualMarket2.ticker,
        new Big(0), // `0` sample for ETH-USD
        redisClient,
      ),
    ]);

    await Promise.all([
      PerpetualMarketTable.update({
        ...perpMarketUpdate1,
        // `nextFundingRate` equal to default funding rate.
        // With a funding sample of `0` we don't expect a mssage to be sent for BTC-USD
        nextFundingRate: '0.0000125',
        defaultFundingRate1H: '0.0000125',
      }),
      PerpetualMarketTable.update({
        ...perpMarketUpdate2,
        nextFundingRate: '0',
        // Differ from current `nextFundingRate` of 0.
        // Expect a message for this.
        defaultFundingRate1H: '0.0000125',
      }),
      PerpetualMarketTable.update(perpMarketUpdate3),
    ]);

    const perpetualMarkets: PerpetualMarketFromDatabase[] = await PerpetualMarketTable.findAll(
      {},
      [],
    );
    const perpetualMarketMap: _.Dictionary<PerpetualMarketFromDatabase> = _.keyBy(
      perpetualMarkets,
      PerpetualMarketColumns.id,
    );
    const liquidityTiers:
    LiquidityTiersFromDatabase[] = await LiquidityTiersTable.findAll({}, []);

    const liquidityTiersMap: LiquidityTiersMap = _.keyBy(
      liquidityTiers,
      LiquidityTiersColumns.id,
    );

    await marketUpdaterTask();

    const newPerpetualMarketMap: _.Dictionary<PerpetualMarketFromDatabase> = {};
    newPerpetualMarketMap[testConstants.defaultPerpetualPosition.perpetualId] = {
      ...perpetualMarketMap[testConstants.defaultPerpetualPosition.perpetualId],
    };
    newPerpetualMarketMap[testConstants.defaultPerpetualMarket2.id] = {
      ...perpetualMarketMap[testConstants.defaultPerpetualMarket2.id],
      nextFundingRate: '0.0000125',
    };
    newPerpetualMarketMap[testConstants.defaultPerpetualMarket3.id] = {
      ...perpetualMarketMap[testConstants.defaultPerpetualMarket3.id],
    };

    const contents: string = JSON.stringify(
      getUpdatedMarkets(perpetualMarketMap, newPerpetualMarketMap, liquidityTiersMap),
    );
    expectMarketWebsocketMessage(producerSendSpy, contents);
    expect(producerSendSpy).toHaveBeenCalledTimes(1);
  });

  it('update sent if position and fills update, but funding was not', async () => {
    const producerSendSpy: jest.SpyInstance = jest.spyOn(producer, 'send');
    synchronizeWrapBackgroundTask(wrapBackgroundTask);
    await OrderTable.create(testConstants.defaultOrder);
    await Promise.all([
      FillTable.create(testConstants.defaultFill),
      PerpetualPositionTable.create(testConstants.defaultPerpetualPosition),
    ]);
    await NextFundingCache.addFundingSample(
      testConstants.defaultPerpetualMarket.ticker,
      new Big(perpMarketUpdate1.nextFundingRate!),
      redisClient,
    );
    // Set funding to the rate returned by the cache
    await PerpetualMarketTable.update({
      id: perpMarketUpdate1.id,
      nextFundingRate: perpMarketUpdate1.nextFundingRate,
    });

    const perpetualMarkets: PerpetualMarketFromDatabase[] = await PerpetualMarketTable.findAll(
      {},
      [],
    );

    const perpetualMarketMap: _.Dictionary<PerpetualMarketFromDatabase> = _.keyBy(
      perpetualMarkets,
      PerpetualMarketColumns.id,
    );
    const liquidityTiers:
    LiquidityTiersFromDatabase[] = await LiquidityTiersTable.findAll({}, []);

    const liquidityTiersMap: LiquidityTiersMap = _.keyBy(
      liquidityTiers,
      LiquidityTiersColumns.id,
    );

    await marketUpdaterTask();

    const newPerpetualMarketMap: _.Dictionary<PerpetualMarketFromDatabase> = {};
    newPerpetualMarketMap[testConstants.defaultPerpetualPosition.perpetualId] = {
      ...perpetualMarketMap[testConstants.defaultPerpetualPosition.perpetualId],
      ...perpMarketUpdate1,
    };
    newPerpetualMarketMap[testConstants.defaultPerpetualMarket2.id] = {
      ...perpetualMarketMap[testConstants.defaultPerpetualMarket2.id],
      ...perpMarketUpdate2,
    };
    newPerpetualMarketMap[testConstants.defaultPerpetualMarket3.id] = {
      ...perpetualMarketMap[testConstants.defaultPerpetualMarket3.id],
      ...perpMarketUpdate3,
    };

    const contents: string = JSON.stringify(
      getUpdatedMarkets(perpetualMarketMap, newPerpetualMarketMap, liquidityTiersMap),
    );
    await expectMarketWebsocketMessage(producerSendSpy, contents);
  });

  it('update sent if funding updates, but positions and fills do not', async () => {
    const producerSendSpy: jest.SpyInstance = jest.spyOn(producer, 'send');
    synchronizeWrapBackgroundTask(wrapBackgroundTask);
    await OrderTable.create(testConstants.defaultOrder);
    await Promise.all([
      FillTable.create(testConstants.defaultFill),
      PerpetualPositionTable.create(testConstants.defaultPerpetualPosition),
    ]);
    await NextFundingCache.addFundingSample(
      testConstants.defaultPerpetualMarket.ticker,
      new Big(perpMarketUpdate1.nextFundingRate!),
      redisClient,
    );
    // Set up funding to be the only updated property
    await PerpetualMarketTable.update({
      ...perpMarketUpdate1,
      nextFundingRate: '0',
    });

    const perpetualMarkets: PerpetualMarketFromDatabase[] = await PerpetualMarketTable.findAll(
      {},
      [],
    );

    const perpetualMarketMap: _.Dictionary<PerpetualMarketFromDatabase> = _.keyBy(
      perpetualMarkets,
      PerpetualMarketColumns.id,
    );

    const liquidityTiers:
    LiquidityTiersFromDatabase[] = await LiquidityTiersTable.findAll({}, []);

    const liquidityTiersMap: LiquidityTiersMap = _.keyBy(
      liquidityTiers,
      LiquidityTiersColumns.id,
    );

    await marketUpdaterTask();

    const newPerpetualMarketMap: _.Dictionary<PerpetualMarketFromDatabase> = {};
    newPerpetualMarketMap[testConstants.defaultPerpetualPosition.perpetualId] = {
      ...perpetualMarketMap[testConstants.defaultPerpetualPosition.perpetualId],
      ...perpMarketUpdate1,
    };
    newPerpetualMarketMap[testConstants.defaultPerpetualMarket2.id] = {
      ...perpetualMarketMap[testConstants.defaultPerpetualMarket2.id],
      ...perpMarketUpdate2,
    };
    newPerpetualMarketMap[testConstants.defaultPerpetualMarket3.id] = {
      ...perpetualMarketMap[testConstants.defaultPerpetualMarket3.id],
      ...perpMarketUpdate3,
    };

    const contents: string = JSON.stringify(
      getUpdatedMarkets(perpetualMarketMap, newPerpetualMarketMap, liquidityTiersMap),
    );
    await expectMarketWebsocketMessage(producerSendSpy, contents);
  });

  it('no message sent if funding is cleared', async () => {
    const producerSendSpy: jest.SpyInstance = jest.spyOn(producer, 'send');
    synchronizeWrapBackgroundTask(wrapBackgroundTask);
    await OrderTable.create(testConstants.defaultOrder);
    await Promise.all([
      FillTable.create(testConstants.defaultFill),
      PerpetualPositionTable.create(testConstants.defaultPerpetualPosition),
    ]);
    await NextFundingCache.addFundingSample(
      testConstants.defaultPerpetualMarket.ticker,
      new Big(perpMarketUpdate1.nextFundingRate!),
      redisClient,
    );

    // Run the task once to update the markets
    await marketUpdaterTask();
    jest.clearAllMocks();

    // Clear all funding samples
    await NextFundingCache.clearFundingSamples(
      testConstants.defaultPerpetualMarket.ticker,
      redisClient,
    );

    await marketUpdaterTask();
    expect(producerSendSpy).toHaveBeenCalledTimes(0);
  });
});
