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
  defaultTendermintEventId5,
  defaultFill,
  defaultOrder,
  defaultPerpetualMarket2,
  defaultTendermintEventId6,
} from '@dydxprotocol-indexer/postgres/build/__tests__/helpers/constants';

describe('update-pnl', () => {
  // Common date helpers
  const JUNE_1 = DateTime.utc(2022, 6, 1).toISO();
  const JUNE_2 = DateTime.utc(2022, 6, 2).toISO();
  const JUNE_3 = DateTime.utc(2022, 6, 3).toISO();

  beforeAll(async () => {
    await dbHelpers.migrate();
  });

  beforeEach(async () => {
    await dbHelpers.clearData();
    await testMocks.seedData();
  });

  afterAll(async () => {
    await dbHelpers.teardown();
    jest.resetAllMocks();
  });

  afterEach(async () => {
    await dbHelpers.clearData();
  });

  /**
   * Helper function to create blocks at multiple heights
   */
  async function createBlocks(heights: { [height: string]: string }) {
    const promises = Object.entries(heights).map(([height, time]) => 
      BlockTable.create({
        blockHeight: height,
        time,
      })
    );
    await Promise.all(promises);
  }

  /**
   * Helper function to create funding payments
   */
  async function createFundingPayments(
    payments: {
      height: string,
      time: string,
      payment: string,
      subaccountId?: string,
      size?: string,
      oraclePrice?: string
    }[]
  ) {
    const promises = payments.map(({ height, time, payment, subaccountId, size, oraclePrice }) =>
      FundingPaymentsTable.create({
        ...defaultFundingPayment,
        subaccountId: subaccountId || defaultFundingPayment.subaccountId,
        createdAtHeight: height,
        createdAt: time,
        payment,
        size: size || defaultFundingPayment.size,
        oraclePrice: oraclePrice || defaultFundingPayment.oraclePrice,
      })
    );
    await Promise.all(promises);
  }

  /**
   * Helper to verify PNL record
   */
  function verifyPnlRecord(
    record: any,
    expectedValues: {
      equity: string,
      netTransfers: string,
      totalPnl: string,
    }
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
      PersistentCacheKeys.PNL_LAST_PROCESSED_HEIGHT
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
  // Create order for the fill
  await OrderTable.create(defaultOrder);
  
  // Create funding payments for both subaccounts at both heights
  // For the first subaccount (defaultSubaccountId)
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
  
  // Create a transfer at height 1 (between subaccounts)
  await TransferTable.create({
    ...defaultTransfer,
    size: '20000',
    createdAtHeight: '1',
    createdAt: JUNE_1,
    senderSubaccountId: defaultTransfer.recipientSubaccountId,
    recipientSubaccountId: defaultTransfer.senderSubaccountId,
  });
  
  // Create a fill to represent buying the BTC position
  await FillTable.create({
    ...defaultFill,
    subaccountId: defaultSubaccountId,
    side: OrderSide.BUY,
    size: '2',  // 2 BTC
    price: '10000',  // at $10,000 each
    quoteAmount: '20000',  // Total $20,000
    createdAtHeight: '1',
    createdAt: JUNE_1,
  });
  
  // Run the PNL update task
  await updatePnlTask();
  
  // Check that PNL entries were created
  const pnlRecords = await PnlTable.findAll({}, []);
  
  // First check height 1
  const { recordsAtHeight: recordsAtHeight1, subaccount1Pnl: subaccount1PnlAtHeight1, subaccount2Pnl: subaccount2PnlAtHeight1 } = 
    findPnlRecords(pnlRecords.results, '1');
  
  expect(recordsAtHeight1.length).toBe(2); // Both subaccounts should have PNL
  expect(subaccount1PnlAtHeight1).toBeDefined();
  expect(subaccount2PnlAtHeight1).toBeDefined();
    
  // At height 1:
  // - Buy fill: -20000 cash flow
  // - Position value: +20000
  // - Funding payment: +10
  // - Transfer: +20000
  // - totalPnl: -20000 + 20000 + 10 = 10
  // - equity: 20000 + 10 = 20010
  verifyPnlRecord(subaccount1PnlAtHeight1, {
    netTransfers: '20000',
    totalPnl: '10',
    equity: '20010',
  });
  
  // Subaccount 2 at height 1: Transfer -20000 (sending), No position, No funding
  verifyPnlRecord(subaccount2PnlAtHeight1, {
    netTransfers: '-20000',
    totalPnl: '0',
    equity: '-20000',
  });
  
  // Now check height 2
  const { recordsAtHeight: recordsAtHeight2, subaccount1Pnl: subaccount1PnlAtHeight2, subaccount2Pnl: subaccount2PnlAtHeight2 } = 
    findPnlRecords(pnlRecords.results, '2');
  
  expect(recordsAtHeight2.length).toBe(2); // Both subaccounts should have PNL
  expect(subaccount1PnlAtHeight2).toBeDefined();
  expect(subaccount2PnlAtHeight2).toBeDefined();
    
  // At height 2:
  // - Previous totalPnl: 10
  // - Position value change: 22000 - 20000 = 2000
  // - Additional funding: +5
  // - totalPnl: 10 + 2000 + 5 = 2015
  // - netTransfers: unchanged at +20000
  // - equity: 20000 + 2015 = 22015
  verifyPnlRecord(subaccount1PnlAtHeight2, {
    netTransfers: '20000',
    totalPnl: '2015', // 10 + 2000 + 5
    equity: '22015',   // 20000 + 2015
  });
  
  // Subaccount 2 at height 2: Transfer -20000 (unchanged), No position, No funding
  verifyPnlRecord(subaccount2PnlAtHeight2, {
    netTransfers: '-20000',
    totalPnl: '0',
    equity: '-20000', // -20000 + 0
  });
  
  // Verify cache was updated
  await verifyCache('2');
});

it('correctly sums funding payments across multiple positions', async () => {
  // Create order for the fills
  await OrderTable.create(defaultOrder);
  
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
  
  await TransferTable.create({
    ...defaultTransfer,
    size: '0', // Zero-sized transfer just to include the subaccount
    createdAtHeight: '1',
    createdAt: JUNE_1,
  });
  
  // Create fills to represent buying the positions
  // BTC position
  await FillTable.create({
    ...defaultFill,
    subaccountId: defaultSubaccountId,
    side: OrderSide.BUY,
    size: '2',  // 2 BTC
    price: '10000',  // at $10,000 each
    quoteAmount: '20000',  // Total $20,000
    createdAtHeight: '1',
    createdAt: JUNE_1,
    clobPairId: '1', // BTC market
    eventId: defaultTendermintEventId,
  });
  
  // ETH position
  await FillTable.create({
    ...defaultFill,
    subaccountId: defaultSubaccountId,
    side: OrderSide.BUY,
    size: '3',  // 3 ETH
    price: '1000',  // at $1,000 each
    quoteAmount: '3000',  // Total $3,000
    createdAtHeight: '1',
    createdAt: JUNE_1,
    clobPairId: '2', // ETH market
    eventId: defaultTendermintEventId2,
  });
  
  // Run the PNL update task
  await updatePnlTask();
  
  // Check that PNL entries were created
  const pnlRecords = await PnlTable.findAll({}, []);
  
  // Look at height 1
  const { subaccount1Pnl: subaccount1PnlAtHeight1 } = 
    findPnlRecords(pnlRecords.results, '1');
    
  // At height 1:
  // - BTC: Buy -$20,000, Position +$20,000, Funding +$10
  // - ETH: Buy -$3,000, Position +$3,000, Funding +$5
  // - totalPnl: (-$20,000 + $20,000 + $10) + (-$3,000 + $3,000 + $5) = $15
  // - netTransfers: 0 (zero-sized transfer)
  verifyPnlRecord(subaccount1PnlAtHeight1, {
    netTransfers: '0',
    totalPnl: '15', // 10 + 5
    equity: '15',   // 0 + 15
  });
  
  // Verify cache was updated
  await verifyCache('1');
});

it('calculates PNL for a long position that is sold between blocks', async () => {
  // Create the necessary block heights
  await BlockTable.create({
    ...defaultBlock,
    blockHeight: '0',
    time: DateTime.utc(2022, 5, 31).toISO(),
  });
    
  await BlockTable.create({
    ...defaultBlock,
    blockHeight: '5',
    time: DateTime.utc(2022, 6, 5).toISO(),
  });
  
  await BlockTable.create({
    ...defaultBlock,
    blockHeight: '7',
    time: DateTime.utc(2022, 6, 7).toISO(),
  });
  
  await BlockTable.create({
    ...defaultBlock,
    blockHeight: '10',
    time: DateTime.utc(2022, 6, 10).toISO(),
  });

  // Create order for the fills
  await OrderTable.create(defaultOrder);
  
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
  
  // Transfer at height 1
  await TransferTable.create({
    ...defaultTransfer,
    size: '30000',
    createdAtHeight: '1',
    createdAt: DateTime.utc(2022, 6, 1).toISO(),
    senderSubaccountId: defaultTransfer.recipientSubaccountId,
    recipientSubaccountId: defaultTransfer.senderSubaccountId,
  });
  
  // Create a fill to buy at block 1
  await FillTable.create({
    ...defaultFill,
    subaccountId: defaultSubaccountId,
    side: OrderSide.BUY,
    size: '3',  // 3 BTC
    price: '9000',  // at $9,000 each
    quoteAmount: '27000',  // Total $27,000
    createdAtHeight: '1',
    createdAt: DateTime.utc(2022, 6, 1).toISO(),
    eventId: defaultTendermintEventId,
  });
  
  // Create a fill to sell between blocks 5-10 (at block 7)
  await FillTable.create({
    ...defaultFill,
    subaccountId: defaultSubaccountId,
    side: OrderSide.SELL,
    size: '3',  // Sell 3 BTC
    price: '10500',  // at $10,500 each
    quoteAmount: '31500',  // Total $31,500
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
  
  // At block 5:
  // - Transfer at block 1: +30000
  // - Buy fill at block 1: -27000 cash flow
  // - Position value at block 5: 3 * 10000 = 30000
  // - Funding payment: +10 (block 5)
  // - totalPnl: -27000 + 30000 + 10 = 3010
  // - equity: 30000 + 3010 = 33010
  verifyPnlRecord(subaccount1PnlAtHeight5, {
    netTransfers: '30000',
    totalPnl: '3010',
    equity: '33010',
  });
  
  // Check block 10
  const { subaccount1Pnl: subaccount1PnlAtHeight10 } = findPnlRecords(pnlRecords.results, '10');
  
  // At block 10:
  // - Previous totalPnl: 3010
  // - Funding payment: +2 (block 10)
  // - Position value at block 5: 3 * 10000 = 30000
  // - Position value at block 10: 0 (position closed)
  // - Position value change: 0 - 30000 = -30000
  // - Cash flow from sells: +31500 (from selling 3 BTC at 10500)
  // - Cash flow from buys: 0 (no buys in this period)
  // - Net position effect: -30000 + 31500 = 1500
  // - totalPnl: 3010 + 2 + 1500 = 4512
  // - netTransfers: still 30000 (no additional transfers)
  // - equity: 30000 + 4512 = 34512
  verifyPnlRecord(subaccount1PnlAtHeight10, {
    netTransfers: '30000',
    totalPnl: '4512',
    equity: '34512',
  });
  
  // Verify cache was updated to the latest height
  await verifyCache('10');
});


it('calculates PNL for a short position that is closed between blocks', async () => {
  // Create the necessary block heights
  await BlockTable.create({
    ...defaultBlock,
    blockHeight: '0',
    time: DateTime.utc(2022, 5, 31).toISO(),
  });
    
  await BlockTable.create({
    ...defaultBlock,
    blockHeight: '5',
    time: DateTime.utc(2022, 6, 5).toISO(),
  });
  
  await BlockTable.create({
    ...defaultBlock,
    blockHeight: '8',
    time: DateTime.utc(2022, 6, 8).toISO(),
  });
  
  await BlockTable.create({
    ...defaultBlock,
    blockHeight: '10',
    time: DateTime.utc(2022, 6, 10).toISO(),
  });

  // Create order for the fills
  await OrderTable.create(defaultOrder);
  
  // Funding payments at blocks 5 and 10
  await FundingPaymentsTable.create({
    ...defaultFundingPayment,
    createdAtHeight: '5',
    createdAt: DateTime.utc(2022, 6, 5).toISO(),
    payment: '-8', // Negative funding for short position
    size: '-4',    // Negative size for short position
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
  
  // Transfer at height 1 to ensure subaccounts are included in calculations
  await TransferTable.create({
    ...defaultTransfer,
    size: '40000',
    createdAtHeight: '1',
    createdAt: DateTime.utc(2022, 6, 1).toISO(),
    senderSubaccountId: defaultTransfer.recipientSubaccountId,
    recipientSubaccountId: defaultTransfer.senderSubaccountId,
  });
  
  // Create a fill to sell at block 3
  await FillTable.create({
    ...defaultFill,
    subaccountId: defaultSubaccountId,
    side: OrderSide.SELL,
    size: '4',  // 4 BTC
    price: '9500',  // at $9,500 each
    quoteAmount: '38000',  // Total $38,000
    createdAtHeight: '3',
    createdAt: DateTime.utc(2022, 6, 3).toISO(),
    eventId: defaultTendermintEventId,
  });
  
  // Create a fill to buy back between blocks 5-10 (at block 8)
  await FillTable.create({
    ...defaultFill,
    subaccountId: defaultSubaccountId,
    side: OrderSide.BUY,
    size: '4',  // Buy 4 BTC
    price: '9200',  // at $9,200 each
    quoteAmount: '36800',  // Total $36,800
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
  
  // At block 5:
  // - Transfer at block 1: +40000
  // - Sell fill at block 3: +38000 cash flow
  // - Short position value at block 5: -4 * 10000 = -40000 (negative for short position)
  // - Funding payment: -8 (negative for short position)
  // - Position value change: -40000 (no previous position value)
  // - Net cash flow: +38000
  // - Net position effect: -40000 + 38000 = -2000 (loss)
  // - totalPnl: 0 - 2000 - 8 = -2008
  // - equity: 40000 - 2008 = 37992
  verifyPnlRecord(subaccount1PnlAtHeight5, {
    netTransfers: '40000',
    totalPnl: '-2008',
    equity: '37992',
  });
  
  // Check block 10
  const { subaccount1Pnl: subaccount1PnlAtHeight10 } = findPnlRecords(pnlRecords.results, '10');
  
  // At block 10:
  // - Previous totalPnl: -2008
  // - Funding payment: +2 (block 10)
  // - Position value at block 5: -40000 (short position)
  // - Position value at block 10: 0 (position closed)
  // - Position value change: 0 - (-40000) = +40000
  // - Cash flow from buys: -36800 (from buying 4 BTC at 9200)
  // - Cash flow from sells: 0 (no sells in this period)
  // - Net position effect: 40000 - 36800 = 3200 (profit)
  // - totalPnl: -2008 + 2 + 3200 = 1194
  // - netTransfers: still 40000 (no additional transfers)
  // - equity: 40000 + 1194 = 41194
  verifyPnlRecord(subaccount1PnlAtHeight10, {
    netTransfers: '40000',
    totalPnl: '1194',
    equity: '41194',
  });
  
  // Verify cache was updated to the latest height
  await verifyCache('10');
});

it('calculates PNL for a position that changes from long to short', async () => {
  // Create the necessary block heights
  await BlockTable.create({
    ...defaultBlock,
    blockHeight: '0',
    time: DateTime.utc(2022, 5, 31).toISO(),
  });
  
  await BlockTable.create({
    ...defaultBlock,
    blockHeight: '5',
    time: DateTime.utc(2022, 6, 5).toISO(),
  });
  
  await BlockTable.create({
    ...defaultBlock,
    blockHeight: '7',
    time: DateTime.utc(2022, 6, 7).toISO(),
  });
  
  await BlockTable.create({
    ...defaultBlock,
    blockHeight: '10',
    time: DateTime.utc(2022, 6, 10).toISO(),
  });

  // Create order for the fills
  await OrderTable.create(defaultOrder);
  
  // Funding payments at blocks 5 and 10
  await FundingPaymentsTable.create({
    ...defaultFundingPayment,
    createdAtHeight: '5',
    createdAt: DateTime.utc(2022, 6, 5).toISO(),
    payment: '15', // Positive funding for long position
    size: '5',     // Positive size for long position
    oraclePrice: '9800', // Position value at block 5: 5 * 9800 = 49000
  });
  
  await FundingPaymentsTable.create({
    ...defaultFundingPayment,
    createdAtHeight: '10',
    createdAt: DateTime.utc(2022, 6, 10).toISO(),
    payment: '-12', // Negative funding for short position
    size: '-2',     // Negative size for short position
    oraclePrice: '9600', // Position value at block 10: -2 * 9600 = -19200
  });
  
  // Transfer at height 1 to ensure subaccounts are included in calculations
  await TransferTable.create({
    ...defaultTransfer,
    size: '50000',
    createdAtHeight: '1',
    createdAt: DateTime.utc(2022, 6, 1).toISO(),
    senderSubaccountId: defaultTransfer.recipientSubaccountId,
    recipientSubaccountId: defaultTransfer.senderSubaccountId,
  });
  
  // Create a fill to buy at block 2
  await FillTable.create({
    ...defaultFill,
    subaccountId: defaultSubaccountId,
    side: OrderSide.BUY,
    size: '5',  // 5 BTC
    price: '9500',  // at $9,500 each
    quoteAmount: '47500',  // Total $47,500
    createdAtHeight: '2',
    createdAt: DateTime.utc(2022, 6, 2).toISO(),
    eventId: defaultTendermintEventId,
  });
  
  // Create a fill to sell more between blocks 5-10 (at block 7) - shifting to a short position
  await FillTable.create({
    ...defaultFill,
    subaccountId: defaultSubaccountId,
    side: OrderSide.SELL,
    size: '7',  // Sell 7 BTC (5 to close + 2 to short)
    price: '9700',  // at $9,700 each
    quoteAmount: '67900',  // Total $67,900
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
  
  // At block 5:
  // - Transfer at block 1: +50000
  // - Buy fill at block 2: -47500 cash flow
  // - Long position value at block 5: 5 * 9800 = 49000
  // - Funding payment: +15
  // - Position value change: 49000 (no previous position value)
  // - Net cash flow: -47500
  // - Net position effect: 49000 - 47500 = 1500 (profit)
  // - totalPnl: 0 + 1500 + 15 = 1515
  // - equity: 50000 + 1515 = 51515
  verifyPnlRecord(subaccount1PnlAtHeight5, {
    netTransfers: '50000',
    totalPnl: '1515',
    equity: '51515',
  });
  
  // Check block 10
  const { subaccount1Pnl: subaccount1PnlAtHeight10 } = findPnlRecords(pnlRecords.results, '10');
  
  // At block 10:
  // - Previous totalPnl: 1515
  // - Funding payment: -12 (block 10)
  // - Position value at block 5: 49000 (long position)
  // - Position value at block 10: -19200 (short position)
  // - Position value change: -19200 - 49000 = -68200
  // - Cash flow from sells: +67900 (from selling 7 BTC at 9700)
  // - Net position effect: -68200 + 67900 = -300 (small loss)
  // - totalPnl: 1515 - 12 - 300 = 1203
  // - netTransfers: still 50000 (no additional transfers)
  // - equity: 50000 + 1203 = 51203
  verifyPnlRecord(subaccount1PnlAtHeight10, {
    netTransfers: '50000',
    totalPnl: '1203',
    equity: '51203',
  });
  
  // Verify cache was updated to the latest height
  await verifyCache('10');
});

it('calculates PNL for operations on multiple different tickers', async () => {
  // Create the necessary block heights
  await BlockTable.create({
    ...defaultBlock,
    blockHeight: '0',
    time: DateTime.utc(2022, 5, 31).toISO(),
  });
  
  await BlockTable.create({
    ...defaultBlock,
    blockHeight: '4',
    time: DateTime.utc(2022, 6, 4).toISO(),
  });
  
  await BlockTable.create({
    ...defaultBlock,
    blockHeight: '5',
    time: DateTime.utc(2022, 6, 5).toISO(),
  });
  
  await BlockTable.create({
    ...defaultBlock,
    blockHeight: '7',
    time: DateTime.utc(2022, 6, 7).toISO(),
  });
  
  await BlockTable.create({
    ...defaultBlock,
    blockHeight: '8',
    time: DateTime.utc(2022, 6, 8).toISO(),
  });
  
  await BlockTable.create({
    ...defaultBlock,
    blockHeight: '10',
    time: DateTime.utc(2022, 6, 10).toISO(),
  });

  // Create order for the fills
  await OrderTable.create(defaultOrder);
  
  // Create funding payments for different tickers
  const btcFundingPayment = defaultFundingPayment;
  const ethFundingPayment = { ...defaultFundingPayment, ticker: 'ETH-USD' };
  
  // Funding payments for BTC  
  await FundingPaymentsTable.create({
    ...btcFundingPayment,
    createdAtHeight: '10',
    createdAt: DateTime.utc(2022, 6, 10).toISO(),
    payment: '-7', // Negative funding for short BTC position
    size: '-1',    // Negative size for short position
    oraclePrice: '9600', // BTC Position value at block 10: -1 * 9600 = -9600
  });
  
  // Funding payments for ETH
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
  
  // Transfer at height 1 to ensure subaccounts are included in calculations
  await TransferTable.create({
    ...defaultTransfer,
    size: '30000',
    createdAtHeight: '1',
    createdAt: DateTime.utc(2022, 6, 1).toISO(),
    senderSubaccountId: defaultTransfer.recipientSubaccountId,
    recipientSubaccountId: defaultTransfer.senderSubaccountId,
  });
  
  // BTC Operations:
  // Block 2: Buy 2 BTC
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
  
  // Block 4: Sell 2 BTC (close position)
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
  
  // Block 7: Sell 1 BTC (open short)
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
  
  // ETH Operations:
  // Block 3: Buy 10 ETH
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
  
  // Block 8: Buy 5 more ETH
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
  
  // At block 5:
  // - Transfer at block 1: +30000
  // 
  // BTC operations:
  // - Buy 2 BTC at block 2: -18000 cash flow
  // - Sell 2 BTC at block 4: +18600 cash flow
  // - BTC position value at block 5: 0 (closed)
  // - BTC position effect: 0 (from position value) + 18600 - 18000 = 600 (profit)
  // - BTC funding: 0 (no position at block 5)
  //
  // ETH operations:
  // - Buy 10 ETH at block 3: -5100 cash flow
  // - ETH position value at block 5: 10 * 520 = 5200
  // - ETH position effect: 5200 - 5100 = 100 (profit)
  // - ETH funding: +3
  //
  // Total position effect: 600 + 100 = 700 (profit)
  // Total funding: 0 + 3 = 3
  // - totalPnl: 0 + 700 + 3 = 703
  // - equity: 30000 + 703 = 30703
  verifyPnlRecord(subaccount1PnlAtHeight5, {
    netTransfers: '30000',
    totalPnl: '703',
    equity: '30703',
  });
  
  // Check block 10
  const { subaccount1Pnl: subaccount1PnlAtHeight10 } = findPnlRecords(pnlRecords.results, '10');
  
  // At block 10:
  // - Previous totalPnl: 703
  // 
  // BTC operations:
  // - Sell 1 BTC at block 7: +9550 cash flow
  // - BTC position value at block 10: -1 * 9600 = -9600 (short)
  // - BTC position value change: -9600 - 0 = -9600
  // - BTC position effect: -9600 + 9550 = -50 (small loss)
  // - BTC funding: -7
  //
  // ETH operations:
  // - Buy 5 more ETH at block 8: -2650 cash flow
  // - ETH position value at block 5: 10 * 520 = 5200
  // - ETH position value at block 10: 15 * 550 = 8250
  // - ETH position value change: 8250 - 5200 = 3050
  // - ETH position effect: 3050 - 2650 = 400 (profit)
  // - ETH funding: +8
  //
  // Total position effect: -50 + 400 = 350 (profit)
  // Total funding: -7 + 8 = 1
  // - totalPnl: 703 + 350 + 1 = 1054
  // - equity: 30000 + 1054 = 31054
  verifyPnlRecord(subaccount1PnlAtHeight10, {
    netTransfers: '30000',
    totalPnl: '1054',
    equity: '31054',
  });
  
  // Verify cache was updated to the latest height
  await verifyCache('10');
});

it('calculates PNL across multiple funding periods with various trades and transfers', async () => {
  // Create the necessary block heights
  await BlockTable.create({
    ...defaultBlock,
    blockHeight: '0',
    time: DateTime.utc(2022, 5, 31).toISO(),
  });
  
  await BlockTable.create({
    ...defaultBlock,
    blockHeight: '5',
    time: DateTime.utc(2022, 6, 5).toISO(),
  });
  
  await BlockTable.create({
    ...defaultBlock,
    blockHeight: '7',
    time: DateTime.utc(2022, 6, 7).toISO(),
  });
  
  await BlockTable.create({
    ...defaultBlock,
    blockHeight: '9',
    time: DateTime.utc(2022, 6, 9).toISO(),
  });
  
  await BlockTable.create({
    ...defaultBlock,
    blockHeight: '10',
    time: DateTime.utc(2022, 6, 10).toISO(),
  });
  
  await BlockTable.create({
    ...defaultBlock,
    blockHeight: '12',
    time: DateTime.utc(2022, 6, 12).toISO(),
  });
  
  await BlockTable.create({
    ...defaultBlock,
    blockHeight: '15',
    time: DateTime.utc(2022, 6, 15).toISO(),
  });

  // Create order for the fills
  await OrderTable.create(defaultOrder);
  
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
    size: '-2',    // Short 2 BTC
    oraclePrice: '9600', // BTC Position value at block 10: -2 * 9600 = -19200
  });
  
  await FundingPaymentsTable.create({
    ...ethFundingPayment,
    createdAtHeight: '10',
    createdAt: DateTime.utc(2022, 6, 10).toISO(),
    payment: '7',
    size: '10',   // Reduced ETH position
    oraclePrice: '540', // ETH Position value at block 10: 10 * 540 = 5400
  });
  
  // Third period (10-15)
  await FundingPaymentsTable.create({
    ...btcFundingPayment,
    createdAtHeight: '15',
    createdAt: DateTime.utc(2022, 6, 15).toISO(),
    payment: '3', // Back to positive funding (long position)
    size: '3',    // Long 3 BTC
    oraclePrice: '9900', // BTC Position value at block 15: 3 * 9900 = 29700
  });
  
  await FundingPaymentsTable.create({
    ...ethFundingPayment,
    createdAtHeight: '15',
    createdAt: DateTime.utc(2022, 6, 15).toISO(),
    payment: '8',
    size: '15',   // Increased ETH position
    oraclePrice: '560', // ETH Position value at block 15: 15 * 560 = 8400
  });
  
  // Multiple transfers across different periods
  // Initial transfer
  await TransferTable.create({
    ...defaultTransfer,
    size: '50000',
    createdAtHeight: '1',
    createdAt: DateTime.utc(2022, 6, 1).toISO(),
    senderSubaccountId: defaultTransfer.recipientSubaccountId,
    recipientSubaccountId: defaultTransfer.senderSubaccountId,
  });
  
  // Second transfer in first period
  await TransferTable.create({
    ...defaultTransfer,
    size: '10000',
    createdAtHeight: '3',
    createdAt: DateTime.utc(2022, 6, 3).toISO(),
    eventId: defaultTendermintEventId2,
    senderSubaccountId: defaultTransfer.recipientSubaccountId,
    recipientSubaccountId: defaultTransfer.senderSubaccountId,
  });
  
  // Transfer in second period (withdrawal)
  await TransferTable.create({
    ...defaultTransfer,
    size: '5000',
    createdAtHeight: '7',
    createdAt: DateTime.utc(2022, 6, 7).toISO(),
    eventId: defaultTendermintEventId3,
    senderSubaccountId: defaultTransfer.senderSubaccountId,
    recipientSubaccountId: defaultTransfer.recipientSubaccountId,
  });
  
  // Transfer in third period
  await TransferTable.create({
    ...defaultTransfer,
    size: '15000',
    createdAtHeight: '12',
    createdAt: DateTime.utc(2022, 6, 12).toISO(),
    eventId: defaultTendermintEventId4,
    senderSubaccountId: defaultTransfer.recipientSubaccountId,
    recipientSubaccountId: defaultTransfer.senderSubaccountId,
  });
  
  // First period (0-5)
  // Buy BTC at block 3
  await FillTable.create({
    ...defaultFill,
    subaccountId: defaultSubaccountId,
    side: OrderSide.BUY,
    size: '5',  // 5 BTC
    price: '9500',  // at $9,500 each
    quoteAmount: '47500',  // Total $47,500
    createdAtHeight: '3',
    createdAt: DateTime.utc(2022, 6, 3).toISO(),
    eventId: defaultTendermintEventId,
  });
  
  // Buy ETH at block 3
  await FillTable.create({
    ...defaultFill,
    subaccountId: defaultSubaccountId,
    side: OrderSide.BUY,
    size: '20',  // 20 ETH
    price: '500',  // at $500 each
    quoteAmount: '10000',  // Total $10,000
    createdAtHeight: '3',
    createdAt: DateTime.utc(2022, 6, 3).toISO(),
    eventId: defaultTendermintEventId4,
  });
  
  // Second period (5-10)
  // Sell more BTC than owned at block 7 (flipping to short)
  await FillTable.create({
    ...defaultFill,
    subaccountId: defaultSubaccountId,
    side: OrderSide.SELL,
    size: '7',  // Sell 7 BTC (5 owned + 2 short)
    price: '9700',  // at $9,700 each
    quoteAmount: '67900',  // Total $67,900
    createdAtHeight: '7',
    createdAt: DateTime.utc(2022, 6, 7).toISO(),
    eventId: defaultTendermintEventId2,
  });
  
  // Sell some ETH at block 9 (partial position reduction)
  await FillTable.create({
    ...defaultFill,
    subaccountId: defaultSubaccountId,
    side: OrderSide.SELL,
    size: '10',  // Sell 10 ETH (keeping 10)
    price: '530',  // at $530 each
    quoteAmount: '5300',  // Total $5,300
    createdAtHeight: '9',
    createdAt: DateTime.utc(2022, 6, 9).toISO(),
    eventId: defaultTendermintEventId5,
  });
  
  // Third period (10-15)
  // Buy back all short BTC and go long at block 12
  await FillTable.create({
    ...defaultFill,
    subaccountId: defaultSubaccountId,
    side: OrderSide.BUY,
    size: '5',  // Buy 5 BTC (2 to cover short + 3 to go long)
    price: '9800',  // at $9,800 each
    quoteAmount: '49000',  // Total $49,000
    createdAtHeight: '12',
    createdAt: DateTime.utc(2022, 6, 12).toISO(),
    eventId: defaultTendermintEventId3,
  });
  
  // Buy more ETH at block 12
  await FillTable.create({
    ...defaultFill,
    subaccountId: defaultSubaccountId,
    side: OrderSide.BUY,
    size: '5',  // Buy 5 more ETH (to 15 total)
    price: '550',  // at $550 each
    quoteAmount: '2750',  // Total $2,750
    createdAtHeight: '12',
    createdAt: DateTime.utc(2022, 6, 12).toISO(),
    eventId: defaultTendermintEventId6,
  });
  
  // Run the PNL update task
  await updatePnlTask();
  
  // Check that PNL entries were created
  const pnlRecords = await PnlTable.findAll({}, []);
  
  // Check block 5
  const { subaccount1Pnl: subaccount1PnlAtHeight5 } = findPnlRecords(pnlRecords.results, '5');
  
  // At block 5:
  // - Transfers: +50000 (block 1) + 10000 (block 3) = 60000
  // 
  // BTC:
  // - Buy 5 BTC at block 3: -47500 cash flow
  // - BTC position value at block 5: 5 * 9800 = 49000
  // - BTC position effect: 49000 - 47500 = 1500 (profit)
  // - BTC funding: +15
  //
  // ETH:
  // - Buy 20 ETH at block 3: -10000 cash flow
  // - ETH position value at block 5: 20 * 520 = 10400
  // - ETH position effect: 10400 - 10000 = 400 (profit)
  // - ETH funding: +4
  //
  // Total position effect: 1500 + 400 = 1900 (profit)
  // Total funding: 15 + 4 = 19
  // - totalPnl: 0 + 1900 + 19 = 1919
  // - equity: 60000 + 1919 = 61919
  verifyPnlRecord(subaccount1PnlAtHeight5, {
    netTransfers: '60000',
    totalPnl: '1919',
    equity: '61919',
  });
  
  // Check block 10
  const { subaccount1Pnl: subaccount1PnlAtHeight10 } = findPnlRecords(pnlRecords.results, '10');
  
  // At block 10:
  // - Previous totalPnl: 1919
  // - Transfers: +60000 (previous) - 5000 (block 7) = 55000
  // 
  // BTC:
  // - Sell 7 BTC at block 7: +67900 cash flow
  // - BTC position value at block 5: 5 * 9800 = 49000 (long)
  // - BTC position value at block 10: -2 * 9600 = -19200 (short)
  // - BTC position value change: -19200 - 49000 = -68200
  // - BTC position effect: -68200 + 67900 = -300 (small loss)
  // - BTC funding: -6
  //
  // ETH:
  // - Sell 10 ETH at block 9: +5300 cash flow
  // - ETH position value at block 5: 20 * 520 = 10400
  // - ETH position value at block 10: 10 * 540 = 5400
  // - ETH position value change: 5400 - 10400 = -5000
  // - ETH position effect: -5000 + 5300 = 300 (profit)
  // - ETH funding: +7
  //
  // Position effect: -300 + 300 = 0 (neutral)
  // Funding: -6 + 7 = 1
  // - totalPnl: 1919 + 0 + 1 = 1920
  // - equity: 55000 + 1920 = 56920
  verifyPnlRecord(subaccount1PnlAtHeight10, {
    netTransfers: '55000',
    totalPnl: '1920',
    equity: '56920',
  });
  
  // Check block 15
  const { subaccount1Pnl: subaccount1PnlAtHeight15 } = findPnlRecords(pnlRecords.results, '15');
  
  // At block 15:
  // - Previous totalPnl: 1920
  // - Transfers: +55000 (previous) + 15000 (block 12) = 70000
  // 
  // BTC:
  // - Buy 5 BTC at block 12: -49000 cash flow
  // - BTC position value at block 10: -2 * 9600 = -19200 (short)
  // - BTC position value at block 15: 3 * 9900 = 29700 (long)
  // - BTC position value change: 29700 - (-19200) = 48900
  // - BTC position effect: 48900 - 49000 = -100 (small loss)
  // - BTC funding: +3
  //
  // ETH:
  // - Buy 5 ETH at block 12: -2750 cash flow
  // - ETH position value at block 10: 10 * 540 = 5400
  // - ETH position value at block 15: 15 * 560 = 8400
  // - ETH position value change: 8400 - 5400 = 3000
  // - ETH position effect: 3000 - 2750 = 250 (profit)
  // - ETH funding: +8
  //
  // Position effect: -100 + 250 = 150 (profit)
  // Funding: 3 + 8 = 11
  // - totalPnl: 1920 + 150 + 11 = 2081
  // - equity: 70000 + 2081 = 72081
  verifyPnlRecord(subaccount1PnlAtHeight15, {
    netTransfers: '70000',
    totalPnl: '2081',
    equity: '72081',
  });
  
  // Verify cache was updated to the latest height
  await verifyCache('15');
});
});