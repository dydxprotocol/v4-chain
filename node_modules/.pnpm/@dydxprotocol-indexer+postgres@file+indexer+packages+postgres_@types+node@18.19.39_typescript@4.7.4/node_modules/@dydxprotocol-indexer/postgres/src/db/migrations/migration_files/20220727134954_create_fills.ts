import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  return knex
    .schema
    .createTable('fills', (table) => {
      table.uuid('id').primary();
      table.uuid('subaccountId').notNullable();
      table.enum('side', [
        'BUY',
        'SELL',
      ]).notNullable();
      table.enum('liquidity', [
        'TAKER',
        'MAKER',
      ]).notNullable();
      table.enum('type', [
        'MARKET',
        'LIMIT',
        'LIQUIDATED',
        'LIQUIDATION',
      ]).notNullable();
      table.bigInteger('clobPairId').notNullable();
      table.uuid('orderId').nullable();
      table.decimal('size', null).notNullable();
      table.decimal('price', null).notNullable();
      table.decimal('quoteAmount', null).notNullable();
      table.binary('eventId', 96).notNullable();
      table.string('transactionHash').notNullable();
      table.timestamp('createdAt').notNullable();
      table.bigInteger('createdAtHeight').notNullable();

      // Foreign
      table.foreign('subaccountId').references('subaccounts.id');
      table.foreign('orderId').references('orders.id');

      // Indices
      table.index(['subaccountId', 'createdAt']);
      table.index(['subaccountId', 'createdAtHeight']);
      table.index(['subaccountId', 'clobPairId']);
      table.index(['clobPairId', 'createdAt']);
      table.index(['orderId']);
    });
}

export async function down(knex: Knex): Promise<void> {
  return knex.schema.dropTableIfExists('fills');
}
