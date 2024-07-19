import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  // eslint-disable-next-line @typescript-eslint/quotes
  await knex.raw(`
    CREATE INDEX CONCURRENTLY IF NOT EXISTS "perpetual_positions_subaccountId_perpetualId_openEventId_index" ON "perpetual_positions" ("subaccountId", "perpetualId", "openEventId" DESC);
  `);
}

export async function down(knex: Knex): Promise<void> {
  await knex.raw(`
    DROP INDEX CONCURRENTLY IF EXISTS "perpetual_positions_subaccountId_perpetualId_openEventId_index";
  `);
}

export const config = {
  transaction: false,
};
