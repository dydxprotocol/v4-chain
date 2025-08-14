import {
  dbHelpers,
  testMocks,
  SubaccountTable,
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
  defaultBlock,
  defaultBlock2,
  defaultSubaccount,
  defaultSubaccount2,
  defaultSubaccountId,
  defaultSubaccountId2,
  defaultTransfer,
  defaultMarket,
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
});