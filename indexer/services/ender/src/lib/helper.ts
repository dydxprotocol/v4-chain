import { logger, ParseMessageError } from '@dydxprotocol-indexer/base';
import {
  OrderSide,
  PerpetualMarketFromDatabase,
  PositionSide,
  protocolTranslations,
} from '@dydxprotocol-indexer/postgres';
import {
  IndexerTendermintEvent,
  IndexerTendermintEvent_BlockEvent,
  Timestamp,
  OrderFillEventV1,
  MarketEventV1,
  SubaccountUpdateEventV1,
  TransferEventV1,
  IndexerOrder,
  StatefulOrderEventV1,
  FundingEventV1,
  AssetCreateEventV1,
  PerpetualMarketCreateEventV1,
  LiquidityTierUpsertEventV1,
  UpdatePerpetualEventV1,
  UpdateClobPairEventV1,
  SubaccountMessage,
  DeleveragingEventV1,
} from '@dydxprotocol-indexer/v4-protos';
import Big from 'big.js';
import _ from 'lodash';
import { DateTime } from 'luxon';

import {
  MILLIS_IN_NANOS,
  SECONDS_IN_MILLIS,
} from '../constants';
import {
  AnnotatedSubaccountMessage,
  DydxIndexerSubtypes,
  EventProtoWithTypeAndVersion,
} from './types';

export function indexerTendermintEventToTransactionIndex(
  event: IndexerTendermintEvent,
): number {
  if (event.transactionIndex !== undefined) {
    return event.transactionIndex;
  } else if (event.blockEvent !== undefined) {
    switch (event.blockEvent) {
      case IndexerTendermintEvent_BlockEvent.BLOCK_EVENT_BEGIN_BLOCK:
        return -2;
      case IndexerTendermintEvent_BlockEvent.BLOCK_EVENT_END_BLOCK:
        return -1;
      default:
        throw new ParseMessageError(`Received V4 event with invalid block event type: ${event.blockEvent}`);
    }
  }

  throw new ParseMessageError(
    'Either transactionIndex or blockEvent must be defined in IndexerTendermintEvent',
  );
}

export function convertToSubaccountMessage(
  annotatedMessage: AnnotatedSubaccountMessage,
): SubaccountMessage {
  const subaccountMessage: SubaccountMessage = _.omit(
    annotatedMessage,
    ['orderId', 'isFill', 'subaccountMessageContents'],
  );
  return subaccountMessage;
}

export function protoTimestampToDate(
  protoTime: Timestamp,
): Date {
  const timeInMillis: number = Number(protoTime.seconds) * SECONDS_IN_MILLIS +
    Math.floor(protoTime.nanos / MILLIS_IN_NANOS);

  return new Date(timeInMillis);
}

export function dateToDateTime(
  protoTime: Date,
): DateTime {
  return DateTime.fromJSDate(protoTime);
}

/**
 * Determines the event subtype and parses the IndexerTendermintEvent
 * to the correct EventProto and returns it all in an object.
 * @param event
 * @returns
 */
export function indexerTendermintEventToEventProtoWithType(
  blockEventIndex: number,
  event: IndexerTendermintEvent,
): EventProtoWithTypeAndVersion | undefined {
  const eventDataBinary: Uint8Array = event.dataBytes;
  // set the default version to 1
  const version: number = event.version === 0 ? 1 : event.version;
  switch (event.subtype) {
    case (DydxIndexerSubtypes.ORDER_FILL.toString()): {
      return {
        type: DydxIndexerSubtypes.ORDER_FILL,
        eventProto: OrderFillEventV1.decode(eventDataBinary),
        indexerTendermintEvent: event,
        version,
        blockEventIndex,
      };
    }
    case (DydxIndexerSubtypes.SUBACCOUNT_UPDATE.toString()): {
      return {
        type: DydxIndexerSubtypes.SUBACCOUNT_UPDATE,
        eventProto: SubaccountUpdateEventV1.decode(eventDataBinary),
        indexerTendermintEvent: event,
        version,
        blockEventIndex,
      };
    }
    case (DydxIndexerSubtypes.TRANSFER.toString()): {
      return {
        type: DydxIndexerSubtypes.TRANSFER,
        eventProto: TransferEventV1.decode(eventDataBinary),
        indexerTendermintEvent: event,
        version,
        blockEventIndex,
      };
    }
    case (DydxIndexerSubtypes.MARKET.toString()): {
      return {
        type: DydxIndexerSubtypes.MARKET,
        eventProto: MarketEventV1.decode(eventDataBinary),
        indexerTendermintEvent: event,
        version,
        blockEventIndex,
      };
    }
    case (DydxIndexerSubtypes.STATEFUL_ORDER.toString()): {
      return {
        type: DydxIndexerSubtypes.STATEFUL_ORDER,
        eventProto: StatefulOrderEventV1.decode(eventDataBinary),
        indexerTendermintEvent: event,
        version,
        blockEventIndex,
      };
    }
    case (DydxIndexerSubtypes.FUNDING.toString()): {
      return {
        type: DydxIndexerSubtypes.FUNDING,
        eventProto: FundingEventV1.decode(eventDataBinary),
        indexerTendermintEvent: event,
        version,
        blockEventIndex,
      };
    }
    case (DydxIndexerSubtypes.ASSET.toString()): {
      return {
        type: DydxIndexerSubtypes.ASSET,
        eventProto: AssetCreateEventV1.decode(eventDataBinary),
        indexerTendermintEvent: event,
        version,
        blockEventIndex,
      };
    }
    case (DydxIndexerSubtypes.PERPETUAL_MARKET.toString()): {
      return {
        type: DydxIndexerSubtypes.PERPETUAL_MARKET,
        eventProto: PerpetualMarketCreateEventV1.decode(eventDataBinary),
        indexerTendermintEvent: event,
        version,
        blockEventIndex,
      };
    }
    case (DydxIndexerSubtypes.LIQUIDITY_TIER.toString()): {
      return {
        type: DydxIndexerSubtypes.LIQUIDITY_TIER,
        eventProto: LiquidityTierUpsertEventV1.decode(eventDataBinary),
        indexerTendermintEvent: event,
        version,
        blockEventIndex,
      };
    }
    case (DydxIndexerSubtypes.UPDATE_PERPETUAL.toString()): {
      return {
        type: DydxIndexerSubtypes.UPDATE_PERPETUAL,
        eventProto: UpdatePerpetualEventV1.decode(eventDataBinary),
        indexerTendermintEvent: event,
        version,
        blockEventIndex,
      };
    }
    case (DydxIndexerSubtypes.UPDATE_CLOB_PAIR.toString()): {
      return {
        type: DydxIndexerSubtypes.UPDATE_CLOB_PAIR,
        eventProto: UpdateClobPairEventV1.decode(eventDataBinary),
        indexerTendermintEvent: event,
        version,
        blockEventIndex,
      };
    }
    case (DydxIndexerSubtypes.DELEVERAGING.toString()): {
      return {
        type: DydxIndexerSubtypes.DELEVERAGING,
        eventProto: DeleveragingEventV1.decode(eventDataBinary),
        indexerTendermintEvent: event,
        version,
        blockEventIndex,
      };
    }
    default: {
      const message: string = `Unable to parse event subtype: ${event.subtype}`;
      logger.error({
        at: 'helpers#indexerTendermintEventToEventWithType',
        message,
      });
      return undefined;
    }
  }
}

/**
 * Returns the size of an order in human readable form.
 * @param order
 * @param perpetualMarket
 */
export function getSize(
  order: IndexerOrder,
  perpetualMarket: PerpetualMarketFromDatabase,
): string {
  return protocolTranslations.quantumsToHumanFixedString(
    order.quantums.toString(10),
    perpetualMarket.atomicResolution,
  );
}

/**
 * Returns the price of an order in human readable form.
 *
 * @param order
 * @param perpetualMarket
 */
export function getPrice(
  order: IndexerOrder,
  perpetualMarket: PerpetualMarketFromDatabase,
): string {
  return protocolTranslations.subticksToPrice(
    order.subticks.toString(10),
    perpetualMarket,
  );
}

/**
 * Returns the trigger price of an order in human readable form.
 *
 * @param order
 * @param perpetualMarket
 * @returns
 */
export function getTriggerPrice(
  order: IndexerOrder,
  perpetualMarket: PerpetualMarketFromDatabase,
): string {
  return protocolTranslations.subticksToPrice(
    order.conditionalOrderTriggerSubticks.toString(10),
    perpetualMarket,
  );
}

/**
 * Returns the weighted average between two prices
 * @param firstPrice
 * @param firstWeight
 * @param secondPrice
 * @param secondWeight
 * @returns
 */
export function getWeightedAverage(
  firstPrice: string,
  firstWeight: string,
  secondPrice: string,
  secondWeight: string,
): string {
  return Big(firstPrice).times(firstWeight).plus(
    Big(secondPrice).times(secondWeight),
  )
    .div(Big(firstWeight).plus(secondWeight))
    .toFixed();
}

/**
 * Returns true if perpetualPostionSide is LONG and orderSide is BUY or
 * if perpetualPostionSide is SHORT and orderSide is SELL
 * @param perpetualPositionSide
 * @param orderSide
 */
export function perpetualPositionAndOrderSideMatching(
  perpetualPositionSide: PositionSide,
  orderSide: OrderSide,
): boolean {
  return (perpetualPositionSide === PositionSide.LONG && orderSide === OrderSide.BUY) ||
    (perpetualPositionSide === PositionSide.SHORT && orderSide === OrderSide.SELL);
}
