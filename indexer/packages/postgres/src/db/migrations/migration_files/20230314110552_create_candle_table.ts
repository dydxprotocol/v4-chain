import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  return knex
    .schema
    .createTable('candles', (table) => {
      table.uuid('id').primary();
      table.timestamp('startedAt').notNullable();
      table.string('ticker').notNullable();
      table.enum('resolution', [
        '1MIN',
        '5MINS',
        '15MINS',
        '30MINS',
        '1HOUR',
        '4HOURS',
        '1DAY',
      ]).notNullable();
      table.decimal('low', null).notNullable();
      table.decimal('high', null).notNullable();
      table.decimal('open', null).notNullable();
      table.decimal('close', null).notNullable();
      table.decimal('baseTokenVolume', null).notNullable();
      table.decimal('usdVolume', null).notNullable();
      table.integer('trades').notNullable();
      table.decimal('startingOpenInterest', null).notNullable();
    });
}

export async function down(knex: Knex): Promise<void> {
  return knex.schema.dropTableIfExists('candles');
}
