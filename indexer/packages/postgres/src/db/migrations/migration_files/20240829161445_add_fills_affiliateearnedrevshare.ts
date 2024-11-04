import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  await knex.schema.table('fills', (table) => {
    table.string('affiliateEarnedRevShare');
  });
}

export async function down(knex: Knex): Promise<void> {
  await knex.schema.table('fills', (table) => {
    table.dropColumn('affiliateEarnedRevShare');
  });
}
