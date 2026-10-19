import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  await knex.schema.alterTable('fills', (table) => {
    // Realized price-PnL of the closing portion of this fill vs the
    // position's average entry price at fill time; excludes fees and
    // funding. NULL for opening/increasing fills and for fills that
    // predate this column (no backfill).
    table.decimal('closedPnl', null).nullable();
  });
}

export async function down(knex: Knex): Promise<void> {
  await knex.schema.alterTable('fills', (table) => {
    table.dropColumn('closedPnl');
  });
}
