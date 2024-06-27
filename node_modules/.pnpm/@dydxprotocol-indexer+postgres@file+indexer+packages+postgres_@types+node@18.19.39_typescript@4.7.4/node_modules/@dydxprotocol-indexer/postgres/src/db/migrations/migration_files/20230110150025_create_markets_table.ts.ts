import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  return knex
    .schema
    .createTable('markets', (table) => {
      table.integer('id').primary();
      table.string('pair').unique().notNullable();
      table.integer('exponent').notNullable();
      table.integer('minPriceChangePpm').notNullable();
      table.decimal('oraclePrice', null).nullable();

      // Indices
      table.index(['pair']);
    });
}

export async function down(knex: Knex): Promise<void> {
  return knex.schema.dropTableIfExists('markets');
}
