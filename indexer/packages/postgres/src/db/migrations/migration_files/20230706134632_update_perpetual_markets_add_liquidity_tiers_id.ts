import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  await knex.schema.alterTable('perpetual_markets', (table) => {
    table
      .integer('liquidityTierId').notNullable().defaultTo(0)
      .references('id')
      .inTable('liquidity_tiers');
  });
}

export async function down(knex: Knex): Promise<void> {
  return knex.schema.alterTable('perpetual_markets', (table) => {
    table.dropForeign(['liquidityTierId']);
    table.dropColumn('liquidityTierId');
  });
}
