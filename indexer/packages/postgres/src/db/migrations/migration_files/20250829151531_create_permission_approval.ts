import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  await knex.schema.createTable('permission_approval', (table) => {
    table.text('suborg_id').primary();
    table.text('arbitrum_approval').nullable();
    table.text('base_approval').nullable();
    table.text('avalanche_approval').nullable();
    table.text('optimism_approval').nullable();
    table.text('ethereum_approval').nullable();
  });
}

export async function down(knex: Knex): Promise<void> {
  await knex.schema.dropTable('permission_approval');
}
