import {
  CANDLES_WEBSOCKET_MESSAGE_VERSION,
  MARKETS_WEBSOCKET_MESSAGE_VERSION,
  ORDERBOOKS_WEBSOCKET_MESSAGE_VERSION,
  SUBACCOUNTS_WEBSOCKET_MESSAGE_VERSION,
  TRADES_WEBSOCKET_MESSAGE_VERSION,
} from '@dydxprotocol-indexer/kafka';
import {
  CandleMessage,
  CandleMessage_Resolution,
  MarketMessage,
  OrderbookMessage,
  SubaccountId,
  SubaccountMessage,
  TradeMessage,
} from '@dydxprotocol-indexer/v4-protos';

export const btcClobPairId: string = '1';
export const ethClobPairId: string = '2';
export const btcTicker: string = 'BTC-USD';
export const ethTicker: string = 'ETH-USD';
export const invalidTicker: string = 'INVALID-INVALID';
export const invalidClobPairId: string = '4125';
export const invalidTopic: string = 'invalidTopic';
export const invalidChannel: string = 'invalidChannel';
export const defaultBlockHeight: string = '0';
export const defaultTxIndex: number = 1;
export const defaultEventIndex: number = 3;
export const defaultOwner: string = 'owner';
export const defaultAccNumber: number = 4;
export const defaultSubaccountId: SubaccountId = {
  owner: defaultOwner,
  number: defaultAccNumber,
};
export const defaultContents: Object = {
  prop: 'property',
  field: 'field',
};
export const defaultContentsString: string = JSON.stringify(defaultContents);

const commonMsgProps: {
  blockHeight: string,
  transactionIndex: number,
  eventIndex: number,
} = {
  blockHeight: defaultBlockHeight,
  transactionIndex: defaultTxIndex,
  eventIndex: defaultEventIndex,
};

export const subaccountMessage: SubaccountMessage = {
  ...commonMsgProps,
  subaccountId: defaultSubaccountId,
  contents: defaultContentsString,
  version: SUBACCOUNTS_WEBSOCKET_MESSAGE_VERSION,
};

export const candlesMessage: CandleMessage = {
  ...commonMsgProps,
  clobPairId: btcClobPairId,
  resolution: CandleMessage_Resolution.ONE_MINUTE,
  contents: defaultContentsString,
  version: CANDLES_WEBSOCKET_MESSAGE_VERSION,
};

export const marketsMessage: MarketMessage = {
  ...commonMsgProps,
  contents: defaultContentsString,
  version: MARKETS_WEBSOCKET_MESSAGE_VERSION,
};

export const orderbookMessage: OrderbookMessage = {
  ...commonMsgProps,
  clobPairId: btcClobPairId,
  contents: defaultContentsString,
  version: ORDERBOOKS_WEBSOCKET_MESSAGE_VERSION,
};

export const tradesMessage: TradeMessage = {
  ...commonMsgProps,
  clobPairId: btcClobPairId,
  contents: defaultContentsString,
  version: TRADES_WEBSOCKET_MESSAGE_VERSION,
};
