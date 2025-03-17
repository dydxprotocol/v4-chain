import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  await knex.raw(`
    CREATE INDEX CONCURRENTLY IF NOT EXISTS "blocks_time_since_march2025_idx"
    ON "blocks" ("time")
    WHERE "time" >= '2025-03-01 00:00:00+00';
  `);
}

export async function down(knex: Knex): Promise<void> {
  await knex.raw(`
    DROP INDEX CONCURRENTLY IF EXISTS 
    "blocks_time_since_march2025_idx";
  `);
}

export const config = {
  transaction: false,
};
