import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  return knex
    .schema
    .createTable('oracle_prices', (table) => {
      table.uuid('id').primary();
      table.integer('marketId').notNullable();
      table.decimal('price', null).notNullable();
      table.timestamp('effectiveAt').notNullable();
      table.bigInteger('effectiveAtHeight').notNullable();

      // Foreign
      table.foreign('marketId').references('markets.id');
      table.foreign('effectiveAtHeight').references('blocks.blockHeight');

      // Indices
      table.index(['marketId', 'effectiveAt']);
      table.index(['marketId']);
    });
}

export async function down(knex: Knex): Promise<void> {
  return knex.schema.dropTableIfExists('oracle_prices');
}
