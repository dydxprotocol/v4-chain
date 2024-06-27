import Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  return knex
    .schema
    .createTable('blocks', (table) => {
      table.bigInteger('blockHeight').primary();
      table.timestamp('time').notNullable();
    });
}

export async function down(knex: Knex): Promise<void> {
  return knex.schema.dropTableIfExists('blocks');
}
