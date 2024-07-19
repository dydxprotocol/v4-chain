import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  return knex
    .schema
    .createTable('assets', (table) => {
      table.string('id').primary();
      table.string('denom').unique().notNullable();
      table.integer('atomicResolution').notNullable();
      table.boolean('hasMarket').notNullable();
      table.integer('marketId').nullable();

      // Indices
      table.index(['denom']);
      table.index(['marketId']);
    });
}

export async function down(knex: Knex): Promise<void> {
  return knex.schema.dropTableIfExists('assets');
}
