import { AssetCreateObject } from '../../src';
import { clearData, migrate, teardown } from '../../src/helpers/db-helpers';
import { getAssetFromId, updateAssets } from '../../src/loops/asset-refresher';
import { defaultAsset, defaultAsset2, defaultAsset3 } from '../helpers/constants';
import { seedData } from '../helpers/mock-generators';

describe('assetRefresher', () => {
  beforeAll(async () => {
    await migrate();
    await seedData();
    await updateAssets();
  });

  afterAll(async () => {
    await clearData();
    await teardown();
  });

  describe('getAssetFromId', () => {
    it.each([
      [defaultAsset],
      [defaultAsset2],
      [defaultAsset3],
    ])('successfully get an asset from id', (asset: AssetCreateObject) => {
      expect(getAssetFromId(asset.id)).toEqual(expect.objectContaining(asset));
    });

    it('returns undefined if asset does not exist', () => {
      expect(() => getAssetFromId('invalid')).toThrowError('Unable to find asset with assetId: invalid');
    });
  });
});
