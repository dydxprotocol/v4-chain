import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  return knex
    .schema
    .alterTable('wallets', (table) => {
      table.decimal('totalTradingRewards', null).defaultTo('0').notNullable();
    });
}

export async function down(knex: Knex): Promise<void> {
  return knex
    .schema
    .alterTable('wallets', (table) => {
      table.dropColumn('totalTradingRewards');
    });
}
