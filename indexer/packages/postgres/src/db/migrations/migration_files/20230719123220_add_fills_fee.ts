import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  await knex.schema.alterTable('fills', (table) => {
    table
      .string('fee').notNullable().defaultTo('0');
  });
}

export async function down(knex: Knex): Promise<void> {
  return knex.schema.alterTable('fills', (table) => {
    table.dropColumn('fee');
  });
}
