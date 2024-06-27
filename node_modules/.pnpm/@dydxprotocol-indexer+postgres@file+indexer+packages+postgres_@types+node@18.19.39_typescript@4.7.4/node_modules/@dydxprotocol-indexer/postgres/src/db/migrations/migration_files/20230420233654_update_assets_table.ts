import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  return knex
    .schema
    .alterTable('assets', (table) => {
      table.dropColumn('denom');
      table.string('symbol').unique().notNullable();

      // Indices
      table.index(['symbol']);
    });
}

export async function down(knex: Knex): Promise<void> {
  return knex.schema.alterTable('assets', (table) => {
    table.dropColumn('symbol');
    table.string('denom').unique().notNullable();

    // Indices
    table.index(['denom']);
  });
}
