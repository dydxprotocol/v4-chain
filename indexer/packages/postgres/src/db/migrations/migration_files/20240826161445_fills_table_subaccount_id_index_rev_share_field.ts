import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  await knex.schema.table('fills', (table) => {
    table.decimal('affiliateEarnedRevShare').notNullable().defaultTo('0');
    table.index('subaccountId');
  });
}

export async function down(knex: Knex): Promise<void> {
  await knex.schema.table('fills', (table) => {
    table.dropIndex('subaccountId');
    table.dropColumn('affiliateEarnedRevShare');
  });
}
