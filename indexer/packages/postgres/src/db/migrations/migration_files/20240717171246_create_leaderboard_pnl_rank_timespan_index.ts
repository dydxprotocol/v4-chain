import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  await knex.raw(`
    CREATE INDEX CONCURRENTLY IF NOT EXISTS "leaderboard_pnl_rank_timespan_index" ON leaderboard_pnl("rank", "timeSpan");
  `);
}

export async function down(knex: Knex): Promise<void> {
  await knex.raw(`
    DROP INDEX CONCURRENTLY IF EXISTS "leaderboard_pnl_rank_timespan_index";
  `);
}

export const config = {
  transaction: false,
};
