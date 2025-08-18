import {
  dbHelpers,
  testMocks,
  TransferTable,
  OraclePriceTable,
  BlockTable,
  PnlTable,
  PerpetualPositionTable,
  PersistentCacheTable,
  PersistentCacheKeys,
  FundingPaymentsTable,
  PerpetualPositionStatus,
  PositionSide,
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
  defaultPerpetualMarket,
  defaultTendermintEventId,
  defaultTendermintEventId2,
  defaultBlock10,
  defaultPerpetualMarket2,
  defaultMarket2,
  defaultTendermintEventId4,
  defaultTendermintEventId3,
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

  it('calculates position effects correctly for open positions in the same subaccount', async () => {
    // Create blocks for heights 0, 1, 7, and 10
    await BlockTable.create({
      blockHeight: '0',
      time: DateTime.utc(2022, 5, 31).toISO(),
    });
    // Height 1 is created by seedData() as defaultBlock
    await BlockTable.create({
      blockHeight: '7',
      time: DateTime.utc(2022, 6, 7).toISO(),
    });
    await BlockTable.create(defaultBlock10);

    // Create oracle prices with increasing prices
    await OraclePriceTable.create({
      ...defaultOraclePrice,
      effectiveAtHeight: '0',
      effectiveAt: DateTime.utc(2022, 5, 31).toISO(),
      price: '10000', // Starting price $10,000
    });

    await OraclePriceTable.create({
      ...defaultOraclePrice,
      effectiveAtHeight: '1',
      effectiveAt: DateTime.utc(2022, 6, 1).toISO(),
      price: '10000', // Price at height 1: $10,000
    });

    await OraclePriceTable.create({
      ...defaultOraclePrice,
      effectiveAtHeight: '7',
      effectiveAt: DateTime.utc(2022, 6, 7).toISO(),
      price: '11000', // Price at height 7: $11,000
    });

    await OraclePriceTable.create({
      ...defaultOraclePrice,
      effectiveAtHeight: '10',
      effectiveAt: DateTime.utc(2022, 6, 10).toISO(),
      price: '12000', // Price at height 10: $12,000
    });

    // Create a different oracle price for market2
    await OraclePriceTable.create({
      marketId: defaultMarket2.id,
      price: '1000', // ETH price at $1000
      effectiveAt: DateTime.utc(2022, 5, 31).toISO(),
      effectiveAtHeight: '0',
    });

    await OraclePriceTable.create({
      marketId: defaultMarket2.id,
      price: '1000', // ETH price at $1000
      effectiveAt: DateTime.utc(2022, 6, 1).toISO(),
      effectiveAtHeight: '1',
    });

    await OraclePriceTable.create({
      marketId: defaultMarket2.id,
      price: '900', // ETH price decreased to $900
      effectiveAt: DateTime.utc(2022, 6, 7).toISO(),
      effectiveAtHeight: '7',
    });

    await OraclePriceTable.create({
      marketId: defaultMarket2.id,
      price: '800', // ETH price decreased to $800
      effectiveAt: DateTime.utc(2022, 6, 10).toISO(),
      effectiveAtHeight: '10',
    });

    // Create transfers at height 1 for the subaccount
    await TransferTable.create({
      ...defaultTransfer,
      createdAtHeight: '1',
      createdAt: defaultBlock.time,
    });

    // Create LONG position for BTC at height 1
    await PerpetualPositionTable.create({
      subaccountId: defaultSubaccountId,
      perpetualId: defaultPerpetualMarket.id,
      side: PositionSide.LONG,
      status: PerpetualPositionStatus.OPEN,
      size: '2', // 2 BTC long
      maxSize: '2',
      entryPrice: '10000', // Entry at $10,000
      sumOpen: '2',
      sumClose: '0',
      createdAt: defaultBlock.time,
      createdAtHeight: '1',
      openEventId: defaultTendermintEventId,
      lastEventId: defaultTendermintEventId,
      settledFunding: '0',
    });

    // Create SHORT position for ETH at height 7
    await PerpetualPositionTable.create({
      subaccountId: defaultSubaccountId,
      perpetualId: defaultPerpetualMarket2.id,
      side: PositionSide.SHORT,
      status: PerpetualPositionStatus.OPEN,
      size: '-5', // 5 ETH short
      maxSize: '5',
      entryPrice: '900', // Entry at $900
      sumOpen: '5',
      sumClose: '0',
      createdAt: DateTime.utc(2022, 6, 7).toISO(),
      createdAtHeight: '7',
      openEventId: defaultTendermintEventId2,
      lastEventId: defaultTendermintEventId2,
      settledFunding: '0',
    });

    // Create funding payments at heights 1 and 10
    await FundingPaymentsTable.create({
      ...defaultFundingPayment,
      createdAtHeight: '1',
      createdAt: defaultBlock.time,
      payment: '2',
    });

    await FundingPaymentsTable.create({
      ...defaultFundingPayment,
      createdAtHeight: '10',
      createdAt: defaultBlock10.time,
      payment: '5',
    });

    await updatePnlTask();

    // Check that PNL entries were created
    const pnlRecords = await PnlTable.findAll({}, []);

    // Find records at height 10
    const recordsAtHeight10 = pnlRecords.results.filter((r) => r.createdAtHeight === '10');
    expect(recordsAtHeight10.length).toBe(2); // One for each subaccount with transfer history

    // Find PNL for the subaccount with positions
    const subaccountWithPositionsPnl = recordsAtHeight10.find((r) => r.subaccountId === defaultSubaccountId);
    expect(subaccountWithPositionsPnl).toBeDefined();

    // Calculate expected position effects
    // BTC: LONG 2 BTC, entry at $10000, current price $12000
    // BTC effect = (12000 - 10000) * 2 = $4000

    // ETH: SHORT 5 ETH, entry at $900, current price $800
    // ETH effect = (900 - 800) * 5 = $500

    // Total position effect = $4000 + $500 = $4500
    expect(subaccountWithPositionsPnl?.deltaPositionEffects).toBe('4500');

    // Total PNL should match position effects + funding
    expect(subaccountWithPositionsPnl?.totalPnl).toBe('4507'); // 4500 from positions + 2 + 5 funding payments

    // Find PNL for the subaccount without positions
    const subaccountWithoutPositionsPnl = recordsAtHeight10.find((r) => r.subaccountId === defaultSubaccountId2);
    expect(subaccountWithoutPositionsPnl).toBeDefined();

    // Subaccount without positions should have zero position effects
    expect(subaccountWithoutPositionsPnl?.deltaPositionEffects).toBe('0');
    expect(subaccountWithoutPositionsPnl?.totalPnl).toBe('0');

    // Verify persistent cache was updated
    const cache = await PersistentCacheTable.findById(
      PersistentCacheKeys.PNL_LAST_PROCESSED_HEIGHT
    );
    expect(cache).toBeDefined();
    expect(cache?.value).toBe('10');
  });

  it('calculates position effects correctly for closed positions', async () => {
    // Create blocks for heights 0, 1, 5, 8, and 10
    await BlockTable.create({
      blockHeight: '0',
      time: DateTime.utc(2022, 5, 31).toISO(),
    });
    // Height 1 is created by seedData() as defaultBlock
    await BlockTable.create({
      blockHeight: '5',
      time: DateTime.utc(2022, 6, 5).toISO(),
    });
    await BlockTable.create({
      blockHeight: '8',
      time: DateTime.utc(2022, 6, 8).toISO(),
    });
    await BlockTable.create(defaultBlock10);

    // Create oracle prices with changing prices
    await OraclePriceTable.create({
      ...defaultOraclePrice,
      effectiveAtHeight: '0',
      effectiveAt: DateTime.utc(2022, 5, 31).toISO(),
      price: '10000', // Starting price $10,000
    });

    await OraclePriceTable.create({
      ...defaultOraclePrice,
      effectiveAtHeight: '1',
      effectiveAt: DateTime.utc(2022, 6, 1).toISO(),
      price: '10000', // Price at height 1: $10,000
    });

    await OraclePriceTable.create({
      ...defaultOraclePrice,
      effectiveAtHeight: '5',
      effectiveAt: DateTime.utc(2022, 6, 5).toISO(),
      price: '11000', // Price at height 5: $11,000
    });

    await OraclePriceTable.create({
      ...defaultOraclePrice,
      effectiveAtHeight: '8',
      effectiveAt: DateTime.utc(2022, 6, 8).toISO(),
      price: '9000', // Price at height 8: $9,000
    });

    await OraclePriceTable.create({
      ...defaultOraclePrice,
      effectiveAtHeight: '10',
      effectiveAt: DateTime.utc(2022, 6, 10).toISO(),
      price: '12000', // Price at height 10: $12,000
    });

    // Create a different oracle price for market2
    await OraclePriceTable.create({
      marketId: defaultMarket2.id,
      price: '1000', // ETH price at $1000
      effectiveAt: DateTime.utc(2022, 5, 31).toISO(),
      effectiveAtHeight: '0',
    });

    await OraclePriceTable.create({
      marketId: defaultMarket2.id,
      price: '1000', // ETH price at $1000
      effectiveAt: DateTime.utc(2022, 6, 1).toISO(),
      effectiveAtHeight: '1',
    });

    await OraclePriceTable.create({
      marketId: defaultMarket2.id,
      price: '1200', // ETH price increased to $1200
      effectiveAt: DateTime.utc(2022, 6, 5).toISO(),
      effectiveAtHeight: '5',
    });

    await OraclePriceTable.create({
      marketId: defaultMarket2.id,
      price: '900', // ETH price dropped to $900
      effectiveAt: DateTime.utc(2022, 6, 8).toISO(),
      effectiveAtHeight: '8',
    });

    await OraclePriceTable.create({
      marketId: defaultMarket2.id,
      price: '800', // ETH price decreased to $800
      effectiveAt: DateTime.utc(2022, 6, 10).toISO(),
      effectiveAtHeight: '10',
    });

    // Create transfers at height 1 for the subaccount
    await TransferTable.create({
      ...defaultTransfer,
      createdAtHeight: '1',
      createdAt: defaultBlock.time,
    });

    // Position 1: Created at height 1, closed at height 8
    // BTC LONG position opened at height 1 with price $10,000 and closed at height 8 with price $9,500
    await PerpetualPositionTable.create({
      subaccountId: defaultSubaccountId,
      perpetualId: defaultPerpetualMarket.id,
      side: PositionSide.LONG,
      status: PerpetualPositionStatus.CLOSED,
      size: '3', // 3 BTC long
      maxSize: '3',
      entryPrice: '10000', // Entry at $10,000
      exitPrice: '9500', // Exit at $9,500
      sumOpen: '3',
      sumClose: '3',
      createdAt: defaultBlock.time,
      createdAtHeight: '1',
      closedAt: DateTime.utc(2022, 6, 8).toISO(),
      closedAtHeight: '8',
      openEventId: defaultTendermintEventId,
      lastEventId: defaultTendermintEventId,
      settledFunding: '0',
    });

    // Position 2: Created at height 5, closed at height 10
    // ETH SHORT position opened at height 5 with price $1200 and closed at height 10 with price $850
    await PerpetualPositionTable.create({
      subaccountId: defaultSubaccountId,
      perpetualId: defaultPerpetualMarket2.id,
      side: PositionSide.SHORT,
      status: PerpetualPositionStatus.CLOSED,
      size: '-4', // 4 ETH short
      maxSize: '4',
      entryPrice: '1200', // Entry at $1,200
      exitPrice: '850', // Exit at $850
      sumOpen: '4',
      sumClose: '4',
      createdAt: DateTime.utc(2022, 6, 5).toISO(),
      createdAtHeight: '5',
      closedAt: defaultBlock10.time,
      closedAtHeight: '10',
      openEventId: defaultTendermintEventId2,
      lastEventId: defaultTendermintEventId2,
      settledFunding: '0',
    });

    // Create funding payments at heights 1 and 10
    await FundingPaymentsTable.create({
      ...defaultFundingPayment,
      createdAtHeight: '1',
      createdAt: defaultBlock.time,
      payment: '2',
    });

    await FundingPaymentsTable.create({
      ...defaultFundingPayment,
      createdAtHeight: '10',
      createdAt: defaultBlock10.time,
      payment: '5',
    });

    // Run the task - this should process both heights 1 and 10
    await updatePnlTask();

    // Check that PNL entries were created
    const pnlRecords = await PnlTable.findAll({}, []);
    console.log('PNL Records:', pnlRecords.results);

    // Find records at height 10
    const recordsAtHeight10 = pnlRecords.results.filter((r) => r.createdAtHeight === '10');
    console.log('Records at Height 10:', recordsAtHeight10);

    expect(recordsAtHeight10.length).toBe(2); // One for each subaccount with transfer history

    // Find PNL for the subaccount with positions
    const subaccountWithPositionsPnl = recordsAtHeight10.find((r) => r.subaccountId === defaultSubaccountId);
    expect(subaccountWithPositionsPnl).toBeDefined();

    // Calculate expected position effects for closed positions
    // Position 1: BTC LONG, entry at $10,000, exit at $9,500, size 3
    // PnL = (9500 - 10000) * 3 = -$1,500 (loss)

    // Position 2: ETH SHORT, entry at $1,200, exit at $850, size 4
    // PnL = (1200 - 850) * 4 = $1,400 (profit)

    // Total position effect = -$1,500 + $1,400 = -$100
    expect(subaccountWithPositionsPnl?.deltaPositionEffects).toBe('-100');

    // Total PNL should match position effects + funding
    // Position effects: -$100
    // Funding payments: $2 + $5 = $7
    // Total: -$100 + $7 = -$93
    expect(subaccountWithPositionsPnl?.totalPnl).toBe('-93');

    // Find PNL for the subaccount without positions
    const subaccountWithoutPositionsPnl = recordsAtHeight10.find((r) => r.subaccountId === defaultSubaccountId2);
    expect(subaccountWithoutPositionsPnl).toBeDefined();

    // Subaccount without positions should have zero position effects
    expect(subaccountWithoutPositionsPnl?.deltaPositionEffects).toBe('0');
    expect(subaccountWithoutPositionsPnl?.totalPnl).toBe('0');

    // Verify persistent cache was updated
    const cache = await PersistentCacheTable.findById(
      PersistentCacheKeys.PNL_LAST_PROCESSED_HEIGHT
    );
    expect(cache).toBeDefined();
    expect(cache?.value).toBe('10');
  });

  it('calculates comprehensive PNL with funding payments, open and closed positions across multiple periods', async () => {
    // Create blocks for heights 0, 1, 5, 8, and 10
    await BlockTable.create({
      blockHeight: '0',
      time: DateTime.utc(2022, 5, 31).toISO(),
    });
    // Height 1 is created by seedData() as defaultBlock
    await BlockTable.create({
      blockHeight: '5',
      time: DateTime.utc(2022, 6, 5).toISO(),
    });
    await BlockTable.create({
      blockHeight: '8',
      time: DateTime.utc(2022, 6, 8).toISO(),
    });
    await BlockTable.create(defaultBlock10);

    // Create oracle prices with changing prices
    await OraclePriceTable.create({
      ...defaultOraclePrice,
      effectiveAtHeight: '0',
      effectiveAt: DateTime.utc(2022, 5, 31).toISO(),
      price: '10000', // Starting price $10,000
    });

    await OraclePriceTable.create({
      ...defaultOraclePrice,
      effectiveAtHeight: '1',
      effectiveAt: DateTime.utc(2022, 6, 1).toISO(),
      price: '10000', // Price at height 1: $10,000
    });

    await OraclePriceTable.create({
      ...defaultOraclePrice,
      effectiveAtHeight: '5',
      effectiveAt: DateTime.utc(2022, 6, 5).toISO(),
      price: '11000', // Price at height 5: $11,000
    });

    await OraclePriceTable.create({
      ...defaultOraclePrice,
      effectiveAtHeight: '8',
      effectiveAt: DateTime.utc(2022, 6, 8).toISO(),
      price: '9000', // Price at height 8: $9,000
    });

    await OraclePriceTable.create({
      ...defaultOraclePrice,
      effectiveAtHeight: '10',
      effectiveAt: DateTime.utc(2022, 6, 10).toISO(),
      price: '12000', // Price at height 10: $12,000
    });

    // Create a different oracle price for market2
    await OraclePriceTable.create({
      marketId: defaultMarket2.id,
      price: '1000', // ETH price at $1000
      effectiveAt: DateTime.utc(2022, 5, 31).toISO(),
      effectiveAtHeight: '0',
    });

    await OraclePriceTable.create({
      marketId: defaultMarket2.id,
      price: '1000', // ETH price at $1000
      effectiveAt: DateTime.utc(2022, 6, 1).toISO(),
      effectiveAtHeight: '1',
    });

    await OraclePriceTable.create({
      marketId: defaultMarket2.id,
      price: '1200', // ETH price increased to $1200
      effectiveAt: DateTime.utc(2022, 6, 5).toISO(),
      effectiveAtHeight: '5',
    });

    await OraclePriceTable.create({
      marketId: defaultMarket2.id,
      price: '900', // ETH price dropped to $900
      effectiveAt: DateTime.utc(2022, 6, 8).toISO(),
      effectiveAtHeight: '8',
    });

    await OraclePriceTable.create({
      marketId: defaultMarket2.id,
      price: '800', // ETH price decreased to $800
      effectiveAt: DateTime.utc(2022, 6, 10).toISO(),
      effectiveAtHeight: '10',
    });

    // Create transfers to ensure subaccounts are included
    await TransferTable.create({
      ...defaultTransfer,
      createdAtHeight: '1',
      createdAt: defaultBlock.time,
    });

    // ====== POSITION 1: BTC LONG, OPEN ======
    // Created at height 1, still open at height 10
    await PerpetualPositionTable.create({
      subaccountId: defaultSubaccountId,
      perpetualId: defaultPerpetualMarket.id,
      side: PositionSide.LONG,
      status: PerpetualPositionStatus.OPEN,
      size: '2', // 2 BTC long
      maxSize: '2',
      entryPrice: '10000', // Entry at $10,000
      sumOpen: '2',
      sumClose: '0',
      createdAt: defaultBlock.time,
      createdAtHeight: '1',
      openEventId: defaultTendermintEventId,
      lastEventId: defaultTendermintEventId,
      settledFunding: '0',
    });

    // ====== POSITION 2: ETH SHORT, CLOSED ======
    // Created at height 1, closed at height 8
    await PerpetualPositionTable.create({
      subaccountId: defaultSubaccountId,
      perpetualId: defaultPerpetualMarket2.id,
      side: PositionSide.SHORT,
      status: PerpetualPositionStatus.CLOSED,
      size: '-3', // 3 ETH short
      maxSize: '3',
      entryPrice: '1000', // Entry at $1,000
      exitPrice: '900', // Exit at $900
      sumOpen: '3',
      sumClose: '3',
      createdAt: defaultBlock.time,
      createdAtHeight: '1',
      closedAt: DateTime.utc(2022, 6, 8).toISO(),
      closedAtHeight: '8',
      openEventId: defaultTendermintEventId2,
      lastEventId: defaultTendermintEventId2,
      settledFunding: '0',
    });

    // ====== POSITION 3: BTC SHORT, OPEN ======
    // Created at height 5, still open at height 10
    await PerpetualPositionTable.create({
      subaccountId: defaultSubaccountId2,
      perpetualId: defaultPerpetualMarket.id,
      side: PositionSide.SHORT,
      status: PerpetualPositionStatus.OPEN,
      size: '-4', // 4 BTC short
      maxSize: '4',
      entryPrice: '11000', // Entry at $11,000
      sumOpen: '4',
      sumClose: '0',
      createdAt: DateTime.utc(2022, 6, 5).toISO(),
      createdAtHeight: '5',
      openEventId: defaultTendermintEventId3,
      lastEventId: defaultTendermintEventId3,
      settledFunding: '0',
    });

    // ====== POSITION 4: ETH LONG, CLOSED ======
    // Created at height 5, closed at height 10
    await PerpetualPositionTable.create({
      subaccountId: defaultSubaccountId2,
      perpetualId: defaultPerpetualMarket2.id,
      side: PositionSide.LONG,
      status: PerpetualPositionStatus.CLOSED,
      size: '5', // 5 ETH long
      maxSize: '5',
      entryPrice: '1200', // Entry at $1,200
      exitPrice: '800', // Exit at $800
      sumOpen: '5',
      sumClose: '5',
      createdAt: DateTime.utc(2022, 6, 5).toISO(),
      createdAtHeight: '5',
      closedAt: defaultBlock10.time,
      closedAtHeight: '10',
      openEventId: defaultTendermintEventId4,
      lastEventId: defaultTendermintEventId4,
      settledFunding: '0',
    });

    // Create funding payments at heights 5 and 10 to ensure we process those heights
    await FundingPaymentsTable.create({
      ...defaultFundingPayment,
      createdAtHeight: '5',
      createdAt: DateTime.utc(2022, 6, 5).toISO(),
      payment: '10',
    });

    await FundingPaymentsTable.create({
      ...defaultFundingPayment,
      createdAtHeight: '10',
      createdAt: defaultBlock10.time,
      payment: '15',
    });

    await updatePnlTask();

    const pnlRecords = await PnlTable.findAll({}, []);
    console.log('PNL Records:', pnlRecords.results);

    // Group records by height
    const recordsAtHeight5 = pnlRecords.results.filter((r) => r.createdAtHeight === '5');
    const recordsAtHeight10 = pnlRecords.results.filter((r) => r.createdAtHeight === '10');

    expect(recordsAtHeight5.length).toBe(2); // One for each subaccount
    expect(recordsAtHeight10.length).toBe(2); // One for each subaccount

    // ===== PNL CALCULATION AT HEIGHT 5 =====

    // Find PNL for subaccount1 at height 5
    const subaccount1PnlAtHeight5 = recordsAtHeight5.find((r) => r.subaccountId === defaultSubaccountId);
    expect(subaccount1PnlAtHeight5).toBeDefined();

    // At height 5, the position effects are calculated from height 0 to height 5
    // Open position: BTC LONG, entry at $10,000, current price $11,000, size 2
    // Position effect = (11000 - 10000) * 2 = $2,000
    expect(subaccount1PnlAtHeight5?.deltaPositionEffects).toBe('2000');

    // Funding payment at height 5 = $10
    expect(subaccount1PnlAtHeight5?.deltaFundingPayments).toBe('10');

    // Total PNL at height 5 = $2,000 + $10 = $2,010
    expect(subaccount1PnlAtHeight5?.totalPnl).toBe('2010');

    // Find PNL for subaccount2 at height 5
    const subaccount2PnlAtHeight5 = recordsAtHeight5.find((r) => r.subaccountId === defaultSubaccountId2);
    expect(subaccount2PnlAtHeight5).toBeDefined();

    // Subaccount2 has just opened positions at height 5, so no PnL yet
    expect(subaccount2PnlAtHeight5?.deltaPositionEffects).toBe('0');
    expect(subaccount2PnlAtHeight5?.deltaFundingPayments).toBe('0');
    expect(subaccount2PnlAtHeight5?.totalPnl).toBe('0');

    // ===== PNL CALCULATION AT HEIGHT 10 =====

    // Find PNL for subaccount1 at height 10
    const subaccount1PnlAtHeight10 = recordsAtHeight10.find((r) => r.subaccountId === defaultSubaccountId);
    expect(subaccount1PnlAtHeight10).toBeDefined();

    // For the BTC LONG position (created at height 1):
    // - Oracle price at height 5: $11,000
    // - Oracle price at height 10: $12,000
    // - Size: 2 BTC
    // - PNL = (12000 - 11000) * 2 = $2,000

    // For the ETH SHORT position (created at height 1, closed at height 8):
    // - Oracle price at height 5: $1,200
    // - Exit price: $900
    // - Size: -3 (short)
    // - PNL = (900 - 1200) * -3 = (1200 - 900) * 3 = $900

    // Total deltaPositionEffects = $2,000 + $900 = $2,900
    expect(subaccount1PnlAtHeight10?.deltaPositionEffects).toBe('2900');

    // Funding payment at height 10 = $15
    expect(subaccount1PnlAtHeight10?.deltaFundingPayments).toBe('15');

    // Total PNL accumulates: previous totalPnl + current deltaPositionEffects + current deltaFundingPayments
    // Total PNL = $2,010 + $2,900 + $15 = $4,925
    expect(subaccount1PnlAtHeight10?.totalPnl).toBe('4925');

    // Find PNL for subaccount2 at height 10
    const subaccount2PnlAtHeight10 = recordsAtHeight10.find((r) => r.subaccountId === defaultSubaccountId2);
    expect(subaccount2PnlAtHeight10).toBeDefined();

    // At height 5:
    // No position effects yet

    // At height 10:
    // BTC SHORT: (11000 - 12000) * 4 = -$4,000 (new since height 5)
    // ETH LONG closed: (800 - 1200) * 5 = -$2,000 (new since height 5)

    // So deltaPositionEffects = -$4,000 - $2,000 = -$6,000
    expect(subaccount2PnlAtHeight10?.deltaPositionEffects).toBe('-6000');

    // Funding payment = $0 (no funding payments for this subaccount)
    expect(subaccount2PnlAtHeight10?.deltaFundingPayments).toBe('0');

    // Total PNL = Previous PNL + Current Position Effects + Current Funding
    // Total PNL = $0 + (-$6,000) + $0 = -$6,000
    expect(subaccount2PnlAtHeight10?.totalPnl).toBe('-6000');

    // Verify persistent cache was updated to the highest height
    const cache = await PersistentCacheTable.findById(
      PersistentCacheKeys.PNL_LAST_PROCESSED_HEIGHT
    );
    expect(cache).toBeDefined();
    expect(cache?.value).toBe('10');
  });
});