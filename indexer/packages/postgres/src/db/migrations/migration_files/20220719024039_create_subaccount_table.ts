import Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  return knex
    .schema
    .createTable('subaccounts', (table) => {
      table.uuid('id').primary();
      table.string('address').notNullable();
      table.integer('subaccountNumber').notNullable();
      table.timestamp('updatedAt').notNullable();
      table.bigInteger('updatedAtHeight').notNullable();

      // Indices
      table.index(['address']);
    });
}

export async function down(knex: Knex): Promise<void> {
  return knex.schema.dropTableIfExists('subaccounts');
}
