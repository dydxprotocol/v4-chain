import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  await knex.schema.alterTable('orders', (table) => {
    table.string('orderRouterAddress').nullable();
  });

  await knex.schema.alterTable('fills', (table) => {
    table.string('orderRouterAddress').nullable();
    table.string('orderRouterFee').nullable();
  });
}

export async function down(knex: Knex): Promise<void> {
  await knex.schema.alterTable('orders', (table) => {
    table.dropColumn('orderRouterAddress');
  });

  await knex.schema.alterTable('fills', (table) => {
    table.dropColumn('orderRouterAddress');
    table.dropColumn('orderRouterFee');
  });
}
