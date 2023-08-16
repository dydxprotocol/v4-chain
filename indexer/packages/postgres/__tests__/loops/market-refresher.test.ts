import { MarketCreateObject } from '../../src';
import { clearData, migrate, teardown } from '../../src/helpers/db-helpers';
import { getMarketFromId, updateMarkets } from '../../src/loops/market-refresher';
import { defaultMarket, defaultMarket2 } from '../helpers/constants';
import { seedData } from '../helpers/mock-generators';

describe('marketRefresher', () => {
  beforeAll(async () => {
    await migrate();
    await seedData();
    await updateMarkets();
  });

  afterAll(async () => {
    await clearData();
    await teardown();
  });

  describe('getMarketFromId', () => {
    it.each([
      [defaultMarket],
      [defaultMarket2],
    ])('successfully get an market from id', (market: MarketCreateObject) => {
      expect(getMarketFromId(market.id)).toEqual(expect.objectContaining(market));
    });

    it('returns undefined if market does not exist', () => {
      expect(() => getMarketFromId(1000)).toThrowError('Unable to find market with id: 1000');
    });
  });
});
