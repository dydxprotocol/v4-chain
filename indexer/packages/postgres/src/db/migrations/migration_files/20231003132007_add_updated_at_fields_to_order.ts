import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  return knex.schema.alterTable('orders', (table) => {
    table.timestamp('updatedAt');
    table.bigInteger('updatedAtHeight');
  });
}

export async function down(knex: Knex): Promise<void> {
  return knex.schema.alterTable('orders', (table) => {
    table.dropColumn('updatedAt');
    table.dropColumn('updatedAtHeight');
  });
}
