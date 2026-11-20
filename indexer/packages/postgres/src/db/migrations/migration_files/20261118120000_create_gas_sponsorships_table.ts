import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  return knex.schema.createTable('gas_sponsorships', (table) => {
    table.string('id').primary();
    table.string('chain_id').notNullable();
    table.string('granter_address').notNullable();
    table.string('grantee_address').notNullable();
    table.string('dydx_address').notNullable();
    // Both amounts are in the smallest unit of the chain's fee denom.
    table.decimal('spend_limit', null).notNullable();
    table.decimal('fee_amount', null).notNullable();
    table.string('status').notNullable();
    table.string('transaction_hash').nullable().unique();
    table.timestamp('expires_at').notNullable();
    table.timestamp('created_at').notNullable();

    table.index(
      ['chain_id', 'grantee_address', 'created_at'],
      'idx_gas_sponsorships_chain_grantee_created_at',
    );
    table.index(['chain_id', 'created_at'], 'idx_gas_sponsorships_chain_created_at');
  });
}

export async function down(knex: Knex): Promise<void> {
  return knex.schema.dropTableIfExists('gas_sponsorships');
}
