import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  return knex
    .schema
    .alterTable('perpetual_positions', (table) => {
      table.foreign('openEventId').references('tendermint_events.id');
      table.foreign('closeEventId').references('tendermint_events.id');
      table.foreign('lastEventId').references('tendermint_events.id');
    });
}

export async function down(knex: Knex): Promise<void> {
  return knex
    .schema
    .alterTable('perpetual_positions', (table) => {
      table.dropForeign(['openEventId']);
      table.dropForeign(['closeEventId']);
      table.dropForeign(['lastEventId']);
    });
}
