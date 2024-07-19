import { TRADES_WEBSOCKET_MESSAGE_VERSION, KafkaTopics } from '@dydxprotocol-indexer/kafka';
import { testConstants, TradeContent, TradeMessageContents } from '@dydxprotocol-indexer/postgres';
import { TradeMessage } from '@dydxprotocol-indexer/v4-protos';

import { AnnotatedSubaccountMessage, ConsolidatedKafkaEvent, SingleTradeMessage } from '../../src/lib/types';

export function contentToTradeMessage(
  tradeContent: TradeContent,
  clobPairId: string,
): TradeMessage {
  const contents: TradeMessageContents = {
    trades: [tradeContent],
  };
  return {
    blockHeight: testConstants.defaultBlock.blockHeight,
    contents: JSON.stringify(contents),
    clobPairId,
    version: TRADES_WEBSOCKET_MESSAGE_VERSION,
  };
}

export function contentToSingleTradeMessage(
  tradeContent: TradeContent,
  clobPairId: string,
  transactionIndex: number = 1,
  eventIndex: number = 1,
): SingleTradeMessage {
  return {
    ...contentToTradeMessage(tradeContent, clobPairId),
    transactionIndex,
    eventIndex,
  };
}

export function createConsolidatedKafkaEventFromTrade(
  trade: SingleTradeMessage,
): ConsolidatedKafkaEvent {
  return {
    topic: KafkaTopics.TO_WEBSOCKETS_TRADES,
    message: trade,
  };
}

export function createConsolidatedKafkaEventFromSubaccount(
  subaccount: AnnotatedSubaccountMessage,
): ConsolidatedKafkaEvent {
  return {
    topic: KafkaTopics.TO_WEBSOCKETS_SUBACCOUNTS,
    message: subaccount,
  };
}
