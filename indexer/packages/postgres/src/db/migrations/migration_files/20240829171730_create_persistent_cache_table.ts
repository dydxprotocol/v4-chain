import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  return knex.schema.createTable('persistent_cache', (table) => {
    table.string('key').primary().notNullable();
    table.string('value').notNullable();
  });
}

export async function down(knex: Knex): Promise<void> {
  return knex.schema.dropTable('persistent_cache');
}
