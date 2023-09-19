import * as Knex from 'knex';

// Use raw SQL to ensure index creation / deletion does not lock the table
export async function up(knex: Knex): Promise<void> {
  await knex.raw(`
    CREATE INDEX CONCURRENTLY IF NOT EXISTS "subaccounts_updatedat_index" on "subaccounts" ("updatedAt");
  `);
}

export async function down(knex: Knex): Promise<void> {
  await knex.raw(`
    DROP INDEX CONCURRENTLY IF EXISTS "subaccounts_updatedat_index";
  `);
}

// `CONCURRENTLY` cannot be used within a transaction
export const config = {
  transaction: false,
};
