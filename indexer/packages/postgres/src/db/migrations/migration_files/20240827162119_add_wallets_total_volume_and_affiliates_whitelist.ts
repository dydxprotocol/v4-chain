import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  return knex
    .schema
    .alterTable('wallets', (table) => {
      table.decimal('totalVolume', null).defaultTo(0).notNullable();
      table.boolean('isWhitelistAffiliate').defaultTo(false).notNullable();
    });
}

export async function down(knex: Knex): Promise<void> {
  return knex
    .schema
    .alterTable('wallets', (table) => {
      table.dropColumn('totalVolume');
      table.dropColumn('isWhitelistAffiliate');
    });
}
