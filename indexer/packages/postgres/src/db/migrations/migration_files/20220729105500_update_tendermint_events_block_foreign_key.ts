import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  return knex
    .schema
    .alterTable('tendermint_events', (table) => {
      table.foreign('blockHeight').references('blocks.blockHeight');
    });
}

export async function down(knex: Knex): Promise<void> {
  return knex
    .schema
    .alterTable('tendermint_events', (table) => {
      table.dropForeign(['blockHeight']);
    });
}
