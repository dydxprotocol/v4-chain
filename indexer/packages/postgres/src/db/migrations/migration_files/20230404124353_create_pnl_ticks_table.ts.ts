import Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  return knex
    .schema
    .createTable('pnl_ticks', (table) => {
      table.uuid('id').primary();
      table.uuid('subaccountId').notNullable();
      table.decimal('equity', null).notNullable();
      table.decimal('totalPnl', null).notNullable();
      table.decimal('netTransfers', null).notNullable();
      table.timestamp('createdAt').notNullable();
      table.bigInteger('blockHeight').notNullable();
      table.timestamp('blockTime').notNullable();

      // Foreign
      table.foreign('subaccountId').references('subaccounts.id');
      table.foreign(['blockHeight', 'blockTime']).references(['blocks.blockHeight', 'blocks.time']);

      // Indices
      table.index(['subaccountId', 'createdAt']);
      table.index(['createdAt']);
    });
}

export async function down(knex: Knex): Promise<void> {
  return knex.schema.dropTableIfExists('pnl_ticks');
}
