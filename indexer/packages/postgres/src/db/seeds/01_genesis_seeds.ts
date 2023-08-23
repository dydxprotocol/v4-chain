import * as Knex from 'knex';

import {
  getPerpetualMarketLiquidityTierUpdateSql,
  getSeedBlocksSql,
  getSeedLiquidityTiersSql,
  getSeedMarketsSql,
  getSeedPerpetualMarketsSql,
} from '../helpers';

// TODO(DEC-760): Seed `PerpetualMarkets`, `Assets` in unit tests.
export async function seed(knex: Knex): Promise<void> {
  await knex.raw(getSeedBlocksSql());
  await knex.raw(getSeedMarketsSql());
  await knex.raw(getSeedLiquidityTiersSql());
  await knex.raw(getSeedPerpetualMarketsSql());

  // Update perpetual_markets table to add liquidityTierId column
  const updateSql: string[] = getPerpetualMarketLiquidityTierUpdateSql();
  for (const sql of updateSql) {
    await knex.raw(sql);
  }
}
