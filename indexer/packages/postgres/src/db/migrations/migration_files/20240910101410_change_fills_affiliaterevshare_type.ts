import * as Knex from 'knex';

// No data has been stored added at time of commit
export async function up(knex: Knex): Promise<void> {
  return knex.schema.alterTable('fills', (table) => {
    // decimal('columnName') has is 8,2 precision and scale
    // decimal('columnName', null) has variable precision and scale
    table.decimal('affiliateRevShare', null).notNullable().defaultTo(0).alter();
  });
}

export async function down(knex: Knex): Promise<void> {
  return knex.schema.alterTable('fills', (table) => {
    table.string('affiliateRevShare').notNullable().defaultTo('0').alter();
  });
}
