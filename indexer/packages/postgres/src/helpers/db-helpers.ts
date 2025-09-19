import { logger } from '@dydxprotocol-indexer/base';

import config from '../config';
import { SQL_TO_JSON_DEFINED_MODELS } from '../constants';
import { knexPrimary, knexReadReplica } from './knex';
import { rawQuery } from './stores-helpers';

const layer2Tables = [
  'perpetual_positions',
  'fills',
  'leaderboard_pnl',
  'funding_payments',
  'subaccounts',
  'turnkey_users',
  'permission_approval',
  'bridge_information',
];

const layer1Tables = [
  'subaccount_usernames',
  'markets',
  'orders',
  'perpetual_markets',
  'tendermint_events',
  'transactions',
  'blocks',
  'assets',
  'candles',
  'liquidity_tiers',
  'wallets',
  'compliance_data',
  'trading_rewards',
  'trading_reward_aggregations',
  'compliance_status',
  'affiliate_referred_users',
  'persistent_cache',
  'affiliate_info',
  'vaults',
];

/**
 * Returns the SQL statement that would convert the provided field and type to the type expected
 * by the model.
 *
 * Raises an error if an unknown conversion is requested.
 */
function getSqlConversionForDydxModelTypes(fieldName: string, type: string): string {
  switch (type) {
    case 'integer':
      return `row_t."${fieldName}"::int`;
    case 'string':
      return `row_t."${fieldName}"::text`;
    case 'boolean':
      return `row_t."${fieldName}"::bool`;
    case 'hex-string':
      return `encode(row_t."${fieldName}", 'hex')`;
    case 'date-time':
      return `to_char(row_t."${fieldName}" at time zone 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS.MS"Z"')`;
    default:
      throw new Error(`Unknown type conversion for ${type}`);
  }
}

/**
 * Defines a `dydx_to_jsonb` function for each of the models in SQL_TO_JSON_DEFINED_MODELS and
 * loads them in Postgres. This allows for plpgsql functions to invoke `dydx_to_jsonb` on the
 * associated models table row type and convert the record into a JSON representation which
 * conforms to the models schema allowing conversion to the model type via the models `fromJson`
 * method.
 */
export async function createModelToJsonFunctions(): Promise<void> {
  await Promise.all(
    SQL_TO_JSON_DEFINED_MODELS.map(async (model) => {
      const sqlProperties: string[] = Object.entries(model.sqlToJsonConversions)
        .map(([key, value]) => `'${key}', ${getSqlConversionForDydxModelTypes(key, value)}`);
      const sqlFn: string = `CREATE OR REPLACE FUNCTION dydx_to_jsonb(row_t ${model.tableName}) RETURNS jsonb AS $$
BEGIN
    RETURN jsonb_build_object(
        ${sqlProperties.join(',\n        ')}
        );
    END;
$$ LANGUAGE plpgsql IMMUTABLE PARALLEL SAFE;`;
      return rawQuery(sqlFn, {}).catch((error) => {
        logger.error({
          at: 'dbHelpers#createModelToJsonFunctions',
          message: `Failed to create or replace function dydx_to_jsonb for model ${model.tableName}.`,
          error,
        });
        throw error;
      });
    }),
  );
}

async function dropData() {
  await Promise.all(
    layer2Tables.map(
      async (table) => {
        return knexPrimary(table).del();
      },
    ),
  );

  // need to use for... of to ensure tables are removed sequentially
  for (const table of layer1Tables) {
    await knexPrimary(table).del();
  }
}

/**
 * Drops all functions named dydx_.* from the database.
 */
async function dropAllDydxFunctions() {
  await knexPrimary.raw(`DO
$do$
DECLARE
   _sql text;
BEGIN
   SELECT INTO _sql
          string_agg(format('DROP %s %s;'
                          , CASE prokind
                              WHEN 'f' THEN 'FUNCTION'
                              WHEN 'a' THEN 'AGGREGATE'
                              WHEN 'p' THEN 'PROCEDURE'
                              WHEN 'w' THEN 'FUNCTION'  -- window function (rarely applicable)
                              -- ELSE NULL              -- not possible in pg 11
                            END
                          , oid::regprocedure)
                   , E'\\n')
   FROM   pg_proc
   WHERE  pronamespace = 'public'::regnamespace  -- schema name here!
   AND proname LIKE 'dydx_%';

   IF _sql IS NOT NULL THEN
      EXECUTE _sql;         -- uncomment payload once you are sure
   ELSE
      RAISE NOTICE 'No fuctions found in schema %', quote_ident('public');
   END IF;
END
$do$;`);
}

export async function clearData() {
  for (const table of layer1Tables) {
    const tableExists = await knexPrimary.schema.hasTable(table);
    if (tableExists) {
      await knexPrimary.raw(`truncate table "${table}" cascade`);
    }
  }

  await Promise.all(
    layer2Tables.map(async (table) => {
      const tableExists = await knexPrimary.schema.hasTable(table);
      if (tableExists) {
        return knexPrimary.raw(`truncate table "${table}" cascade`);
      }
    }),
  );
}

export async function clearSchema() {
  await knexPrimary.schema.raw('DROP SCHEMA public CASCADE');
  await knexPrimary.schema.raw('CREATE SCHEMA public');
}

export async function reset() {
  await dropData();
  await dropAllDydxFunctions();
  await rollback();
}

export async function rollback() {
  await knexPrimary.migrate.rollback({ loadExtensions: ['.js'] });
}

export async function migrate() {
  return knexPrimary.migrate.latest({ loadExtensions: ['.js'] });
}

export async function teardown() {
  await dropAllDydxFunctions();
  await knexPrimary.destroy();
  if (config.IS_USING_DB_READONLY) {
    await knexReadReplica.getConnection().destroy();
  }
}
