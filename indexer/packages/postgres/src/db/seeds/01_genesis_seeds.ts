import * as Knex from 'knex';

import {
  getSeedAssetPositionsSql,
  getSeedAssetsSql,
  getSeedPerpetualMarketsSql,
  getSeedSubaccountsSql,
  getSeedMarketsSql,
  getSeedBlocksSql,
  getSeedLiquidityTiersSql,
  getPerpetualMarketLiquidityTierUpdateSql,
} from '../helpers';

// TODO(DEC-760): Seed `PerpetualMarkets`, `Assets` in unit tests.
export async function seed(knex: Knex): Promise<void> {
  await knex.raw(getSeedBlocksSql());
  await knex.raw(getSeedMarketsSql());
  await knex.raw(getSeedLiquidityTiersSql());
  await knex.raw(getSeedPerpetualMarketsSql());
  await knex.raw(getSeedAssetsSql());
  await knex.raw(getSeedSubaccountsSql());
  // AssetPosition seeding needs to be run after subaccounts/assets due to foreign key
  // dependencies.
  await knex.raw(await getSeedAssetPositionsSql());

  // Update perpetual_markets table to add liquidityTierId column
  const updateSql: string[] = getPerpetualMarketLiquidityTierUpdateSql();
  for (const sql of updateSql) {
    await knex.raw(sql);
  }
}
