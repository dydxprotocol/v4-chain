import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  return knex
    .schema
    .alterTable('perpetual_positions', (table) => {
      table.decimal('cumulativeEntryPrice', null).defaultTo(0).notNullable();
    });
}

export async function down(knex: Knex): Promise<void> {
  return knex
    .schema
    .alterTable('perpetual_positions', (table) => {
      table.dropColumn('cumulativeEntryPrice');
    });
}
