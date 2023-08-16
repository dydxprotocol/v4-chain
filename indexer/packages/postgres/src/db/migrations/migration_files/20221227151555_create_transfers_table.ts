import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  return knex
    .schema
    .createTable('transfers', (table) => {
      table.uuid('id').primary();
      table.uuid('senderSubaccountId').notNullable();
      table.uuid('recipientSubaccountId').notNullable();
      table.string('assetId').notNullable();
      table.decimal('size', null).notNullable();
      table.binary('eventId', 96).notNullable();
      table.string('transactionHash').notNullable();
      table.timestamp('createdAt').notNullable();
      table.bigInteger('createdAtHeight').notNullable();

      // Foreign
      table.foreign('senderSubaccountId').references('subaccounts.id');
      table.foreign('recipientSubaccountId').references('subaccounts.id');
      table.foreign('assetId').references('assets.id');
      table.foreign('eventId').references('tendermint_events.id');

      // Indices
      table.index(['senderSubaccountId', 'createdAt']);
      table.index(['senderSubaccountId', 'createdAtHeight']);
      table.index(['recipientSubaccountId', 'createdAt']);
      table.index(['recipientSubaccountId', 'createdAtHeight']);
      table.index(['assetId', 'createdAt']);
    });
}

export async function down(knex: Knex): Promise<void> {
  return knex.schema.dropTableIfExists('transfers');
}
