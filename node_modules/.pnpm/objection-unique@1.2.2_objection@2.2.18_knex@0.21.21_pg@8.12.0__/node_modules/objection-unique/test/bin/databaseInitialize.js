#!/usr/bin/env node

/**
 * Module dependencies.
 */

import knex from '../utils/knex';

/**
 * Exit process.
 */

const exit = process.exit;

/**
 * Initialize database.
 */

(async () => {
  // Drop `Test` table.
  await knex.schema.dropTableIfExists('Test');
  await knex.schema.dropTableIfExists('CompoundTest');

  // Create `Test` table.
  await knex.schema
   .createTableIfNotExists('Test', table => {
     table.increments('id').primary();
     table.string('foo').unique();
     table.string('bar').unique();
     table.string('biz');
   });

  // Create `CompoundTest` table.
  await knex.schema
    .createTableIfNotExists('CompoundTest', table => {
      table.increments('id').primary();
      table.string('foo');
      table.string('bar');
      table.string('biz');
      table.unique(['foo', 'bar']);
    });

  exit();
})();
