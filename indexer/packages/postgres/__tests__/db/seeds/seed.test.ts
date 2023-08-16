import { knexPrimary } from '../../../src/helpers/knex';
import { seed } from '../../../src/db/seeds/01_genesis_seeds';
import { clearData, migrate, teardown } from '../../../src/helpers/db-helpers';
import {
  AssetFromDatabase,
  LiquidityTiersFromDatabase,
  MarketFromDatabase,
  PerpetualMarketFromDatabase,
} from '../../../src/types';
import * as AssetTable from '../../../src/stores/asset-table';
import * as PerpetualMarketTable from '../../../src/stores/perpetual-market-table';
import * as MarketTable from '../../../src/stores/market-table';
import * as LiquidityTiersTable from '../../../src/stores/liquidity-tiers-table';
import {
  expectAsset, expectLiquidityTier, expectMarketParamAndPrice, expectPerpetualMarket,
} from '../helpers';
import {
  getAssetsFromGenesis,
  getClobPairsFromGenesis,
  getLiquidityTiersFromGenesis,
  getMarketParamsFromGenesis,
  getMarketPricesFromGenesis,
  getPerpetualsFromGenesis,
} from '../../../src/db/helpers';
import {
  defaultLiquidityTier,
  defaultLiquidityTier2,
  defaultMarket,
  defaultMarket2,
  defaultPerpetualMarket,
} from '../../helpers/constants';

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

    const perpetualMarkets: PerpetualMarketFromDatabase[] = await PerpetualMarketTable.findAll(
      {},
      [],
      { readReplica: true },
    );

    const assets: AssetFromDatabase[] = await AssetTable.findAll(
      {},
      [],
      { readReplica: true },
    );

    const liquidityTiers: LiquidityTiersFromDatabase[] = await LiquidityTiersTable.findAll(
      {},
      [],
      { readReplica: true },
    );

    expect(perpetualMarkets).toHaveLength(33);
    perpetualMarkets.forEach((perpetualMarket: PerpetualMarketFromDatabase, index: number) => {
      expectPerpetualMarket(
        perpetualMarket,
        getPerpetualsFromGenesis()[index],
        getClobPairsFromGenesis()[index],
      );
    });

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

    expect(assets).toHaveLength(1);
    expectAsset(assets[0],
      getAssetsFromGenesis()[0]);
  });

  it('seed should update the liquidityTierId for existing Perpetual Markets', async () => {
    await Promise.all([
      MarketTable.create(defaultMarket),
      MarketTable.create(defaultMarket2),
    ]);
    await Promise.all([
      LiquidityTiersTable.create(defaultLiquidityTier),
      LiquidityTiersTable.create(defaultLiquidityTier2),
    ]);
    await PerpetualMarketTable.create({
      ...defaultPerpetualMarket,
      liquidityTierId: 1,
    });

    await seed(knexPrimary);

    const perpetualMarkets: PerpetualMarketFromDatabase[] = await PerpetualMarketTable.findAll(
      {},
      [],
      { readReplica: true },
    );
    expect(perpetualMarkets[0].liquidityTierId).toEqual(0);
  });

  it('can be run multiple times', async () => {
    await seed(knexPrimary);
    await seed(knexPrimary);

    const perpetualMarkets: PerpetualMarketFromDatabase[] = await PerpetualMarketTable.findAll(
      {},
      [],
      { readReplica: true },
    );

    expect(perpetualMarkets).toHaveLength(33);

    const assets: AssetFromDatabase[] = await AssetTable.findAll(
      {},
      [],
      { readReplica: true },
    );
    expect(assets).toHaveLength(1);

    const markets: MarketFromDatabase[] = await MarketTable.findAll(
      {},
      [],
      { readReplica: true },
    );
    expect(markets).toHaveLength(33);
  });
});
