import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  return knex
    .schema
    .createTable('trading_reward_aggregations', (table) => {
      table.uuid('id').primary();
      table.string('address').notNullable();
      table.timestamp('startedAt').notNullable();
      table.bigInteger('startedAtHeight').notNullable();
      table.timestamp('endedAt').nullable();
      table.bigInteger('endedAtHeight').nullable();
      table.enum(
        'period',
        [
          'DAILY',
          'WEEKLY',
          'MONTHLY',
        ],
      ).notNullable();
      table.decimal('amount').notNullable();

      // Foreign
      table.foreign('address').references('wallets.address');
      table.foreign('startedAtHeight').references('blocks.blockHeight');

      // Indices
      table.index(['address', 'startedAtHeight']);
      table.index(['period', 'startedAtHeight']);
    });
}

export async function down(knex: Knex): Promise<void> {
  return knex.schema.dropTableIfExists('trading_reward_aggregations');
}
