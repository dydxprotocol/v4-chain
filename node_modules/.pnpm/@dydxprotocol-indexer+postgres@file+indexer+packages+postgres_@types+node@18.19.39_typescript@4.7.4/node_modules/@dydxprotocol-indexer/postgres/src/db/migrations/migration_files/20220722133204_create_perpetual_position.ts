import Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  return knex
    .schema
    .createTable('perpetual_positions', (table) => {
      table.uuid('id').primary();
      table.uuid('subaccountId').notNullable();
      table.bigInteger('perpetualId').notNullable();
      table.enum('side', [
        'LONG',
        'SHORT',
      ]).notNullable();
      table.enum('status', [
        'OPEN',
        'CLOSED',
        'LIQUIDATED',
      ]).notNullable();
      // The size of the position. Positive for long, negative for short.
      table.decimal('size', null).notNullable();
      table.decimal('maxSize', null).notNullable();
      table.decimal('entryPrice', null).notNullable();
      table.decimal('exitPrice', null).nullable();
      table.decimal('sumOpen', null).notNullable();
      table.decimal('sumClose', null).notNullable();
      table.timestamp('createdAt').notNullable();
      table.timestamp('closedAt').nullable();
      table.bigInteger('createdAtHeight').notNullable();
      table.bigInteger('closedAtHeight').nullable();
      table.binary('openEventId', 96).notNullable();
      table.binary('closeEventId', 96).nullable();
      table.binary('lastEventId', 96).notNullable();
      // This is the sum of all settled funding payments made to the subaccount. A positive value
      // indicates the subaccount received funding, negative indicates the subaccount paid funding.
      table.decimal('settledFunding', null).notNullable();

      // Foreign
      table.foreign('perpetualId').references('perpetual_markets.id');
      table.foreign('subaccountId').references('subaccounts.id');

      // Indices
      table.index(['subaccountId', 'status']);
      table.index(['subaccountId', 'createdAt']);
      table.index(['subaccountId', 'createdAtHeight']);
      table.index(['perpetualId', 'status']);
    });
}

export async function down(knex: Knex): Promise<void> {
  return knex.schema.dropTableIfExists('perpetual_positions');
}
