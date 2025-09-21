import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  return knex.schema.createTable('bridge_information', (table) => {
    table.string('id').primary();
    table.string('from_address').notNullable();
    table.string('chain_id').notNullable();
    table.string('amount').notNullable();
    table.string('transaction_hash').nullable().unique();
    table.timestamp('created_at').notNullable();

    // Index
    table.index('from_address', 'idx_bridge_information_from_address');
  });
}

export async function down(knex: Knex): Promise<void> {
  return knex.schema.dropTableIfExists('bridge_information');
}
