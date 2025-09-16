import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  return knex.schema.createTable('affiliate_referee_stats', (table) => {
    table.string('refereeAddress').primary().notNullable();
    table.string('affiliateAddress').notNullable();
    table.decimal('affiliateEarnings', null).notNullable().defaultTo(0);
    table.integer('referredMakerTrades').notNullable().defaultTo(0);
    table.integer('referredTakerTrades').notNullable().defaultTo(0);
    table.decimal('referredTotalVolume', null).notNullable().defaultTo(0);
    table.bigInteger('referralBlockHeight').notNullable();
    table.decimal('referredTakerFees', null).notNullable().defaultTo(0);
    table.decimal('referredMakerFees', null).notNullable().defaultTo(0);
    table.decimal('referredMakerRebates', null).notNullable().defaultTo(0);
    table.decimal('referredLiquidationFees', null).notNullable().defaultTo(0);

    // Indices
    table.index(['affiliateAddress'], 'idx_affiliate_address');
  });
}

export async function down(knex: Knex): Promise<void> {
  return knex.schema.dropTable('affiliate_referee_stats');
}
