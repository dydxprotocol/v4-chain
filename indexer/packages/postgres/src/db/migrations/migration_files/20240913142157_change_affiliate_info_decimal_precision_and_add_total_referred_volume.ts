import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  return knex.schema.alterTable('affiliate_info', (table) => {
    // null indicates variable precision whereas not specifying will result in 8,2 precision,scale
    table.decimal('affiliateEarnings', null).alter();
    table.decimal('totalReferredFees', null).alter();
    table.decimal('referredNetProtocolEarnings', null).alter();

    table.decimal('referredTotalVolume', null).notNullable();
  });
}

export async function down(knex: Knex): Promise<void> {
  return knex.schema.alterTable('affiliate_info', (table) => {
    table.decimal('affiliateEarnings').alter();
    table.decimal('totalReferredFees').alter();
    table.decimal('referredNetProtocolEarnings').alter();

    table.dropColumn('referredTotalVolume');
  });
}
