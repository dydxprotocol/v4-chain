import Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  return knex
    .schema
    .createTable('liquidity_tiers', (table) => {
      table.integer('id').primary();
      table.string('name').nullable();
      table.bigInteger('initialMarginPpm').notNullable();  // in ppm
      table.bigInteger('maintenanceFractionPpm').notNullable();  // in ppm. Needs to be multiplied by initialMarginFraction to get the actual maintenanceMarginFraction.
      table.decimal('basePositionNotional', null).notNullable();  // in human-readable form.
    });
}

export async function down(knex: Knex): Promise<void> {
  return knex.schema.dropTableIfExists('liquidity_tiers');
}
