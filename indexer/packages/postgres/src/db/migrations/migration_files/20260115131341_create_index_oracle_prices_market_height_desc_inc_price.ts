import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  await knex.raw(`
    CREATE INDEX oracle_prices_market_height_desc_inc_price
    ON "oracle_prices" ("marketId", "effectiveAtHeight" DESC)
    INCLUDE ("price");
  `);
}

export async function down(knex: Knex): Promise<void> {
  await knex.raw(`
    DROP INDEX IF EXISTS "oracle_prices_market_height_desc_inc_price";
  `);
}

export const config = {
  transaction: false,
};
