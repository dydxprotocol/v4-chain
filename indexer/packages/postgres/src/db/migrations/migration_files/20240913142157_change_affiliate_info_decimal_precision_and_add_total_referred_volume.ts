import * as Knex from 'knex';

// No data has been stored added at time of commit
export async function up(knex: Knex): Promise<void> {
  return knex.schema.alterTable('affiliate_info', (table) => {
    // null indicates variable precision whereas not specifying will result in 8,2 precision,scale
    table.decimal('affiliateEarnings', null).notNullable().defaultTo(0).alter();
    table.decimal('totalReferredFees', null).notNullable().defaultTo(0).alter();
    table.decimal('referredNetProtocolEarnings', null).notNullable().defaultTo(0).alter();

    table.decimal('referredTotalVolume', null).notNullable().defaultTo(0);
  });
}

export async function down(knex: Knex): Promise<void> {
  return knex.schema.alterTable('affiliate_info', (table) => {
    table.decimal('affiliateEarnings').notNullable().defaultTo(0).alter();
    table.decimal('totalReferredFees').notNullable().defaultTo(0).alter();
    table.decimal('referredNetProtocolEarnings').notNullable().defaultTo(0).alter();

    table.dropColumn('referredTotalVolume');
  });
}
