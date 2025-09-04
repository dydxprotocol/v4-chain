import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  await knex.schema.createTable('permission_approval', (table) => {
    table.text('suborg_id').notNullable();
    table.text('chain_id').notNullable();
    table.text('approval').notNullable();
    table.primary(['suborg_id', 'chain_id']);
  });
}

export async function down(knex: Knex): Promise<void> {
  await knex.schema.dropTable('permission_approval');
}
