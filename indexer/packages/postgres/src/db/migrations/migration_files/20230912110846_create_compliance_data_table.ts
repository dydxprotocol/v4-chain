import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  return knex
    .schema
    .createTable('compliance_data', (table) => {
      table.string('address').primary();
      table.string('chain').nullable();
      table.boolean('sanctioned').notNullable();
      table.decimal('riskScore').nullable();
      table.timestamp('updatedAt').notNullable();

      // Index
      table.index(['updatedAt']);
    });
}

export async function down(knex: Knex): Promise<void> {
  return knex.schema.dropTableIfExists('compliance_data');
}
