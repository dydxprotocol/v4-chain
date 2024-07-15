import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  await knex.schema.createTable('subaccount_usernames', (table) => {
    // username should be unique across the table
    table.string('username').notNullable().unique();
    table.string('subaccountId').notNullable().primary();
  });
}

export async function down(knex: Knex): Promise<void> {
  await knex.schema.dropTable('subaccount_usernames');
}
