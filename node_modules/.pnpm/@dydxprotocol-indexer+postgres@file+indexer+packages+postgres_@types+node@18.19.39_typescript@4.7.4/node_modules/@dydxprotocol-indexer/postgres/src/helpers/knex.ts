import Knex from 'knex';
import { Model } from 'objection';

import config from '../config';
import {
  knexPrimaryConfigForEnv,
  knexReadReplicaConfigForEnv,
} from '../db/knexfile';

// TODO: add specific logger for the logging of knex
const logging = {
  log: {
    warn() {},
    error() {},
    deprecate() {},
    debug() {},
  },
};

class KnexReadReplica {
  private knexReadReplica: Knex | null = null;

  constructor(isUsingReadOnlyDB: boolean, knexConfig: Knex.Config) {
    if (isUsingReadOnlyDB) {
      this.knexReadReplica = Knex({
        ...knexConfig,
        ...logging,
      });
    }
  }

  public getConnection(): Knex {
    if (!this.knexReadReplica) {
      throw new Error('Service is not configured to use read only DB');
    }

    return this.knexReadReplica;
  }
}
export const knexPrimary: Knex = Knex({
  ...knexPrimaryConfigForEnv,
  ...logging,
});

export const knexReadReplica: KnexReadReplica = new KnexReadReplica(
  config.IS_USING_DB_READONLY,
  knexReadReplicaConfigForEnv,
);

// Bind all Models to the primary knex instance. You only
// need to do this once before you use any of
// your model classes.
Model.knex(knexPrimary);
