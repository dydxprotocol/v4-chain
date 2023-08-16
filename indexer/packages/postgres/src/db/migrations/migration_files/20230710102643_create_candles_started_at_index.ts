import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  // eslint-disable-next-line @typescript-eslint/quotes
  await knex.raw(`
    CREATE INDEX CONCURRENTLY IF NOT EXISTS "candles_started_at_index" ON "candles" ("startedAt");
  `);
}

export async function down(knex: Knex): Promise<void> {
  await knex.raw(`
    DROP INDEX CONCURRENTLY IF EXISTS "candles_started_at_index";
  `);
}

export const config = {
  transaction: false,
};
