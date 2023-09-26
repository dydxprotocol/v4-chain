import { readFileSync } from 'fs';
import path from 'path';

import { logger } from '@dydxprotocol-indexer/base';
import { dbHelpers, storeHelpers } from '@dydxprotocol-indexer/postgres';

export type PostgresFunction = {
  // The name of the script
  readonly name: string;
  // The contents of the script
  readonly script: string;
};

/**
 * Loads a named script from the specified path.
 *
 * @param name The name of the script.
 * @param scriptPath The path to the script.
 * @returns The created script object
 */
function newScript(name: string, scriptPath: string): PostgresFunction {
  const script: string = readFileSync(path.resolve(__dirname, scriptPath)).toString();
  return {
    name,
    script,
  };
}

const scripts: string[] = [
  'create_extension_pg_stat_statements.sql',
  'create_extension_uuid_ossp.sql',
  'dydx_event_id_from_parts.sql',
  'dydx_event_to_transaction_index.sql',
  'dydx_from_jsonlib_long.sql',
  'dydx_from_protocol_order_side.sql',
  'dydx_from_protocol_time_in_force.sql',
  'dydx_from_serializable_int.sql',
  'dydx_get_fee_from_liquidity.sql',
  'dydx_get_order_status.sql',
  'dydx_get_total_filled_from_liquidity.sql',
  'dydx_get_weighted_average.sql',
  'dydx_order_fill_handler_per_order.sql',
  'dydx_perpetual_position_and_order_side_matching.sql',
  'dydx_subaccount_update_handler.sql',
  'dydx_trim_scale.sql',
  'dydx_uuid.sql',
  'dydx_uuid_from_asset_position_parts.sql',
  'dydx_uuid_from_fill_event_parts.sql',
  'dydx_uuid_from_order_id.sql',
  'dydx_uuid_from_order_id_parts.sql',
  'dydx_uuid_from_perpetual_position_parts.sql',
  'dydx_uuid_from_subaccount_id.sql',
  'dydx_uuid_from_subaccount_id_parts.sql',
  'dydx_uuid_from_transaction_parts.sql',
  'dydx_create_transaction.sql',
  'dydx_create_initial_rows_for_tendermint_block.sql',
  'dydx_create_tendermint_event.sql',
  'dydx_tendermint_event_to_transaction_index.sql',
];

export async function createPostgresFunctions(): Promise<void> {
  await Promise.all([
    dbHelpers.createModelToJsonFunctions(),
    ...scripts.map((script: string) => storeHelpers.rawQuery(newScript(script, `../../scripts/${script}`).script, {})
      .catch((error) => {
        logger.error({
          at: 'dbHelpers#createModelToJsonFunctions',
          message: `Failed to create or replace function contained in ${script}`,
          error,
        });
        throw error;
      }),
    ),
  ]);
}
