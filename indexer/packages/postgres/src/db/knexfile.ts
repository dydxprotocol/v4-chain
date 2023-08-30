import {
  NodeEnv,
} from '@dydxprotocol-indexer/base';
import Knex from 'knex';

import config from '../config';
import './pg-config';

const SUPPORTED_ENVIRONMENTS: string[] = Object.values(NodeEnv);
const environment: string = config.NODE_ENV || NodeEnv.DEVELOPMENT;

if (!SUPPORTED_ENVIRONMENTS.includes(environment)) {
  throw new Error(`Unknown node environment ${environment}`);
}

export function getConfigForHost(host: string) : Knex.Config {
  return {
    client: 'pg',
    connection: {
      host,
      port: config.DB_PORT,
      database: config.DB_NAME,
      user: config.DB_USERNAME,
      password: config.DB_PASSWORD,
    },
    pool: {
      min: config.PG_POOL_MIN,
      max: config.PG_POOL_MAX,
      acquireTimeoutMillis: config.PG_ACQUIRE_CONNECTION_TIMEOUT_MS,
      createTimeoutMillis: config.PG_ACQUIRE_CONNECTION_TIMEOUT_MS,
    },
    acquireConnectionTimeout: config.PG_ACQUIRE_CONNECTION_TIMEOUT_MS,
    migrations: {
      directory: `${__dirname}/migrations/migration_files/`,
    },
  };
}

export const knexPrimaryConfigForEnv: Knex.Config = getConfigForHost(config.DB_HOSTNAME);

export const knexReadReplicaConfigForEnv:
Knex.Config = getConfigForHost(config.DB_READONLY_HOSTNAME);
