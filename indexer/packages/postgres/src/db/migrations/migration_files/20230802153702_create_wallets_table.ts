import Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  return knex
    .schema
    .createTable('wallets', (table) => {
      table.string('address').primary();
    });
}

export async function down(knex: Knex): Promise<void> {
  return knex.schema.dropTableIfExists('wallets');
}
