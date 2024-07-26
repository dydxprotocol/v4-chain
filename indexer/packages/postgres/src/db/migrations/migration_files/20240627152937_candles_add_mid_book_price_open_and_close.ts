import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  return knex
    .schema
    .alterTable('candles', (table) => {
      table.decimal('orderbookMidPriceOpen', null);
      table.decimal('orderbookMidPriceClose', null);
    });
}

export async function down(knex: Knex): Promise<void> {
  return knex
    .schema
    .alterTable('candles', (table) => {
      table.dropColumn('orderbookMidPriceOpen');
      table.dropColumn('orderbookMidPriceClose');
    });
}
