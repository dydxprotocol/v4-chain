import {
  dbHelpers,
  testMocks,
  TransferTable,
  BlockTable,
  PnlTable,
  PersistentCacheTable,
  PersistentCacheKeys,
  FundingPaymentsTable,
  FillTable,
  OrderSide,
  OrderTable,
  TendermintEventTable,
  TendermintEventCreateObject,
} from '@dydxprotocol-indexer/postgres';
import { DateTime } from 'luxon';
import updatePnlTask from '../../src/tasks/update-pnl';
import {
  defaultSubaccountId,
  defaultSubaccountId2,
  defaultBlock,
  defaultTransfer,
  defaultFundingPayment,
  defaultPerpetualMarket,
  defaultTendermintEventId,
  defaultTendermintEventId2,
  defaultTendermintEventId3,
  defaultTendermintEventId4,
  defaultFill,
  defaultOrder,
  defaultPerpetualMarket2,
} from '@dydxprotocol-indexer/postgres/build/__tests__/helpers/constants';

describe('update-pnl', () => {
  // Common date helpers
  const JUNE_1 = DateTime.utc(2022, 6, 1).toISO();
  const JUNE_2 = DateTime.utc(2022, 6, 2).toISO();

  let defaultTendermintEventId5: Buffer;
  let defaultTendermintEventId6: Buffer;

  beforeAll(async () => {
    await dbHelpers.migrate();
  });

  beforeEach(async () => {
    await dbHelpers.clearData();
    await testMocks.seedData();

    // Create common blocks used by all tests
    const blockHeights: string[] = ['0', '3', '4', '5', '6', '7', '8', '9', '10', '12', '15'];

    await Promise.all(blockHeights.map((height) => BlockTable.create({
      ...defaultBlock,
      blockHeight: height,
    }),
    ));

    // Create events to be used for fills
    const defaultTendermintEvent5: TendermintEventCreateObject = {
      blockHeight: '3',
      transactionIndex: 0,
      eventIndex: 0,
    };
    const defaultTendermintEvent6: TendermintEventCreateObject = {
      blockHeight: '3',
      transactionIndex: 1,
      eventIndex: 1,
    };

    defaultTendermintEventId5 = await TendermintEventTable.createEventId('3', 0, 0);
    defaultTendermintEventId6 = await TendermintEventTable.createEventId('3', 1, 1);
    await Promise.all([
      TendermintEventTable.create(defaultTendermintEvent5),
      TendermintEventTable.create(defaultTendermintEvent6),
    ]);

    // Create order for the fills
    await OrderTable.create(defaultOrder);

    // Create a common transfer at height 1 for all tests to use
    await TransferTable.create({
      ...defaultTransfer,
      size: '30000',
      createdAtHeight: '1',
      createdAt: JUNE_1,
      senderSubaccountId: defaultTransfer.recipientSubaccountId,
      recipientSubaccountId: defaultTransfer.senderSubaccountId,
    });
  });

  afterAll(async () => {
    await dbHelpers.teardown();
    jest.resetAllMocks();
  });

  afterEach(async () => {
    await dbHelpers.clearData();
  });

  /**
   * Helper to verify PNL record
   */
  function verifyPnlRecord(
    record: any,
    expectedValues: {
      equity: string,
      netTransfers: string,
      totalPnl: string,
    },
  ) {
    expect(record.equity).toBe(expectedValues.equity);
    expect(record.netTransfers).toBe(expectedValues.netTransfers);
    expect(record.totalPnl).toBe(expectedValues.totalPnl);
  }

  /**
   * Helper to verify cache was updated correctly
   */
  async function verifyCache(expectedHeight: string) {
    const cache = await PersistentCacheTable.findById(
      PersistentCacheKeys.PNL_LAST_PROCESSED_HEIGHT,
    );
    expect(cache).toBeDefined();
    expect(cache?.value).toBe(expectedHeight);
  }

  /**
   * Helper to find PNL records by height
   */
  function findPnlRecords(records: any[], height: string) {
    const recordsAtHeight = records.filter((r) => r.createdAtHeight === height);
    const subaccount1Pnl = recordsAtHeight.find((r) => r.subaccountId === defaultSubaccountId);
    const subaccount2Pnl = recordsAtHeight.find((r) => r.subaccountId === defaultSubaccountId2);
    return {
      recordsAtHeight,
      subaccount1Pnl,
      subaccount2Pnl,
    };
  }

  it('calculates initial PNL state with transfers and checks multiple heights', async () => {
    // Create funding payments for both subaccounts at both heights
    await FundingPaymentsTable.create({
      ...defaultFundingPayment,
      createdAtHeight: '1',
      createdAt: JUNE_1,
      payment: '10',
      size: '2',
      oraclePrice: '10000', // Position value: 2 * 10000 = 20000
    });

    await FundingPaymentsTable.create({
      ...defaultFundingPayment,
      createdAtHeight: '2',
      createdAt: JUNE_2,
      payment: '5',
      size: '2',
      oraclePrice: '11000', // Position value: 2 * 11000 = 22000
    });

    // Create a fill to represent buying the BTC position
    await FillTable.create({
      ...defaultFill,
      subaccountId: defaultSubaccountId,
      side: OrderSide.BUY,
      size: '2', // 2 BTC
      price: '10000', // at $10,000 each
      quoteAmount: '20000', // Total $20,000
      createdAtHeight: '1',
      createdAt: JUNE_1,
    });

    // Run the PNL update task
    await updatePnlTask();

    // Check that PNL entries were created
    const pnlRecords = await PnlTable.findAll({}, []);

    // First check height 1
    const { recordsAtHeight: recordsAtHeight1, subaccount1Pnl: subaccount1PnlAtHeight1, subaccount2Pnl: subaccount2PnlAtHeight1 } = findPnlRecords(pnlRecords.results, '1');

    expect(recordsAtHeight1.length).toBe(2); // Both subaccounts should have PNL
    expect(subaccount1PnlAtHeight1).toBeDefined();
    expect(subaccount2PnlAtHeight1).toBeDefined();

    verifyPnlRecord(subaccount1PnlAtHeight1, {
      netTransfers: '30000',
      totalPnl: '10',
      equity: '30010',
    });

    verifyPnlRecord(subaccount2PnlAtHeight1, {
      netTransfers: '-30000',
      totalPnl: '0',
      equity: '-30000',
    });

    // Now check height 2
    const { subaccount1Pnl: subaccount1PnlAtHeight2, subaccount2Pnl: subaccount2PnlAtHeight2 } = findPnlRecords(pnlRecords.results, '2');

    verifyPnlRecord(subaccount1PnlAtHeight2, {
      netTransfers: '30000',
      totalPnl: '2015',
      equity: '32015',
    });

    verifyPnlRecord(subaccount2PnlAtHeight2, {
      netTransfers: '-30000',
      totalPnl: '0',
      equity: '-30000',
    });

    // Verify cache was updated
    await verifyCache('2');
  });

  it('correctly sums funding payments across multiple positions', async () => {
    // Create funding payments for the first position (BTC)
    await FundingPaymentsTable.create({
      ...defaultFundingPayment,
      createdAtHeight: '1',
      createdAt: JUNE_1,
      payment: '10',
      size: '2',
      oraclePrice: '10000', // Position value: 2 * 10000 = 20000
      perpetualId: defaultPerpetualMarket.id,
      ticker: defaultPerpetualMarket.ticker, // 'BTC-USD'
    });

    // Create funding payments for the second position (ETH)
    await FundingPaymentsTable.create({
      ...defaultFundingPayment,
      createdAtHeight: '1',
      createdAt: JUNE_1,
      payment: '5',
      size: '3',
      oraclePrice: '1000', // Position value: 3 * 1000 = 3000
      perpetualId: defaultPerpetualMarket2.id,
      ticker: defaultPerpetualMarket2.ticker, // 'ETH-USD'
      subaccountId: defaultSubaccountId,
    });

    // Create fills to represent buying the positions
    // BTC position
    await FillTable.create({
      ...defaultFill,
      subaccountId: defaultSubaccountId,
      side: OrderSide.BUY,
      size: '2', // 2 BTC
      price: '10000', // at $10,000 each
      quoteAmount: '20000', // Total $20,000
      createdAtHeight: '1',
      createdAt: JUNE_1,
      eventId: defaultTendermintEventId,
    });

    // ETH position
    await FillTable.create({
      ...defaultFill,
      subaccountId: defaultSubaccountId,
      side: OrderSide.BUY,
      size: '3', // 3 ETH
      price: '1000', // at $1,000 each
      quoteAmount: '3000', // Total $3,000
      createdAtHeight: '1',
      createdAt: JUNE_1,
      eventId: defaultTendermintEventId2,
    });

    // Run the PNL update task
    await updatePnlTask();

    // Check that PNL entries were created
    const pnlRecords = await PnlTable.findAll({}, []);

    // Look at height 1
    const { subaccount1Pnl: subaccount1PnlAtHeight1 } = findPnlRecords(pnlRecords.results, '1');

    verifyPnlRecord(subaccount1PnlAtHeight1, {
      netTransfers: '30000',
      totalPnl: '15', // 10 + 5
      equity: '30015', // 30000 + 15
    });

    // Verify cache was updated
    await verifyCache('1');
  });

  it('calculates PNL for a long position that is sold between blocks', async () => {
    // Funding payments at blocks 5 and 10 only (nothing at time 0)
    await FundingPaymentsTable.create({
      ...defaultFundingPayment,
      createdAtHeight: '5',
      createdAt: DateTime.utc(2022, 6, 5).toISO(),
      payment: '10',
      size: '3',
      oraclePrice: '10000', // Position value at block 5: 3 * 10000 = 30000
    });

    await FundingPaymentsTable.create({
      ...defaultFundingPayment,
      createdAtHeight: '10',
      createdAt: DateTime.utc(2022, 6, 10).toISO(),
      payment: '2',
      size: '0', // Position is closed by this point
      oraclePrice: '11000',
    });

    // Create a fill to buy at block 1
    await FillTable.create({
      ...defaultFill,
      subaccountId: defaultSubaccountId,
      side: OrderSide.BUY,
      size: '3', // 3 BTC
      price: '9000', // at $9,000 each
      quoteAmount: '27000', // Total $27,000
      createdAtHeight: '1',
      createdAt: JUNE_1,
      eventId: defaultTendermintEventId,
    });

    // Create a fill to sell between blocks 5-10 (at block 7)
    await FillTable.create({
      ...defaultFill,
      subaccountId: defaultSubaccountId,
      side: OrderSide.SELL,
      size: '3', // Sell 3 BTC
      price: '10500', // at $10,500 each
      quoteAmount: '31500', // Total $31,500
      createdAtHeight: '7',
      createdAt: DateTime.utc(2022, 6, 7).toISO(),
      eventId: defaultTendermintEventId2,
    });

    // Run the PNL update task
    await updatePnlTask();

    // Check that PNL entries were created
    const pnlRecords = await PnlTable.findAll({}, []);

    // Check block 5
    const { subaccount1Pnl: subaccount1PnlAtHeight5 } = findPnlRecords(pnlRecords.results, '5');

    verifyPnlRecord(subaccount1PnlAtHeight5, {
      netTransfers: '30000',
      totalPnl: '3010',
      equity: '33010',
    });

    // Check block 10
    const { subaccount1Pnl: subaccount1PnlAtHeight10 } = findPnlRecords(pnlRecords.results, '10');

    verifyPnlRecord(subaccount1PnlAtHeight10, {
      netTransfers: '30000',
      totalPnl: '4512',
      equity: '34512',
    });

    // Verify cache was updated to the latest height
    await verifyCache('10');
  });

  it('calculates PNL for a short position that is closed between blocks', async () => {
    // Funding payments at blocks 5 and 10
    await FundingPaymentsTable.create({
      ...defaultFundingPayment,
      createdAtHeight: '5',
      createdAt: DateTime.utc(2022, 6, 5).toISO(),
      payment: '-8', // Negative funding for short position
      size: '-4', // Negative size for short position
      oraclePrice: '10000', // Position value at block 5: -4 * 10000 = -40000
    });

    await FundingPaymentsTable.create({
      ...defaultFundingPayment,
      createdAtHeight: '10',
      createdAt: DateTime.utc(2022, 6, 10).toISO(),
      payment: '2',
      size: '0',
      oraclePrice: '11000',
    });

    // Create a fill to sell at block 3
    await FillTable.create({
      ...defaultFill,
      subaccountId: defaultSubaccountId,
      side: OrderSide.SELL,
      size: '4', // 4 BTC
      price: '9500', // at $9,500 each
      quoteAmount: '38000', // Total $38,000
      createdAtHeight: '3',
      createdAt: DateTime.utc(2022, 6, 3).toISO(),
      eventId: defaultTendermintEventId,
    });

    // Create a fill to buy back between blocks 5-10 (at block 8)
    await FillTable.create({
      ...defaultFill,
      subaccountId: defaultSubaccountId,
      side: OrderSide.BUY,
      size: '4', // Buy 4 BTC
      price: '9200', // at $9,200 each
      quoteAmount: '36800', // Total $36,800
      createdAtHeight: '8',
      createdAt: DateTime.utc(2022, 6, 8).toISO(),
      eventId: defaultTendermintEventId2,
    });

    // Run the PNL update task
    await updatePnlTask();

    // Check that PNL entries were created
    const pnlRecords = await PnlTable.findAll({}, []);

    // Check block 5
    const { subaccount1Pnl: subaccount1PnlAtHeight5 } = findPnlRecords(pnlRecords.results, '5');

    verifyPnlRecord(subaccount1PnlAtHeight5, {
      netTransfers: '30000',
      totalPnl: '-2008',
      equity: '27992',
    });

    // Check block 10
    const { subaccount1Pnl: subaccount1PnlAtHeight10 } = findPnlRecords(pnlRecords.results, '10');

    verifyPnlRecord(subaccount1PnlAtHeight10, {
      netTransfers: '30000',
      totalPnl: '1194',
      equity: '31194',
    });

    // Verify cache was updated to the latest height
    await verifyCache('10');
  });

  it('calculates PNL for a position that changes from long to short', async () => {
    // Funding payments at blocks 5 and 10
    await FundingPaymentsTable.create({
      ...defaultFundingPayment,
      createdAtHeight: '5',
      createdAt: DateTime.utc(2022, 6, 5).toISO(),
      payment: '15', // Positive funding for long position
      size: '5', // Positive size for long position
      oraclePrice: '9800', // Position value at block 5: 5 * 9800 = 49000
    });

    await FundingPaymentsTable.create({
      ...defaultFundingPayment,
      createdAtHeight: '10',
      createdAt: DateTime.utc(2022, 6, 10).toISO(),
      payment: '-12', // Negative funding for short position
      size: '-2', // Negative size for short position
      oraclePrice: '9600', // Position value at block 10: -2 * 9600 = -19200
    });

    // Need extra transfer for this test as we need more funds
    await TransferTable.create({
      ...defaultTransfer,
      eventId: defaultTendermintEventId4,
      size: '20000',
      createdAtHeight: '1',
      createdAt: JUNE_1,
      senderSubaccountId: defaultTransfer.recipientSubaccountId,
      recipientSubaccountId: defaultTransfer.senderSubaccountId,
    });

    // Create a fill to buy at block 2
    await FillTable.create({
      ...defaultFill,
      subaccountId: defaultSubaccountId,
      side: OrderSide.BUY,
      size: '5', // 5 BTC
      price: '9500', // at $9,500 each
      quoteAmount: '47500', // Total $47,500
      createdAtHeight: '2',
      createdAt: DateTime.utc(2022, 6, 2).toISO(),
      eventId: defaultTendermintEventId,
    });

    // Create a fill to sell more between blocks 5-10 (at block 7) - shifting to a short position
    await FillTable.create({
      ...defaultFill,
      subaccountId: defaultSubaccountId,
      side: OrderSide.SELL,
      size: '7', // Sell 7 BTC (5 to close + 2 to short)
      price: '9700', // at $9,700 each
      quoteAmount: '67900', // Total $67,900
      createdAtHeight: '7',
      createdAt: DateTime.utc(2022, 6, 7).toISO(),
      eventId: defaultTendermintEventId2,
    });

    // Run the PNL update task
    await updatePnlTask();

    // Check that PNL entries were created
    const pnlRecords = await PnlTable.findAll({}, []);

    // Check block 5
    const { subaccount1Pnl: subaccount1PnlAtHeight5 } = findPnlRecords(pnlRecords.results, '5');

    verifyPnlRecord(subaccount1PnlAtHeight5, {
      netTransfers: '50000',
      totalPnl: '1515',
      equity: '51515',
    });

    // Check block 10
    const { subaccount1Pnl: subaccount1PnlAtHeight10 } = findPnlRecords(pnlRecords.results, '10');

    verifyPnlRecord(subaccount1PnlAtHeight10, {
      netTransfers: '50000',
      totalPnl: '1203',
      equity: '51203',
    });

    // Verify cache was updated to the latest height
    await verifyCache('10');
  });

  it('calculates PNL for operations on multiple different tickers', async () => {
    // Create funding payments for different tickers
    const btcFundingPayment = defaultFundingPayment;
    const ethFundingPayment = { ...defaultFundingPayment, ticker: 'ETH-USD' };

    // Funding payments for ETH and BTC
    await FundingPaymentsTable.create({
      ...btcFundingPayment,
      createdAtHeight: '10',
      createdAt: DateTime.utc(2022, 6, 10).toISO(),
      payment: '-7', // Negative funding for short BTC position
      size: '-1', // Negative size for short position
      oraclePrice: '9600', // BTC Position value at block 10: -1 * 9600 = -9600
    });

    await FundingPaymentsTable.create({
      ...ethFundingPayment,
      createdAtHeight: '5',
      createdAt: DateTime.utc(2022, 6, 5).toISO(),
      payment: '3',
      size: '10',
      oraclePrice: '520', // ETH Position value at block 5: 10 * 520 = 5200
    });

    await FundingPaymentsTable.create({
      ...ethFundingPayment,
      createdAtHeight: '10',
      createdAt: DateTime.utc(2022, 6, 10).toISO(),
      payment: '8',
      size: '15',
      oraclePrice: '550', // ETH Position value at block 10: 15 * 550 = 8250
    });

    // Trading operations for BTC and ETH
    await FillTable.create({
      ...defaultFill,
      subaccountId: defaultSubaccountId,
      side: OrderSide.BUY,
      size: '2',
      price: '9000',
      quoteAmount: '18000',
      createdAtHeight: '2',
      createdAt: DateTime.utc(2022, 6, 2).toISO(),
      eventId: defaultTendermintEventId,
    });

    await FillTable.create({
      ...defaultFill,
      subaccountId: defaultSubaccountId,
      side: OrderSide.SELL,
      size: '2',
      price: '9300',
      quoteAmount: '18600',
      createdAtHeight: '4',
      createdAt: DateTime.utc(2022, 6, 4).toISO(),
      eventId: defaultTendermintEventId2,
    });

    await FillTable.create({
      ...defaultFill,
      subaccountId: defaultSubaccountId,
      side: OrderSide.SELL,
      size: '1',
      price: '9550',
      quoteAmount: '9550',
      createdAtHeight: '7',
      createdAt: DateTime.utc(2022, 6, 7).toISO(),
      eventId: defaultTendermintEventId3,
    });

    await FillTable.create({
      ...defaultFill,
      subaccountId: defaultSubaccountId,
      side: OrderSide.BUY,
      size: '10',
      price: '510',
      quoteAmount: '5100',
      createdAtHeight: '3',
      createdAt: DateTime.utc(2022, 6, 3).toISO(),
      eventId: defaultTendermintEventId4,
    });

    await FillTable.create({
      ...defaultFill,
      subaccountId: defaultSubaccountId,
      side: OrderSide.BUY,
      size: '5',
      price: '530',
      quoteAmount: '2650',
      createdAtHeight: '8',
      createdAt: DateTime.utc(2022, 6, 8).toISO(),
      eventId: defaultTendermintEventId5,
    });

    // Run the PNL update task
    await updatePnlTask();

    // Check that PNL entries were created
    const pnlRecords = await PnlTable.findAll({}, []);

    // Check block 5
    const { subaccount1Pnl: subaccount1PnlAtHeight5 } = findPnlRecords(pnlRecords.results, '5');

    verifyPnlRecord(subaccount1PnlAtHeight5, {
      netTransfers: '30000',
      totalPnl: '703',
      equity: '30703',
    });

    // Check block 10
    const { subaccount1Pnl: subaccount1PnlAtHeight10 } = findPnlRecords(pnlRecords.results, '10');

    verifyPnlRecord(subaccount1PnlAtHeight10, {
      netTransfers: '30000',
      totalPnl: '1054',
      equity: '31054',
    });

    // Verify cache was updated to the latest height
    await verifyCache('10');
  });

  it('calculates PNL across multiple funding periods with various trades and transfers', async () => {
    // Funding payments for BTC across three periods
    const btcFundingPayment = defaultFundingPayment;
    const ethFundingPayment = { ...defaultFundingPayment, ticker: 'ETH-USD' };

    // First period (0-5)
    await FundingPaymentsTable.create({
      ...btcFundingPayment,
      createdAtHeight: '5',
      createdAt: DateTime.utc(2022, 6, 5).toISO(),
      payment: '15',
      size: '5',
      oraclePrice: '9800', // BTC Position value at block 5: 5 * 9800 = 49000
    });

    await FundingPaymentsTable.create({
      ...ethFundingPayment,
      createdAtHeight: '5',
      createdAt: DateTime.utc(2022, 6, 5).toISO(),
      payment: '4',
      size: '20',
      oraclePrice: '520', // ETH Position value at block 5: 20 * 520 = 10400
    });

    // Second period (5-10)
    await FundingPaymentsTable.create({
      ...btcFundingPayment,
      createdAtHeight: '10',
      createdAt: DateTime.utc(2022, 6, 10).toISO(),
      payment: '-6', // Negative funding for short position
      size: '-2', // Short 2 BTC
      oraclePrice: '9600', // BTC Position value at block 10: -2 * 9600 = -19200
    });

    await FundingPaymentsTable.create({
      ...ethFundingPayment,
      createdAtHeight: '10',
      createdAt: DateTime.utc(2022, 6, 10).toISO(),
      payment: '7',
      size: '10', // Reduced ETH position
      oraclePrice: '540', // ETH Position value at block 10: 10 * 540 = 5400
    });

    // Third period (10-15)
    await FundingPaymentsTable.create({
      ...btcFundingPayment,
      createdAtHeight: '15',
      createdAt: DateTime.utc(2022, 6, 15).toISO(),
      payment: '3', // Back to positive funding (long position)
      size: '3', // Long 3 BTC
      oraclePrice: '9900', // BTC Position value at block 15: 3 * 9900 = 29700
    });

    await FundingPaymentsTable.create({
      ...ethFundingPayment,
      createdAtHeight: '15',
      createdAt: DateTime.utc(2022, 6, 15).toISO(),
      payment: '8',
      size: '15', // Increased ETH position
      oraclePrice: '560', // ETH Position value at block 15: 15 * 560 = 8400
    });

    // Additional transfers with unique event IDs
    await TransferTable.create({
      ...defaultTransfer,
      eventId: defaultTendermintEventId2,
      size: '30000',
      createdAtHeight: '3',
      createdAt: DateTime.utc(2022, 6, 3).toISO(),
      senderSubaccountId: defaultTransfer.recipientSubaccountId,
      recipientSubaccountId: defaultTransfer.senderSubaccountId,
    });

    await TransferTable.create({
      ...defaultTransfer,
      eventId: defaultTendermintEventId3,
      size: '5000',
      createdAtHeight: '7',
      createdAt: DateTime.utc(2022, 6, 7).toISO(),
      senderSubaccountId: defaultTransfer.senderSubaccountId,
      recipientSubaccountId: defaultTransfer.recipientSubaccountId,
    });

    await TransferTable.create({
      ...defaultTransfer,
      eventId: defaultTendermintEventId4,
      size: '15000',
      createdAtHeight: '12',
      createdAt: DateTime.utc(2022, 6, 12).toISO(),
      senderSubaccountId: defaultTransfer.recipientSubaccountId,
      recipientSubaccountId: defaultTransfer.senderSubaccountId,
    });

    // First period (0-5) trades
    await FillTable.create({
      ...defaultFill,
      subaccountId: defaultSubaccountId,
      side: OrderSide.BUY,
      size: '5', // 5 BTC
      price: '9500', // at $9,500 each
      quoteAmount: '47500', // Total $47,500
      createdAtHeight: '3',
      createdAt: DateTime.utc(2022, 6, 3).toISO(),
      eventId: defaultTendermintEventId,
    });

    await FillTable.create({
      ...defaultFill,
      subaccountId: defaultSubaccountId,
      side: OrderSide.BUY,
      size: '20', // 20 ETH
      price: '500', // at $500 each
      quoteAmount: '10000', // Total $10,000
      createdAtHeight: '3',
      createdAt: DateTime.utc(2022, 6, 3).toISO(),
      eventId: defaultTendermintEventId5,
    });

    // Second period (5-10) trades
    await FillTable.create({
      ...defaultFill,
      subaccountId: defaultSubaccountId,
      side: OrderSide.SELL,
      size: '7', // Sell 7 BTC (5 owned + 2 short)
      price: '9700', // at $9,700 each
      quoteAmount: '67900', // Total $67,900
      createdAtHeight: '7',
      createdAt: DateTime.utc(2022, 6, 7).toISO(),
      eventId: defaultTendermintEventId2,
    });

    await FillTable.create({
      ...defaultFill,
      subaccountId: defaultSubaccountId,
      side: OrderSide.SELL,
      size: '10', // Sell 10 ETH (keeping 10)
      price: '530', // at $530 each
      quoteAmount: '5300', // Total $5,300
      createdAtHeight: '9',
      createdAt: DateTime.utc(2022, 6, 9).toISO(),
      eventId: defaultTendermintEventId6,
    });

    // Third period (10-15) trades
    await FillTable.create({
      ...defaultFill,
      subaccountId: defaultSubaccountId,
      side: OrderSide.BUY,
      size: '5', // Buy 5 BTC (2 to cover short + 3 to go long)
      price: '9800', // at $9,800 each
      quoteAmount: '49000', // Total $49,000
      createdAtHeight: '12',
      createdAt: DateTime.utc(2022, 6, 12).toISO(),
      eventId: defaultTendermintEventId3,
    });

    await FillTable.create({
      ...defaultFill,
      subaccountId: defaultSubaccountId,
      side: OrderSide.BUY,
      size: '5', // Buy 5 more ETH (to 15 total)
      price: '550', // at $550 each
      quoteAmount: '2750', // Total $2,750
      createdAtHeight: '12',
      createdAt: DateTime.utc(2022, 6, 12).toISO(),
      eventId: defaultTendermintEventId4,
    });

    // Run the PNL update task
    await updatePnlTask();

    // Check that PNL entries were created
    const pnlRecords = await PnlTable.findAll({}, []);

    // Check block 5
    const { subaccount1Pnl: subaccount1PnlAtHeight5 } = findPnlRecords(pnlRecords.results, '5');

    verifyPnlRecord(subaccount1PnlAtHeight5, {
      netTransfers: '60000',
      totalPnl: '1919',
      equity: '61919',
    });

    // Check block 10
    const { subaccount1Pnl: subaccount1PnlAtHeight10 } = findPnlRecords(pnlRecords.results, '10');

    verifyPnlRecord(subaccount1PnlAtHeight10, {
      netTransfers: '55000',
      totalPnl: '1920',
      equity: '56920',
    });

    // Check block 15
    const { subaccount1Pnl: subaccount1PnlAtHeight15 } = findPnlRecords(pnlRecords.results, '15');

    verifyPnlRecord(subaccount1PnlAtHeight15, {
      netTransfers: '70000',
      totalPnl: '2081',
      equity: '72081',
    });

    // Verify cache was updated to the latest height
    await verifyCache('15');
  });
});
