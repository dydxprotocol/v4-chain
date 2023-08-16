import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  return knex.schema.alterTable('fills', (table) => {
    table.bigInteger('clientMetadata').nullable();
  });
}

export async function down(knex: Knex): Promise<void> {
  return knex.schema.alterTable('fills', (table) => {
    table.dropColumn('clientMetadata');
  });
}
