import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  return knex
    .schema
    .createTable('orders', (table) => {
      table.uuid('id').primary();
      table.uuid('subaccountId').notNullable();
      table.bigInteger('clientId').notNullable();
      table.bigInteger('clobPairId').notNullable();
      table.enum(
        'side',
        [
          'BUY',
          'SELL',
        ],
      ).notNullable();
      table.decimal('size', null).notNullable();
      table.decimal('totalFilled', null).notNullable();
      table.decimal('price', null).notNullable();
      table.enum(
        'type',
        [
          'LIMIT',
          'MARKET',
          'STOP_LIMIT',
          'STOP_MARKET',
          'TRAILING_STOP',
          'TAKE_PROFIT',
          'TAKE_PROFIT_MARKET',
          'LIQUIDATED',
          'LIQUIDATION',
          'HARD_TRADE',
          'FAILED_HARD_TRADE',
          'TRANSFER_PLACEHOLDER',
        ],
      ).notNullable();
      table.enum(
        'status',
        [
          'OPEN',
          'FILLED',
          'CANCELED',
          'BEST_EFFORT_CANCELED',
        ],
      ).notNullable();
      table.enum(
        'timeInForce',
        [
          'GTT',
          'FOK',
          'IOC',
          'POST_ONLY',
        ],
      ).notNullable();
      table.boolean('reduceOnly').notNullable();
      // `orderFlags` indicate the type of the order on the protocol, and is one of the fields that
      // is used to generate the order UUID.
      table.bigInteger('orderFlags').notNullable();
      // One of `goodTilBlock` or `goodTilBlockTime` must be set, but not both. This is enforced
      // by a CHECK constraint in a later migration (20230201173734).
      table.bigInteger('goodTilBlock').nullable();
      table.timestamp('goodTilBlockTime').nullable();
      table.bigInteger('createdAtHeight').nullable();
      table.bigInteger('clientMetadata').notNullable();

      // Foreign
      table.foreign('subaccountId').references('subaccounts.id');

      // Indices
      table.index(['subaccountId']);
      table.index(['clobPairId', 'side', 'price']);
    });
}

export async function down(knex: Knex): Promise<void> {
  return knex.schema.dropTableIfExists('orders');
}
