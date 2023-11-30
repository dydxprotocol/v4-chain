import { TradingRewardFromDatabase } from '../../src/types';
import { clearData, migrate, teardown } from '../../src/helpers/db-helpers';
import { defaultTradingReward, defaultWallet } from '../helpers/constants';
import * as TradingRewardTable from '../../src/stores/trading-reward-table';
import { WalletTable } from '../../src';

describe('TradingReward store', () => {
  beforeAll(async () => {
    await migrate();
  });

  beforeEach(async () => {
    await WalletTable.create(defaultWallet);
  });

  afterEach(async () => {
    await clearData();
  });

  afterAll(async () => {
    await teardown();
  });

  it('Successfully creates a TradingReward', async () => {
    await TradingRewardTable.create(defaultTradingReward);
  });

  it('Successfully finds all TradingRewards', async () => {
    await Promise.all([
      TradingRewardTable.create(defaultTradingReward),
      TradingRewardTable.create({
        ...defaultTradingReward,
        blockHeight: '20',
      }),
    ]);

    const tradingRewards: TradingRewardFromDatabase[] = await TradingRewardTable.findAll(
      {},
      [],
      { readReplica: true },
    );

    expect(tradingRewards.length).toEqual(2);
    expect(tradingRewards[0]).toEqual(expect.objectContaining({
      ...defaultTradingReward,
      blockHeight: '20',
    }));
    expect(tradingRewards[1]).toEqual(expect.objectContaining(defaultTradingReward));
  });

  it('Successfully finds a TradingReward', async () => {
    await TradingRewardTable.create(defaultTradingReward);

    const tradingReward: TradingRewardFromDatabase | undefined = await TradingRewardTable.findById(
      TradingRewardTable.uuid(defaultTradingReward.address, defaultTradingReward.blockHeight),
    );

    expect(tradingReward).toEqual(expect.objectContaining(defaultTradingReward));
  });
});
