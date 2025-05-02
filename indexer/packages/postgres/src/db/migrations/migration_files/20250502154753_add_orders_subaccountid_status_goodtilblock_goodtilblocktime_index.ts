import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  await knex.raw(`
    CREATE INDEX CONCURRENTLY IF NOT EXISTS orders_subaccountid_status_goodtilblock_goodtilblocktime_index
      ON orders (
        "subaccountId",
        status,
        "goodTilBlock" DESC,
        "goodTilBlockTime" DESC
      );
  `);
}

export async function down(knex: Knex): Promise<void> {
  await knex.raw(`
    DROP INDEX CONCURRENTLY IF EXISTS orders_subaccountid_status_goodtilblock_goodtilblocktime_index;
  `);
}

export const config = {
  transaction: false,
};
