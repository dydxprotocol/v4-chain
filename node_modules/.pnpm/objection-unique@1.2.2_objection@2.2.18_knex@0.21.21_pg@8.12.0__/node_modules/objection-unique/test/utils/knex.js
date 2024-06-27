
/**
 * Module dependencies.
 */

import knex from 'knex';

// Knex configuration.
const configuration = {
  client: 'sqlite3',
  connection: {
    filename: './test.db'
  },
  useNullAsDefault: true
};

/**
 * Export `knex`.
 */

export default knex(configuration);
