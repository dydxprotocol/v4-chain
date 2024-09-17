import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  await knex.schema.table('fills', (table) => {
    table.renameColumn('affiliateEarnedRevShare', 'affiliateRevShare');
  });
}

export async function down(knex: Knex): Promise<void> {
  await knex.schema.table('fills', (table) => {
    table.renameColumn('affiliateRevShare', 'affiliateEarnedRevShare');
  });
}
