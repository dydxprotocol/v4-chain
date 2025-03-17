import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  await knex.raw(`
    CREATE INDEX CONCURRENTLY IF NOT EXISTS "orders_open_shortterm_idx" 
    ON "orders" ("goodTilBlock")
    WHERE status = 'OPEN' AND "orderFlags" = 0 AND "goodTilBlock" IS NOT NULL;
  `);
}

export async function down(knex: Knex): Promise<void> {
  await knex.raw(`
    DROP INDEX CONCURRENTLY IF EXISTS "orders_open_shortterm_idx";
  `);
}

export const config = {
  transaction: false,
};
