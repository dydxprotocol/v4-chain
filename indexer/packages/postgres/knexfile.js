require('ts-node/register');
require('dotenv-flow').config();

// TODO pull from knexfile.ts
module.exports = {
  migrations: {
    directory: './src/db/migrations/migration_files/',
  },
  client: 'pg',
  connection: {
    host: process.env.DB_HOSTNAME,
    port: parseInt(process.env.DB_PORT, 10),
    database: process.env.DB_NAME,
    user: process.env.DB_USERNAME,
    password: process.env.DB_PASSWORD,
  },
  pool: {
    min: parseInt(process.env.PG_POOL_MIN, 10),
    max: parseInt(process.env.PG_POOL_MAX, 10),
  },
};
