import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  return knex
    .schema
    .alterTable('liquidity_tiers', (table) => {
      table.dropColumn('basePositionNotional');
    });
}

export async function down(knex: Knex): Promise<void> {
  return knex
    .schema
    .alterTable('liquidity_tiers', (table) => {
      table.decimal('basePositionNotional', null).notNullable().defaultTo('0');
    });
}
