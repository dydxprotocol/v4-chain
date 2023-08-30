import {
  IndexerTendermintEvent,
  IndexerOrder_Side,
  IndexerSubaccountId,
  IndexerOrder_TimeInForce,
  IndexerOrderId, IndexerTendermintEvent_BlockEvent,
} from '@dydxprotocol-indexer/v4-protos';
import {
  BUFFER_ENCODING_UTF_8,
  dbHelpers,
  AssetPositionTable,
  PerpetualPositionTable,
  Liquidity,
  OrderSide,
  PositionSide,
  TendermintEventTable,
  FillTable,
  OrderTable,
  protocolTranslations,
  SubaccountTable,
  storeHelpers,
  testConstants,
  uuid,
} from '@dydxprotocol-indexer/postgres';

import { createPostgresFunctions } from '../../src/helpers/postgres/postgres-functions';
import Long from 'long';
import {
  getWeightedAverage,
  indexerTendermintEventToTransactionIndex,
  perpetualPositionAndOrderSideMatching,
} from '../../src/lib/helper';
import { bigIntToBytes } from '@dydxprotocol-indexer/v4-proto-parser';

describe('SQL Function Tests', () => {
  beforeAll(async () => {
    await dbHelpers.migrate();
    await createPostgresFunctions();
  });

  afterAll(async () => {
    await dbHelpers.teardown();
    jest.resetAllMocks();
  });

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
    { transactionIndex: 5 } as IndexerTendermintEvent,
    {
      blockEvent: IndexerTendermintEvent_BlockEvent.BLOCK_EVENT_BEGIN_BLOCK,
    } as IndexerTendermintEvent,
    {
      blockEvent: IndexerTendermintEvent_BlockEvent.BLOCK_EVENT_END_BLOCK,
    } as IndexerTendermintEvent,
  ])('dydx_event_to_transaction_index (%s)', async (event: IndexerTendermintEvent) => {
    const result = await getSingleRawQueryResultRow(`SELECT dydx_event_to_transaction_index('${JSON.stringify(event)}') AS result;`);
    expect(result).toEqual(indexerTendermintEventToTransactionIndex(event));
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
  ])('dydx_uuid_from_perpetual_position_parts (%s)', async (subaccountId: IndexerSubaccountId, blockHeight: number, transactionIndex: number, eventIndex: number) => {
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
});

async function getSingleRawQueryResultRow(query: string): Promise<object> {
  const queryResult = await storeHelpers.rawQuery(query, {}).catch((error) => {
    throw error;
  });
  return queryResult.rows[0].result;
}
