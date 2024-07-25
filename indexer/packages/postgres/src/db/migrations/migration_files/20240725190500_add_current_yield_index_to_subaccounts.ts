import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  await knex.schema.alterTable('subaccounts', (table) => {
    table.string('assetYieldIndex').notNullable();
  });
}

export async function down(knex: Knex): Promise<void> {
  await knex.schema.alterTable('subaccounts', (table) => {
    table.dropColumn('assetYieldIndex');
  });
}