import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  return knex.schema.createTable('affiliate_referred_users', (table) => {
    table.string('refereeAddress').primary().notNullable();
    table.string('affiliateAddress').notNullable();
    table.bigInteger('referredAtBlock').notNullable();

    // Index on affiliateAddress for faster queries
    table.index(['affiliateAddress']);
  });
}

export async function down(knex: Knex): Promise<void> {
  return knex.schema.dropTable('affiliate_referred_users');
}
