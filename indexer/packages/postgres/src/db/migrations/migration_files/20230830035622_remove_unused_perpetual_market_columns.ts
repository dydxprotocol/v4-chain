import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  return knex
    .schema
    .alterTable('perpetual_markets', (table) => {
      table.dropColumn('baseAsset');
      table.dropColumn('quoteAsset');
      table.dropColumn('basePositionSize');
      table.dropColumn('incrementalPositionSize');
      table.dropColumn('maxPositionSize');
    });
}

export async function down(knex: Knex): Promise<void> {
  return knex
    .schema
    .alterTable('perpetual_markets', (table) => {
      table.string('baseAsset').notNullable();
      table.string('quoteAsset').notNullable();
      table.decimal('basePositionSize', null).notNullable();
      table.decimal('incrementalPositionSize', null).notNullable();
      table.decimal('maxPositionSize', null).notNullable();
    });
}
