import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  return knex.schema.alterTable('turnkey_users', (table) => {
    // The Injective account the user's deposits are bridged to. The user chooses it, so unlike
    // the other addresses it is not unique: a user naming an address must not be able to keep
    // another user from naming the same one.
    table.text('injective_address').nullable();
  });
}

export async function down(knex: Knex): Promise<void> {
  return knex.schema.alterTable('turnkey_users', (table) => {
    table.dropColumn('injective_address');
  });
}
