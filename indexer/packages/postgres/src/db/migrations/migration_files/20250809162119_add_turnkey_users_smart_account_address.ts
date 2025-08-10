import * as Knex from "knex";


export async function up(knex: Knex): Promise<void> {
  return knex.schema.alterTable('turnkey_users', (table) => {
    table.text('smart_account_address').nullable().unique();
    table.index(['smart_account_address'], 'idx_turnkey_users_smart_account_address');
  });
}


export async function down(knex: Knex): Promise<void> {
  return knex.schema.alterTable('turnkey_users', (table) => {
    table.dropIndex(['smart_account_address'], 'idx_turnkey_users_smart_account_address');
    table.dropColumn('smart_account_address');
  });
}

