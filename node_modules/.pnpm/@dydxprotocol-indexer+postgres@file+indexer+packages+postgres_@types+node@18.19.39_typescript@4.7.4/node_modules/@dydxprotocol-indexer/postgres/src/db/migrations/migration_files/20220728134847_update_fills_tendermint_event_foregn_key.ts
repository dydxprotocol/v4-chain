import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  return knex
    .schema
    .alterTable('fills', (table) => {
      table.foreign('eventId').references('tendermint_events.id');
    });
}

export async function down(knex: Knex): Promise<void> {
  return knex
    .schema
    .alterTable('fills', (table) => {
      table.dropForeign(['eventId']);
    });
}
