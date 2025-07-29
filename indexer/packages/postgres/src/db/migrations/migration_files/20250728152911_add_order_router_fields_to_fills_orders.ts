import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  await knex.schema.table('orders', (table) => {
    table.string('orderRouterAddress');
  });
  await knex.schema.table('fills', (table) => {
    table.string('orderRouterAddress');
  });
}


export async function down(knex: Knex): Promise<void> {
  await knex.schema.table('orders', (table) => {
    table.dropColumn('orderRouterAddress');
  });
  await knex.schema.table('fills', (table) => {
    table.dropColumn('orderRouterAddress');
  });
}