import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  return knex.schema.raw(`
    ALTER TABLE "compliance_status"
    DROP CONSTRAINT "compliance_status_status_check",
    ADD CONSTRAINT "compliance_status_status_check" 
    CHECK (status IN ('COMPLIANT', 'FIRST_STRIKE_CLOSE_ONLY', 'FIRST_STRIKE', 'CLOSE_ONLY', 'BLOCKED'))
  `);
}

export async function down(knex: Knex): Promise<void> {
  return knex.schema.raw(`
    ALTER TABLE "compliance_status"
    DROP CONSTRAINT "compliance_status_status_check",
    ADD CONSTRAINT "compliance_status_status_check" 
    CHECK (status IN ('COMPLIANT', 'FIRST_STRIKE', 'CLOSE_ONLY', 'BLOCKED'))
  `);
}
