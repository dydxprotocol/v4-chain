import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  return knex
    .schema
    .alterTable('perpetual_markets', (table) => {
      table.dropColumn('lastPrice');
    });
}

export async function down(knex: Knex): Promise<void> {
  return knex
    .schema
    .alterTable('perpetual_markets', (table) => {
      table.decimal('lastPrice', null).defaultTo('0').notNullable();
    });
}
