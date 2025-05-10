import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  // oracle_prices has other indices on ("marketId", ...) already.
  await knex.raw(`
    DROP INDEX CONCURRENTLY IF EXISTS oracle_prices_marketid_index;
  `);

  // tendermint_events has an index on ("blockHeight", "transactionIndex") already.
  await knex.raw(`
    DROP INDEX CONCURRENTLY IF EXISTS tendermint_events_blockheight_index;
  `);

  // transactions has an index on ("blockHeight", "transactionIndex") already.
  await knex.raw(`
    DROP INDEX CONCURRENTLY IF EXISTS transactions_blockheight_index;
  `);
}

export async function down(knex: Knex): Promise<void> {
  await knex.raw(`
    CREATE INDEX CONCURRENTLY IF NOT EXISTS oracle_prices_marketid_index
      ON oracle_prices("marketId");
  `);

  await knex.raw(`
    CREATE INDEX CONCURRENTLY IF NOT EXISTS tendermint_events_blockheight_index
      ON tendermint_events("blockHeight");
  `);

  await knex.raw(`
    CREATE INDEX CONCURRENTLY IF NOT EXISTS transactions_blockheight_index
      ON transactions("blockHeight");
  `);
}

export const config = {
  transaction: false,
};
