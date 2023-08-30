import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  await knex.schema.alterTable('perpetual_markets', (table) => {
    table.dropColumn('initialMarginFraction');
    table.dropColumn('incrementalInitialMarginFraction');
    table.dropColumn('maintenanceMarginFraction');
  });
}

export async function down(knex: Knex): Promise<void> {
  await knex.schema.alterTable('perpetual_markets', (table) => {
    table.decimal('initialMarginFraction', null).notNullable().defaultTo('0');
    table.decimal('incrementalInitialMarginFraction', null).notNullable().defaultTo('0');
    table.decimal('maintenanceMarginFraction', null).notNullable().defaultTo('0');
  });
}
