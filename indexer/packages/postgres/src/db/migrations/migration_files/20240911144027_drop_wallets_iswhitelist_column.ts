import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  return knex
    .schema
    .alterTable('wallets', (table) => {
      table.dropColumn('isWhitelistAffiliate');
    });
}

export async function down(knex: Knex): Promise<void> {
  return knex
    .schema
    .alterTable('wallets', (table) => {
      table.boolean('isWhitelistAffiliate').defaultTo(false).notNullable();
    });
}
