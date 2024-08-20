import {
  BlockTable,
  dbHelpers,
  PerpetualPositionTable,
  PnlTicksTable,
  testConstants,
  testMocks,
  WalletTable,
  SubaccountTable,
  LeaderboardPnlTable,
  LeaderboardPnlFromDatabase,
  LeaderboardPnlTimeSpan,
} from '@dydxprotocol-indexer/postgres';
import { LeaderboardPnlProcessedCache, redis } from '@dydxprotocol-indexer/redis';

import generateLeaderboardTaskFromTimespan from '../../src/tasks/create-leaderboard';
import { DateTime } from 'luxon';
import { redisClient } from '../../src/helpers/redis';

describe('create-leaderboard', () => {
  beforeAll(async () => {
    await dbHelpers.migrate();
    await dbHelpers.clearData();
  });

  beforeEach(async () => {
    await testMocks.seedData();
    await WalletTable.create(testConstants.defaultWallet3);
    await SubaccountTable.create(testConstants.defaultSubaccountWithAlternateAddress);
    await setupRankedPnlTicksData();
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

  it('Succeeds in populating the leaderboard with ranked pnl ticks', async () => {
    await Promise.all([
      PerpetualPositionTable.create(testConstants.defaultPerpetualPosition),
      PerpetualPositionTable.create({
        ...testConstants.defaultPerpetualPosition,
        perpetualId: testConstants.defaultPerpetualMarket2.id,
        openEventId: testConstants.defaultTendermintEventId2,
      }),
    ]);
    const task: () => Promise<void> = generateLeaderboardTaskFromTimespan(
      LeaderboardPnlTimeSpan.ALL_TIME);
    await task();
    const { results: pnlTicks } = await PnlTicksTable.findAll(
      {},
      [],
    );
    expect(pnlTicks.length).toEqual(2);
    const leaderboardResults: LeaderboardPnlFromDatabase[] = await LeaderboardPnlTable.findAll(
      {},
      [],
    );
    expect(leaderboardResults.length).toEqual(2);
  });

  it('leaderboard not updated if last processed pnl time < cached leaderboard time', async () => {
    await LeaderboardPnlProcessedCache.setProcessedTime(
      LeaderboardPnlTimeSpan.ALL_TIME,
      DateTime.utc().toISO(),
      redisClient,
    );
    const task: () => Promise<void> = generateLeaderboardTaskFromTimespan(
      LeaderboardPnlTimeSpan.ALL_TIME);
    await task();
    const leaderboardResults: LeaderboardPnlFromDatabase[] = await LeaderboardPnlTable.findAll(
      {},
      [],
    );
    expect(leaderboardResults.length).toEqual(0);
  });
});

async function setupRankedPnlTicksData() {
  await Promise.all([
    BlockTable.create({
      blockHeight: '3',
      time: testConstants.defaultBlock.time,
    }),
    BlockTable.create({
      blockHeight: '5',
      time: testConstants.defaultBlock.time,
    }),
  ]);
  await PnlTicksTable.createMany([
    {
      subaccountId: testConstants.defaultSubaccountId,
      equity: '1100',
      createdAt: DateTime.utc().toISO(),
      totalPnl: '1200',
      netTransfers: '50',
      blockHeight: '9',
      blockTime: testConstants.defaultBlock.time,
    },
    {
      subaccountId: testConstants.defaultSubaccountIdWithAlternateAddress,
      equity: '1090',
      createdAt: DateTime.utc().toISO(),
      totalPnl: '1190',
      netTransfers: '50',
      blockHeight: '7',
      blockTime: testConstants.defaultBlock.time,
    },
  ]);
}
