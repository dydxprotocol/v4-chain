import Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  return knex
    .schema
    .createTable('asset_positions', (table) => {
      table.uuid('id').primary();
      table.string('assetId').notNullable();
      table.uuid('subaccountId').notNullable();
      table.decimal('size', null).notNullable();
      table.boolean('isLong').notNullable();

      // Foreign
      table.foreign('assetId').references('assets.id');
      table.foreign('subaccountId').references('subaccounts.id');

      // Indices
      table.index(['subaccountId']);
      table.index(['subaccountId', 'assetId']);
    });
}

export async function down(knex: Knex): Promise<void> {
  return knex.schema.dropTableIfExists('asset_positions');
}
