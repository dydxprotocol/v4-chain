import * as Knex from 'knex';

// No data has been stored added at time of commit
export async function up(knex: Knex): Promise<void> {
  return knex.schema.alterTable('affiliate_info', (table) => {
    // decimal('columnName') has is 8,2 precision and scale
    // decimal('columnName', null) has variable precision and scale
    table.decimal('affiliateEarnings', null).notNullable().defaultTo(0).alter();
    table.decimal('totalReferredFees', null).notNullable().defaultTo(0).alter();
    table.decimal('referredNetProtocolEarnings', null).notNullable().defaultTo(0).alter();

    table.decimal('totalReferredVolume', null).notNullable().defaultTo(0);
  });
}

export async function down(knex: Knex): Promise<void> {
  return knex.schema.alterTable('affiliate_info', (table) => {
    table.decimal('affiliateEarnings').notNullable().defaultTo(0).alter();
    table.decimal('totalReferredFees').notNullable().defaultTo(0).alter();
    table.decimal('referredNetProtocolEarnings').notNullable().defaultTo(0).alter();

    table.dropColumn('totalReferredVolume');
  });
}
