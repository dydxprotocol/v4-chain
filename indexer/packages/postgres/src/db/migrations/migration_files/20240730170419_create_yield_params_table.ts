import Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  return knex
    .schema
    .createTable('yield_params', (table) => {
      table.uuid("id").primary()
      table.string("sDAIPrice").notNullable()
      table.string("assetYieldIndex").notNullable()
      table.timestamp('createdAt').notNullable();
      table.bigInteger('createdAtHeight').notNullable();
    });
}

export async function down(knex: Knex): Promise<void> {
  return knex.schema.dropTableIfExists('yield_params');
}
