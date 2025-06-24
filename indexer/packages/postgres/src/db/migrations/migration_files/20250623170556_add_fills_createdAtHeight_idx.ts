import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  await knex.raw(`
    CREATE INDEX CONCURRENTLY IF NOT EXISTS fills_createdatheight_index
      ON fills("createdAtHeight");
  `);
}

export async function down(knex: Knex): Promise<void> {
  await knex.raw(`
    DROP INDEX CONCURRENTLY IF EXISTS fills_createdatheight_index;
  `);
}

export const config = {
  transaction: false,
};
