import {
  dbHelpers,
  testMocks,
} from '@dydxprotocol-indexer/postgres';

describe('orderbook-mid-price-memory-cache', () => {

  beforeAll(async () => {
    await dbHelpers.migrate();
    await dbHelpers.clearData();
  });

  beforeEach(async () => {
    await testMocks.seedData();
  });

  afterEach(async () => {
    await dbHelpers.clearData();
  });

  describe('getOrderbookMidPrice', () => {
    it('should return the mid price for a given ticker', async () => {
    });
  });
});
