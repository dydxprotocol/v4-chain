import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  return knex.schema.createTable('tokens', (table) => {
    table.increments('id').primary();
    table.string('token').notNullable().unique();
    table.string('address').notNullable();
    table.foreign('address').references('wallets.address').onDelete('CASCADE');
    table.timestamp('updatedAt').notNullable().defaultTo(knex.fn.now());
  });
}

export async function down(knex: Knex): Promise<void> {
  return knex.schema.dropTable('tokens');
}
