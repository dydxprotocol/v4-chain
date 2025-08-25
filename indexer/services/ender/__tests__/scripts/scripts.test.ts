import {
  AssetPositionTable,
  BlockTable,
  BUFFER_ENCODING_UTF_8,
  CLOB_STATUS_TO_MARKET_STATUS,
  dbHelpers,
  FillTable,
  FundingIndexUpdatesTable,
  Liquidity,
  OraclePriceTable,
  OrderSide,
  OrderTable,
  PerpetualPositionTable,
  PositionSide,
  protocolTranslations,
  storeHelpers,
  SubaccountTable,
  TendermintEventFromDatabase,
  TendermintEventTable,
  testConstants,
  TransactionFromDatabase,
  TransactionTable,
  TransferTable,
  uuid,
} from '@dydxprotocol-indexer/postgres';
import {
  AssetCreateEventV1,
  IndexerOrder_ConditionType,
  IndexerOrder_Side,
  IndexerOrder_TimeInForce,
  IndexerOrderId,
  IndexerSubaccountId,
  IndexerTendermintEvent,
  IndexerTendermintEvent_BlockEvent,
  MarketEventV1,
  SubaccountUpdateEventV1,
} from '@dydxprotocol-indexer/v4-protos';

import { bigIntToBytes } from '@dydxprotocol-indexer/v4-proto-parser';
import Long from 'long';
import { createPostgresFunctions } from '../../src/helpers/postgres/postgres-functions';
import {
  getWeightedAverage,
  indexerTendermintEventToTransactionIndex,
  perpetualPositionAndOrderSideMatching,
} from '../../src/lib/helper';
import { DydxIndexerSubtypes } from '../../src/lib/types';
import { defaultAssetCreateEvent, defaultMarketCreate } from '../helpers/constants';
import { createIndexerTendermintEvent } from '../helpers/indexer-proto-helpers';

describe('SQL Function Tests', () => {
  beforeAll(async () => {
    await dbHelpers.migrate();
    await createPostgresFunctions();
  });

  beforeEach(async () => {
    await dbHelpers.clearData();
  });

  afterEach(async () => {
    await dbHelpers.clearData();
  });

  afterAll(async () => {
    await dbHelpers.teardown();
    jest.resetAllMocks();
  });

  const defaultTxHash: string = '0x32343534306431622d306461302d343831322d613730372d3965613162336162';
  const defaultTxHash2: string = '0x32363534306431622d306461302d343831322d613730372d3965613162336162';
  const defaultSubaccountUpdateEvent: SubaccountUpdateEventV1 = SubaccountUpdateEventV1
    .fromPartial({
      subaccountId: {
        owner: '',
        number: 0,
      },
    });
  const defaultSubaccountUpdateEventBinary: Uint8Array = Uint8Array.from(
    SubaccountUpdateEventV1.encode(
      defaultSubaccountUpdateEvent,
    ).finish(),
  );

  const defaultMarketEventBinary: Uint8Array = Uint8Array.from(MarketEventV1.encode(
    defaultMarketCreate,
  ).finish());

  const defaultAssetEventBinary: Uint8Array = Uint8Array.from(AssetCreateEventV1.encode(
    defaultAssetCreateEvent,
  ).finish());

  const transactionIndex0: number = 0;
  const transactionIndex1: number = 1;
  const eventIndex0: number = 0;
  const eventIndex1: number = 1;

  const events: IndexerTendermintEvent[] = [
    createIndexerTendermintEvent(
      DydxIndexerSubtypes.FUNDING,
      defaultMarketEventBinary,
      -1,
      eventIndex0,
    ),
    createIndexerTendermintEvent(
      DydxIndexerSubtypes.SUBACCOUNT_UPDATE,
      defaultSubaccountUpdateEventBinary,
      transactionIndex0,
      eventIndex0,
    ),
    createIndexerTendermintEvent(
      DydxIndexerSubtypes.ASSET,
      defaultAssetEventBinary,
      transactionIndex0,
      eventIndex1,
    ),
    createIndexerTendermintEvent(
      DydxIndexerSubtypes.MARKET,
      defaultMarketEventBinary,
      transactionIndex1,
      eventIndex0,
    ),
  ];

  it.each([
    [0, 0, 0],
    [1, 2, 3],
    [9, 8, 7],
  ])('dydx_event_id_from_parts (%d, %d, %d)', async (blockHeight: number, transactionIndex: number, eventIndex: number) => {
    const result = await getSingleRawQueryResultRow(`SELECT dydx_event_id_from_parts(${blockHeight}, ${transactionIndex}, ${eventIndex}) AS result;`);
    expect(result).toEqual(TendermintEventTable.createEventId(
      `${blockHeight}`,
      transactionIndex,
      eventIndex,
    ));
  });

  it.each([
    Long.fromNumber(1_000_000_000, true),
    Long.fromNumber(1_000_000_000, false),
    Long.fromNumber(-1_000_000_000, false),
  ])('dydx_from_jsonlib_long (%s)', async (value: Long) => {
    const result = await getSingleRawQueryResultRow(`SELECT dydx_trim_scale(dydx_from_jsonlib_long('${JSON.stringify(value)}')) AS result`);
    expect(result).toEqual(value.toString());
  });

  it.each([
    ['SIDE_BUY', IndexerOrder_Side.SIDE_BUY],
    ['SIDE_SELL', IndexerOrder_Side.SIDE_SELL],
  ])('dydx_from_protocol_order_side(%s)', async (_name: string, value: IndexerOrder_Side) => {
    const result = await getSingleRawQueryResultRow(`SELECT dydx_from_protocol_order_side('${value}') AS result`);
    expect(result).toEqual(protocolTranslations.protocolOrderSideToOrderSide(value));
  });

  it.each([
    ['TIME_IN_FORCE_UNSPECIFIED', IndexerOrder_TimeInForce.TIME_IN_FORCE_UNSPECIFIED],
    ['TIME_IN_FORCE_IOC', IndexerOrder_TimeInForce.TIME_IN_FORCE_IOC],
    ['TIME_IN_FORCE_POST_ONLY', IndexerOrder_TimeInForce.TIME_IN_FORCE_POST_ONLY],
    ['TIME_IN_FORCE_FILL_OR_KILL', IndexerOrder_TimeInForce.TIME_IN_FORCE_FILL_OR_KILL],
    ['UNRECOGNIZED', IndexerOrder_TimeInForce.UNRECOGNIZED],
  ])('dydx_from_protocol_time_in_force (%s)', async (_name: string, value: IndexerOrder_TimeInForce) => {
    const result = await getSingleRawQueryResultRow(`SELECT dydx_from_protocol_time_in_force('${value}') AS result`);
    expect(result).toEqual(protocolTranslations.protocolOrderTIFToTIF(value));
  });

  it.each([
    ['LIMIT', 32, IndexerOrder_ConditionType.UNRECOGNIZED],
    ['LIMIT', 32, IndexerOrder_ConditionType.CONDITION_TYPE_UNSPECIFIED],
    ['TAKE_PROFIT', 32, IndexerOrder_ConditionType.CONDITION_TYPE_TAKE_PROFIT],
    ['STOP_LIMIT', 32, IndexerOrder_ConditionType.CONDITION_TYPE_STOP_LOSS],
    ['LIMIT', 0, IndexerOrder_ConditionType.UNRECOGNIZED],
    ['LIMIT', 64, IndexerOrder_ConditionType.UNRECOGNIZED],
    ['TWAP', 128, IndexerOrder_ConditionType.UNRECOGNIZED],
    ['TWAP_SUBORDER', 256, IndexerOrder_ConditionType.UNRECOGNIZED],
  ])('dydx_protocol_convert_to_order_type (%s)', async (_name: string, orderFlags: number, value: IndexerOrder_ConditionType) => {
    const result = await getSingleRawQueryResultRow(`SELECT dydx_protocol_convert_to_order_type(${orderFlags}, '${value}'::jsonb) AS result`);
    expect(result).toEqual(
      protocolTranslations.protocolConditionTypeToOrderType(value, orderFlags),
    );
  });

  it.each([
    '0', '1', '-1', '10000000000000000000000000000', '-20000000000000000000000000000',
  ])('dydx_from_serializable_int (%s)', async (value: string) => {
    const byteArray = bigIntToBytes(BigInt(value));
    const result = await getSingleRawQueryResultRow(
      `SELECT dydx_trim_scale(dydx_from_serializable_int('${JSON.stringify(byteArray)}')) AS result`);
    expect(result).toEqual(value);
  });

  it.each([
    ['same amount of rounded decimal places (0)', 0, 1, 0, 1],
    ['same amount of rounded decimal places (big.js DP = 20)', 1, 2, 3, 4],
    ['first price is null', null, 1, 10, 2],
    ['second price is null', 3, 1, null, 2],
  ])('dydx_get_weighted_average (%s)', async (_name: string, firstPrice: number | null, firstWeight: number, secondPrice: number | null, secondWeight: number) => {
    const result = await getSingleRawQueryResultRow(`SELECT dydx_get_weighted_average(${JSON.stringify(firstPrice)}, ${firstWeight}, ${JSON.stringify(secondPrice)}, ${secondWeight}) AS result`);
    expect(result).toEqual(getWeightedAverage(firstPrice ? `${firstPrice}` : '0', `${firstWeight}`, secondPrice ? `${secondPrice}` : '0', `${secondWeight}`));
  });

  it.each([
    [PositionSide.LONG, OrderSide.BUY],
    [PositionSide.LONG, OrderSide.SELL],
    [PositionSide.SHORT, OrderSide.BUY],
    [PositionSide.SHORT, OrderSide.SELL],
  ])('dydx_perpetual_position_and_order_side_matching (%s, %s)', async (perpSide: PositionSide, orderSide: OrderSide) => {
    const result = await getSingleRawQueryResultRow(`SELECT dydx_perpetual_position_and_order_side_matching('${perpSide}', '${orderSide}') AS result;`);
    expect(result).toEqual(perpetualPositionAndOrderSideMatching(perpSide, orderSide));
  });

  it.each([
    ['0', '0'],
    ['0.', '0'],
    ['0.000', '0'],
    ['10', '10'],
    ['-10', '-10'],
    ['1.23456789012345678901234567890', '1.2345678901234567890123456789'],
    ['-1.2300', '-1.23'],
  ])('dydx_trim_scale (%s)', async (value: string, expected) => {
    const result = await getSingleRawQueryResultRow(`SELECT dydx_trim_scale('${value}') AS result;`);
    expect(result).toEqual(expected);
  });

  it.each([
    'foo',
    'bar',
  ])('dydx_uuid (%s)', async (value: string) => {
    const result = await getSingleRawQueryResultRow(`SELECT dydx_uuid('${value}') AS result`);
    expect(result).toEqual(uuid.getUuid(Buffer.from(value, BUFFER_ENCODING_UTF_8)));
  });

  it.each([
    [
      {
        owner: testConstants.defaultSubaccount.address,
        number: testConstants.defaultSubaccount.subaccountNumber,
      },
      '0',
    ],
    [
      {
        owner: testConstants.defaultSubaccount2.address,
        number: testConstants.defaultSubaccount2.subaccountNumber,
      },
      '1',
    ],
  ])('dydx_uuid_from_asset_position_parts (%s)', async (subaccountId: IndexerSubaccountId, assetId: string) => {
    const subaccountUuid = SubaccountTable.subaccountIdToUuid(subaccountId);
    const result = await getSingleRawQueryResultRow(
      `SELECT dydx_uuid_from_asset_position_parts('${subaccountUuid}', '${assetId}') AS result`);
    expect(result).toEqual(AssetPositionTable.uuid(subaccountUuid, assetId));
  });

  it.each([
    [Liquidity.TAKER, 1, 2, 3],
    [Liquidity.MAKER, 4, 5, 6],
  ])('dydx_uuid_from_fill_event_parts (%s)', async (liquidity: Liquidity, blockHeight: number, transactionIndex: number, eventIndex: number) => {
    const eventId = TendermintEventTable.createEventId(`${blockHeight}`, transactionIndex, eventIndex);
    const result = await getSingleRawQueryResultRow(
      `SELECT dydx_uuid_from_fill_event_parts('\\x${eventId.toString('hex')}'::bytea, '${liquidity}') AS result`);
    expect(result).toEqual(FillTable.uuid(eventId, liquidity));
  });

  it.each([
    [1, 2, 3, 4],
    [5, 6, 7, 8],
  ])('dydx_uuid_from_funding_index_update_parts (%s, %s, %s, %s)', async (blockHeight: number, transactionIndex: number, eventIndex: number, perpetualId: number) => {
    const eventId = TendermintEventTable.createEventId(`${blockHeight}`, transactionIndex, eventIndex);
    const result = await getSingleRawQueryResultRow(
      `SELECT dydx_uuid_from_funding_index_update_parts('${blockHeight}', '\\x${eventId.toString('hex')}'::bytea, '${perpetualId}') AS result`);
    expect(result).toEqual(FundingIndexUpdatesTable.uuid(`${blockHeight}`, eventId, `${perpetualId}`));
  });

  it.each([
    {
      subaccountId: {
        owner: testConstants.defaultSubaccount.address,
        number: testConstants.defaultSubaccount.subaccountNumber,
      },
      clientId: 3,
      orderFlags: 4,
      clobPairId: 5,
    },
  ])('dydx_uuid_from_order_id and parts (%s)', async (orderId: IndexerOrderId) => {
    let result = await getSingleRawQueryResultRow(`SELECT dydx_uuid_from_order_id('${JSON.stringify(orderId)}') AS result`);
    expect(result).toEqual(OrderTable.orderIdToUuid(orderId));

    result = await getSingleRawQueryResultRow(`SELECT dydx_uuid_from_order_id_parts('${SubaccountTable.subaccountIdToUuid(orderId.subaccountId!)}', '${orderId.clientId}', '${orderId.clobPairId}', '${orderId.orderFlags}') AS result`);
    expect(result).toEqual(OrderTable.orderIdToUuid(orderId));
  });

  it.each([
    [
      {
        owner: testConstants.defaultSubaccount.address,
        number: testConstants.defaultSubaccount.subaccountNumber,
      },
      1,
      2,
      3,
    ],
    [
      {
        owner: testConstants.defaultSubaccount2.address,
        number: testConstants.defaultSubaccount2.subaccountNumber,
      },
      4,
      5,
      6,
    ],
  ])('dydx_uuid_from_perpetual_position_parts (%s, %s, %s, %s)', async (subaccountId: IndexerSubaccountId, blockHeight: number, transactionIndex: number, eventIndex: number) => {
    const subaccountUuid = SubaccountTable.subaccountIdToUuid(subaccountId);
    const eventId = TendermintEventTable.createEventId(`${blockHeight}`, transactionIndex, eventIndex);
    const result = await getSingleRawQueryResultRow(
      `SELECT dydx_uuid_from_perpetual_position_parts('${subaccountUuid}', '\\x${eventId.toString('hex')}'::bytea) AS result`);
    expect(result).toEqual(PerpetualPositionTable.uuid(subaccountUuid, eventId));
  });

  it.each([
    {
      owner: testConstants.defaultSubaccount.address,
      number: testConstants.defaultSubaccount.subaccountNumber,
    },
  ])('dydx_uuid_from_subaccount_id and parts (%s)', async (subaccountId: IndexerSubaccountId) => {
    let result = await getSingleRawQueryResultRow(`SELECT dydx_uuid_from_subaccount_id('${JSON.stringify(subaccountId)}') AS result`);
    expect(result).toEqual(SubaccountTable.subaccountIdToUuid(subaccountId));

    result = await getSingleRawQueryResultRow(`SELECT dydx_uuid_from_subaccount_id_parts('${subaccountId.owner}', '${subaccountId.number}') AS result`);
    expect(result).toEqual(SubaccountTable.subaccountIdToUuid(subaccountId));
  });

  it.each([
    [
      {
        owner: testConstants.defaultSubaccount.address,
        number: testConstants.defaultSubaccount.subaccountNumber,
      },
      {
        owner: testConstants.defaultSubaccount2.address,
        number: testConstants.defaultSubaccount2.subaccountNumber,
      },
      undefined,
      undefined,
    ],
    [
      {
        owner: testConstants.defaultSubaccount2.address,
        number: testConstants.defaultSubaccount2.subaccountNumber,
      },
      undefined,
      'senderWallet',
      undefined,
    ],
    [
      {
        owner: testConstants.defaultSubaccount.address,
        number: testConstants.defaultSubaccount.subaccountNumber,
      },
      undefined,
      undefined,
      'recipientWallet',
    ],
    [
      undefined,
      undefined,
      'senderWallet',
      'recipientWallet',
    ],
  ])('dydx_uuid_from_transfer_parts (%s, %s, %s, %s)', async (
    senderSubaccountId: IndexerSubaccountId | undefined,
    recipientSubaccountId: IndexerSubaccountId | undefined,
    senderWalletAddress: string | undefined,
    recipientWalletAddress: string | undefined) => {
    const eventId: Buffer = TendermintEventTable.createEventId('1', 2, 3);
    const assetId: string = '0';
    const senderSubaccountUuid: string | undefined = senderSubaccountId
      ? SubaccountTable.subaccountIdToUuid(senderSubaccountId) : undefined;
    const recipientSubaccountUuid: string | undefined = recipientSubaccountId
      ? SubaccountTable.subaccountIdToUuid(recipientSubaccountId) : undefined;
    const result = await getSingleRawQueryResultRow(
      `SELECT dydx_uuid_from_transfer_parts('\\x${eventId.toString('hex')}'::bytea, '${assetId}', ${senderSubaccountUuid ? `'${senderSubaccountUuid}'` : 'NULL'}, ${recipientSubaccountUuid ? `'${recipientSubaccountUuid}'` : 'NULL'}, ${senderWalletAddress ? `'${senderWalletAddress}'` : 'NULL'}, ${recipientWalletAddress ? `'${recipientWalletAddress}'` : 'NULL'}) AS result`);
    expect(result).toEqual(TransferTable.uuid(
      eventId,
      assetId,
      senderSubaccountUuid,
      recipientSubaccountUuid,
      senderWalletAddress,
      recipientWalletAddress,
    ));
  });

  it.each([
    {
      event: { transactionIndex: 123 },
      expectedResult: 123,
    },
    {
      event: { blockEvent: IndexerTendermintEvent_BlockEvent.BLOCK_EVENT_BEGIN_BLOCK },
      expectedResult: -2,
    },
    {
      event: { blockEvent: IndexerTendermintEvent_BlockEvent.BLOCK_EVENT_END_BLOCK },
      expectedResult: -1,
    },
    {
      event: { blockEvent: '3' },
      expectedError: 'Received V4 event with invalid block event type: "3"',
    },
    {
      event: {},
      expectedError: 'Either transactionIndex or blockEvent must be defined in IndexerTendermintEvent',
    },
  ])('dydx_tendermint_event_to_transaction_index should return the expected result', async (
    { event, expectedResult, expectedError },
  ) => {
    try {
      const result = await getSingleRawQueryResultRow(
        `SELECT dydx_tendermint_event_to_transaction_index('${JSON.stringify(event)}') AS result`);
      if (expectedError) {
        throw new Error('Expected an error but got a result.');
      }
      expect(result).toEqual(expectedResult);
    } catch (error) {
      if (!expectedError) {
        throw error;
      }
      expect(error.message).toContain(expectedError);
    }
  });

  it.each([
    [
      '123456',
      123,
    ],
  ])('dydx_uuid_from_transaction_parts (%s, %s)', async (blockHeight: string, transactionIndex: number) => {
    const result = await getSingleRawQueryResultRow(
      `SELECT dydx_uuid_from_transaction_parts('${blockHeight}', '${transactionIndex}') AS result`);
    expect(result).toEqual(
      TransactionTable.uuid(blockHeight, transactionIndex),
    );
  });

  it.each([
    [
      123,
      '123456',
    ],
  ])('dydx_uuid_from_oracle_price_parts (%s, %s)', async (marketId: number, blockHeight: string) => {
    const result = await getSingleRawQueryResultRow(
      `SELECT dydx_uuid_from_oracle_price_parts('${marketId}', '${blockHeight}') AS result`);
    expect(result).toEqual(
      OraclePriceTable.uuid(marketId, blockHeight),
    );
  });

  it('dydx_clob_pair_status_to_market_status should convert all statuses', async () => {
    for (const [key, value] of Object.entries(CLOB_STATUS_TO_MARKET_STATUS)) {
      const result = await getSingleRawQueryResultRow(
        `SELECT dydx_clob_pair_status_to_market_status('${key}') AS result`);
      expect(result).toEqual(value);
    }
  });

  it('dydx_create_transaction.sql should insert a transaction and return correct jsonb', async () => {
    const transactionHash: string = 'txnhash';
    const blockHeight: string = '1';
    const transactionIndex: number = 123;

    const returnedJsonb = await getSingleRawQueryResultRow(`SELECT dydx_create_transaction('${transactionHash}', '${blockHeight}', ${transactionIndex}) AS result`);

    const transactions: TransactionFromDatabase[] = await TransactionTable.findAll(
      {},
      [],
      { readReplica: true },
    );

    expect(transactions.length).toEqual(1);
    const transaction = transactions[0];
    expect(transaction).toEqual(expect.objectContaining({
      transactionHash,
      blockHeight,
      transactionIndex,
    }));
    expect(returnedJsonb).toEqual(expect.objectContaining({
      transactionHash,
      blockHeight: Number(blockHeight),
      transactionIndex,
    }));
  });

  it('dydx_create_tendermint_event.sql should insert a tendermint event and return correct jsonb', async () => {
    await BlockTable.create(testConstants.defaultBlock);
    const transactionIndex: number = 0;
    const eventIndex: number = 0;

    const indexerTendermintEvent: IndexerTendermintEvent = createIndexerTendermintEvent(
      DydxIndexerSubtypes.ASSET,
      AssetCreateEventV1.encode(defaultAssetCreateEvent).finish(),
      transactionIndex,
      eventIndex,
    );

    const result = await getSingleRawQueryResultRow(
      `SELECT dydx_create_tendermint_event('${JSON.stringify(indexerTendermintEvent)}', '${testConstants.defaultBlock.blockHeight}') AS result`,
    );
    const tendermintEvents: TendermintEventFromDatabase[] = await TendermintEventTable.findAll(
      {},
      [],
      { readReplica: true },
    );

    expect(tendermintEvents.length).toEqual(1);
    const tendermintEvent = tendermintEvents[0];
    expect(tendermintEvent).toEqual(expect.objectContaining({
      blockHeight: testConstants.defaultBlock.blockHeight,
      transactionIndex,
      eventIndex,
      id: TendermintEventTable.createEventId(
        testConstants.defaultBlock.blockHeight,
        transactionIndex,
        eventIndex,
      ),
    }));
    expect(result).toEqual(expect.objectContaining({
      blockHeight: Number(testConstants.defaultBlock.blockHeight),
      transactionIndex,
      eventIndex,
    }));
  });

  it('dydx_create_initial_rows_for_tendermint_block.sql should insert the initial rows correctly', async () => {
    const blockHeight = '1';
    const txHashes = [defaultTxHash, defaultTxHash2];
    const dateTimeIso = '2020-01-01T00:00:00.000Z';
    await getSingleRawQueryResultRow(
      `SELECT dydx_create_initial_rows_for_tendermint_block(${blockHeight}, '${dateTimeIso}', '${JSON.stringify(txHashes)}', '${JSON.stringify(events)}')`,
    );
    // Validate blocks table
    const blocks = await BlockTable.findAll({}, [], { readReplica: true });
    expect(blocks.length).toEqual(1);
    expect(blocks[0]).toEqual(expect.objectContaining({
      blockHeight,
      time: dateTimeIso,
    }));

    // Validate transactions table
    const transactions: TransactionFromDatabase[] = await TransactionTable.findAll(
      {},
      [],
      { readReplica: true },
    );
    expect(transactions.length).toEqual(txHashes.length);
    txHashes.forEach((txHash, index) => {
      expect(transactions[index].transactionHash).toEqual(txHash);
    });

    // Validate tendermint_events table
    const tendermintEvents: TendermintEventFromDatabase[] = await TendermintEventTable.findAll(
      {},
      [],
      { readReplica: true },
    );
    expect(tendermintEvents.length).toEqual(events.length);
    events.forEach((event, index) => {
      expect(tendermintEvents[index]).toEqual(expect.objectContaining({
        blockHeight,
        transactionIndex: indexerTendermintEventToTransactionIndex(event),
        eventIndex: event.eventIndex,
        id: TendermintEventTable.createEventId(
          blockHeight,
          indexerTendermintEventToTransactionIndex(event),
          event.eventIndex,
        ),
      }));
    });
  });
});

async function getSingleRawQueryResultRow(query: string): Promise<object> {
  const queryResult = await storeHelpers.rawQuery(query, {}).catch((error: Error) => {
    throw error;
  });
  return queryResult.rows[0].result;
}
