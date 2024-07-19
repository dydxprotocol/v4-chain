import { LeaderboardPnlFromDatabase } from '../../src/types';
import * as LeaderboardPnlTable from '../../src/stores/leaderboard-pnl-table';
import { clearData, migrate, teardown } from '../../src/helpers/db-helpers';
import {
  defaultLeaderboardPnl2OneDay,
  defaultLeaderboardPnlOneDay,
  defaultLeaderboardPnl1AllTime,
  defaultLeaderboardPnlOneDayToUpsert,
  defaultWallet3,
} from '../helpers/constants';
import { seedData } from '../helpers/mock-generators';
import { WalletTable } from '../../src';

describe('LeaderboardPnl store', () => {
  beforeEach(async () => {
    await seedData();
    await WalletTable.create(defaultWallet3);
  });

  beforeAll(async () => {
    await migrate();
  });

  afterEach(async () => {
    await clearData();
  });

  afterAll(async () => {
    await teardown();
  });

  it('Successfully creates a LeaderboardPnl', async () => {
    await LeaderboardPnlTable.create(defaultLeaderboardPnlOneDay);
  });

  it('Successfully creates multiple LeaderboardPnls', async () => {
    await Promise.all([
      LeaderboardPnlTable.create(defaultLeaderboardPnlOneDay),
      LeaderboardPnlTable.create(defaultLeaderboardPnl2OneDay),
      LeaderboardPnlTable.create(defaultLeaderboardPnl1AllTime),
    ]);

    const leaderboardPnls: LeaderboardPnlFromDatabase[] = await LeaderboardPnlTable.findAll(
      {},
      [],
    );

    expect(leaderboardPnls.length).toEqual(3);
  });

  it('Successfully finds LeaderboardPnl with address and timespan', async () => {
    await Promise.all([
      LeaderboardPnlTable.create(defaultLeaderboardPnlOneDay),
      LeaderboardPnlTable.create(defaultLeaderboardPnl2OneDay),
      LeaderboardPnlTable.create(defaultLeaderboardPnl1AllTime),
    ]);

    const leaderboardPnl: LeaderboardPnlFromDatabase[] = await LeaderboardPnlTable.findAll(
      {
        address: [defaultLeaderboardPnlOneDay.address],
        timeSpan: [defaultLeaderboardPnlOneDay.timeSpan],
      },
      [],
    );

    expect(leaderboardPnl.length).toEqual(1);
    expect(leaderboardPnl[0]).toEqual(expect.objectContaining(defaultLeaderboardPnlOneDay));
  });

  it('Successfully upserts a LeaderboardPnl', async () => {
    await LeaderboardPnlTable.upsert(defaultLeaderboardPnlOneDay);

    await LeaderboardPnlTable.upsert(defaultLeaderboardPnlOneDayToUpsert);

    const leaderboardPnls: LeaderboardPnlFromDatabase[] = await LeaderboardPnlTable.findAll(
      {},
      [],
    );

    expect(leaderboardPnls.length).toEqual(1);
    expect(leaderboardPnls[0]).toEqual(
      expect.objectContaining(defaultLeaderboardPnlOneDayToUpsert));
  });

  it('Successfully bulk upserts LeaderboardPnls', async () => {
    await LeaderboardPnlTable.bulkUpsert(
      [defaultLeaderboardPnlOneDay, defaultLeaderboardPnl2OneDay]);

    const leaderboardPnls: LeaderboardPnlFromDatabase[] = await LeaderboardPnlTable.findAll(
      {},
      [],
    );

    expect(leaderboardPnls.length).toEqual(2);
    expect(leaderboardPnls[0]).toEqual(expect.objectContaining(defaultLeaderboardPnlOneDay));
    expect(leaderboardPnls[1]).toEqual(expect.objectContaining(defaultLeaderboardPnl2OneDay));
  });
});
