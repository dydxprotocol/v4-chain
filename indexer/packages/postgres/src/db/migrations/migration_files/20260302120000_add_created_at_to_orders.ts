import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  return knex.schema.alterTable('orders', (table) => {
    table.timestamp('createdAt').nullable();
  });
}

export async function down(knex: Knex): Promise<void> {
  return knex.schema.alterTable('orders', (table) => {
    table.dropColumn('createdAt');
  });
}
