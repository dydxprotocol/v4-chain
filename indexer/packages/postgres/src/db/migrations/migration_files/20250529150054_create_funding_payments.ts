import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  return knex
    .schema
    .createTable('funding_payments', (table) => {
      table.uuid('subaccountId').notNullable();
      table.timestamp('createdAt').notNullable();
      table.bigInteger('createdAtHeight').notNullable();
      table.bigInteger('perpetualId').notNullable();
      table.text('ticker').notNullable();
      table.decimal('oraclePrice', null).notNullable();
      table.decimal('size', null).notNullable();
      table.enum('side', [
        'LONG',
        'SHORT',
      ]).notNullable();
      table.decimal('rate', null).notNullable();
      table.decimal('payment', null).notNullable();
      table.decimal('fundingIndex', null).notNullable();

      // Primary key
      table.primary(['subaccountId', 'createdAt', 'ticker']);

      // Foreign key
      table.foreign('subaccountId').references('subaccounts.id');
      table.foreign('perpetualId').references('perpetual_markets.id');

      // Index
      table.index(['createdAtHeight']);
    });
}

export async function down(knex: Knex): Promise<void> {
  return knex.schema.dropTableIfExists('funding_payments');
}
