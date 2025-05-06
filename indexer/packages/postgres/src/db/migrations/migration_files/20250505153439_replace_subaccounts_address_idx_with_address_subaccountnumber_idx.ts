import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  await knex.raw(`
    DROP INDEX CONCURRENTLY IF EXISTS subaccounts_address_index;
  `);

  await knex.raw(`
    CREATE INDEX CONCURRENTLY IF NOT EXISTS subaccounts_address_subaccountnumber_index
      ON subaccounts(address, "subaccountNumber")
      INCLUDE (id);
  `);
}

export async function down(knex: Knex): Promise<void> {
  await knex.raw(`
    DROP INDEX CONCURRENTLY IF EXISTS subaccounts_address_subaccountnumber_index;
  `);

  await knex.raw(`
    CREATE INDEX CONCURRENTLY IF NOT EXISTS subaccounts_address_index
      ON subaccounts (address);
  `);
}

export const config = {
  transaction: false,
};
