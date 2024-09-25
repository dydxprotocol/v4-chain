import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  return knex.schema.alterTable('affiliate_info', (table) => {
    table.decimal('affiliateEarnings', null).notNullable().defaultTo(0).alter();
    table.decimal('totalReferredFees', null).notNullable().defaultTo(0).alter();
    table.decimal('referredNetProtocolEarnings', null).notNullable().defaultTo(0).alter();
  });
}

export async function down(knex: Knex): Promise<void> {
  return knex.schema.alterTable('affiliate_info', (table) => {
    table.decimal('affiliateEarnings', null).alter();
    table.decimal('totalReferredFees', null).alter();
    table.decimal('referredNetProtocolEarnings', null).alter();
  });
}
