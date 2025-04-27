import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  await knex.raw(`
    DROP INDEX CONCURRENTLY IF EXISTS candles_ticker_resolution_index;
  `);

  await knex.raw(`
    CREATE INDEX CONCURRENTLY IF NOT EXISTS candles_ticker_resolution_startedat_index
      ON candles (ticker, resolution, "startedAt" DESC);
  `);
}

export async function down(knex: Knex): Promise<void> {
  await knex.raw(`
    DROP INDEX CONCURRENTLY IF EXISTS candles_ticker_resolution_startedat_index;
  `);

  await knex.raw(`
    CREATE INDEX CONCURRENTLY IF NOT EXISTS candles_ticker_resolution_index
      ON candles (ticker, resolution);
  `);
}

export const config = {
  transaction: false,
};
