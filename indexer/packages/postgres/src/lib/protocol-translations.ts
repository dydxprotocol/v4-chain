import {
  bytesToBigInt,
  ORDER_FLAG_CONDITIONAL,
  ORDER_FLAG_LONG_TERM,
  ORDER_FLAG_SHORT_TERM,
  ORDER_FLAG_TWAP,
  ORDER_FLAG_TWAP_SUBORDER,
} from '@dydxprotocol-indexer/v4-proto-parser';
import {
  ClobPairStatus,
  IndexerOrder,
  IndexerOrder_ConditionType,
  IndexerOrder_Side,
  IndexerOrder_TimeInForce,
} from '@dydxprotocol-indexer/v4-protos';
import Big from 'big.js';
import { DateTime } from 'luxon';

import {
  CLOB_STATUS_TO_MARKET_STATUS,
  FUNDING_RATE_FROM_PROTOCOL_IN_HOURS,
  PPM_EXPONENT,
  QUOTE_CURRENCY_ATOMIC_RESOLUTION,
} from '../constants';
import {
  IsoString, OrderSide, OrderType, PerpetualMarketFromDatabase, PerpetualMarketStatus, TimeInForce,
} from '../types';
import { InvalidClobPairStatusError } from './errors';

// Mapping from the TimeInForce enum from the protocol to the TimeInForce enum in the Indexer
const PROTOCOL_TIF_TO_INDEXER_TIF_MAP: Record<IndexerOrder_TimeInForce, TimeInForce> = {
  [IndexerOrder_TimeInForce.TIME_IN_FORCE_FILL_OR_KILL]: TimeInForce.FOK,
  [IndexerOrder_TimeInForce.TIME_IN_FORCE_IOC]: TimeInForce.IOC,
  // Default behavior with UNSPECIFIED = GTT (Good-Til-Time)
  [IndexerOrder_TimeInForce.TIME_IN_FORCE_UNSPECIFIED]: TimeInForce.GTT,
  [IndexerOrder_TimeInForce.UNRECOGNIZED]: TimeInForce.GTT,
  [IndexerOrder_TimeInForce.TIME_IN_FORCE_POST_ONLY]: TimeInForce.POST_ONLY,
};

// Reverse mapping of above
const INDEXER_TIF_TO_PROTOCOL_TIF_MAP: Record<TimeInForce, IndexerOrder_TimeInForce> = {
  [TimeInForce.FOK]: IndexerOrder_TimeInForce.TIME_IN_FORCE_FILL_OR_KILL,
  [TimeInForce.IOC]: IndexerOrder_TimeInForce.TIME_IN_FORCE_IOC,
  [TimeInForce.GTT]: IndexerOrder_TimeInForce.TIME_IN_FORCE_UNSPECIFIED,
  [TimeInForce.POST_ONLY]: IndexerOrder_TimeInForce.TIME_IN_FORCE_POST_ONLY,
};

// Mapping from Condition type enum from protocol to OrderType enum in the Indexer
const CONDITION_TYPE_TO_ORDER_TYPE_MAP: Record<IndexerOrder_ConditionType, OrderType> = {
  // Default behavior with UNSPECIFIED / UNRECOGNIZED = Limit Order
  [IndexerOrder_ConditionType.CONDITION_TYPE_UNSPECIFIED]: OrderType.LIMIT,
  [IndexerOrder_ConditionType.UNRECOGNIZED]: OrderType.LIMIT,
  [IndexerOrder_ConditionType.CONDITION_TYPE_STOP_LOSS]: OrderType.STOP_LIMIT,
  [IndexerOrder_ConditionType.CONDITION_TYPE_TAKE_PROFIT]: OrderType.TAKE_PROFIT,
};

// Reverse mapping of above
const ORDER_TYPE_TO_CONDITION_TYPE_MAP: Record<OrderType, IndexerOrder_ConditionType> = {
  // Limit orders are unspecified
  [OrderType.LIMIT]: IndexerOrder_ConditionType.CONDITION_TYPE_UNSPECIFIED,

  // Only STOP_LIMIT is used currently
  [OrderType.STOP_LIMIT]: IndexerOrder_ConditionType.CONDITION_TYPE_STOP_LOSS,
  [OrderType.STOP_MARKET]: IndexerOrder_ConditionType.CONDITION_TYPE_STOP_LOSS,

  // Only TAKE_PROFIT is used currently
  [OrderType.TAKE_PROFIT]: IndexerOrder_ConditionType.CONDITION_TYPE_TAKE_PROFIT,
  [OrderType.TAKE_PROFIT_MARKET]: IndexerOrder_ConditionType.CONDITION_TYPE_TAKE_PROFIT,

  // TODO(IND-356): Remove irrelevant order types
  // Unused order types
  [OrderType.MARKET]: IndexerOrder_ConditionType.CONDITION_TYPE_UNSPECIFIED,
  [OrderType.TRAILING_STOP]: IndexerOrder_ConditionType.CONDITION_TYPE_UNSPECIFIED,

  [OrderType.TWAP]: IndexerOrder_ConditionType.CONDITION_TYPE_UNSPECIFIED,
  [OrderType.TWAP_SUBORDER]: IndexerOrder_ConditionType.CONDITION_TYPE_UNSPECIFIED,
};

/**
 * @param subticks - quote quantums/base quantums e.g. (1e-14 USDC/1e-10 BTC)
 * @returns - quote currency / base currency (human readable price)
 */
export function subticksToPrice(
  subticks: string,
  perpetualMarket: PerpetualMarketFromDatabase,
): string {
  return Big(subticks)
    .times(Big(10).pow(perpetualMarket.quantumConversionExponent))
    .times(Big(10).pow(QUOTE_CURRENCY_ATOMIC_RESOLUTION))
    .div(Big(10).pow(perpetualMarket.atomicResolution))
    .toFixed();
}

/**
 * @param price - quote currency / base currency (human readable price)
 * @returns - quote quantums/base quantums e.g. (1e-14 USDC/1e-10 BTC)
 */
export function priceToSubticks(
  price: string,
  perpetualMarket: PerpetualMarketFromDatabase,
): string {
  return Big(price)
    .times(Big(10).pow(perpetualMarket.atomicResolution))
    .div(Big(10).pow(QUOTE_CURRENCY_ATOMIC_RESOLUTION))
    .div(Big(10).pow(perpetualMarket.quantumConversionExponent))
    .toFixed();
}

/**
 * @param perpetualMarket
 * @returns - tick size for a given perpetual market. Tick size is the minimum price movement on
 * the market in quote human.
 */
export function getTickSize(
  perpetualMarket: PerpetualMarketFromDatabase,
): string {
  return Big(perpetualMarket.subticksPerTick)
    .times(Big(10).pow(perpetualMarket.quantumConversionExponent))
    .times(Big(10).pow(QUOTE_CURRENCY_ATOMIC_RESOLUTION))
    .div(Big(10).pow(perpetualMarket.atomicResolution))
    .toFixed();
}

/**
 * Given a funding index from the protocol, convert it to a human readable units such that when
 * multiplied by the position size in human-readable units of the base currency (e.g. 1 ETH, 1 BTC)
 * results in the funding payment in human-readable units of the quote currency (e.g. 1.2 USDC).
 * The funding index from the protocol is in parts-per-million, and of the units
 * quote-quantums / base quantums.
 * Formula for funding payment:
 * F_Q = funding payment in quote quantums
 * F_Q_H = funding payment in human-readable units of the quote currency
 * R_Q = atomic resolution of quote currency
 * P_B = position size in base quantums
 * P_B_H = position size in human-readable units of the base currency
 * R_B = atomic resolution of base currency
 * FI = funding index from protocol
 * FI_H = funding index in human-readable units, such that when multipled by the position size in
 * human-readable units of the base currency, returns the funding payment in human-readable units of
 * the quote currency.
 *
 * F_Q_H = F_Q * (10 ^ R_Q) -> formula to get funding payment in human-readable units
 * P_B_H = P_B * (10 ^ R_B) -> formula to get position size in human-readable units
 * F_Q = P_B * (FI / 10 ^ 6) -> formula to get funding payment from funding index and position size
 * F_Q_H = P_B_H * (FI_H) -> desired formula with funding index in human-readable units
 *
 * Note: in F_Q = P_B * (FI * 10 ^ 6), we divide by 10^6 as funding index from protocol is in
 * parts-per-million.
 * Solving for FI_H, we get
 * FI_H = F_Q_H / P_B_H
 * FI_H = (F_Q * 10 ^ R_Q) / (P_B * 10 ^ R_B)
 * FI_H = ((P_B * (FI / 10 ^ 6)) *  10 ^ R_Q) / (P_B * 10 ^ R_B)
 * FI_H = ((FI / 10 ^ 6) * 10 ^ R_Q) / (10 ^ R_B)
 * FI_H = (FI * 10 ^ -6 * 10 ^ R_Q) / (10 ^ R_B) -> formula for human-readable funding index using
 * resolutions of base and quote currency along with the funding-index from the protocol
 * @param fundingIndex
 * @param perpetualMarket
 */
export function fundingIndexToHumanFixedString(
  fundingIndex: string,
  perpetualMarket: PerpetualMarketFromDatabase,
): string {
  return Big(fundingIndex)
    .times(Big(10).pow(PPM_EXPONENT)) // PPM = parts-per-million
    .times(Big(10).pow(QUOTE_CURRENCY_ATOMIC_RESOLUTION))
    .div(Big(10).pow(perpetualMarket.atomicResolution))
    .toFixed();
}

/**
 * Given a funding value in parts-per-million, convert it to a funding index in human-readable
 * units. This formula requires dividing ppm by 1^6 and also dividing by 8 as funding is returned
 * in 8 hour and we store funding rate in 1 hour increments.
 */
export function funding8HourValuePpmTo1HourRate(
  fundingValuePpm: number,
): string {
  return Big(fundingValuePpm)
    .times(Big(10).pow(PPM_EXPONENT))
    .div(FUNDING_RATE_FROM_PROTOCOL_IN_HOURS)
    .toFixed();
}

/**
 * Returns the step size for the given perpetual market.
 * Step size is the smallest factor allowed for order amounts on the market in human.
 * @param perpetualMarket
 * @returns
 */
export function getStepSize(
  perpetualMarket: PerpetualMarketFromDatabase,
): string {
  return Big(perpetualMarket.stepBaseQuantums)
    .times(Big(10).pow(perpetualMarket.atomicResolution))
    .toFixed();
}

export function quantumsToHumanFixedString(
  baseQuantums: string,
  atomicResolution: number,
): string {
  return quantumsToHuman(
    baseQuantums,
    atomicResolution,
  ).toFixed();
}

/**
 * Returns the absolute human-readable size in string form, from the serialized
 * quantums buffer received from the protocol. This is a temporary function until
 * DEC-1597 (deprecate `isLong`) is completed.
 * @param serializedQuantums The serialized quantums buffer from the protocol
 * @param atomicResolution atomic resolution of the market
 * @returns absolute human-readable size
 */
export function serializedQuantumsToAbsHumanFixedString(
  serializedQuantums: Uint8Array,
  atomicResolution: number,
): string {
  return quantumsToHuman(
    bytesToBigInt(serializedQuantums).toString(),
    atomicResolution,
  ).abs().toFixed();
}

/**
 * @param quantums - the smallest increment of position size, and is determined by atomicResolution.
 * For example, an atomicResolution of 8 means the smallest increment of position size is 1e-8.
 * @returns - human readable position size
 */
export function quantumsToHuman(
  quantums: string,
  atomicResolution: number,
): Big {
  return Big(quantums)
    .times(Big(10).pow(atomicResolution));
}

export function humanToQuantums(
  human: string,
  atomicResolution: number,
): Big {
  return Big(human)
    .div(Big(10).pow(atomicResolution));
}

/**
 * Converts a price from the `Price` module in the V4 protocol to a human readable price.
 * @param protocolPrice Price value from the `Price` module.
 * @param exponent Exponent for the price from the `Price` module.
 * @returns Human readable price as a string.
 */
export function protocolPriceToHuman(
  protocolPrice: string,
  exponent: number,
): string {
  return Big(protocolPrice)
    .times(Big(10).pow(exponent))
    .toFixed();
}

/**
 * Converts the `Order_Side` enum from the protobuf to the `OrderSide` enum in postgres
 * @param protocolOrderSide `IndexerOrder_Side` enum from protobuf
 * @returns `OrderSide` corresponding to the `Order_Side` passed in
 */
export function protocolOrderSideToOrderSide(
  protocolOrderSide: IndexerOrder_Side,
): OrderSide {
  return protocolOrderSide === IndexerOrder_Side.SIDE_BUY ? OrderSide.BUY : OrderSide.SELL;
}

/**
 * Converts the TimeInForce field from an IndexerOrder proto to a TimeInForce enum in the Indexer.
 * Special cased:
 * - UNSPECIFIED -> GTT
 * Throw an error if the input TimeInForce enum is not in the known enum values for TimeInForce.
 * @param protocolOrderTIF
 * @returns
 */
export function protocolOrderTIFToTIF(
  protocolOrderTIF: IndexerOrder_TimeInForce,
): TimeInForce {
  if (!(protocolOrderTIF in PROTOCOL_TIF_TO_INDEXER_TIF_MAP)) {
    throw new Error(`Unexpected TimeInForce from protocol: ${protocolOrderTIF}`);
  }

  return PROTOCOL_TIF_TO_INDEXER_TIF_MAP[protocolOrderTIF];
}

/**
 * Converts TimeInForce enum in the Indexer to the TimeInForce field from an IndexerOrder proto.
 * GTT -> UNSPECIFIED.
 * Throw an error if the input TimeInForce enum is not in the known enum values for TimeInForce.
 * @param timeInForce
 */
export function tifToProtocolOrderTIF(timeInForce: TimeInForce): IndexerOrder_TimeInForce {
  if (!(timeInForce in INDEXER_TIF_TO_PROTOCOL_TIF_MAP)) {
    throw new Error(`Unexpected TimeInForce: ${timeInForce}`);
  }

  return INDEXER_TIF_TO_PROTOCOL_TIF_MAP[timeInForce];
}

// Gets `goodTilBlock` from an `IndexerOrder`, undefined if it does not exist.
export function getGoodTilBlock(order: IndexerOrder): number | undefined {
  return order.goodTilBlock;
}

// Gets `goodTilBlockTime` from an `Order` as an ISO string, undefined if it does not exist.
export function getGoodTilBlockTime(order: IndexerOrder): IsoString | undefined {
  if (order.goodTilBlockTime !== undefined) {
    // `goodTilBlockTime` is the unix timestamp in seconds
    // Reference:
    // https://github.com/dydxprotocol/v4/blob/main/proto/dydxprotocol/clob/order.proto#L138-L144
    return DateTime.fromSeconds(order.goodTilBlockTime).toUTC().toISO();
  }
  return undefined;
}

/**
 * Converts ConditionType enum from an IndexerOrder proto to an OrderType in the Indexer.
 * Special cased:
 * - UNSPECIFIED -> LIMIT
 * Throw an error if the input ConditionType is not in the known enum values for ConditionType.
 * @param protocolConditionType
 * @returns
 */
export function protocolConditionTypeToOrderType(
  protocolConditionType: IndexerOrder_ConditionType,
  orderFlag: number = 32,
): OrderType {
  if (!(protocolConditionType in CONDITION_TYPE_TO_ORDER_TYPE_MAP)) {
    throw new Error(`Unexpected ConditionType: ${protocolConditionType}`);
  }

  switch (orderFlag) {
    case ORDER_FLAG_SHORT_TERM:
      return OrderType.LIMIT;
    case ORDER_FLAG_CONDITIONAL:
      switch (protocolConditionType) {
        case IndexerOrder_ConditionType.UNRECOGNIZED:
        case IndexerOrder_ConditionType.CONDITION_TYPE_UNSPECIFIED:
          return OrderType.LIMIT;
        case IndexerOrder_ConditionType.CONDITION_TYPE_STOP_LOSS:
          return OrderType.STOP_LIMIT;
        case IndexerOrder_ConditionType.CONDITION_TYPE_TAKE_PROFIT:
          return OrderType.TAKE_PROFIT;
        default:
          throw new Error(`Unexpected ConditionType: ${protocolConditionType}`);
      }
    case ORDER_FLAG_LONG_TERM:
      return OrderType.LIMIT;
    case ORDER_FLAG_TWAP:
      return OrderType.TWAP;
    case ORDER_FLAG_TWAP_SUBORDER:
      return OrderType.TWAP_SUBORDER;
    default:
      throw new Error(`Unexpected OrderFlags: ${orderFlag}`);
  }
}

/**
 * Converts OrderType enum to protocol ConditionType.
 * Special cased:
 * - all unused types (not LIMIT / STOP-LIMIT/MARKET / TAKE-PROFIT (MARKET)) default to unspecified
 * - STOP_LIMIT and STOP_MARKET map to STOP_LOSS
 * - TAKE_PROFIT and TAKE_PROFIT_MARKET map to TAKE_PROFIT
 * @param orderType
 * @returns
 */
export function orderTypeToProtocolConditionType(
  orderType: OrderType,
): IndexerOrder_ConditionType {
  if (!(orderType in ORDER_TYPE_TO_CONDITION_TYPE_MAP)) {
    throw new Error(`Unexpected OrderType: ${orderType}`);
  }

  return ORDER_TYPE_TO_CONDITION_TYPE_MAP[orderType];
}

export function clobStatusToMarketStatus(clobPairStatus: ClobPairStatus): PerpetualMarketStatus {
  if (
    clobPairStatus !== ClobPairStatus.CLOB_PAIR_STATUS_UNSPECIFIED &&
    clobPairStatus !== ClobPairStatus.UNRECOGNIZED &&
    clobPairStatus in CLOB_STATUS_TO_MARKET_STATUS
  ) {
    return CLOB_STATUS_TO_MARKET_STATUS[clobPairStatus];
  } else {
    throw new InvalidClobPairStatusError(clobPairStatus);
  }
}
