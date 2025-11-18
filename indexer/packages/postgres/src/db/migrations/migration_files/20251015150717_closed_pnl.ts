import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  await knex.schema.alterTable('perpetual_positions', (table) => {
    table.decimal('totalRealizedPnl', null).nullable();
  });

  await knex.schema.alterTable('fills', (table) => {
    table.decimal('positionSizeBefore', null).nullable();
    table.decimal('entryPriceBefore', null).nullable();
    table.enum('positionSideBefore', [
      'LONG',
      'SHORT',
    ]).nullable().defaultTo(null);
  });
}

export async function down(knex: Knex): Promise<void> {
  await knex.schema.alterTable('perpetual_positions', (table) => {
    table.dropColumn('totalRealizedPnl');
  });

  await knex.schema.alterTable('fills', (table) => {
    table.dropColumn('positionSizeBefore');
    table.dropColumn('entryPriceBefore');
    table.dropColumn('positionSideBefore');
  });
}
