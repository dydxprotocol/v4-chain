import {
  BlockTable,
  dbHelpers,
  FundingIndexMap,
  FundingIndexUpdatesCreateObject,
  FundingIndexUpdatesTable,
  IsoString,
  PerpetualPositionFromDatabase,
  PerpetualPositionTable,
  PnlTicksCreateObject,
  PriceMap,
  SubaccountAssetNetTransferMap,
  SubaccountFromDatabase,
  SubaccountTable,
  testConstants,
  testMocks,
  Transaction,
  TransferTable,
  PositionSide,
  PerpetualMarketStatus,
  PerpetualMarketType,
  MarketTable,
  PerpetualMarketTable,
  perpetualMarketRefresher,
} from '@dydxprotocol-indexer/postgres';
import {
  calculateEquity,
  calculateTotalPnl,
  getBlockHeightToFundingIndexMap,
  getNewPnlTick,
  getPnlTicksCreateObjects,
  getUsdcTransfersSinceLastPnlTick,
  getAccountsToUpdate,
  normalizeStartTime,
} from '../../src/helpers/pnl-ticks-helper';
import { defaultPnlTickForSubaccounts } from '../../src/helpers/constants';
import Big from 'big.js';
import { DateTime } from 'luxon';
import { LatestAccountPnlTicksCache, PnlTickForSubaccounts, redis } from '@dydxprotocol-indexer/redis';
import { redisClient } from '../../src/helpers/redis';
import { ZERO } from '../../src/lib/constants';
import { SubaccountUsdcTransferMap } from '../../src/helpers/types';
import config from '../../src/config';
import _ from 'lodash';
import { ONE_HOUR_IN_MILLISECONDS } from '@dydxprotocol-indexer/base';

describe('pnl-ticks-helper', () => {
  const positions: PerpetualPositionFromDatabase[] = [
    {
      ...testConstants.defaultPerpetualPosition,
      entryPrice: '20000',
      sumOpen: '10',
      sumClose: '0',
      id: testConstants.defaultPerpetualPositionId,
    },
  ];
  const lastUpdatedFundingIndexMap: FundingIndexMap = {
    [testConstants.defaultPerpetualMarket.id]: Big('10050'),
    [testConstants.defaultPerpetualMarket2.id]: Big('5'),
    5: Big('5'),
  };
  const currentFundingIndexMap: FundingIndexMap = {
    [testConstants.defaultPerpetualMarket.id]: Big('11000'),
    [testConstants.defaultPerpetualMarket2.id]: Big('8'),
    5: Big('8'),
  };
  const marketPrices: PriceMap = {
    [testConstants.defaultPerpetualMarket.id]: '20000',
    [testConstants.defaultPerpetualMarket2.id]: '1000',
  };
  const pnlTickForSubaccounts: PnlTickForSubaccounts = {
    [testConstants.defaultSubaccountId]: testConstants.defaultPnlTick,
  };
  const dateTime: DateTime = DateTime.utc();

  beforeAll(async () => {
    await dbHelpers.migrate();
    await dbHelpers.clearData();
  });

  beforeEach(async () => {
    await testMocks.seedData();
    await perpetualMarketRefresher.updatePerpetualMarkets();
  });

  afterAll(async () => {
    await dbHelpers.teardown();
  });

  afterEach(async () => {
    await dbHelpers.clearData();
    await redis.deleteAllAsync(redisClient);
    jest.resetAllMocks();
  });

  it('getUsdcTransfersSinceLastPnlTick no transfers', async () => {
    const subaccountIds: string[] = [
      testConstants.defaultSubaccountId,
      testConstants.defaultSubaccountId2,
    ];
    const blockHeight: string = '5';
    const netUsdcTransfers: SubaccountUsdcTransferMap = await getUsdcTransfersSinceLastPnlTick(
      subaccountIds,
      defaultPnlTickForSubaccounts,
      blockHeight,
    );
    expect(netUsdcTransfers).toEqual({});
  });

  it('getBlockHeightToFundingIndexMap', async () => {
    await Promise.all([
      BlockTable.create({
        blockHeight: '3',
        time: testConstants.defaultBlock.time,
      }),
      SubaccountTable.create({
        ...testConstants.defaultSubaccount3,
        updatedAtHeight: '3',
      }),
      SubaccountTable.update({
        id: testConstants.defaultSubaccountId,
        updatedAtHeight: '1',
        updatedAt: testConstants.defaultSubaccount.updatedAt,
      }),
      SubaccountTable.update({
        id: testConstants.defaultSubaccountId2,
        updatedAtHeight: '2',
        updatedAt: testConstants.defaultSubaccount.updatedAt,
      }),
    ]);

    const fundingIndexUpdate2: FundingIndexUpdatesCreateObject = {
      ...testConstants.defaultFundingIndexUpdate,
      perpetualId: testConstants.defaultPerpetualMarket2.id,
      fundingIndex: '5',
      effectiveAtHeight: '2',
    };
    const fundingIndexUpdate3: FundingIndexUpdatesCreateObject = {
      ...testConstants.defaultFundingIndexUpdate,
      eventId: testConstants.defaultTendermintEventId2,
      fundingIndex: '100',
      effectiveAtHeight: '1',
    };
    const fundingIndexUpdate4: FundingIndexUpdatesCreateObject = {
      ...testConstants.defaultFundingIndexUpdate,
      perpetualId: testConstants.defaultPerpetualMarket2.id,
      eventId: testConstants.defaultTendermintEventId2,
      fundingIndex: '2',
      effectiveAtHeight: '1',
    };
    await Promise.all([
      FundingIndexUpdatesTable.create(testConstants.defaultFundingIndexUpdate),
      FundingIndexUpdatesTable.create(fundingIndexUpdate2),
      FundingIndexUpdatesTable.create(fundingIndexUpdate3),
      FundingIndexUpdatesTable.create(fundingIndexUpdate4),
    ]);
    const subaccountsWithTransfers: SubaccountFromDatabase[] = await SubaccountTable.findAll(
      {}, [], {},
    );
    const accountsToUpdate1: string[] = [testConstants.defaultSubaccountId];
    const heightToFundingIndices1:
    _.Dictionary<FundingIndexMap> = await getBlockHeightToFundingIndexMap(
      subaccountsWithTransfers, accountsToUpdate1,
    );

    expect(heightToFundingIndices1[1][testConstants.defaultPerpetualMarket.id]).toEqual(Big('100'));
    expect(heightToFundingIndices1[1][testConstants.defaultPerpetualMarket2.id]).toEqual(Big('2'));

    const accountsToUpdate2: string[] = [testConstants.defaultSubaccountId2];
    const heightToFundingIndices2:
    _.Dictionary<FundingIndexMap> = await getBlockHeightToFundingIndexMap(
      subaccountsWithTransfers, accountsToUpdate2,
    );

    expect(heightToFundingIndices2[2][testConstants.defaultPerpetualMarket.id]).toEqual(Big('10050'));
    expect(heightToFundingIndices2[2][testConstants.defaultPerpetualMarket2.id]).toEqual(Big('5'));
    const accountsToUpdate3: string[] = [
      testConstants.defaultSubaccountId,
      testConstants.defaultSubaccountId2,
      testConstants.defaultSubaccountId3,
    ];
    const heightToFundingIndices3:
    _.Dictionary<FundingIndexMap> = await getBlockHeightToFundingIndexMap(
      subaccountsWithTransfers, accountsToUpdate3,
    );

    expect(heightToFundingIndices3[1][testConstants.defaultPerpetualMarket.id]).toEqual(Big('100'));
    expect(heightToFundingIndices3[1][testConstants.defaultPerpetualMarket2.id]).toEqual(Big('2'));
    expect(heightToFundingIndices3[2][testConstants.defaultPerpetualMarket.id]).toEqual(Big('10050'));
    expect(heightToFundingIndices3[2][testConstants.defaultPerpetualMarket2.id]).toEqual(Big('5'));
    expect(heightToFundingIndices3[3][testConstants.defaultPerpetualMarket.id]).toEqual(Big('10050'));
    expect(heightToFundingIndices3[3][testConstants.defaultPerpetualMarket2.id]).toEqual(Big('5'));
  });

  it('getUsdcTransfersSinceLastPnlTick with transfers', async () => {
    const subaccountIds: string[] = [
      testConstants.defaultSubaccountId,
      testConstants.defaultSubaccountId2,
    ];
    const blockHeight: string = '5';
    await Promise.all([
      TransferTable.create({
        ...testConstants.defaultTransfer,
        createdAtHeight: '3',
      }),
      TransferTable.create({
        ...testConstants.defaultTransfer,
        size: '10.5',
        createdAtHeight: '4',
        eventId: testConstants.defaultTendermintEventId2,
      }),
    ]);
    const netUsdcTransfers: SubaccountUsdcTransferMap = await getUsdcTransfersSinceLastPnlTick(
      subaccountIds,
      defaultPnlTickForSubaccounts,
      blockHeight,
    );
    expect(netUsdcTransfers).toEqual(expect.objectContaining({
      [testConstants.defaultSubaccountId]: new Big('-20.5'),
      [testConstants.defaultSubaccountId2]: new Big('20.5'),
    }));
  });

  it('normalizeStartTime', () => {
    const time: Date = new Date('2021-01-09T20:00:50.000Z');
    // 1 hour
    config.PNL_TICK_UPDATE_INTERVAL_MS = 1000 * 60 * 60;
    const result1: Date = normalizeStartTime(time);
    expect(result1.toISOString()).toBe('2021-01-09T20:00:00.000Z');
    // 1 day
    config.PNL_TICK_UPDATE_INTERVAL_MS = 1000 * 60 * 60 * 24;
    const result2: Date = normalizeStartTime(time);
    expect(result2.toISOString()).toBe('2021-01-09T00:00:00.000Z');
  });

  it('getAccountsToUpdate', () => {
    const accountToLastUpdatedBlockTime: _.Dictionary<IsoString> = {
      account1: '2024-01-01T10:00:00Z',
      account2: '2024-01-01T11:00:00Z',
      account3: '2024-01-01T11:01:00Z',
      account4: '2024-01-01T11:10:00Z',
      account5: '2024-01-01T12:00:00Z',
      account6: '2024-01-01T12:00:10Z',
    };
    const blockTime: IsoString = '2024-01-01T12:01:00Z';
    config.PNL_TICK_UPDATE_INTERVAL_MS = ONE_HOUR_IN_MILLISECONDS;

    const expectedAccountsToUpdate: string[] = ['account1', 'account2', 'account3', 'account4'];
    const accountsToUpdate: string[] = getAccountsToUpdate(
      accountToLastUpdatedBlockTime,
      blockTime,
    );
    expect(accountsToUpdate).toEqual(expectedAccountsToUpdate);
  });

  it('calculateEquity', () => {
    const usdcPosition: Big = new Big('100');
    const equity: Big = calculateEquity(
      usdcPosition,
      positions,
      marketPrices,
      lastUpdatedFundingIndexMap,
      currentFundingIndexMap,
    );
    expect(equity).toEqual(new Big('190600'));
  });

  it('calculateEquity with no positions', () => {
    const usdcPosition: Big = new Big('100');
    const equity: Big = calculateEquity(
      usdcPosition,
      [],
      marketPrices,
      {},
      {},
    );
    expect(equity).toEqual(usdcPosition);
  });

  it('calculateEquity with LONG position', () => {
    const longPosition: PerpetualPositionFromDatabase = {
      ...testConstants.defaultPerpetualPosition,
      perpetualId: testConstants.defaultPerpetualMarket2.id,
      entryPrice: '20000',
      sumOpen: '10',
      sumClose: '0',
      openEventId: testConstants.defaultTendermintEventId2,
      id: PerpetualPositionTable.uuid(
        testConstants.defaultPerpetualPosition.subaccountId,
        testConstants.defaultTendermintEventId2,
      ),
    };
    const usdcPosition: Big = new Big('10000');
    const equity: Big = calculateEquity(
      usdcPosition,
      [longPosition],
      marketPrices,
      lastUpdatedFundingIndexMap,
      currentFundingIndexMap,
    );
    expect(equity).toEqual(new Big('19970'));
  });

  it('calculateEquity with SHORT position', () => {
    const shortPosition: PerpetualPositionFromDatabase = {
      ...testConstants.defaultPerpetualPosition,
      perpetualId: testConstants.defaultPerpetualMarket2.id,
      side: PositionSide.SHORT,
      entryPrice: '20000',
      size: '-10',
      sumOpen: '10',
      sumClose: '0',
      openEventId: testConstants.defaultTendermintEventId2,
      id: PerpetualPositionTable.uuid(
        testConstants.defaultPerpetualPosition.subaccountId,
        testConstants.defaultTendermintEventId2,
      ),
    };
    const usdcPosition: Big = new Big('10000');
    const equity: Big = calculateEquity(
      usdcPosition,
      [shortPosition],
      marketPrices,
      lastUpdatedFundingIndexMap,
      currentFundingIndexMap,
    );
    expect(equity).toEqual(new Big('30'));
  });

  it('calculateEquity with multiple positions', () => {
    const positions2: PerpetualPositionFromDatabase[] = [
      ...positions,
      {
        ...testConstants.defaultPerpetualPosition,
        side: PositionSide.SHORT,
        perpetualId: testConstants.defaultPerpetualMarket2.id,
        entryPrice: '20000',
        sumOpen: '10',
        size: '-10',
        sumClose: '0',
        openEventId: testConstants.defaultTendermintEventId2,
        id: PerpetualPositionTable.uuid(
          testConstants.defaultPerpetualPosition.subaccountId,
          testConstants.defaultTendermintEventId2,
        ),
      },
    ];
    const usdcPosition: Big = new Big('10000');
    const equity: Big = calculateEquity(
      usdcPosition,
      positions2,
      marketPrices,
      lastUpdatedFundingIndexMap,
      currentFundingIndexMap,
    );
    expect(equity).toEqual(new Big('190530'));
  });

  it('calculateEquity with perpetualId/marketId mismatch', async () => {
    await MarketTable.create({
      id: 1000,
      pair: 'TEST-USD',
      exponent: -12,
      minPriceChangePpm: 50,
      oraclePrice: '1.00',
    });
    await PerpetualMarketTable.create({
      id: '5',
      clobPairId: '5',
      ticker: 'TEST-USD',
      marketId: 1000,
      status: PerpetualMarketStatus.ACTIVE,
      priceChange24H: '0.000000001',
      volume24H: '10000000',
      trades24H: 200,
      nextFundingRate: '1.2',
      openInterest: '40000',
      quantumConversionExponent: -16,
      atomicResolution: -2,
      subticksPerTick: 10,
      stepBaseQuantums: 1,
      liquidityTierId: 0,
      marketType: PerpetualMarketType.ISOLATED,
      baseOpenInterest: '100000',
      defaultFundingRate1H: '0.0001',
    });
    await perpetualMarketRefresher.updatePerpetualMarkets();
    const positions2: PerpetualPositionFromDatabase[] = [
      ...positions,
      {
        ...testConstants.defaultPerpetualPosition,
        side: PositionSide.SHORT,
        perpetualId: '5',
        entryPrice: '20000',
        sumOpen: '10',
        size: '-10',
        sumClose: '0',
        openEventId: testConstants.defaultTendermintEventId2,
        id: PerpetualPositionTable.uuid(
          testConstants.defaultPerpetualPosition.subaccountId,
          testConstants.defaultTendermintEventId2,
        ),
      },
    ];
    const marketPricesMissingPrice: PriceMap = {
      [testConstants.defaultPerpetualMarket.id]: '20000',
      1000: '1000',
    };
    const usdcPosition: Big = new Big('10000');
    const equity: Big = calculateEquity(
      usdcPosition,
      positions2,
      marketPricesMissingPrice,
      lastUpdatedFundingIndexMap,
      currentFundingIndexMap,
    );
    expect(equity).toEqual(new Big('190530'));
  });

  it('calculateTotalPnl', () => {
    const equity: Big = new Big('200100');
    const transfers: string = '-20.5';
    const totalPnl: Big = calculateTotalPnl(
      equity,
      transfers,
    );
    expect(totalPnl).toEqual(new Big('200120.5'));
  });

  it('calculateTotalPnl with 0 equity', () => {
    const equity: Big = ZERO;
    const transfers: string = '-20.5';
    const totalPnl: Big = calculateTotalPnl(
      equity,
      transfers,
    );
    expect(totalPnl).toEqual(new Big('20.5'));
  });

  it('getNewPnlTick', () => {
    const subaccountAssetNetTransferMap: SubaccountAssetNetTransferMap = {
      [testConstants.defaultSubaccountId]: {
        [testConstants.defaultAsset.id]: '-20.5',
        [testConstants.defaultAsset2.id]: '30.5',
      },
      [testConstants.defaultSubaccountId2]: {
        [testConstants.defaultAsset.id]: '10',
      },
    };
    const usdcPosition: Big = new Big('100');
    const usdcNetTransfersSinceLastPnlTick: Big = new Big('-5.5');
    const latestBlockHeight: string = '5';
    const latestBlockTime: IsoString = DateTime.utc(2022, 6, 2, 0, 30).toISO();
    const pnlTick: PnlTicksCreateObject = getNewPnlTick(
      testConstants.defaultSubaccountId,
      subaccountAssetNetTransferMap,
      marketPrices,
      positions,
      usdcPosition,
      usdcNetTransfersSinceLastPnlTick,
      dateTime,
      latestBlockHeight,
      latestBlockTime,
      pnlTickForSubaccounts,
      lastUpdatedFundingIndexMap,
      currentFundingIndexMap,
    );
    expect(pnlTick).toEqual({
      subaccountId: testConstants.defaultSubaccountId,
      equity: '190600.000000',
      totalPnl: '190620.500000',
      netTransfers: '-5.500000',
      createdAt: dateTime.toISO(),
      blockHeight: latestBlockHeight,
      blockTime: latestBlockTime,
    });
  });

  it('getNewPnlTicks with prior pnl ticks', async () => {
    config.PNL_TICK_UPDATE_INTERVAL_MS = 3_600_000;
    const blockHeight: string = '5';
    const blockTime: IsoString = DateTime.utc(2022, 6, 2, 0, 30).toISO();
    const ticksForSubaccounts: PnlTickForSubaccounts = {
      [testConstants.defaultSubaccountId]: {
        ...testConstants.defaultPnlTick,
        createdAt: DateTime.utc(2022, 6, 2).toISO(),
      },
    };
    await LatestAccountPnlTicksCache.set(
      ticksForSubaccounts,
      redisClient,
    );
    await BlockTable.create({
      blockHeight,
      time: blockTime,
    });
    await TransferTable.create(testConstants.defaultTransfer);
    const txId: number = await Transaction.start();
    jest.spyOn(DateTime, 'utc').mockImplementation(() => dateTime);
    const newTicksToCreate: PnlTicksCreateObject[] = await
    getPnlTicksCreateObjects(txId);
    await Transaction.rollback(txId);
    expect(newTicksToCreate.length).toEqual(1);
    expect(newTicksToCreate).toEqual(
      expect.arrayContaining([
        {
          createdAt: dateTime.toISO(),
          blockHeight,
          blockTime,
          equity: '0.000000',
          netTransfers: '10.000000',
          subaccountId: testConstants.defaultSubaccountId2,
          totalPnl: '-10.000000',
        },
      ]),
    );
  });

  it('getNewPnlTicks without prior pnl ticks', async () => {
    jest.spyOn(DateTime, 'utc').mockImplementation(() => dateTime);
    const blockHeight: string = '5';
    const blockTime: IsoString = DateTime.utc(2022, 6, 2, 0, 30).toISO();
    await TransferTable.create(testConstants.defaultTransfer);
    await BlockTable.create({
      blockHeight,
      time: blockTime,
    });
    const txId: number = await Transaction.start();
    const newTicksToCreate: PnlTicksCreateObject[] = await
    getPnlTicksCreateObjects(txId);
    await Transaction.rollback(txId);
    expect(newTicksToCreate.length).toEqual(2);
    expect(newTicksToCreate).toEqual(
      expect.arrayContaining([
        {
          createdAt: dateTime.toISO(),
          blockHeight,
          blockTime,
          equity: '0.000000',
          netTransfers: '10.000000',
          subaccountId: testConstants.defaultSubaccountId2,
          totalPnl: '-10.000000',
        },
        {
          createdAt: dateTime.toISO(),
          blockHeight,
          blockTime,
          equity: '0.000000',
          netTransfers: '-10.000000',
          subaccountId: testConstants.defaultSubaccountId,
          totalPnl: '10.000000',
        },
      ]),
    );
  });
});
