import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  await knex.schema.createTable('subaccount_usernames', (table) => {
    // username is primary key and is unique across the table
    table.string('username').notNullable().primary();
    // subaccounts is a foreign key to the subaccounts table subaccounts.id
    table.uuid('subaccountId').notNullable().references('id').inTable('subaccounts');
  });
}

export async function down(knex: Knex): Promise<void> {
  await knex.schema.dropTable('subaccount_usernames');
}
