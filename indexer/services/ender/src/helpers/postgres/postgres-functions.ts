import { readFileSync } from 'fs';
import path from 'path';

import { logger } from '@dydxprotocol-indexer/base';
import { dbHelpers, storeHelpers } from '@dydxprotocol-indexer/postgres';

export type PostgresFunction = {
  // The name of the script
  readonly name: string,
  // The contents of the script
  readonly script: string,
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
  'dydx_asset_create_handler.sql',
  'dydx_block_processor_ordered_handlers.sql',
  'dydx_block_processor_unordered_handlers.sql',
  'dydx_deleveraging_handler.sql',
  'dydx_funding_handler.sql',
  'dydx_liquidity_tier_handler.sql',
  'dydx_market_create_handler.sql',
  'dydx_market_modify_handler.sql',
  'dydx_market_price_update_handler.sql',
  'dydx_perpetual_market_v1_handler.sql',
  'dydx_perpetual_market_v2_handler.sql',
  'dydx_perpetual_market_v3_handler.sql',
  'dydx_register_affiliate_handler.sql',
  'dydx_stateful_order_handler.sql',
  'dydx_subaccount_update_handler.sql',
  'dydx_trading_rewards_handler.sql',
  'dydx_transfer_handler.sql',
  'dydx_update_clob_pair_handler.sql',
  'dydx_update_perpetual_v1_handler.sql',
  'dydx_update_perpetual_v2_handler.sql',
  'dydx_update_perpetual_v3_handler.sql',
  'dydx_vault_upsert_handler.sql',
];

const DB_SETUP_SCRIPTS: string[] = [
  'create_extension_pg_stat_statements.sql',
  'create_extension_uuid_ossp.sql',
];

const HELPER_SCRIPTS: string[] = [
  'dydx_clob_pair_status_to_market_status.sql',
  'dydx_create_initial_rows_for_tendermint_block.sql',
  'dydx_create_tendermint_event.sql',
  'dydx_create_transaction.sql',
  'dydx_event_id_from_parts.sql',
  'dydx_from_jsonlib_long.sql',
  'dydx_from_protocol_order_side.sql',
  'dydx_from_protocol_time_in_force.sql',
  'dydx_from_serializable_int.sql',
  'dydx_get_fee_from_liquidity.sql',
  'dydx_get_builder_fee_from_liquidity.sql',
  'dydx_get_builder_address_from_liquidity.sql',
  'dydx_get_order_status.sql',
  'dydx_get_perpetual_market_for_clob_pair.sql',
  'dydx_get_market_for_id.sql',
  'dydx_get_total_filled_from_liquidity.sql',
  'dydx_get_weighted_average.sql',
  'dydx_get_order_router_address_from_liquidity.sql',
  'dydx_get_order_router_fee_from_liquidity.sql',
  'dydx_liquidation_fill_handler_per_order.sql',
  'dydx_order_fill_handler_per_order.sql',
  'dydx_perpetual_position_and_order_side_matching.sql',
  'dydx_process_trading_reward_event.sql',
  'dydx_protocol_convert_to_order_type.sql',
  'dydx_tendermint_event_to_transaction_index.sql',
  'dydx_trim_scale.sql',
  'dydx_update_perpetual_position_aggregate_fields.sql',
  'dydx_uuid.sql',
  'dydx_uuid_from_asset_position_parts.sql',
  'dydx_uuid_from_fill_event_parts.sql',
  'dydx_uuid_from_funding_index_update_parts.sql',
  'dydx_uuid_from_oracle_price_parts.sql',
  'dydx_uuid_from_order_id.sql',
  'dydx_uuid_from_order_id_parts.sql',
  'dydx_uuid_from_perpetual_position_parts.sql',
  'dydx_uuid_from_subaccount_id.sql',
  'dydx_uuid_from_subaccount_id_parts.sql',
  'dydx_uuid_from_trading_rewards_parts.sql',
  'dydx_uuid_from_transaction_parts.sql',
  'dydx_uuid_from_transfer_parts.sql',
  'dydx_protocol_market_type_to_perpetual_market_type.sql',
  'dydx_protocol_vault_status_to_vault_status.sql',
  'dydx_order_flags.sql',
  'dydx_apply_fill_realized_effects.sql',
];

const MAIN_SCRIPTS: string[] = [
  'dydx_block_processor.sql',
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
