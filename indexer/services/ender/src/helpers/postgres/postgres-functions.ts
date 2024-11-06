import { readFileSync } from 'fs';
import path from 'path';

import { logger } from '@klyraprotocol-indexer/base';
import { dbHelpers, storeHelpers } from '@klyraprotocol-indexer/postgres';

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

const HANDLER_SCRIPTS: string[] = [
  'klyra_asset_create_handler.sql',
  'klyra_block_processor_ordered_handlers.sql',
  'klyra_block_processor_unordered_handlers.sql',
  'klyra_deleveraging_handler.sql',
  'klyra_funding_handler.sql',
  'klyra_liquidity_tier_handler.sql',
  'klyra_market_create_handler.sql',
  'klyra_market_modify_handler.sql',
  'klyra_market_price_update_handler.sql',
  'klyra_open_interest_update_handler.sql',
  'klyra_perpetual_market_v1_handler.sql',
  'klyra_perpetual_market_v2_handler.sql',
  'klyra_stateful_order_handler.sql',
  'klyra_subaccount_update_handler.sql',
  'klyra_transfer_handler.sql',
  'klyra_update_clob_pair_handler.sql',
  'klyra_update_perpetual_handler.sql',
  'klyra_yield_params_handler.sql',
];

const DB_SETUP_SCRIPTS: string[] = [
  'create_extension_pg_stat_statements.sql',
  'create_extension_uuid_ossp.sql',
];

const HELPER_SCRIPTS: string[] = [
  'klyra_clob_pair_status_to_market_status.sql',
  'klyra_create_initial_rows_for_tendermint_block.sql',
  'klyra_create_tendermint_event.sql',
  'klyra_create_transaction.sql',
  'klyra_event_id_from_parts.sql',
  'klyra_from_jsonlib_long.sql',
  'klyra_from_protocol_order_side.sql',
  'klyra_from_protocol_time_in_force.sql',
  'klyra_from_serializable_int.sql',
  'klyra_get_fee_from_liquidity.sql',
  'klyra_get_order_status.sql',
  'klyra_get_perpetual_market_for_clob_pair.sql',
  'klyra_get_total_filled_from_liquidity.sql',
  'klyra_get_weighted_average.sql',
  'klyra_liquidation_fill_handler_per_order.sql',
  'klyra_order_fill_handler_per_order.sql',
  'klyra_perpetual_position_and_order_side_matching.sql',
  'klyra_protocol_condition_type_to_order_type.sql',
  'klyra_tendermint_event_to_transaction_index.sql',
  'klyra_trim_scale.sql',
  'klyra_update_perpetual_position_aggregate_fields.sql',
  'klyra_uuid.sql',
  'klyra_uuid_from_asset_position_parts.sql',
  'klyra_uuid_from_fill_event_parts.sql',
  'klyra_uuid_from_funding_index_update_parts.sql',
  'klyra_uuid_from_oracle_price_parts.sql',
  'klyra_uuid_from_order_id.sql',
  'klyra_uuid_from_order_id_parts.sql',
  'klyra_uuid_from_perpetual_position_parts.sql',
  'klyra_uuid_from_subaccount_id.sql',
  'klyra_uuid_from_subaccount_id_parts.sql',
  'klyra_uuid_from_transaction_parts.sql',
  'klyra_uuid_from_transfer_parts.sql',
  'klyra_uuid_from_yield_params_parts.sql',
  'klyra_protocol_market_type_to_perpetual_market_type.sql',
];

const MAIN_SCRIPTS: string[] = [
  'klyra_block_processor.sql',
];

const SCRIPTS: string[] = [
  ...HANDLER_SCRIPTS.map((script: string) => `handlers/${script}`),
  ...HELPER_SCRIPTS.map((script: string) => `helpers/${script}`),
  ...DB_SETUP_SCRIPTS.map((script: string) => `setup/${script}`),
  ...MAIN_SCRIPTS,
];

export async function createPostgresFunctions(): Promise<void> {
  await Promise.all([
    dbHelpers.createModelToJsonFunctions(),
    ...SCRIPTS.map((script: string) => storeHelpers.rawQuery(newScript(script, `../../scripts/${script}`).script, {})
      .catch((error: Error) => {
        logger.error({
          at: 'postgres-functions#createPostgresFunctions',
          message: `Failed to create or replace function contained in ${script}`,
          error,
        });
        throw error;
      }),
    ),
  ]);
}
