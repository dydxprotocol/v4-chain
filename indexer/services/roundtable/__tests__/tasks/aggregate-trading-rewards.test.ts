import {
  dbHelpers, testMocks,
} from '@dydxprotocol-indexer/postgres';

describe('aggregate-trading-rewards', () => {
  beforeAll(async () => {
    await dbHelpers.migrate();
    await dbHelpers.clearData();
  });

  beforeEach(async () => {
    await testMocks.seedData();
  });

  afterEach(async () => {
    await dbHelpers.clearData();
    jest.resetAllMocks();
  });

  afterAll(async () => {
    await dbHelpers.teardown();
  });
});
