import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  await knex.schema.alterTable('perpetual_positions', (table) => {
    table.dropColumn('totalRealizedPnl');
  });
}

export async function down(knex: Knex): Promise<void> {
  await knex.schema.alterTable('perpetual_positions', (table) => {
    table.decimal('totalRealizedPnl', null).nullable();
  });
}
