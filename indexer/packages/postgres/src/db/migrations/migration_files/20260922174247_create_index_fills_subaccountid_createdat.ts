import * as Knex from 'knex';

/**
 * Adds the index backing the parent-subaccount fills query
 *   SELECT * FROM fills WHERE "subaccountId" IN (...) ORDER BY "createdAt" ASC, "eventId" ASC, ...
 * Without it the planner scans fills_createdat_index from the oldest row forward, filtering by
 * subaccountId, which walks a huge fraction of the table for subaccounts without old fills.
 */
export async function up(knex: Knex): Promise<void> {
  await knex.raw(`
    CREATE INDEX CONCURRENTLY IF NOT EXISTS "fills_subaccountid_createdat_eventid_index"
      ON "fills" ("subaccountId", "createdAt", "eventId", "createdAtHeight" DESC);
  `);
}

export async function down(knex: Knex): Promise<void> {
  await knex.raw(`
    DROP INDEX CONCURRENTLY IF EXISTS "fills_subaccountid_createdat_eventid_index";
  `);
}

export const config = {
  transaction: false,
};
