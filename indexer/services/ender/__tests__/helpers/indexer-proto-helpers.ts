import { randomUUID } from 'crypto';

import {
  createKafkaMessage,
  MARKETS_WEBSOCKET_MESSAGE_VERSION,
  SUBACCOUNTS_WEBSOCKET_MESSAGE_VERSION,
  TRADES_WEBSOCKET_MESSAGE_VERSION,
  KafkaTopics,
} from '@dydxprotocol-indexer/kafka';
import {
  FillFromDatabase,
  FillTable,
  Liquidity,
  OrderFromDatabase,
  OrderSide,
  OrderStatus,
  OrderTable,
  OrderType,
  PerpetualPositionFromDatabase,
  PerpetualPositionTable,
  SubaccountMessageContents,
  testConstants,
  TimeInForce,
  TradeMessageContents,
  apiTranslations,
  FillType,
  perpetualMarketRefresher,
  PerpetualMarketStatus,
  PerpetualMarketFromDatabase,
  PerpetualMarketTable,
  IsoString,
  fillTypeToTradeType,
  OrderSubaccountMessageContents,
  MarketFromDatabase,
  MarketTable,
  MarketsMap,
  MarketColumns,
  UpdatedPerpetualPositionSubaccountKafkaObject,
} from '@dydxprotocol-indexer/postgres';
import { getOrderIdHash, ORDER_FLAG_CONDITIONAL, ORDER_FLAG_TWAP } from '@dydxprotocol-indexer/v4-proto-parser';
import {
  LiquidationOrderV1,
  MarketMessage,
  IndexerOrder,
  IndexerOrder_Side,
  OrderFillEventV1,
  IndexerSubaccountId,
  SubaccountMessage,
  TradeMessage,
  IndexerOrder_TimeInForce,
  IndexerTendermintBlock,
  IndexerTendermintEvent,
  IndexerTendermintEvent_BlockEvent,
  Timestamp,
  OffChainUpdateV1,
  IndexerOrderId,
  PerpetualMarketCreateEventV1,
  PerpetualMarketCreateEventV2,
  PerpetualMarketCreateEventV3,
  DeleveragingEventV1,
  protoTimestampToDate,
  PerpetualMarketType,
} from '@dydxprotocol-indexer/v4-protos';
import { IHeaders, Message, ProducerRecord } from 'kafkajs';
import _ from 'lodash';

import {
  annotateWithPnl,
  convertPerpetualPosition,
  generateFillSubaccountMessage,
  generatePerpetualMarketMessage,
  generatePerpetualPositionsContents,
} from '../../src/helpers/kafka-helper';
import { DydxIndexerSubtypes, VulcanMessage } from '../../src/lib/types';

// TX Hash is SHA256, so is of length 64 hexadecimal without the '0x'.
// This can be generated from 32 characters, because each character generates two
// hexadecimal characters.
const NUM_CHARS_IN_TX_HASH: number = 32;
const defaultPerpetualMarketTicker: string = testConstants.defaultPerpetualMarket.ticker;

/**
 * Creates an IndexerTendermintEvent, if transactionIndex < 0, creates a block event,
 * otherwise creates a transaction event.
 * @param subtype
 * @param dataBytes
 * @param transactionIndex
 * @param eventIndex
 * @param version
 * @returns
 */
export function createIndexerTendermintEvent(
  subtype: string,
  dataBytes: Uint8Array,
  transactionIndex: number,
  eventIndex: number,
  version: number = 1,
): IndexerTendermintEvent {
  if (transactionIndex < 0) {
    // blockEvent
    return {
      subtype,
      dataBytes,
      blockEvent: IndexerTendermintEvent_BlockEvent.BLOCK_EVENT_END_BLOCK,
      eventIndex,
      version,
    };
  }
  // transactionIndex
  return {
    subtype,
    dataBytes,
    transactionIndex,
    eventIndex,
    version,
  };
}

export function createTxHashes(numTx: number): string[] {
  const txHashes: string[] = [];
  _.times(numTx, () => {
    const txHashArray: string[] = [];
    const randomString: string = randomUUID().substring(0, NUM_CHARS_IN_TX_HASH);
    _.times(randomString.length, (index: number) => {
      txHashArray.push(randomString.charCodeAt(index).toString(16));
    });
    txHashes.push('0x'.concat(txHashArray.join('')));
  });
  return txHashes;
}

export function createIndexerTendermintBlock(
  height: number,
  time: Timestamp,
  events: IndexerTendermintEvent[],
  txHashes: string[],
): IndexerTendermintBlock {
  return {
    height,
    time: protoTimestampToDate(time),
    events,
    txHashes,
  };
}

export function expectSubaccountKafkaMessage({
  producerSendMock,
  blockHeight,
  transactionIndex,
  eventIndex,
  contents,
  subaccountIdProto,
}: {
  producerSendMock: jest.SpyInstance,
  blockHeight: string,
  transactionIndex: number,
  eventIndex: number,
  contents: string,
  subaccountIdProto: IndexerSubaccountId,
}): void {
  expect(producerSendMock.mock.calls.length).toBeGreaterThanOrEqual(1);
  expect(producerSendMock.mock.calls[0].length).toBeGreaterThanOrEqual(1);

  const subaccountProducerRecords: ProducerRecord[] = _.filter(
    _.flatten(producerSendMock.mock.calls) as ProducerRecord[],
    (producerRecord: ProducerRecord) => {
      return producerRecord.topic === KafkaTopics.TO_WEBSOCKETS_SUBACCOUNTS.toString();
    },
  );
  expect(subaccountProducerRecords.length).toEqual(1);

  const subaccountProducerRecord: ProducerRecord = subaccountProducerRecords[0];
  const subaccountMessages: SubaccountMessage[] = _.map(
    subaccountProducerRecord.messages,
    (message: Message) => {
      expect(Buffer.isBuffer(message.value));

      const messageValueBinary: Uint8Array = new Uint8Array(
        // Can assume Buffer, since we check above that it is a buffer
        message.value as Buffer,
      );
      return SubaccountMessage.decode(messageValueBinary);
    },
  );

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const subaccountMessageJsons: any[] = _.map(subaccountMessages, (message: SubaccountMessage) => {
    return {
      ...message,
      contents: JSON.parse(message.contents),
    };
  });
  const expectedSubaccountMessage: SubaccountMessage = SubaccountMessage.fromPartial({
    blockHeight,
    transactionIndex,
    eventIndex,
    subaccountId: subaccountIdProto,
    contents,
    version: SUBACCOUNTS_WEBSOCKET_MESSAGE_VERSION,
  });
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const expectedSubaccountMessageJson: any = {
    ...expectedSubaccountMessage,
    contents: JSON.parse(expectedSubaccountMessage.contents),
  };
  expect(subaccountMessageJsons).toContainEqual(expectedSubaccountMessageJson);
}

export function expectPerpetualMarketKafkaMessage(
  producerSendMock: jest.SpyInstance,
  perpetualMarkets: PerpetualMarketFromDatabase[],
) {
  expectMarketKafkaMessage({
    producerSendMock,
    contents: JSON.stringify(generatePerpetualMarketMessage(perpetualMarkets)),
  });
}

export function expectMarketKafkaMessage({
  producerSendMock,
  contents,
}: {
  producerSendMock: jest.SpyInstance,
  contents: string,
}): void {
  expect(producerSendMock.mock.calls.length).toBeGreaterThanOrEqual(1);
  expect(producerSendMock.mock.calls[0].length).toBeGreaterThanOrEqual(1);

  const marketProducerRecords: ProducerRecord[] = _.filter(
    _.flatten(producerSendMock.mock.calls) as ProducerRecord[],
    (producerRecord: ProducerRecord) => {
      return producerRecord.topic === KafkaTopics.TO_WEBSOCKETS_MARKETS.toString();
    },
  );
  expect(marketProducerRecords.length).toEqual(1);

  const marketProducerRecord: ProducerRecord = marketProducerRecords[0];
  const marketMessages: MarketMessage[] = _.map(
    marketProducerRecord.messages,
    (message: Message) => {
      expect(Buffer.isBuffer(message.value));

      const messageValueBinary: Uint8Array = new Uint8Array(
        // Can assume Buffer, since we check above that it is a buffer
        message.value as Buffer,
      );
      return MarketMessage.decode(messageValueBinary);
    },
  );

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const marketMessageJsons: any[] = _.map(marketMessages, (message: MarketMessage) => {
    return {
      ...message,
      contents: JSON.parse(message.contents),
    };
  });
  const expectedMarketMessage: MarketMessage = MarketMessage.fromPartial({
    contents,
    version: MARKETS_WEBSOCKET_MESSAGE_VERSION,
  });
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const expectedMarketMessageJson: any = {
    ...expectedMarketMessage,
    contents: JSON.parse(expectedMarketMessage.contents),
  };
  expect(marketMessageJsons).toContainEqual(expectedMarketMessageJson);
}

export function expectTradeKafkaMessage({
  producerSendMock,
  blockHeight,
  contents,
  clobPairId,
}: {
  producerSendMock: jest.SpyInstance,
  blockHeight: string,
  contents: string,
  clobPairId: string,
}): void {
  expect(producerSendMock.mock.calls.length).toBeGreaterThanOrEqual(1);
  expect(producerSendMock.mock.calls[0].length).toBeGreaterThanOrEqual(1);

  const tradeProducerRecords: ProducerRecord[] = _.filter(
    _.flatten(producerSendMock.mock.calls) as ProducerRecord[],
    (producerRecord: ProducerRecord) => {
      return producerRecord.topic === KafkaTopics.TO_WEBSOCKETS_TRADES.toString();
    },
  );
  expect(tradeProducerRecords.length).toEqual(1);

  const tradeProducerRecord: ProducerRecord = tradeProducerRecords[0];
  const tradeMessages: TradeMessage[] = _.map(
    tradeProducerRecord.messages,
    (message: Message) => {
      expect(Buffer.isBuffer(message.value));

      const messageValueBinary: Uint8Array = new Uint8Array(
        // Can assume Buffer, since we check above that it is a buffer
        message.value as Buffer,
      );
      return TradeMessage.decode(messageValueBinary);
    },
  );

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const tradeMessageJsons: any[] = _.map(tradeMessages, (message: TradeMessage) => {
    return {
      ...message,
      contents: JSON.parse(message.contents),
    };
  });
  const expectedTradeMessage: TradeMessage = TradeMessage.fromPartial({
    blockHeight,
    contents,
    clobPairId,
    version: TRADES_WEBSOCKET_MESSAGE_VERSION,
  });
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const expectedTradeMessageJson: any = {
    ...expectedTradeMessage,
    contents: JSON.parse(expectedTradeMessage.contents),
  };
  expect(tradeMessageJsons).toContainEqual(expectedTradeMessageJson);
}

export function expectVulcanKafkaMessage({
  producerSendMock,
  orderId,
  offchainUpdate,
  headers,
}: {
  producerSendMock: jest.SpyInstance,
  orderId: IndexerOrderId,
  offchainUpdate: OffChainUpdateV1,
  headers?: IHeaders,
}): void {
  expect(producerSendMock.mock.calls.length).toBeGreaterThanOrEqual(1);
  expect(producerSendMock.mock.calls[0].length).toBeGreaterThanOrEqual(1);

  const vulcanProducerRecords: ProducerRecord[] = _.filter(
    _.flatten(producerSendMock.mock.calls) as ProducerRecord[],
    (producerRecord: ProducerRecord) => {
      return producerRecord.topic === KafkaTopics.TO_VULCAN.toString();
    },
  );
  expect(vulcanProducerRecords.length).toEqual(1);

  const vulcanProducerRecord: ProducerRecord = vulcanProducerRecords[0];
  const vulcanMessages: VulcanMessage[] = _.map(
    vulcanProducerRecord.messages,
    (message: Message): VulcanMessage => {
      expect(Buffer.isBuffer(message.value));
      const messageValueBinary: Uint8Array = new Uint8Array(
        // Can assume Buffer, since we check above that it is a buffer
        message.value as Buffer,
      );

      return {
        key: message.key as Buffer,
        value: OffChainUpdateV1.decode(messageValueBinary),
        headers: message.headers,
      };
    },
  );

  expect(vulcanMessages).toContainEqual({
    key: getOrderIdHash(orderId),
    value: offchainUpdate,
    headers,
  });
}

export function createLiquidationOrder({
  subaccountId,
  clobPairId,
  perpetualId,
  quantums,
  isBuy,
  subticks,
}: {
  subaccountId: IndexerSubaccountId,
  clobPairId: string,
  perpetualId: string,
  quantums: number,
  isBuy: boolean,
  subticks: number,
}): LiquidationOrderV1 {
  return LiquidationOrderV1.fromPartial({
    liquidated: subaccountId,
    clobPairId: parseInt(clobPairId, 10),
    perpetualId: parseInt(perpetualId, 10),
    totalSize: quantums,
    isBuy,
    subticks,
  });
}

export function createOrder({
  subaccountId,
  clientId,
  side,
  quantums,
  subticks,
  goodTilOneof,
  clobPairId,
  orderFlags,
  timeInForce,
  reduceOnly,
  clientMetadata,
  builderAddress,
  feePpm,
  orderRouterAddress,
  duration,
  interval,
  priceTolerance,
}: {
  subaccountId: IndexerSubaccountId,
  clientId: number,
  side: IndexerOrder_Side,
  quantums: number | Long,
  subticks: number | Long,
  goodTilOneof: Partial<IndexerOrder>,
  clobPairId: string,
  orderFlags: string,
  timeInForce: IndexerOrder_TimeInForce,
  reduceOnly: boolean,
  clientMetadata: number,
  builderAddress?: string,
  feePpm?: number,
  orderRouterAddress?: string,
  duration?: number,
  interval?: number,
  priceTolerance?: number,
}): IndexerOrder {
  // eslint-disable-next-line  @typescript-eslint/no-explicit-any
  let orderJSON: any = {
    orderId: {
      subaccountId,
      clientId,
      clobPairId: parseInt(clobPairId, 10),
      orderFlags: parseInt(orderFlags, 10),
    },
    side,
    quantums,
    subticks,
    timeInForce,
    reduceOnly,
    clientMetadata,
  };

  if (builderAddress !== undefined && feePpm !== undefined) {
    orderJSON = {
      ...orderJSON,
      builderCodeParams: {
        builderAddress,
        feePpm,
      },
    };
  }

  if (duration !== undefined && interval !== undefined && priceTolerance !== undefined) {
    orderJSON = {
      ...orderJSON,
      twapParameters: {
        duration,
        interval,
        priceTolerance,
      },
    };
  }

  if (goodTilOneof.goodTilBlock !== undefined) {
    orderJSON = {
      ...orderJSON,
      goodTilBlock: goodTilOneof.goodTilBlock,
    };
  } else if (goodTilOneof.goodTilBlockTime !== undefined) {
    orderJSON = {
      ...orderJSON,
      goodTilBlockTime: goodTilOneof.goodTilBlockTime,
    };
  }

  if (orderRouterAddress !== undefined) {
    orderJSON = {
      ...orderJSON,
      orderRouterAddress,
    };
  }

  return IndexerOrder.fromPartial(orderJSON);
}

export function createKafkaMessageFromOrderFillEvent({
  orderFillEvent,
  transactionIndex,
  eventIndex,
  height,
  time,
  txHash,
}: {
  orderFillEvent: OrderFillEventV1,
  transactionIndex: number,
  eventIndex: number,
  height: number,
  time: Timestamp,
  txHash: string,
}) {
  const events: IndexerTendermintEvent[] = [
    createIndexerTendermintEvent(
      DydxIndexerSubtypes.ORDER_FILL,
      Uint8Array.from(OrderFillEventV1.encode(orderFillEvent).finish()),
      transactionIndex,
      eventIndex,
    ),
  ];

  const block: IndexerTendermintBlock = createIndexerTendermintBlock(
    height,
    time,
    events,
    [txHash],
  );

  const binaryBlock: Uint8Array = Uint8Array.from(IndexerTendermintBlock.encode(block).finish());
  return createKafkaMessage(Buffer.from(binaryBlock));
}

export function createKafkaMessageFromDeleveragingEvent({
  deleveragingEvent,
  transactionIndex,
  eventIndex,
  height,
  time,
  txHash,
}: {
  deleveragingEvent: DeleveragingEventV1,
  transactionIndex: number,
  eventIndex: number,
  height: number,
  time: Timestamp,
  txHash: string,
}) {
  const events: IndexerTendermintEvent[] = [
    createIndexerTendermintEvent(
      DydxIndexerSubtypes.DELEVERAGING,
      Uint8Array.from(DeleveragingEventV1.encode(deleveragingEvent).finish()),
      transactionIndex,
      eventIndex,
    ),
  ];

  const block: IndexerTendermintBlock = createIndexerTendermintBlock(
    height,
    time,
    events,
    [txHash],
  );

  const binaryBlock: Uint8Array = Uint8Array.from(IndexerTendermintBlock.encode(block).finish());
  return createKafkaMessage(Buffer.from(binaryBlock));
}

export function liquidationOrderToOrderSide(
  liquidationOrder: LiquidationOrderV1,
): OrderSide {
  return liquidationOrder.isBuy ? OrderSide.BUY : OrderSide.SELL;
}

export async function expectFillInDatabase({
  subaccountId,
  clientId,
  liquidity,
  size,
  price,
  quoteAmount,
  eventId,
  transactionHash,
  createdAt,
  createdAtHeight,
  type,
  clobPairId,
  side,
  orderFlags,
  clientMetadata,
  fee,
  affiliateRevShare,
  hasOrderId = true,
  builderAddress = null,
  builderFee = null,
  orderRouterAddress = null,
  orderRouterFee = null,
  positionSideBefore,
  entryPriceBefore,
  positionSizeBefore,
}: {
  subaccountId: string,
  clientId: string,
  liquidity: Liquidity,
  size: string,
  price: string,
  quoteAmount: string,
  eventId: Buffer,
  transactionHash: string,
  createdAt: string,
  createdAtHeight: string,
  type: FillType,
  clobPairId: string,
  side: OrderSide,
  orderFlags: string,
  clientMetadata: string | null,
  fee: string,
  affiliateRevShare: string,
  hasOrderId?: boolean,
  builderAddress?: string | null,
  builderFee?: string | null,
  orderRouterAddress?: string | null,
  orderRouterFee?: string | null,
  positionSizeBefore?: string | null,
  entryPriceBefore?: string | null,
  positionSideBefore?: string | null,
}): Promise<void> {
  const fillId: string = FillTable.uuid(eventId, liquidity);
  const fill: FillFromDatabase | undefined = await FillTable.findById(fillId);

  expect(fill).not.toEqual(undefined);
  expect(fill).toEqual(expect.objectContaining({
    subaccountId,
    side,
    liquidity,
    type,
    clobPairId,
    orderId: hasOrderId ? OrderTable.uuid(subaccountId, clientId, clobPairId, orderFlags) : null,
    size,
    price,
    quoteAmount,
    eventId,
    transactionHash,
    createdAt,
    createdAtHeight,
    clientMetadata,
    fee,
    affiliateRevShare,
    builderAddress,
    builderFee,
    orderRouterAddress,
    orderRouterFee,
    ...(positionSideBefore ? { positionSideBefore } : {}),
    ...(positionSizeBefore ? { positionSizeBefore } : {}),
    ...(entryPriceBefore ? { entryPriceBefore } : {}),
  }));
}

export async function expectNoOrdersExistForSubaccountClobPairId({
  subaccountId,
  clobPairId,
}: {
  subaccountId: string,
  clobPairId: string,
}): Promise<void> {
  const ordersFromDatabase: OrderFromDatabase[] = await
  OrderTable.findBySubaccountIdAndClobPair(subaccountId, clobPairId);
  expect(ordersFromDatabase).toHaveLength(0);
}

export async function expectOrderInDatabase({
  subaccountId,
  clientId,
  size,
  totalFilled,
  price,
  status,
  clobPairId,
  side,
  orderFlags,
  timeInForce,
  reduceOnly,
  goodTilBlock,
  goodTilBlockTime,
  clientMetadata,
  updatedAt,
  updatedAtHeight,
  builderAddress,
  feePpm,
  duration,
  interval,
  priceTolerance,
  orderType,
}: {
  subaccountId: string,
  clientId: string,
  size: string,
  totalFilled: string,
  price: string,
  status: OrderStatus,
  clobPairId: string,
  side: OrderSide,
  orderFlags: string,
  timeInForce: TimeInForce,
  reduceOnly: boolean,
  goodTilBlock?: string,
  goodTilBlockTime?: string,
  clientMetadata: string,
  updatedAt: IsoString,
  updatedAtHeight: string,
  builderAddress?: string,
  feePpm?: number,
  duration?: number,
  interval?: number,
  priceTolerance?: number,
  orderType?: OrderType,
}): Promise<void> {
  const orderId: string = OrderTable.uuid(subaccountId, clientId, clobPairId, orderFlags);
  const orderFromDatabase: OrderFromDatabase | undefined = await
  OrderTable.findById(orderId);

  expect(orderFromDatabase).not.toEqual(undefined);
  expect(orderFromDatabase).toEqual(expect.objectContaining({
    subaccountId,
    clientId,
    clobPairId,
    side,
    size,
    totalFilled,
    price,
    type: orderType ?? OrderType.LIMIT, // TODO: Add additional order types
    status,
    timeInForce,
    reduceOnly,
    orderFlags,
    goodTilBlock: goodTilBlock ?? null,
    goodTilBlockTime: goodTilBlockTime ?? null,
    clientMetadata,
    updatedAt,
    updatedAtHeight,
    builderAddress: builderAddress ?? null,
    feePpm: feePpm ?? null,
    duration: duration ?? null,
    interval: interval ?? null,
    priceTolerance: priceTolerance ?? null,
  }));
}

export async function expectFillSubaccountKafkaMessageFromLiquidationEvent(
  producerSendMock: jest.SpyInstance,
  subaccountIdProto: IndexerSubaccountId,
  fillId: string,
  positionId: string,
  blockHeight: string = '3',
  transactionIndex: number = 0,
  eventIndex: number = 0,
  ticker: string = defaultPerpetualMarketTicker,
) {
  const [
    fill,
    position,
  ]: [
    FillFromDatabase | undefined,
    PerpetualPositionFromDatabase | undefined,
  ] = await Promise.all([
    FillTable.findById(fillId),
    PerpetualPositionTable.findById(positionId),
  ]);
  expect(fill).toBeDefined();
  expect(position).toBeDefined();

  const markets: MarketFromDatabase[] = await MarketTable.findAll(
    {},
    [],
  );
  const marketIdToMarket: MarketsMap = _.keyBy(
    markets,
    MarketColumns.id,
  );
  const positionUpdate = annotateWithPnl(
    convertPerpetualPosition(position!),
    perpetualMarketRefresher.getPerpetualMarketsMap(),
    marketIdToMarket[parseInt(position!.perpetualId, 10)],
  );

  const contents: SubaccountMessageContents = {
    fills: [
      generateFillSubaccountMessage(fill!, ticker),
    ],
    perpetualPositions: generatePerpetualPositionsContents(
      subaccountIdProto,
      [positionUpdate],
      perpetualMarketRefresher.getPerpetualMarketsMap(),
    ),
    blockHeight,
  };

  expectSubaccountKafkaMessage({
    producerSendMock,
    blockHeight,
    transactionIndex,
    eventIndex,
    contents: JSON.stringify(contents),
    subaccountIdProto,
  });
}

function isConditionalOrder(order: OrderFromDatabase): boolean {
  return Number(order.orderFlags) === ORDER_FLAG_CONDITIONAL;
}

function isTwapOrder(order: OrderFromDatabase): boolean {
  return Number(order.orderFlags) === ORDER_FLAG_TWAP;
}

export function expectOrderSubaccountKafkaMessage(
  producerSendMock: jest.SpyInstance,
  subaccountIdProto: IndexerSubaccountId,
  order: OrderFromDatabase,
  blockHeight: string = '3',
  transactionIndex: number = 0,
  eventIndex: number = 0,
  ticker: string = defaultPerpetualMarketTicker,
): void {
  const {
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    triggerPrice, totalFilled, goodTilBlock, ...orderWithoutUnwantedFields
  } = order!;
  let orderObject: OrderSubaccountMessageContents;

  if (isConditionalOrder(order)) {
    orderObject = {
      ...order!,
      timeInForce: apiTranslations.orderTIFToAPITIF(order!.timeInForce),
      postOnly: apiTranslations.isOrderTIFPostOnly(order!.timeInForce),
      goodTilBlock: order!.goodTilBlock,
      goodTilBlockTime: order!.goodTilBlockTime,
      ticker,
    };
  } else if (isTwapOrder(order)) {
    orderObject = {
      ...order!,
      timeInForce: apiTranslations.orderTIFToAPITIF(order!.timeInForce),
      postOnly: apiTranslations.isOrderTIFPostOnly(order!.timeInForce),
      goodTilBlock: order!.goodTilBlock,
      goodTilBlockTime: order!.goodTilBlockTime,
      ticker,
    };
  } else {
    orderObject = {
      ...orderWithoutUnwantedFields!,
      timeInForce: apiTranslations.orderTIFToAPITIF(order!.timeInForce),
      postOnly: apiTranslations.isOrderTIFPostOnly(order!.timeInForce),
      goodTilBlockTime: order!.goodTilBlockTime,
      ticker,
    };
  }

  const contents: SubaccountMessageContents = {
    orders: [
      orderObject,
    ],
    blockHeight,
  };

  expectSubaccountKafkaMessage({
    producerSendMock,
    blockHeight,
    transactionIndex,
    eventIndex,
    contents: JSON.stringify(contents),
    subaccountIdProto,
  });
}

export async function expectOrderFillAndPositionSubaccountKafkaMessageFromIds(
  producerSendMock: jest.SpyInstance,
  subaccountIdProto: IndexerSubaccountId,
  orderId: string,
  fillId: string,
  positionId: string,
  blockHeight: string = '3',
  transactionIndex: number = 0,
  eventIndex: number = 0,
) {
  const [
    order,
    fill,
    position,
  ]: [
    OrderFromDatabase | undefined,
    FillFromDatabase | undefined,
    PerpetualPositionFromDatabase | undefined,
  ] = await Promise.all([
    OrderTable.findById(orderId),
    FillTable.findById(fillId),
    PerpetualPositionTable.findById(positionId),
  ]);

  const perpetualMarket: PerpetualMarketFromDatabase | undefined = await PerpetualMarketTable
    .findById(
      position!.perpetualId,
    );

  expect(order).toBeDefined();
  expect(fill).toBeDefined();
  expect(perpetualMarket).toBeDefined();

  const contents: SubaccountMessageContents = {
    orders: [
      {
        ...order!,
        timeInForce: apiTranslations.orderTIFToAPITIF(order!.timeInForce),
        postOnly: apiTranslations.isOrderTIFPostOnly(order!.timeInForce),
        goodTilBlock: order!.goodTilBlock,
        goodTilBlockTime: order!.goodTilBlockTime,
        ticker: perpetualMarket!.ticker,
      },
    ],
    fills: [
      generateFillSubaccountMessage(fill!, perpetualMarket!.ticker),
    ],
    blockHeight,
  };

  if (position !== undefined) {
    const markets: MarketFromDatabase[] = await MarketTable.findAll(
      {},
      [],
    );
    const marketIdToMarket: MarketsMap = _.keyBy(
      markets,
      MarketColumns.id,
    );
    const positionUpdate: UpdatedPerpetualPositionSubaccountKafkaObject = annotateWithPnl(
      convertPerpetualPosition(position),
      perpetualMarketRefresher.getPerpetualMarketsMap(),
      marketIdToMarket[parseInt(position.perpetualId, 10)],
    );
    contents.perpetualPositions = generatePerpetualPositionsContents(
      subaccountIdProto,
      [positionUpdate],
      perpetualMarketRefresher.getPerpetualMarketsMap(),
    );
  }

  expectSubaccountKafkaMessage({
    producerSendMock,
    blockHeight,
    transactionIndex,
    eventIndex,
    contents: JSON.stringify(contents),
    subaccountIdProto,
  });
}

export async function expectDefaultTradeKafkaMessageFromTakerFillId(
  producerSendMock: jest.SpyInstance,
  eventId: Buffer,
  blockHeight: string = '3',
) {
  const takerFillId: string = FillTable.uuid(eventId, Liquidity.TAKER);
  const takerFill: FillFromDatabase | undefined = await FillTable.findById(takerFillId);
  expect(takerFill).toBeDefined();

  const contents: TradeMessageContents = {
    trades: [
      {
        id: eventId.toString('hex'),
        size: takerFill!.size,
        price: takerFill!.price,
        side: takerFill!.side.toString(),
        createdAt: takerFill!.createdAt,
        type: fillTypeToTradeType(takerFill!.type),
      },
    ],
  };

  expectTradeKafkaMessage({
    producerSendMock,
    blockHeight,
    contents: JSON.stringify(contents),
    clobPairId: takerFill!.clobPairId,
  });
}

export async function expectPerpetualPosition(
  perpetualPositionId: string,
  fields: {
    sumOpen?: string,
    sumClose?: string,
    entryPrice?: string,
    exitPrice?: string | null,
    totalRealizedPnl?: string | null,
  },
) {
  const perpetualPosition:
  PerpetualPositionFromDatabase | undefined = await PerpetualPositionTable.findById(
    perpetualPositionId,
  );

  expect(perpetualPosition).toBeDefined();

  if (fields.sumOpen !== undefined) {
    expect(perpetualPosition!.sumOpen).toEqual(fields.sumOpen);
  }
  if (fields.sumClose !== undefined) {
    expect(perpetualPosition!.sumClose).toEqual(fields.sumClose);
  }
  if (fields.entryPrice !== undefined) {
    expect(perpetualPosition!.entryPrice).toEqual(fields.entryPrice);
  }
  if (fields.exitPrice !== undefined) {
    expect(perpetualPosition!.exitPrice).toEqual(fields.exitPrice);
  }
  if (fields.totalRealizedPnl !== undefined) {
    expect(perpetualPosition!.totalRealizedPnl).toEqual(fields.totalRealizedPnl);
  }
}

// Values of the `PerpetualMarketCreateObject` which are hard-coded and not derived
// from PerpetualMarketCreate events.
export const HARDCODED_PERPETUAL_MARKET_VALUES: Object = {
  priceChange24H: '0',
  trades24H: 0,
  volume24H: '0',
  nextFundingRate: '0',
  status: PerpetualMarketStatus.ACTIVE,
  openInterest: '0',
};

export function expectPerpetualMarketV1(
  perpetualMarket: PerpetualMarketFromDatabase,
  perpetual: PerpetualMarketCreateEventV1,
): void {
  // TODO(IND-219): Set initialMarginFraction/maintenanceMarginFraction using LiquidityTier
  expect(perpetualMarket).toEqual(expect.objectContaining({
    ...HARDCODED_PERPETUAL_MARKET_VALUES,
    id: perpetual.id.toString(),
    status: PerpetualMarketStatus.INITIALIZING,
    clobPairId: perpetual.clobPairId.toString(),
    ticker: perpetual.ticker,
    marketId: perpetual.marketId,
    quantumConversionExponent: perpetual.quantumConversionExponent,
    atomicResolution: perpetual.atomicResolution,
    subticksPerTick: perpetual.subticksPerTick,
    stepBaseQuantums: Number(perpetual.stepBaseQuantums),
    liquidityTierId: perpetual.liquidityTier,
    marketType: 'CROSS',
  }));
}

export function expectPerpetualMarketV2(
  perpetualMarket: PerpetualMarketFromDatabase,
  perpetual: PerpetualMarketCreateEventV2,
): void {
  // TODO(IND-219): Set initialMarginFraction/maintenanceMarginFraction using LiquidityTier
  expect(perpetualMarket).toEqual(expect.objectContaining({
    ...HARDCODED_PERPETUAL_MARKET_VALUES,
    id: perpetual.id.toString(),
    status: PerpetualMarketStatus.INITIALIZING,
    clobPairId: perpetual.clobPairId.toString(),
    ticker: perpetual.ticker,
    marketId: perpetual.marketId,
    quantumConversionExponent: perpetual.quantumConversionExponent,
    atomicResolution: perpetual.atomicResolution,
    subticksPerTick: perpetual.subticksPerTick,
    stepBaseQuantums: Number(perpetual.stepBaseQuantums),
    liquidityTierId: perpetual.liquidityTier,
    marketType: eventPerpetualMarketTypeToIndexerPerpetualMarketType(
      perpetual.marketType,
    ),
  }));
}

export function expectPerpetualMarketV3(
  perpetualMarket: PerpetualMarketFromDatabase,
  perpetual: PerpetualMarketCreateEventV3,
): void {
  // TODO(IND-219): Set initialMarginFraction/maintenanceMarginFraction using LiquidityTier
  expect(perpetualMarket).toEqual(expect.objectContaining({
    ...HARDCODED_PERPETUAL_MARKET_VALUES,
    id: perpetual.id.toString(),
    status: PerpetualMarketStatus.INITIALIZING,
    clobPairId: perpetual.clobPairId.toString(),
    ticker: perpetual.ticker,
    marketId: perpetual.marketId,
    quantumConversionExponent: perpetual.quantumConversionExponent,
    atomicResolution: perpetual.atomicResolution,
    subticksPerTick: perpetual.subticksPerTick,
    stepBaseQuantums: Number(perpetual.stepBaseQuantums),
    liquidityTierId: perpetual.liquidityTier,
    marketType: eventPerpetualMarketTypeToIndexerPerpetualMarketType(
      perpetual.marketType,
    ),
    defaultFundingRate1H: ((perpetual.defaultFunding8hrPpm / 1_000_000) / 8).toString(),
  }));
}

export function eventPerpetualMarketTypeToIndexerPerpetualMarketType(
  perpetualMarketType: PerpetualMarketType,
): string {
  switch (perpetualMarketType) {
    case PerpetualMarketType.PERPETUAL_MARKET_TYPE_CROSS:
      return 'CROSS';
    case PerpetualMarketType.PERPETUAL_MARKET_TYPE_ISOLATED:
      return 'ISOLATED';
    default:
      throw new Error(`Unknown perpetual market type: ${perpetualMarketType}`);
  }
}
