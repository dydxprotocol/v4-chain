import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  await knex.schema.createTable('subaccount_usernames', (table) => {
    // username should be unique across the table
    table.string('username').notNullable().unique();
    // subaccounts is a foreign key to the subaccounts table subaccounts.id
    table.uuid('subaccountId').notNullable().primary().references('id').inTable('subaccounts');
  });
}

export async function down(knex: Knex): Promise<void> {
  await knex.schema.dropTable('subaccount_usernames');
}
