import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  await knex.schema.alterTable('perpetual_markets', (table) => {
    table.decimal('baseOpenInterest', null).defaultTo(0);
  });
}

export async function down(knex: Knex): Promise<void> {
  await knex.schema.alterTable('perpetual_markets', (table) => {
    table.dropColumn('baseOpenInterest');
  });
}
