import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  await knex.raw(`
    DROP INDEX CONCURRENTLY IF EXISTS "transfers_sendersubaccountid_createdatheight_index";

    CREATE INDEX CONCURRENTLY IF NOT EXISTS transfers_sender_id_height_nn
    ON transfers ("senderSubaccountId", "createdAtHeight")
    WHERE "senderSubaccountId" IS NOT NULL;

    DROP INDEX CONCURRENTLY IF EXISTS "transfers_recipientsubaccountid_createdatheight_index";

    CREATE INDEX CONCURRENTLY IF NOT EXISTS transfers_recipient_id_height_nn
    ON transfers ("recipientSubaccountId", "createdAtHeight")
    WHERE "recipientSubaccountId" IS NOT NULL;
  `);
}

export async function down(knex: Knex): Promise<void> {
  await knex.raw(`
    DROP INDEX CONCURRENTLY IF EXISTS "transfers_sender_id_height_nn";

    CREATE INDEX CONCURRENTLY IF NOT EXISTS transfers_sendersubaccountid_createdatheight_index
    ON transfers ("senderSubaccountId", "createdAtHeight");

    DROP INDEX CONCURRENTLY IF EXISTS "transfers_recipient_id_height_nn";

    CREATE INDEX CONCURRENTLY IF NOT EXISTS transfers_recipientsubaccountid_createdatheight_index
    ON transfers ("recipientSubaccountId", "createdAtHeight");
  `);
}

export const config = {
  transaction: false,
};
