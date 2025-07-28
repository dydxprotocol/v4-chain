import * as Knex from "knex";

const ORDERS_TABLE = 'orders';
const FILLS_TABLE = 'fills';
const ORDER_ROUTER_ADDRESS_COLUMN = 'orderRouterAddress';

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