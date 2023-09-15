import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  return knex
    .schema
    .createTable('compliance_data', (table) => {
      table.string('address').notNullable();
      table.enum('provider', [
        'ELLIPTIC',
      ]).notNullable();
      table.string('chain').nullable();
      table.boolean('blocked').notNullable();
      table.decimal('riskScore').nullable();
      table.timestamp('updatedAt').notNullable();

      // Composite primary key
      table.primary(['address', 'provider']);

      // Index
      table.index(['updatedAt']);
      table.index(['blocked']);
    });
}

export async function down(knex: Knex): Promise<void> {
  return knex.schema.dropTableIfExists('compliance_data');
}
