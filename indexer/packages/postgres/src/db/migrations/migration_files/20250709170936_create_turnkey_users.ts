import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  return knex
    .schema
    .createTable('turnkey_users', (table) => {
      table.text('suborg_id').primary(); // Primary key
      table.text('email').nullable().unique(); // optional but unique
      table.text('svm_address').notNullable().unique(); // indexed
      table.text('evm_address').notNullable().unique(); // indexed
      table.text('salt').notNullable(); // used to generate dydx keypair when combining the onboarding signature
      table.text('dydx_address').nullable().unique();
      table.timestamp('created_at').notNullable();
      table.text('smart_account_address').nullable().unique();

      // Indexes
      table.index(['svm_address'], 'idx_turnkey_users_svm_address');
      table.index(['evm_address'], 'idx_turnkey_users_evm_address');
      // used to validate that user is already a user.
      table.index(['email'], 'idx_turnkey_users_email');
    });
}

export async function down(knex: Knex): Promise<void> {
  return knex.schema.dropTableIfExists('turnkey_users');
}
