import { AssetPositionFromDatabase } from '../../src/types';
import * as AssetPositionTable from '../../src/stores/asset-position-table';
import { findUsdcPositionForSubaccounts } from '../../src/stores/asset-position-table';

import { clearData, migrate, teardown } from '../../src/helpers/db-helpers';
import { seedData } from '../helpers/mock-generators';
import { randomUUID } from 'crypto';
import {
  defaultAsset,
  defaultAssetPosition,
  defaultAssetPosition2,
  defaultAssetPositionId,
  defaultSubaccountId,
  defaultSubaccountId2,
} from '../helpers/constants';
import Big from 'big.js';

describe('Asset position store', () => {
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

  it('Successfully finds all asset positions', async () => {
    await Promise.all([
      AssetPositionTable.upsert(defaultAssetPosition),
      AssetPositionTable.upsert(defaultAssetPosition2),
    ]);

    const assetPositions: AssetPositionFromDatabase[] = await AssetPositionTable.findAll(
      {},
      [],
      { readReplica: true },
    );

    expect(assetPositions.length).toEqual(2);
    expect(assetPositions[0]).toEqual(expect.objectContaining(defaultAssetPosition));
    expect(assetPositions[1]).toEqual(expect.objectContaining(defaultAssetPosition2));
  });

  it('Successfully finds Asset Position with subaccount id', async () => {
    await Promise.all([
      AssetPositionTable.upsert(defaultAssetPosition),
      AssetPositionTable.upsert(defaultAssetPosition2),
    ]);

    const assetPositions: AssetPositionFromDatabase[] = await AssetPositionTable.findAll(
      {
        subaccountId: [defaultSubaccountId],
      },
      [],
      { readReplica: true },
    );

    expect(assetPositions.length).toEqual(1);
    expect(assetPositions[0]).toEqual(expect.objectContaining(defaultAssetPosition));
  });

  it('Successfully finds an asset position with id', async () => {
    const assetPositionFromDatabase: AssetPositionFromDatabase | undefined = await
    AssetPositionTable.upsert(defaultAssetPosition);

    const assetPosition: AssetPositionFromDatabase | undefined = await AssetPositionTable.findById(
      assetPositionFromDatabase.id,
    );

    expect(assetPosition).toEqual(expect.objectContaining(defaultAssetPosition));
  });

  it('Successfully finds USDC positions for subaccountIds', async () => {
    await Promise.all([
      AssetPositionTable.upsert(defaultAssetPosition),
      AssetPositionTable.upsert({
        ...defaultAssetPosition2,
        assetId: defaultAsset.id,
      }),
    ]);

    const assetPositions: { [subaccountId: string]: Big } = await findUsdcPositionForSubaccounts([
      defaultSubaccountId,
      defaultSubaccountId2,
    ]);

    expect(assetPositions).toEqual(expect.objectContaining({
      [defaultSubaccountId]: Big(defaultAssetPosition.size),
      [defaultSubaccountId2]: Big(0).minus(defaultAssetPosition2.size),
    }));
  });

  it('Unable finds an asset position', async () => {
    const assetPosition: AssetPositionFromDatabase | undefined = await AssetPositionTable.findById(
      randomUUID(),
    );
    expect(assetPosition).toEqual(undefined);
  });

  it('Successfully upserts an asset position', async () => {
    await AssetPositionTable.upsert(defaultAssetPosition);
    await expect(
      AssetPositionTable.findById(defaultAssetPositionId),
    ).resolves.toEqual(expect.objectContaining(defaultAssetPosition));
  });

  it('Successfully upserts a preexisting asset position', async () => {
    await AssetPositionTable.upsert(defaultAssetPosition);
    await AssetPositionTable.upsert({
      ...defaultAssetPosition,
      size: '20000',
    });
    await expect(
      AssetPositionTable.findById(defaultAssetPositionId),
    ).resolves.toEqual(expect.objectContaining({
      ...defaultAssetPosition,
      size: '20000',
    }));
  });
});
