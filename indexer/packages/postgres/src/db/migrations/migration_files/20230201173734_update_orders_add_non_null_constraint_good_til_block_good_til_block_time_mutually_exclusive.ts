import * as Knex from 'knex';

// Add a constraint to the `orders` table to ensure that one of `goodTilBlock` and
// `goodTilBlockTime` is non-null for each row, but not both.
export async function up(knex: Knex): Promise<void> {
  return knex
    .schema
    .raw(`
      ALTER TABLE "orders" ADD CONSTRAINT good_til_block_or_good_til_block_time_non_null CHECK (num_nonnulls("goodTilBlock", "goodTilBlockTime") = 1);
    `);
}

export async function down(knex: Knex): Promise<void> {
  return knex
    .schema
    .raw(`
      ALTER TABLE "orders" DROP CONSTRAINT good_til_block_or_good_til_block_time_non_null;
    `);
}
