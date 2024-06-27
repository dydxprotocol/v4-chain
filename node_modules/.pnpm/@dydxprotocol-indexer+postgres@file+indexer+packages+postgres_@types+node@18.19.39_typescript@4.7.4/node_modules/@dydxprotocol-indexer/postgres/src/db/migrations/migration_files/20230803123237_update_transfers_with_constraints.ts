import * as Knex from 'knex';

// Add a constraint to the `orders` table to ensure that one of `goodTilBlock` and
// `goodTilBlockTime` is non-null for each row, but not both.
export async function up(knex: Knex): Promise<void> {
  return knex
    .schema
    .raw(`
      ALTER TABLE "transfers" ADD CONSTRAINT recipient_non_null CHECK (num_nonnulls("recipientWalletAddress", "recipientSubaccountId") = 1);
    `).raw(`
      ALTER TABLE "transfers" ADD CONSTRAINT sender_non_null CHECK (num_nonnulls("senderWalletAddress", "senderSubaccountId") = 1);
    `);
}

export async function down(knex: Knex): Promise<void> {
  return knex
    .schema
    .raw(`
      ALTER TABLE "transfers" DROP CONSTRAINT recipient_non_null;
    `)
    .raw(`
      ALTER TABLE "transfers" DROP CONSTRAINT sender_non_null;
    `);
}
