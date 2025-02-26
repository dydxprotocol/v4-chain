import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  return knex.schema.createTable('funding_payments', (table) => {
    table.string('address').primary().notNullable(); // address of vault
    table.bigInteger('clobPairId').notNullable(); // clob pair id for vault
    table.bigInteger('size').notNullable(); // size of position
    table.decimal('rate').notNullable(); // rate funding was paid at
    table.decimal('amount').notNullable(); // amount paid
    table.boolean('isLong').notNullable(); // position is Long
    table.timestamp('createdAt').notNullable();
  });
}

export async function down(knex: Knex): Promise<void> {
  return knex.schema.dropTable('funding_payments');
}
