import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  await knex.raw(`
    CREATE INDEX CONCURRENTLY IF NOT EXISTS "oracle_prices_marketid_effectiveatheight_index" ON "oracle_prices" ("marketId", "effectiveAtHeight");
  `);
}

export async function down(knex: Knex): Promise<void> {
  await knex.raw(`
    DROP INDEX CONCURRENTLY IF EXISTS "oracle_prices_marketid_effectiveatheight_index";
  `);
}

export const config = {
  transaction: false,
};
