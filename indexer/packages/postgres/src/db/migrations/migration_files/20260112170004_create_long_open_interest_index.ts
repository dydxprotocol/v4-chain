import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  await knex.raw(`
    CREATE INDEX CONCURRENTLY perpetual_positions_oi_open_idx ON perpetual_positions ("perpetualId", side) INCLUDE (size) WHERE status = 'OPEN';
    `);
}

export async function down(knex: Knex): Promise<void> {
  await knex.raw(`
    DROP INDEX CONCURRENTLY IF EXISTS "perpetual_positions_oi_open_idx";
    `);
}

export const config = {
  transaction: false,
};
