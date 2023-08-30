import Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  return knex
    .schema
    .createTable('transactions', (table) => {
      table.uuid('id').primary();
      table.bigInteger('blockHeight').notNullable();
      table.integer('transactionIndex').notNullable();
      table.string('transactionHash').notNullable();

      // Indices
      table.index(['blockHeight']);
      table.index(['blockHeight', 'transactionIndex']);
    });
}

export async function down(knex: Knex): Promise<void> {
  return knex.schema.dropTableIfExists('transactions');
}
