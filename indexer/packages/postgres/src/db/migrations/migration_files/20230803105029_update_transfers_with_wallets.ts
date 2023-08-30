import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  return knex.schema.alterTable('transfers', (table) => {
    // Make senderSubaccountId and recipientSubaccountId nullable
    table.uuid('senderSubaccountId').nullable().alter();
    table.uuid('recipientSubaccountId').nullable().alter();

    // Add nullable columns recipientWalletAddress and senderWalletAddress
    table.string('recipientWalletAddress').nullable();
    table.string('senderWalletAddress').nullable();

    // Foreign key constraints for new columns
    table.foreign('recipientWalletAddress').references('wallets.address');
    table.foreign('senderWalletAddress').references('wallets.address');
  });
}

export async function down(knex: Knex): Promise<void> {
  return knex.schema.alterTable('transfers', (table) => {
    // Remove the foreign key constraints
    table.dropForeign(['recipientWalletAddress']);
    table.dropForeign(['senderWalletAddress']);

    // Drop the new columns
    table.dropColumn('recipientWalletAddress');
    table.dropColumn('senderWalletAddress');

    // Make senderSubaccountId and recipientSubaccountId not nullable again
    table.uuid('senderSubaccountId').notNullable().alter();
    table.uuid('recipientSubaccountId').notNullable().alter();
  });
}
