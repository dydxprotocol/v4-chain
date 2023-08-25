import { knexPrimary } from '../../../src/helpers/knex';
import { seed } from '../../../src/db/seeds/01_genesis_seeds';
import { clearData, migrate, teardown } from '../../../src/helpers/db-helpers';
import { MarketColumns, MarketFromDatabase, Ordering } from '../../../src/types';
import * as MarketTable from '../../../src/stores/market-table';

describe('seed', () => {
  beforeAll(async () => {
    await migrate();
  });

  afterEach(async () => {
    await clearData();
  });

  afterAll(async () => {
    await clearData();
    await teardown();
  });

  it('seeds database', async () => {
    await seed(knexPrimary);

    const markets: MarketFromDatabase[] = await MarketTable.findAll(
      {},
      [],
      { readReplica: true, orderBy: [[MarketColumns.id, Ordering.ASC]] },
    );

    expect(markets).toHaveLength(35);
  });

  it('can be run multiple times', async () => {
    await seed(knexPrimary);
    await seed(knexPrimary);

    const markets: MarketFromDatabase[] = await MarketTable.findAll(
      {},
      [],
      { readReplica: true },
    );
    expect(markets).toHaveLength(35);
  });
});
