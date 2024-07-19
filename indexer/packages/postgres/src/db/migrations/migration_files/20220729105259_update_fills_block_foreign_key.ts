import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  return knex
    .schema
    .alterTable('fills', (table) => {
      table.foreign('createdAtHeight').references('blocks.blockHeight');
    });
}

export async function down(knex: Knex): Promise<void> {
  return knex
    .schema
    .alterTable('fills', (table) => {
      table.dropForeign(['createdAtHeight']);
    });
}
