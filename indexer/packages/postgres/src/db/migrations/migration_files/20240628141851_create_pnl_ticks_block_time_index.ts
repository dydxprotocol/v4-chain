import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  await knex.raw(`
    CREATE INDEX CONCURRENTLY IF NOT EXISTS "pnl_ticks_blocktime_index" ON "pnl_ticks" ("blockTime");
  `);
}

export async function down(knex: Knex): Promise<void> {
  await knex.raw(`
    DROP INDEX CONCURRENTLY IF EXISTS "pnl_ticks_blocktime_index";
  `);
}

export const config = {
  transaction: false,
};
