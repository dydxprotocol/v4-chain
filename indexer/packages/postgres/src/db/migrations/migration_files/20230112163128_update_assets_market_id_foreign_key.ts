import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  return knex
    .schema
    .alterTable('assets', (table) => {
      table.foreign('marketId').references('markets.id');
    });
}

export async function down(knex: Knex): Promise<void> {
  return knex
    .schema
    .alterTable('assets', (table) => {
      table.dropForeign(['marketId']);
    });
}
