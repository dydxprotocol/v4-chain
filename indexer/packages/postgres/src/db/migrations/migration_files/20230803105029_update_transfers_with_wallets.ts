import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  return knex.schema.alterTable('transfers', (table) => {
    // Make senderSubaccountId and recipientSubaccountId nullable
    table.uuid('senderSubaccountId').nullable().alter();
    table.uuid('recipientSubaccountId').nullable().alter();

    // Add nullable columns recipientWalletAddress and senderWalletAddress
    table.string('recipientWalletAddress').nullable();
    table.string('senderWalletAddress').nullable();
  });
}

export async function down(knex: Knex): Promise<void> {
  return knex.schema.alterTable('transfers', (table) => {
    // Drop the new columns
    table.dropColumn('recipientWalletAddress');
    table.dropColumn('senderWalletAddress');

    // Make senderSubaccountId and recipientSubaccountId not nullable again
    table.uuid('senderSubaccountId').notNullable().alter();
    table.uuid('recipientSubaccountId').notNullable().alter();
  });
}
