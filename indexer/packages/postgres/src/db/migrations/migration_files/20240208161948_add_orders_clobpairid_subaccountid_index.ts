import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  // eslint-disable-next-line @typescript-eslint/quotes
  await knex.raw(`
    CREATE INDEX CONCURRENTLY IF NOT EXISTS orders_clobPairId_subaccountId_index ON orders("clobPairId", "subaccountId");
  `);
}

export async function down(knex: Knex): Promise<void> {
  await knex.raw(`
    DROP INDEX CONCURRENTLY IF EXISTS "orders_clobPairId_subaccountId_index";
  `);
}

export const config = {
  transaction: false,
};
