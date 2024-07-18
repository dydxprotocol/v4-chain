import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  return knex.schema.createTable('leaderboard_pnl', (table) => {
    table.string('address').notNullable().references('address').inTable('wallets');
    table.enum(
      'timeSpan',
      [
        'ONE_DAY',
        'SEVEN_DAYS',
        'THIRTY_DAYS',
        'ONE_YEAR',
        'ALL_TIME',
      ],
    );
    table.string('pnl').notNullable();
    table.string('currentEquity').notNullable();
    table.integer('rank').notNullable();
    table.primary(['address', 'timeSpan']);
  });
}

export async function down(knex: Knex): Promise<void> {
  return knex.schema.dropTable('leaderboard_pnl');
}
