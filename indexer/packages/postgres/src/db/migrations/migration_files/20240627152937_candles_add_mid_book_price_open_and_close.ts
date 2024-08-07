import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  return knex
    .schema
    .alterTable('candles', (table) => {
<<<<<<< HEAD
      table.decimal('orderbookMidPriceOpen', null);
      table.decimal('orderbookMidPriceClose', null);
=======
      table.decimal('orderbookMidPriceOpen', null).nullable();
      table.decimal('orderbookMidPriceClose', null).nullable();
>>>>>>> 73b04dc3 (Adam/add candles hloc (#2047))
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
