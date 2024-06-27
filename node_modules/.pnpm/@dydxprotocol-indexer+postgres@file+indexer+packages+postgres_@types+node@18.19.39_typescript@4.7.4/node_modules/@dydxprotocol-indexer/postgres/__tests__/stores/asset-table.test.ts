import { AssetFromDatabase } from '../../src/types';
import * as AssetTable from '../../src/stores/asset-table';
import {
  clearData,
  migrate,
  teardown,
} from '../../src/helpers/db-helpers';
import { UniqueViolationError } from 'objection';
import {
  defaultAsset,
  defaultAsset2,
  defaultMarket,
  defaultMarket2,
} from '../helpers/constants';
import * as MarketTable from '../../src/stores/market-table';

describe('Asset store', () => {
  beforeAll(async () => {
    await migrate();
  });

  beforeEach(async () => {
    await Promise.all([
      MarketTable.create(defaultMarket),
      MarketTable.create(defaultMarket2),
    ]);
  });

  afterEach(async () => {
    await clearData();
    jest.resetAllMocks();
  });

  afterAll(async () => {
    await teardown();
    jest.clearAllMocks();
  });

  it('Successfully creates a Asset', async () => {
    await AssetTable.create(defaultAsset);
  });

  it('Fails to create second asset with the same ID', async () => {
    try {
      await Promise.all([
        AssetTable.create(defaultAsset),
        AssetTable.create(defaultAsset),
      ]);
    } catch (e) {
      expect(e).toBeInstanceOf(UniqueViolationError);
    }
  });

  it('Successfully finds all Assets', async () => {
    await Promise.all([
      AssetTable.create(defaultAsset),
      AssetTable.create(defaultAsset2),
    ]);

    const assets: AssetFromDatabase[] = await AssetTable.findAll(
      {},
      [],
      { readReplica: true },
    );

    expect(assets.length).toEqual(2);
    expect(assets[0]).toEqual(expect.objectContaining(defaultAsset2));
    expect(assets[1]).toEqual(expect.objectContaining(defaultAsset));
  });

  it('Successfully finds Asset with symbol', async () => {
    await Promise.all([
      AssetTable.create(defaultAsset),
      AssetTable.create(defaultAsset2),
    ]);

    const assets: AssetFromDatabase[] = await AssetTable.findAll(
      {
        symbol: defaultAsset.symbol,
      },
      [],
      { readReplica: true },
    );

    expect(assets.length).toEqual(1);
    expect(assets[0]).toEqual(expect.objectContaining(defaultAsset));
  });

  it('Successfully finds a Asset', async () => {
    await AssetTable.create(defaultAsset);

    const asset: AssetFromDatabase | undefined = await AssetTable.findById(
      defaultAsset.id,
    );

    expect(asset).toEqual(expect.objectContaining(defaultAsset));
  });

  it('Unable finds a Asset', async () => {
    const asset: AssetFromDatabase | undefined = await AssetTable.findById(
      defaultAsset.id,
    );
    expect(asset).toEqual(undefined);
  });

  it('Successfully updates a asset', async () => {
    await AssetTable.create(defaultAsset);

    const asset: AssetFromDatabase | undefined = await AssetTable.update({
      id: defaultAsset.id,
      symbol: 'ETH',
    });

    expect(asset).toEqual(expect.objectContaining({
      ...defaultAsset,
      symbol: 'ETH',
    }));
  });

  it('Fails to update asset to have same symbol as existing asset', async () => {
    try {
      await AssetTable.create(defaultAsset);
      await AssetTable.create({
        id: '1', symbol: 'ETH', atomicResolution: -10, hasMarket: true,
      });
      await AssetTable.update({
        id: defaultAsset.id,
        symbol: 'ETH',
      });
    } catch (e) {
      expect(e).toBeInstanceOf(UniqueViolationError);
    }
  });
});
