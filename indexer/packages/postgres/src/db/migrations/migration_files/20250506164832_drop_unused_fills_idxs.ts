import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  await knex.raw(`
    DROP INDEX CONCURRENTLY IF EXISTS fills_clobPairId_createdAtHeight_partial;
  `);

  await knex.raw(`
    DROP INDEX CONCURRENTLY IF EXISTS fills_clobpairid_createdat_index;
  `);

  await knex.raw(`
    DROP INDEX CONCURRENTLY IF EXISTS fills_orderid_index;
  `);

  await knex.raw(`
    DROP INDEX CONCURRENTLY IF EXISTS fills_subaccountid_createdat_index;
  `);

  await knex.raw(`
    DROP INDEX CONCURRENTLY IF EXISTS fills_subaccountid_index;
  `);

  await knex.raw(`
    DROP INDEX CONCURRENTLY IF EXISTS fills_type_index;
  `);
}

export async function down(knex: Knex): Promise<void> {
  await knex.raw(`
    CREATE INDEX CONCURRENTLY IF NOT EXISTS fills_clobPairId_createdAtHeight_partial
      ON fills("clobPairId", "createdAtHeight")
      WHERE liquidity = 'TAKER';
  `);

  await knex.raw(`
    CREATE INDEX CONCURRENTLY IF NOT EXISTS fills_clobpairid_createdat_index
      ON fills("clobPairId", "createdAt");
  `);

  await knex.raw(`
    CREATE INDEX CONCURRENTLY IF NOT EXISTS fills_orderid_index
      ON fills("orderId");
  `);

  await knex.raw(`
    CREATE INDEX CONCURRENTLY IF NOT EXISTS fills_subaccountid_createdat_index
      ON fills("subaccountId", "createdAt");
  `);

  await knex.raw(`
    CREATE INDEX CONCURRENTLY IF NOT EXISTS fills_subaccountid_index
      ON fills("subaccountId");
  `);

  await knex.raw(`
    CREATE INDEX CONCURRENTLY IF NOT EXISTS fills_type_index
      ON fills(type);
  `);
}

export const config = {
  transaction: false,
};
