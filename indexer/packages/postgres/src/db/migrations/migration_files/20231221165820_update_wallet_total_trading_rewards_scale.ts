import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  return knex
    .schema
    .alterTable('wallets', (table) => {
      // 27 is the max precision and 18 is scale, which means 9 digits before the decimal point
      // and 18 after
      table.decimal('totalTradingRewards', 27, 18).notNullable().alter();
    });
}

export async function down(knex: Knex): Promise<void> {
  return knex
    .schema
    .alterTable('wallets', (table) => {
      table.decimal('totalTradingRewards').notNullable().alter();
    });
}
