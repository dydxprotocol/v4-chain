import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  return knex
    .schema
    .alterTable('trading_reward_aggregations', (table) => {
      // 27 is the max precision and 18 is scale, which means 9 digits before the decimal point
      // and 18 after
      table.decimal('amount', 27, 18).notNullable().alter();
    });
}

export async function down(knex: Knex): Promise<void> {
  return knex
    .schema
    .alterTable('trading_reward_aggregations', (table) => {
      table.decimal('amount').notNullable().alter();
    });
}
