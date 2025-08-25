import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  return knex.schema.alterTable('turnkey_users', (table) => {
    table.text('smart_account_address').nullable().unique();
  });
}

export async function down(knex: Knex): Promise<void> {
  return knex.schema.alterTable('turnkey_users', (table) => {
    table.dropColumn('smart_account_address');
  });
}
