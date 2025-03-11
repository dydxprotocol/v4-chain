import * as Knex from 'knex';

// Disable transactions because CREATE INDEX CONCURRENTLY cannot run inside a transaction block
export const config = {
  transaction: false,
};

export async function up(knex: Knex): Promise<void> {
  await knex.raw(`
    CREATE INDEX CONCURRENTLY IF NOT EXISTS candles_resolution_started_at_1_4_hour_feb2025_idx
    ON candles (resolution, "startedAt")
    WHERE resolution IN ('1HOUR', '4HOURS')
      AND "startedAt" > '2025-02-01 00:00:00+00'::timestamp with time zone;
  `);
}

export async function down(knex: Knex): Promise<void> {
  await knex.raw(`
    DROP INDEX IF EXISTS candles_resolution_started_at_1_4_hour_feb2025_idx;
  `);
}
