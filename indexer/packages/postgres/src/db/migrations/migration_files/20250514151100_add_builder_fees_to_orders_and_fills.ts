import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  await knex.schema.alterTable('orders', (table) => {
    table.string('builderAddress').nullable().defaultTo(null);
    table.string('feePpm').nullable().defaultTo(null);
  });

  await knex.schema.alterTable('fills', (table) => {
    table.string('builderAddress').nullable().defaultTo(null);
    table.string('builderFee').nullable().defaultTo(null);
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
