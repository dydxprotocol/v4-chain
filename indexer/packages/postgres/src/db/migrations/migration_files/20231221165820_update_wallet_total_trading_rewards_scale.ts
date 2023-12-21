import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  return knex
    .schema
    .alterTable('wallets', (table) => {
      table.decimal('totalTradingRewards', 18, 18).notNullable().alter();
    });
}

export async function down(knex: Knex): Promise<void> {
  return knex
    .schema
    .alterTable('wallets', (table) => {
      table.decimal('totalTradingRewards').notNullable().alter();
    });
}
