import * as Knex from 'knex';

import { getSeedBlocksSql, getSeedLiquidityTiersSql, getSeedMarketsSql } from '../helpers';

// TODO(DEC-760): Seed `PerpetualMarkets`, `Assets` in unit tests.
export async function seed(knex: Knex): Promise<void> {
  await knex.raw(getSeedBlocksSql());
  await knex.raw(getSeedMarketsSql());
  await knex.raw(getSeedLiquidityTiersSql());
}
