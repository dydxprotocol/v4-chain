import * as Knex from 'knex';

// Disable transactions because CREATE INDEX CONCURRENTLY cannot run inside a transaction block
export const config = {
  transaction: false,
};

// Index for efficiently retrieving recent 1 and 4 hour candles (used by sparklines endpoint).
export async function up(knex: Knex): Promise<void> {
  await knex.raw(`
    CREATE INDEX CONCURRENTLY IF NOT EXISTS add_candles_sparklines_feb2025_idx
    ON candles (resolution, "startedAt")
    WHERE resolution IN ('1HOUR', '4HOURS')
      AND "startedAt" > '2025-02-01 00:00:00+00'::timestamp with time zone;
  `);
}

export async function down(knex: Knex): Promise<void> {
  await knex.raw(`
    DROP INDEX IF EXISTS add_candles_sparklines_feb2025_idx;
  `);
}
