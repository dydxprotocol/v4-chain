import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  return knex.schema.createTable('affiliate_per_referee_stats', (table) => {
    table.string('affiliateAddress').primary().notNullable();
    table.string('refereeAddress').primary().notNullable();
    table.decimal('affiliateEarnings', null).notNullable().defaultTo(0);
    table.integer('referredMakerTrades').notNullable().defaultTo(0);
    table.integer('referredTakerTrades').notNullable().defaultTo(0);
    table.decimal('referredTotalVolume', null).notNullable().defaultTo(0);
    table.bigInteger('firstReferralBlockHeight').notNullable();
    table.decimal('totalReferredTakerFees', null).notNullable().defaultTo(0);
    table.decimal('totalReferredMakerFees', null).notNullable().defaultTo(0);
    table.decimal('totalReferredMakerRebates', null).notNullable().defaultTo(0);
    table.decimal('totalReferredLiquidationfees', null).notNullable().defaultTo(0);
  });
}

export async function down(knex: Knex): Promise<void> {
  return knex.schema.dropTable('affiliate_per_referee_stats');
}
