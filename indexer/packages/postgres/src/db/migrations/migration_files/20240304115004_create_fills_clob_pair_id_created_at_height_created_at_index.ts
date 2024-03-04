import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  // Partial index only when `liquidity` is 'TAKER' as this index is meant to speed up the query to
  // fetch all trades from the database, and we arbitrarily only use 'TAKER' fills for trades to
  // avoid double counting.
  // eslint-disable-next-line @typescript-eslint/quotes
  await knex.raw(`
    CREATE INDEX CONCURRENTLY IF NOT EXISTS "fills_clobPairId_createdAtHeight_createdAt_partial" ON "fills" ("clobPairId", "createdAtHeight", "createdAt") WHERE "liquidity" = 'TAKER';
  `);
}

export async function down(knex: Knex): Promise<void> {
  await knex.raw(`
    DROP INDEX CONCURRENTLY IF EXISTS "fills_clobPairId_createdAtHeight_createdAt_partial";
  `);
}

export const config = {
  transaction: false,
};
