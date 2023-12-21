import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  return knex
    .schema
    .alterTable('trading_rewards', (table) => {
      table.decimal('amount', 18, 18).notNullable().alter();
    });
}

export async function down(knex: Knex): Promise<void> {
  return knex
    .schema
    .alterTable('trading_rewards', (table) => {
      table.decimal('amount').notNullable().alter();
    });
}
