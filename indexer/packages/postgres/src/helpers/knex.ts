import Knex from 'knex';
import { Model } from 'objection';

import config from '../config';
import {
  knexPrimaryConfigForEnv,
  knexReadReplicaConfigForEnvs,
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
  private knexReadReplicas: Knex[] = [];

  constructor(isUsingReadOnlyDB: boolean, knexConfigs: Knex.Config[]) {
    if (isUsingReadOnlyDB) {
      this.knexReadReplicas = knexConfigs.map((knexConfig: Knex.Config) => Knex({
        ...knexConfig,
        ...logging,
      }),
      );
    }
  }

  public getConnection(): Knex {
    if (this.knexReadReplicas.length === 0) {
      throw new Error('Service is not configured to use a read replica');
    }
    const randomIndex = Math.floor(Math.random() * this.knexReadReplicas.length);
    return this.knexReadReplicas[randomIndex];
  }
}
export const knexPrimary: Knex = Knex({
  ...knexPrimaryConfigForEnv,
  ...logging,
});

export const knexReadReplica: KnexReadReplica = new KnexReadReplica(
  config.IS_USING_DB_READONLY,
  knexReadReplicaConfigForEnvs,
);

// Bind all Models to the primary knex instance. You only
// need to do this once before you use any of
// your model classes.
Model.knex(knexPrimary);
