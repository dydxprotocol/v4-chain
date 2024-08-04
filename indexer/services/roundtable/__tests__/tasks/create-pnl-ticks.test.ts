import {
  BlockTable,
  dbHelpers,
  OraclePriceTable,
  PerpetualPositionTable,
  PnlTicksTable,
  testConstants,
  testMocks,
  TransferTable,
  FundingIndexUpdatesTable,
} from '@dydxprotocol-indexer/postgres';

import createPnlTicksTask from '../../src/tasks/create-pnl-ticks';
import { LatestAccountPnlTicksCache, PnlTickForSubaccounts, redis } from '@dydxprotocol-indexer/redis';
import { DateTime } from 'luxon';
import config from '../../src/config';
import { redisClient } from '../../src/helpers/redis';
import { logger } from '@dydxprotocol-indexer/base';

describe('create-pnl-ticks', () => {

  const pnlTickForSubaccounts: PnlTickForSubaccounts = {
    [testConstants.defaultSubaccountId]: {
      ...testConstants.defaultPnlTick,
      createdAt: DateTime.utc(2022, 6, 1, 0, 0, 0).toISO(),
      blockTime: DateTime.utc(2022, 6, 1, 0, 0, 0).toISO(),
    },
  };
  const existingPnlTicksNeedsUpdate: PnlTickForSubaccounts = {
    [testConstants.defaultSubaccountId]: {
      ...testConstants.defaultPnlTick,
      createdAt: DateTime.utc(2022, 5, 31, 23, 59, 0).toISO(),
      blockTime: DateTime.utc(2022, 5, 31, 23, 59, 0).toISO(),
    },
  };
  const dateTime: DateTime = DateTime.utc(2022, 6, 1, 0, 30, 0);

  beforeAll(async () => {
    await dbHelpers.migrate();
    await dbHelpers.clearData();
  });

  beforeEach(async () => {
    await testMocks.seedData();
    await Promise.all([
      BlockTable.create({
        ...testConstants.defaultBlock,
        blockHeight: '3',
      }),
      BlockTable.create({
        ...testConstants.defaultBlock,
        blockHeight: '4',
      }),
      BlockTable.create({
        ...testConstants.defaultBlock,
        blockHeight: '5',
      }),
    ]);
    await Promise.all([
      OraclePriceTable.create(testConstants.defaultOraclePrice),
      OraclePriceTable.create(testConstants.defaultOraclePrice2),
    ]);

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
    await Promise.all([
      FundingIndexUpdatesTable.create(testConstants.defaultFundingIndexUpdate),
      FundingIndexUpdatesTable.create({
        ...testConstants.defaultFundingIndexUpdate,
        perpetualId: testConstants.defaultPerpetualMarket2.id,
      }),
    ]);
  });

  afterAll(async () => {
    await dbHelpers.teardown();
    jest.resetAllMocks();
  });

  afterEach(async () => {
    await dbHelpers.clearData();
    await redis.deleteAllAsync(redisClient);
    jest.resetAllMocks();
  });

  it('succeeds with no prior pnl ticks and no open perpetual positions', async () => {
    const date: number = new Date(2023, 4, 18, 0, 0, 0).valueOf();
    jest.spyOn(Date, 'now').mockImplementation(() => date);
    jest.spyOn(DateTime, 'utc').mockImplementation(() => dateTime);
    await createPnlTicksTask();
    const { results: pnlTicks } = await PnlTicksTable.findAll(
      {},
      [],
      {},
    );
    expect(pnlTicks.length).toEqual(2);
    expect(pnlTicks).toEqual(
      expect.arrayContaining([
        {
          id: PnlTicksTable.uuid(testConstants.defaultSubaccountId2, dateTime.toISO()),
          createdAt: dateTime.toISO(),
          blockHeight: '5',
          blockTime: testConstants.defaultBlock.time,
          equity: '0.000000',
          netTransfers: '20.500000',
          subaccountId: testConstants.defaultSubaccountId2,
          totalPnl: '-20.500000',
        },
        {
          id: PnlTicksTable.uuid(testConstants.defaultSubaccountId, dateTime.toISO()),
          createdAt: dateTime.toISO(),
          blockHeight: '5',
          blockTime: testConstants.defaultBlock.time,
          equity: '0.000000',
          netTransfers: '-20.500000',
          subaccountId: testConstants.defaultSubaccountId,
          totalPnl: '20.500000',
        },
      ]),
    );
  });

  it('succeeds with no prior pnl ticks and open perpetual positions', async () => {
    const date: number = new Date(2023, 4, 18, 0, 0, 0).valueOf();
    jest.spyOn(Date, 'now').mockImplementation(() => date);
    jest.spyOn(DateTime, 'utc').mockImplementation(() => dateTime);
    await Promise.all([
      PerpetualPositionTable.create(testConstants.defaultPerpetualPosition),
      PerpetualPositionTable.create({
        ...testConstants.defaultPerpetualPosition,
        perpetualId: testConstants.defaultPerpetualMarket2.id,
        openEventId: testConstants.defaultTendermintEventId2,
      }),
    ]);
    await createPnlTicksTask();
    const { results: pnlTicks } = await PnlTicksTable.findAll(
      {},
      [],
      {},
    );
    expect(pnlTicks.length).toEqual(2);
    expect(pnlTicks).toEqual(
      expect.arrayContaining([
        {
          id: PnlTicksTable.uuid(testConstants.defaultSubaccountId2, dateTime.toISO()),
          createdAt: dateTime.toISO(),
          blockHeight: '5',
          blockTime: testConstants.defaultBlock.time,
          equity: '0.000000',
          netTransfers: '20.500000',
          subaccountId: testConstants.defaultSubaccountId2,
          totalPnl: '-20.500000',
        },
        {
          id: PnlTicksTable.uuid(testConstants.defaultSubaccountId, dateTime.toISO()),
          createdAt: dateTime.toISO(),
          blockHeight: '5',
          blockTime: testConstants.defaultBlock.time,
          equity: '105000.000000',
          netTransfers: '-20.500000',
          subaccountId: testConstants.defaultSubaccountId,
          totalPnl: '105020.500000',
        },
      ]),
    );
  });

  it(
    'succeeds with prior pnl ticks and open perpetual positions, updates pnl correctly',
    async () => {
      const date: number = new Date(2023, 4, 18, 0, 0, 0).valueOf();
      jest.spyOn(Date, 'now').mockImplementation(() => date);
      config.PNL_TICK_UPDATE_INTERVAL_MS = 3_600_000;
      jest.spyOn(DateTime, 'utc').mockImplementation(() => dateTime);
      await LatestAccountPnlTicksCache.set(
        existingPnlTicksNeedsUpdate,
        redisClient,
      );
      await Promise.all([
        PerpetualPositionTable.create(testConstants.defaultPerpetualPosition),
        PerpetualPositionTable.create({
          ...testConstants.defaultPerpetualPosition,
          perpetualId: testConstants.defaultPerpetualMarket2.id,
          openEventId: testConstants.defaultTendermintEventId2,
        }),
      ]);
      await createPnlTicksTask();
      const { results: pnlTicks } = await PnlTicksTable.findAll(
        {},
        [],
        {},
      );
      expect(pnlTicks.length).toEqual(2);
      expect(pnlTicks).toEqual(
        expect.arrayContaining([
          {
            id: PnlTicksTable.uuid(testConstants.defaultSubaccountId2, dateTime.toISO()),
            createdAt: dateTime.toISO(),
            blockHeight: '5',
            blockTime: testConstants.defaultBlock.time,
            equity: '0.000000',
            netTransfers: '20.500000',
            subaccountId: testConstants.defaultSubaccountId2,
            totalPnl: '-20.500000',
          },
          {
            id: PnlTicksTable.uuid(testConstants.defaultSubaccountId, dateTime.toISO()),
            createdAt: dateTime.toISO(),
            blockHeight: '5',
            blockTime: testConstants.defaultBlock.time,
            equity: '105000.000000',
            netTransfers: '-20.500000',
            subaccountId: testConstants.defaultSubaccountId,
            totalPnl: '105020.500000',
          },
        ]),
      );
    });

  it(
    'succeeds with prior pnl ticks and open perpetual positions, respects PNL_TICK_UPDATE_INTERVAL_MS',
    async () => {
      const date: number = new Date(2023, 4, 18, 0, 0, 0).valueOf();
      jest.spyOn(Date, 'now').mockImplementation(() => date);
      config.PNL_TICK_UPDATE_INTERVAL_MS = 3_600_000;
      jest.spyOn(DateTime, 'utc').mockImplementation(() => dateTime);
      await LatestAccountPnlTicksCache.set(
        pnlTickForSubaccounts,
        redisClient,
      );
      await Promise.all([
        PerpetualPositionTable.create(testConstants.defaultPerpetualPosition),
        PerpetualPositionTable.create({
          ...testConstants.defaultPerpetualPosition,
          perpetualId: testConstants.defaultPerpetualMarket2.id,
          openEventId: testConstants.defaultTendermintEventId2,
        }),
      ]);
      await createPnlTicksTask();
      const { results: pnlTicks } = await PnlTicksTable.findAll(
        {},
        [],
        {},
      );
      expect(pnlTicks.length).toEqual(1);
      expect(pnlTicks).toEqual(
        expect.arrayContaining([
          {
            id: PnlTicksTable.uuid(testConstants.defaultSubaccountId2, dateTime.toISO()),
            createdAt: dateTime.toISO(),
            blockHeight: '5',
            blockTime: testConstants.defaultBlock.time,
            equity: '0.000000',
            netTransfers: '20.500000',
            subaccountId: testConstants.defaultSubaccountId2,
            totalPnl: '-20.500000',
          },
        ]),
      );
    });

  it(
    'no-op if PNL_TICK_UPDATE_INTERVAL_MS has not been reached',
    async () => {
      config.PNL_TICK_UPDATE_INTERVAL_MS = 3_600_000;
      await PnlTicksTable.create({
        ...testConstants.defaultPnlTick,
        blockTime: testConstants.defaultBlock.time,
      });
      const {
        maxBlockTime,
        // eslint-disable-next-line @typescript-eslint/no-unused-vars
        count,
      }: {
        maxBlockTime: string,
        count: number,
      } = await PnlTicksTable.findLatestProcessedBlocktimeAndCount();

      const date: number = Date.parse(maxBlockTime).valueOf();
      jest.spyOn(Date, 'now').mockImplementation(() => date);
      jest.spyOn(DateTime, 'utc').mockImplementation(() => dateTime);
      jest.spyOn(logger, 'info');
      await LatestAccountPnlTicksCache.set(
        pnlTickForSubaccounts,
        redisClient,
      );
      await Promise.all([
        PerpetualPositionTable.create(testConstants.defaultPerpetualPosition),
        PerpetualPositionTable.create({
          ...testConstants.defaultPerpetualPosition,
          perpetualId: testConstants.defaultPerpetualMarket2.id,
          openEventId: testConstants.defaultTendermintEventId2,
        }),
      ]);
      await createPnlTicksTask();
      const { results: pnlTicks } = await PnlTicksTable.findAll(
        {},
        [],
        {},
      );
      // no new pnl ticks should be created.
      expect(pnlTicks.length).toEqual(1);
      expect(logger.info).toHaveBeenCalledWith(
        expect.objectContaining({
          message: 'Skipping run because update interval has not been reached and all subaccounts have been processed',
        }),
      );
    });

  it(
    'calculates pnl per subaccount if last run hit subaccount limit',
    async () => {
      const pnlTicksHelper = require('../../src/helpers/pnl-ticks-helper');
      const getPnlTicksCreateObjectsSpy = jest.spyOn(pnlTicksHelper, 'getPnlTicksCreateObjects');
      config.PNL_TICK_UPDATE_INTERVAL_MS = 3_600_000;
      config.PNL_TICK_MAX_ACCOUNTS_PER_RUN = 1;
      await PnlTicksTable.create({
        ...testConstants.defaultPnlTick,
        blockTime: testConstants.defaultBlock.time,
      });
      const {
        maxBlockTime,
        // eslint-disable-next-line @typescript-eslint/no-unused-vars
        count,
      }: {
        maxBlockTime: string,
        count: number,
      } = await PnlTicksTable.findLatestProcessedBlocktimeAndCount();

      const date: number = Date.parse(maxBlockTime).valueOf();
      jest.spyOn(Date, 'now').mockImplementation(() => date);
      jest.spyOn(DateTime, 'utc').mockImplementation(() => dateTime);
      jest.spyOn(logger, 'info');
      await LatestAccountPnlTicksCache.set(
        pnlTickForSubaccounts,
        redisClient,
      );
      await Promise.all([
        PerpetualPositionTable.create(testConstants.defaultPerpetualPosition),
        PerpetualPositionTable.create({
          ...testConstants.defaultPerpetualPosition,
          perpetualId: testConstants.defaultPerpetualMarket2.id,
          openEventId: testConstants.defaultTendermintEventId2,
        }),
      ]);
      await createPnlTicksTask();
      expect(getPnlTicksCreateObjectsSpy).toHaveBeenCalledTimes(1);
    });
});
