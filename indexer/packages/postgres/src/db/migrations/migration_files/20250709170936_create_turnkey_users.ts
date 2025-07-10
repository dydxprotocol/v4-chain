import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  return knex
    .schema
    .createTable('turnkey_users', (table) => {
      table.text('suborgId').primary(); // Primary key
      table.text('username').nullable(); // optional
      table.text('email').nullable(); // optional
      table.text('svmAddress').notNullable(); // indexed
      table.text('evmAddress').notNullable(); // indexed
      table.text('salt').notNullable(); // used to generate dydx keypair when combining the onboarding signature
      table.text('dydxAddress').nullable();
      table.timestamp('createdAt').notNullable();

      // Indexes
      table.index(['svmAddress'], 'idx_turnkey_users_svm_address');
      table.index(['evmAddress'], 'idx_turnkey_users_evm_address');
    });
}

export async function down(knex: Knex): Promise<void> {
  return knex.schema.dropTableIfExists('turnkey_users');
}
