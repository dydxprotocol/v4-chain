import * as Knex from 'knex';

/**
 * Drops the ("blockHeight", "transactionIndex") index on tendermint_events: no query reads it, and
 * on mainnet it took 93 GiB on every instance. Its only remaining role was the referencing side of
 * the "blockHeight" -> blocks foreign key, which Postgres uses only when a blocks row is deleted or
 * rekeyed; the indexer never does either, so rebuild the index before deleting blocks rows by hand.
 *
 * On mainnet this index was dropped by hand, so there this migration is a no-op.
 */
export async function up(knex: Knex): Promise<void> {
  await knex.raw(`
    DROP INDEX CONCURRENTLY IF EXISTS "tendermint_events_blockheight_transactionindex_index";
  `);
}

export async function down(knex: Knex): Promise<void> {
  await knex.raw(`
    CREATE INDEX CONCURRENTLY IF NOT EXISTS "tendermint_events_blockheight_transactionindex_index"
      ON "tendermint_events" ("blockHeight", "transactionIndex");
  `);
}

export const config = {
  transaction: false,
};
