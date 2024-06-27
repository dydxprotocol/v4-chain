import Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  return knex
    .schema
    .createTable('tendermint_events', (table) => {
      table.binary('id', 96).primary();
      table.bigInteger('blockHeight').notNullable();
      table.integer('transactionIndex').notNullable();
      table.integer('eventIndex').notNullable();

      // Indices
      table.index(['blockHeight']);
      table.index(['blockHeight', 'transactionIndex']);
    });
}

export async function down(knex: Knex): Promise<void> {
  return knex.schema.dropTableIfExists('tendermint_events');
}
