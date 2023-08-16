import Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  return knex
    .schema
    .createTable('funding_index_updates', (table) => {
      table.uuid('id').primary();
      table.bigInteger('perpetualId').notNullable();
      table.binary('eventId', 96).notNullable();
      table.decimal('rate', null).notNullable();
      table.decimal('oraclePrice', null).notNullable();
      table.decimal('fundingIndex', null).notNullable();
      table.timestamp('effectiveAt').notNullable();
      table.bigInteger('effectiveAtHeight').notNullable();

      // Foreign
      table.foreign('eventId').references('tendermint_events.id');
      table.foreign('perpetualId').references('perpetual_markets.id');
      table.foreign('effectiveAtHeight').references('blocks.blockHeight');

      // Indices
      table.index(['perpetualId']);
      table.index(['effectiveAtHeight']);
    });
}

export async function down(knex: Knex): Promise<void> {
  return knex.schema.dropTableIfExists('funding_index_updates');
}
