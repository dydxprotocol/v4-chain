import {
  dbHelpers,
  testMocks,
  TransferTable,
  OraclePriceTable,
  BlockTable,
  PnlTable,
  PersistentCacheTable,
  PersistentCacheKeys,
  FundingPaymentsTable,
} from '@dydxprotocol-indexer/postgres';

import updatePnlTask from '../../src/tasks/update-pnl';

import {
  defaultSubaccountId,
  defaultSubaccountId2,
  defaultBlock,
  defaultBlock2,
  defaultTransfer,
  defaultOraclePrice,
  defaultFundingPayment,
} from '@dydxprotocol-indexer/postgres/build/__tests__/helpers/constants';
import { DateTime } from 'luxon';

describe('update-pnl', () => {
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

  it('computes zero PnL for subaccounts with transfers but no positions', async () => {
    // Create block 0 to establish the timeline
    await BlockTable.create({
        blockHeight: '0',
        time: DateTime.utc(2022, 5, 31).toISO(),
    });

    // Create oracle price at heights 0 and 1
    await OraclePriceTable.create({
      ...defaultOraclePrice,
      effectiveAtHeight: '0',
    });
    await OraclePriceTable.create(defaultOraclePrice);
    
    // Create a transfer between these subaccounts
    await TransferTable.create({
      ...defaultTransfer,
      createdAtHeight: '1',
    });

    // Create a funding payment at height 1
    await FundingPaymentsTable.create(defaultFundingPayment);

    await updatePnlTask();

    // Check that PNL entries were created for both subaccounts
    const pnlRecords = await PnlTable.findAll({}, []);

    // Expect 2 records (one for each subaccount)
    expect(pnlRecords.results.length).toBe(2);
    
    // The subaccount with the funding payment should have non-zero deltaFundingPayments
    const subaccount1Pnl = pnlRecords.results.find((r: { subaccountId: string; }) => r.subaccountId === defaultSubaccountId);
    const subaccount2Pnl = pnlRecords.results.find((r: { subaccountId: string; }) => r.subaccountId === defaultSubaccountId2);
    
    expect(subaccount1Pnl).toBeDefined();
    expect(subaccount2Pnl).toBeDefined();
    
    // Subaccount1 has a funding payment
    expect(subaccount1Pnl?.deltaFundingPayments).toBe('5');
    expect(subaccount1Pnl?.deltaPositionEffects).toBe('0');
    expect(subaccount1Pnl?.totalPnl).toBe('5');
    
    // Subaccount2 has no funding payment or positions
    expect(subaccount2Pnl?.deltaFundingPayments).toBe('0');
    expect(subaccount2Pnl?.deltaPositionEffects).toBe('0');
    expect(subaccount2Pnl?.totalPnl).toBe('0');

    // Verify createdAt and createdAtHeight values
    for (const record of pnlRecords.results) {
      // createdAtHeight should be the end height (1)
      expect(record.createdAtHeight).toBe('1');
      
      // createdAt should match the timestamp of the oracle price at height 1
      expect(record.createdAt).toBe(defaultOraclePrice.effectiveAt);
    }
    
    // Verify persistent cache was updated
    const cache = await PersistentCacheTable.findById(
      PersistentCacheKeys.PNL_LAST_PROCESSED_HEIGHT
    );
    expect(cache).toBeDefined();
    expect(cache?.value).toBe('1');
  });

  it('processes multiple periods of pnl calculations', async () => {
  // Create block 0 to establish the timeline
  await BlockTable.create({
    blockHeight: '0',
    time: DateTime.utc(2022, 5, 31).toISO(),
  });
  
  // Create oracle prices for all heights
  await OraclePriceTable.create({
    ...defaultOraclePrice,
    effectiveAtHeight: '0',
    effectiveAt: DateTime.utc(2022, 5, 31).toISO(),
  });
  await OraclePriceTable.create({
    ...defaultOraclePrice,
    effectiveAtHeight: '1',
    effectiveAt: DateTime.utc(2022, 6, 1).toISO(),
  });
  await OraclePriceTable.create({
    ...defaultOraclePrice,
    effectiveAtHeight: '2',
    effectiveAt: DateTime.utc(2022, 6, 2).toISO(),
  });
  
  // Create transfers
  await TransferTable.create({
    ...defaultTransfer,
    createdAtHeight: '1',
  });
  
  // Create funding payments at heights 1 and 2
  await FundingPaymentsTable.create({
    ...defaultFundingPayment,
    createdAtHeight: '1',
    createdAt: defaultBlock.time,
    payment: '5',
  });
  await FundingPaymentsTable.create({
    ...defaultFundingPayment,
    createdAtHeight: '2',
    createdAt: defaultBlock2.time,
    payment: '10',
  });
  
  // Run the task
  await updatePnlTask();
  
  // Check that PNL entries were created
  const pnlRecords = await PnlTable.findAll({}, []);
  
  // Sort the records by height to make testing easier
  const sortedRecords = pnlRecords.results.sort((a, b) => 
    parseInt(a.createdAtHeight, 10) - parseInt(b.createdAtHeight, 10));
  
  // We should have records for both heights, 2 heights * 2 subaccounts = 4 records total
  expect(sortedRecords.length).toBe(4);
  
  // Group records by height
  const recordsAtHeight1 = sortedRecords.filter((r) => r.createdAtHeight === '1');
  const recordsAtHeight2 = sortedRecords.filter((r) => r.createdAtHeight === '2');
  
  expect(recordsAtHeight1.length).toBe(2);
  expect(recordsAtHeight2.length).toBe(2);
  
  // Find records for each subaccount at height 1
  const subaccount1PnlAtHeight1 = recordsAtHeight1.find((r) => r.subaccountId === defaultSubaccountId);
  const subaccount2PnlAtHeight1 = recordsAtHeight1.find((r) => r.subaccountId === defaultSubaccountId2);
  
  expect(subaccount1PnlAtHeight1).toBeDefined();
  expect(subaccount2PnlAtHeight1).toBeDefined();
  
  // PNL at height 1
  expect(subaccount1PnlAtHeight1?.deltaFundingPayments).toBe('5');
  expect(subaccount1PnlAtHeight1?.deltaPositionEffects).toBe('0');
  expect(subaccount1PnlAtHeight1?.totalPnl).toBe('5');
  
  expect(subaccount2PnlAtHeight1?.deltaFundingPayments).toBe('0');
  expect(subaccount2PnlAtHeight1?.deltaPositionEffects).toBe('0');
  expect(subaccount2PnlAtHeight1?.totalPnl).toBe('0');
  
  // Find records for each subaccount at height 2
  const subaccount1PnlAtHeight2 = recordsAtHeight2.find((r) => r.subaccountId === defaultSubaccountId);
  const subaccount2PnlAtHeight2 = recordsAtHeight2.find((r) => r.subaccountId === defaultSubaccountId2);
  
  expect(subaccount1PnlAtHeight2).toBeDefined();
  expect(subaccount2PnlAtHeight2).toBeDefined();
  
  // PNL at height 2 - subaccount1 should have additional funding payment and accumulated total
  expect(subaccount1PnlAtHeight2?.deltaFundingPayments).toBe('10');
  expect(subaccount1PnlAtHeight2?.deltaPositionEffects).toBe('0');
  expect(subaccount1PnlAtHeight2?.totalPnl).toBe('15'); // 5 from height 1 + 10 from height 2
  
  // PNL at height 2 - subaccount2 still has no funding or positions
  expect(subaccount2PnlAtHeight2?.deltaFundingPayments).toBe('0');
  expect(subaccount2PnlAtHeight2?.deltaPositionEffects).toBe('0');
  expect(subaccount2PnlAtHeight2?.totalPnl).toBe('0');
  
  // Verify createdAt matches oracle price timestamps for the respective heights
  for (const record of recordsAtHeight1) {
    expect(record.createdAt).toBe(DateTime.utc(2022, 6, 1).toISO());
  }
  
  for (const record of recordsAtHeight2) {
    expect(record.createdAt).toBe(DateTime.utc(2022, 6, 2).toISO());
  }
  
  // Verify persistent cache was updated to the highest height
  const cache = await PersistentCacheTable.findById(
    PersistentCacheKeys.PNL_LAST_PROCESSED_HEIGHT
  );
  expect(cache).toBeDefined();
  expect(cache?.value).toBe('2');
});
});