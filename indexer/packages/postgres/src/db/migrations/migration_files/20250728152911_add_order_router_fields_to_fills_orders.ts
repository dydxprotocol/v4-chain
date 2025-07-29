import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  await knex.schema.table('orders', (table) => {
    table.string('orderRouterAddress').nullable();
  });
  await knex.schema.table('fills', (table) => {
    table.string('orderRouterFee').nullable();
    table.string('orderRouterAddress').nullable();
  });
}

export async function down(knex: Knex): Promise<void> {
  await knex.schema.table('orders', (table) => {
    table.dropColumn('orderRouterAddress');
  });
  await knex.schema.table('fills', (table) => {
    table.dropColumn('orderRouterFee');
    table.dropColumn('orderRouterAddress');
  });
}
