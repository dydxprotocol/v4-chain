import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  return knex
    .schema
    .alterTable('perpetual_markets', (table) => {
      table.dropColumn('minOrderBaseQuantums');
    });
}

export async function down(knex: Knex): Promise<void> {
  return knex
    .schema
    .alterTable('perpetual_markets', (table) => {
      table.decimal('minOrderBaseQuantums', null).notNullable();
    });
}
