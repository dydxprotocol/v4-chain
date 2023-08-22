import { knexPrimary } from '../../../src/helpers/knex';
import { seed } from '../../../src/db/seeds/01_genesis_seeds';
import { clearData, migrate, teardown } from '../../../src/helpers/db-helpers';
import { LiquidityTiersFromDatabase, MarketFromDatabase } from '../../../src/types';
import * as MarketTable from '../../../src/stores/market-table';
import * as LiquidityTiersTable from '../../../src/stores/liquidity-tiers-table';
import { expectLiquidityTier, expectMarketParamAndPrice } from '../helpers';
import {
  getLiquidityTiersFromGenesis,
  getMarketParamsFromGenesis,
  getMarketPricesFromGenesis,
} from '../../../src/db/helpers';

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

    const liquidityTiers: LiquidityTiersFromDatabase[] = await LiquidityTiersTable.findAll(
      {},
      [],
      { readReplica: true },
    );

    const markets: MarketFromDatabase[] = await MarketTable.findAll(
      {},
      [],
      { readReplica: true },
    );

    expect(markets).toHaveLength(33);
    markets.forEach((marketFromDb: MarketFromDatabase, index: number) => {
      expectMarketParamAndPrice(
        marketFromDb,
        getMarketParamsFromGenesis()[index],
        getMarketPricesFromGenesis()[index],
      );
    });

    expect(liquidityTiers).toHaveLength(3);
    liquidityTiers.forEach((liquidityTier: LiquidityTiersFromDatabase, index: number) => {
      expectLiquidityTier(
        liquidityTier,
        getLiquidityTiersFromGenesis()[index],
      );
    });
  });

  it('can be run multiple times', async () => {
    await seed(knexPrimary);
    await seed(knexPrimary);

    const markets: MarketFromDatabase[] = await MarketTable.findAll(
      {},
      [],
      { readReplica: true },
    );
    expect(markets).toHaveLength(33);
  });
});
