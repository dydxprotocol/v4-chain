import { LiquidityTiersFromDatabase } from '../../src/types';
import * as LiquidityTierTable from '../../src/stores/liquidity-tiers-table';
import {
  clearData,
  migrate,
  teardown,
} from '../../src/helpers/db-helpers';
import { UniqueViolationError } from 'objection';
import { defaultLiquidityTier, defaultLiquidityTier2 } from '../helpers/constants';

describe('LiquidityTier store', () => {
  beforeAll(async () => {
    await migrate();
  });

  afterEach(async () => {
    await clearData();
    jest.resetAllMocks();
  });

  afterAll(async () => {
    await teardown();
    jest.clearAllMocks();
  });

  it('Successfully creates a liquidity tier', async () => {
    await LiquidityTierTable.create(defaultLiquidityTier);
  });

  it('Fails to create second liquidity tier with the same ID', async () => {
    try {
      await Promise.all([
        LiquidityTierTable.create(defaultLiquidityTier),
        LiquidityTierTable.create(defaultLiquidityTier),
      ]);
    } catch (e) {
      expect(e).toBeInstanceOf(UniqueViolationError);
    }
  });

  it('Successfully finds all liquidity tiers', async () => {
    await Promise.all([
      LiquidityTierTable.create(defaultLiquidityTier),
      LiquidityTierTable.create(defaultLiquidityTier2),
    ]);

    const liquidityTiers: LiquidityTiersFromDatabase[] = await LiquidityTierTable.findAll(
      {},
      [],
      { readReplica: true },
    );

    expect(liquidityTiers.length).toEqual(2);
    expect(liquidityTiers[0]).toEqual(expect.objectContaining(defaultLiquidityTier));
    expect(liquidityTiers[1]).toEqual(expect.objectContaining(defaultLiquidityTier2));
  });

  it('Successfully finds a liquidity tier', async () => {
    await LiquidityTierTable.create(defaultLiquidityTier);

    const liquidityTier: LiquidityTiersFromDatabase | undefined = await LiquidityTierTable.findById(
      defaultLiquidityTier.id,
    );

    expect(liquidityTier).toEqual(expect.objectContaining(defaultLiquidityTier));
  });

  it('Unable to find a liquidity tier', async () => {
    const liquidityTier: LiquidityTiersFromDatabase | undefined = await LiquidityTierTable.findById(
      defaultLiquidityTier.id,
    );
    expect(liquidityTier).toEqual(undefined);
  });

  it('Successfully updates a liquidity tier', async () => {
    await LiquidityTierTable.create(defaultLiquidityTier);

    const liquidityTier: LiquidityTiersFromDatabase | undefined = await LiquidityTierTable.update({
      id: defaultLiquidityTier.id,
      initialMarginPpm: '1000',
    });

    expect(liquidityTier).toEqual(expect.objectContaining({
      ...defaultLiquidityTier,
      initialMarginPpm: '1000',
    }));
  });

  it('Successfully upserts an existing liquidity tier', async () => {
    await LiquidityTierTable.create(defaultLiquidityTier);

    const liquidityTier: LiquidityTiersFromDatabase | undefined = await LiquidityTierTable.upsert({
      ...defaultLiquidityTier,
      initialMarginPpm: '1000',
    });

    expect(liquidityTier).toEqual(expect.objectContaining({
      ...defaultLiquidityTier,
      initialMarginPpm: '1000',
    }));
  });

  it('Successfully upserts a liquidity tier', async () => {
    const liquidityTier: LiquidityTiersFromDatabase | undefined = await
    LiquidityTierTable.upsert(defaultLiquidityTier);

    expect(liquidityTier).toEqual(expect.objectContaining(defaultLiquidityTier));
  });
});
