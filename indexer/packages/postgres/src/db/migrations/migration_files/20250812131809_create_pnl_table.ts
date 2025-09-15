import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  return knex
    .schema
    .createTable('pnl', (table) => {
      table.uuid('subaccountId').notNullable();
      table.timestamp('createdAt').notNullable();
      table.bigInteger('createdAtHeight').notNullable();
      table.decimal('equity', null).notNullable();
      table.decimal('netTransfers', null).notNullable();
      table.decimal('totalPnl', null).notNullable();

      // Primary key
      table.primary(['subaccountId', 'createdAt']);

      // Foreign
      table.foreign('subaccountId').references('subaccounts.id');

      // Indices
      table.index(['createdAtHeight']);
    });
}

export async function down(knex: Knex): Promise<void> {
  return knex.schema.dropTableIfExists('pnl');
}
