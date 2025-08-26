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
  Liquidity,
  FillType,
  OrderTable,
} from '@dydxprotocol-indexer/postgres';
import { DateTime } from 'luxon';

import updatePnlTask from '../../src/tasks/update-pnl';

import {
  defaultSubaccountId,
  defaultSubaccountId2,
  defaultBlock,
  defaultBlock2,
  defaultTransfer,
  defaultFundingPayment,
  defaultPerpetualMarket,
  defaultTendermintEventId,
  defaultTendermintEventId2,
  defaultTendermintEventId3,
  defaultOrderId,
  defaultFill,
  defaultOrder,
  defaultDeposit,
  defaultFundingPayment2,
  defaultPerpetualMarket2,
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
  
//   // For the second subaccount (defaultSubaccountId2)
//   await FundingPaymentsTable.create({
//     ...defaultFundingPayment2, // This uses defaultSubaccountId2
//     createdAtHeight: '1',
//     createdAt: JUNE_1,
//     payment: '0',
//     size: '0',      // Explicitly set to 0 (no position)
//     oraclePrice: '10000',
//   });
  
//   await FundingPaymentsTable.create({
//     ...defaultFundingPayment2,
//     createdAtHeight: '2',
//     createdAt: JUNE_2,
//     payment: '0',
//     size: '0',      // Explicitly set to 0 (no position)
//     oraclePrice: '11000',
//   });
  
  // Create a transfer at height 1 (between subaccounts)
  await TransferTable.create({
    ...defaultTransfer,
    size: '1000',
    createdAtHeight: '1',
    createdAt: JUNE_1,
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
  // - Transfer: -1000
  // - totalPnl: -20000 + 20000 + 10 = 10
  // - equity: -1000 + 10 = -990
  verifyPnlRecord(subaccount1PnlAtHeight1, {
    netTransfers: '-1000',
    totalPnl: '10', // -20000 + 20000 + 10
    equity: '-990', // -1000 + 10
  });
  
  // Subaccount 2 at height 1: Transfer +1000, No position, No funding
  verifyPnlRecord(subaccount2PnlAtHeight1, {
    netTransfers: '1000',
    totalPnl: '0',
    equity: '1000', // 1000 + 0
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
  // - netTransfers: unchanged at -1000
  // - equity: -1000 + 2015 = 1015
  verifyPnlRecord(subaccount1PnlAtHeight2, {
    netTransfers: '-1000',
    totalPnl: '2015', // 10 + 2000 + 5
    equity: '1015',   // -1000 + 2015
  });
  
  // Subaccount 2 at height 2: Transfer +1000 (unchanged), No position, No funding
  verifyPnlRecord(subaccount2PnlAtHeight2, {
    netTransfers: '1000',
    totalPnl: '0',
    equity: '1000', // 1000 + 0
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
  
  // Create a dummy transfer to ensure the subaccount is included in the calculation
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

});