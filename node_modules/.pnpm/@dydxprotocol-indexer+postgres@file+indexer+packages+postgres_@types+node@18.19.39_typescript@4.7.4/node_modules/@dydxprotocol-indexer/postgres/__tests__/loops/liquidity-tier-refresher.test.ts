import { LiquidityTiersCreateObject } from '../../src';
import { clearData, migrate, teardown } from '../../src/helpers/db-helpers';
import { getLiquidityTierFromId, updateLiquidityTiers } from '../../src/loops/liquidity-tier-refresher';
import { defaultLiquidityTier, defaultLiquidityTier2 } from '../helpers/constants';
import { seedData } from '../helpers/mock-generators';

describe('liquidityTierRefresher', () => {
  beforeAll(async () => {
    await migrate();
    await seedData();
    await updateLiquidityTiers();
  });

  afterAll(async () => {
    await clearData();
    await teardown();
  });

  describe('getLiquidityTierFromId', () => {
    it.each([
      [defaultLiquidityTier],
      [defaultLiquidityTier2],
    ])('successfully get an liquidityTier from id',
      (liquidityTier: LiquidityTiersCreateObject) => {
        expect(getLiquidityTierFromId(liquidityTier.id)).toEqual(
          expect.objectContaining(liquidityTier),
        );
      });

    it('throws error if liquidityTier does not exist', () => {
      expect(() => getLiquidityTierFromId(50)).toThrowError('Unable to find liquidity tier with id: 50');
    });
  });
});
