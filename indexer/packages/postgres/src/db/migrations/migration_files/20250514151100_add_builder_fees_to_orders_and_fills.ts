import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  await knex.schema.alterTable('orders', (table) => {
    table.string('builderAddress').nullable();
    table.string('feePpm').nullable();
  });

  await knex.schema.alterTable('fills', (table) => {
    table.string('builderAddress').nullable();
    table.string('builderFee').nullable();
  });
}

export async function down(knex: Knex): Promise<void> {
  await knex.schema.alterTable('orders', (table) => {
    table.dropColumn('builderAddress');
    table.dropColumn('feePpm');
  });

  await knex.schema.alterTable('fills', (table) => {
    table.dropColumn('builderAddress');
    table.dropColumn('builderFee');
  });
}
