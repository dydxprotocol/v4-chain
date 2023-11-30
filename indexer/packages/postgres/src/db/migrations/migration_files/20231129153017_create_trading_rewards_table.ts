import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  return knex
    .schema
    .createTable('trading_rewards', (table) => {
      table.uuid('id').primary();
      table.string('address').notNullable();
      table.timestamp('blockTime').notNullable();
      table.bigInteger('blockHeight').notNullable();
      table.decimal('amount').notNullable();

      // Foreign
      table.foreign('address').references('wallets.address');

      // Indices
      table.index(['address']);
      table.index(['blockTime']);
      table.index(['blockHeight']);
    });
}

export async function down(knex: Knex): Promise<void> {
  return knex.schema.dropTableIfExists('trading_rewards');
}
