import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  return knex.schema.createTable('vaults', (table) => {
    table.string('address').primary().notNullable(); // address of vault
    table.bigInteger('clobPairId').notNullable(); // clob pair id for vault
    table.enum('status', [
      'DEACTIVATED',
      'STAND_BY',
      'QUOTING',
      'CLOSE_ONLY',
    ]).notNullable(); // quoting status of vault
    table.timestamp('createdAt').notNullable();
    table.timestamp('updatedAt').notNullable();
  });
}

export async function down(knex: Knex): Promise<void> {
  return knex.schema.dropTable('vaults');
}
