import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  return knex.schema.alterTable('orders', (table) => {
    table.string('duration').nullable();
    table.string('interval').nullable();
    table.string('priceTolerance').nullable();
  });
}

export async function down(knex: Knex): Promise<void> {
  return knex.schema.alterTable('orders', (table) => {
    table.dropColumn('duration');
    table.dropColumn('interval');
    table.dropColumn('priceTolerance');
  });
}
