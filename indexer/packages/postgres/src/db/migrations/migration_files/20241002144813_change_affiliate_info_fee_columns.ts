import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  return knex
    .schema
    .alterTable('affiliate_info', (table) => {
      table.dropColumn('totalReferredFees');
      table.dropColumn('referredNetProtocolEarnings');
      table.decimal('totalReferredTakerFees', null).notNullable().defaultTo(0);
      table.decimal('totalReferredMakerFees', null).notNullable().defaultTo(0);
      table.decimal('totalReferredMakerRebates', null).notNullable().defaultTo(0);
    });
}

export async function down(knex: Knex): Promise<void> {
  return knex
    .schema
    .alterTable('affiliate_info', (table) => {
      table.decimal('totalReferredFees', null).notNullable().defaultTo(0);
      table.decimal('referredNetProtocolEarnings', null).notNullable().defaultTo(0);
      table.dropColumn('totalReferredTakerFees');
      table.dropColumn('totalReferredMakerFees');
      table.dropColumn('totalReferredMakerRebates');
    });
}
