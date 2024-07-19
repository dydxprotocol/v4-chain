import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  return knex.schema.raw(`
    ALTER TABLE "perpetual_markets"
    DROP CONSTRAINT "perpetual_markets_status_check",
    ADD CONSTRAINT "perpetual_markets_status_check" 
     CHECK (status IN ('ACTIVE', 'PAUSED', 'CANCEL_ONLY', 'POST_ONLY', 'INITIALIZING', 'FINAL_SETTLEMENT'))
  `);
}

export async function down(knex: Knex): Promise<void> {
  return knex.schema.raw(`
    ALTER TABLE "perpetual_markets"
    DROP CONSTRAINT "perpetual_markets_status_check",
    ADD CONSTRAINT "perpetual_markets_status_check" 
     CHECK (status IN ('ACTIVE', 'PAUSED', 'CANCEL_ONLY', 'POST_ONLY', 'INITIALIZING'))
  `);
}
