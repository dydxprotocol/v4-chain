import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  return knex.raw(`
    ALTER TABLE ONLY fills
    DROP CONSTRAINT IF EXISTS fills_type_check;

    ALTER TABLE ONLY fills
    ADD CONSTRAINT fills_type_check
    CHECK (type = ANY (ARRAY['MARKET'::text, 'LIMIT'::text, 'LIQUIDATED'::text, 'LIQUIDATION'::text, 'DELEVERAGED'::text, 'OFFSETTING'::text]));
  `);
}

export async function down(knex: Knex): Promise<void> {
  return knex.raw(`
    ALTER TABLE ONLY fills
    DROP CONSTRAINT IF EXISTS fills_type_check;

    ALTER TABLE ONLY fills
    ADD CONSTRAINT fills_type_check
    CHECK (type = ANY (ARRAY['MARKET'::text, 'LIMIT'::text, 'LIQUIDATED'::text, 'LIQUIDATION'::text]));
  `);
}
