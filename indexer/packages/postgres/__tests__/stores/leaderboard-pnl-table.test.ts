import { LeaderboardPNLFromDatabase } from '../../src/types';
import * as LeaderboardPNLTable from '../../src/stores/leaderboard-pnl-table';
import { clearData, migrate, teardown } from '../../src/helpers/db-helpers';
import {
  defaultLeaderboardPnl2OneDay,
  defaultLeaderboardPnlOneDay,
  defaultLeaderboardPnl1AllTime,
  defaultLeaderboardPnlOneDayToUpsert,
} from '../helpers/constants';
import { seedData } from '../helpers/mock-generators';

describe('LeaderboardPNL store', () => {
  beforeEach(async () => {
    await seedData();
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

  it('Successfully creates a LeaderboardPNL', async () => {
    await LeaderboardPNLTable.create(defaultLeaderboardPnlOneDay);
  });

  it('Successfully creates multiple LeaderboardPNLs', async () => {
    await Promise.all([
      LeaderboardPNLTable.create(defaultLeaderboardPnlOneDay),
      LeaderboardPNLTable.create(defaultLeaderboardPnl2OneDay),
      LeaderboardPNLTable.create(defaultLeaderboardPnl1AllTime),
    ]);

    const leaderboardPNLs: LeaderboardPNLFromDatabase[] = await LeaderboardPNLTable.findAll(
      {},
      [],
      { readReplica: true },
    );

    expect(leaderboardPNLs.length).toEqual(3);
  });

  it('Successfully finds LeaderboardPNL with subaccountId and timespan', async () => {
    await Promise.all([
      LeaderboardPNLTable.create(defaultLeaderboardPnlOneDay),
      LeaderboardPNLTable.create(defaultLeaderboardPnl2OneDay),
      LeaderboardPNLTable.create(defaultLeaderboardPnl1AllTime),
    ]);

    const leaderboardPNL: LeaderboardPNLFromDatabase[] = await LeaderboardPNLTable.findAll(
      {
        subaccountId: [defaultLeaderboardPnlOneDay.subaccountId],
        timeSpan: [defaultLeaderboardPnlOneDay.timeSpan],
      },
      [],
      { readReplica: true },
    );

    expect(leaderboardPNL.length).toEqual(1);
    expect(leaderboardPNL[0]).toEqual(expect.objectContaining(defaultLeaderboardPnlOneDay));
  });

  it('Successfully upserts a LeaderboardPNL', async () => {
    await LeaderboardPNLTable.upsert(defaultLeaderboardPnlOneDay);

    await LeaderboardPNLTable.upsert(defaultLeaderboardPnlOneDayToUpsert);

    const leaderboardPNLs: LeaderboardPNLFromDatabase[] = await LeaderboardPNLTable.findAll(
      {},
      [],
      { readReplica: true },
    );

    expect(leaderboardPNLs.length).toEqual(1);
    expect(leaderboardPNLs[0]).toEqual(
      expect.objectContaining(defaultLeaderboardPnlOneDayToUpsert));
  });

  it('Successfully bulk upserts LeaderboardPNLs', async () => {
    await LeaderboardPNLTable.bulkUpsert(
      [defaultLeaderboardPnlOneDay, defaultLeaderboardPnl2OneDay]);

    const leaderboardPNLs: LeaderboardPNLFromDatabase[] = await LeaderboardPNLTable.findAll(
      {},
      [],
      { readReplica: true },
    );

    expect(leaderboardPNLs.length).toEqual(2);
    expect(leaderboardPNLs[0]).toEqual(expect.objectContaining(defaultLeaderboardPnlOneDay));
    expect(leaderboardPNLs[1]).toEqual(expect.objectContaining(defaultLeaderboardPnl2OneDay));
  });
});
