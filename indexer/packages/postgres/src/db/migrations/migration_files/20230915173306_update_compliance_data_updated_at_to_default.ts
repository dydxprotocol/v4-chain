import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  await knex.schema.alterTable('compliance_data', (table) => {
    table
      .timestamp('updatedAt').notNullable().defaultTo(knex.fn.now()).alter();
  });
}

export async function down(knex: Knex): Promise<void> {
  await knex.schema.alterTable('compliance_data', (table) => {
    table
      .timestamp('updatedAt').notNullable().alter();
  });
}
