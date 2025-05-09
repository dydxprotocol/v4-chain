import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  await knex.raw(`
    DROP INDEX CONCURRENTLY IF EXISTS orders_clobpairid_side_price_index;
  `);
}

export async function down(knex: Knex): Promise<void> {
  await knex.raw(`
    CREATE INDEX CONCURRENTLY IF NOT EXISTS orders_clobpairid_side_price_index
      ON orders("clobPairId", side, price);
  `);
}

export const config = {
  transaction: false,
};
