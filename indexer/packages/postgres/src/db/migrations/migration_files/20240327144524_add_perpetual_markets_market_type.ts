import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  return knex.schema.table('perpetual_markets', (table) => {
    table.string('marketType').notNullable().defaultTo('CROSS');
  });
}

export async function down(knex: Knex): Promise<void> {
  return knex.schema.table('perpetual_markets', (table) => {
    table.dropColumn('marketType');
  });
}
