import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  return knex.schema.createTable('affiliate_info', (table) => {
    table.string('address').primary().notNullable();
    table.decimal('affiliateEarnings').notNullable();
    table.integer('referredMakerTrades').notNullable();
    table.integer('referredTakerTrades').notNullable();
    table.decimal('totalReferredFees').notNullable();
    table.integer('totalReferredUsers').notNullable();
    table.decimal('referredNetProtocolEarnings').notNullable();
    table.bigInteger('firstReferralBlockHeight').notNullable();
  });
}

export async function down(knex: Knex): Promise<void> {
  return knex.schema.dropTable('affiliate_info');
}
