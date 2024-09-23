import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  await knex.schema.alterTable('liquidity_tiers', (table) => {
    table.string('openInterestLowerCap').nullable().defaultTo(null);
    table.string('openInterestUpperCap').nullable().defaultTo(null);
  });

}

export async function down(knex: Knex): Promise<void> {
  await knex.schema.alterTable('liquidity_tiers', (table) => {
    table.dropColumn('openInterestLowerCap');
    table.dropColumn('openInterestUpperCap');
  });
}
