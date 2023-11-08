/**
 * Environment variables required by Ender.
 */

import {
  parseSchema,
  baseConfigSchema,
  parseBoolean,
} from '@dydxprotocol-indexer/base';
import {
  kafkaConfigSchema,
} from '@dydxprotocol-indexer/kafka';
import {
  postgresConfigSchema,
} from '@dydxprotocol-indexer/postgres';
import { redisConfigSchema } from '@dydxprotocol-indexer/redis';

export const configSchema = {
  ...baseConfigSchema,
  ...postgresConfigSchema,
  ...redisConfigSchema,
  ...kafkaConfigSchema,
  SEND_WEBSOCKET_MESSAGES: parseBoolean({
    default: true,
  }),
  USE_ASSET_CREATE_HANDLER_SQL_FUNCTION: parseBoolean({
    default: true,
  }),
  USE_FUNDING_HANDLER_SQL_FUNCTION: parseBoolean({
    default: true,
  }),
  USE_LIQUIDATION_HANDLER_SQL_FUNCTION: parseBoolean({
    default: true,
  }),
  USE_LIQUIDITY_TIER_HANDLER_SQL_FUNCTION: parseBoolean({
    default: true,
  }),
  USE_MARKET_CREATE_HANDLER_SQL_FUNCTION: parseBoolean({
    default: true,
  }),
  USE_MARKET_MODIFY_HANDLER_SQL_FUNCTION: parseBoolean({
    default: true,
  }),
  USE_MARKET_PRICE_UPDATE_HANDLER_SQL_FUNCTION: parseBoolean({
    default: true,
  }),
  USE_ORDER_HANDLER_SQL_FUNCTION: parseBoolean({
    default: true,
  }),
  USE_PERPETUAL_MARKET_HANDLER_SQL_FUNCTION: parseBoolean({
    default: true,
  }),
  USE_STATEFUL_ORDER_HANDLER_SQL_FUNCTION: parseBoolean({
    default: true,
  }),
  USE_SUBACCOUNT_UPDATE_SQL_FUNCTION: parseBoolean({
    default: true,
  }),
  USE_TRANSFER_HANDLER_SQL_FUNCTION: parseBoolean({
    default: true,
  }),
  USE_UPDATE_CLOB_PAIR_HANDLER_SQL_FUNCTION: parseBoolean({
    default: true,
  }),
  USE_UPDATE_PERPETUAL_HANDLER_SQL_FUNCTION: parseBoolean({
    default: true,
  }),
  USE_SQL_FUNCTION_TO_CREATE_INITIAL_ROWS: parseBoolean({
    default: true,
  }),
};

export default parseSchema(configSchema);
