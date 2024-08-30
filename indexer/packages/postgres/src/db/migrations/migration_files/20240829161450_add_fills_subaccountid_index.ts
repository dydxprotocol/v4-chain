import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  await knex.raw(`
    CREATE INDEX CONCURRENTLY IF NOT EXISTS "fills_subaccountid_index" ON "fills" ("subaccountId");
    `);
}

export async function down(knex: Knex): Promise<void> {
  await knex.raw(`
    DROP INDEX CONCURRENTLY IF EXISTS "fills_subaccountid_index";
    `);
}

export const config = {
  transaction: false,
};
