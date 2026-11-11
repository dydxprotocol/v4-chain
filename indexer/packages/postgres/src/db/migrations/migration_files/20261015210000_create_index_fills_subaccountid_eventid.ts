import * as Knex from 'knex';

/**
 * Adds the index backing the paginated fills query
 *   SELECT * FROM fills WHERE "subaccountId" IN ($1) ORDER BY "eventId" ASC, "createdAtHeight" DESC
 *   LIMIT ... OFFSET ...
 * Without it the planner sorts every fill of the subaccount, which for the largest subaccounts is
 * a parallel sequential scan of the whole table on each request.
 *
 * On mainnet this index was built by hand (it takes far longer than bazooka's Lambda timeout), so
 * there this migration only records it.
 */
export async function up(knex: Knex): Promise<void> {
  await knex.raw(`
    CREATE INDEX CONCURRENTLY IF NOT EXISTS "fills_subaccountid_eventid_createdatheight_index"
      ON "fills" ("subaccountId", "eventId", "createdAtHeight" DESC);
  `);
}

export async function down(knex: Knex): Promise<void> {
  await knex.raw(`
    DROP INDEX CONCURRENTLY IF EXISTS "fills_subaccountid_eventid_createdatheight_index";
  `);
}

export const config = {
  transaction: false,
};
