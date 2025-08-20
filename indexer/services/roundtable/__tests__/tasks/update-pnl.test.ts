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
import { DateTime } from 'luxon';

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
  defaultPerpetualMarket2,
  defaultMarket2,
  defaultTendermintEventId,
  defaultTendermintEventId2,
  defaultTendermintEventId3,
  defaultTendermintEventId4,
  defaultBlock10,
} from '@dydxprotocol-indexer/postgres/build/__tests__/helpers/constants';

describe('update-pnl', () => {
  // Common date helpers
  const MAY_31 = DateTime.utc(2022, 5, 31).toISO();
  const JUNE_1 = DateTime.utc(2022, 6, 1).toISO();
  const JUNE_5 = DateTime.utc(2022, 6, 5).toISO();
  const JUNE_7 = DateTime.utc(2022, 6, 7).toISO();
  const JUNE_8 = DateTime.utc(2022, 6, 8).toISO();
  const JUNE_10 = DateTime.utc(2022, 6, 10).toISO();

  // Common price points
  const BTC_PRICES = {
    HEIGHT_0: '10000',
    HEIGHT_1: '10000',
    HEIGHT_5: '11000',
    HEIGHT_7: '11000',
    HEIGHT_8: '9000',
    HEIGHT_10: '12000',
  };

  const ETH_PRICES = {
    HEIGHT_0: '1000',
    HEIGHT_1: '1000',
    HEIGHT_5: '1200',
    HEIGHT_7: '900',
    HEIGHT_8: '900',
    HEIGHT_10: '800',
  };

  beforeAll(async () => {
    await dbHelpers.migrate();
  });

  beforeEach(async () => {
    await dbHelpers.clearData();
    await testMocks.seedData();

    // Common setup for all tests:

    // 1. Create block 0
    await BlockTable.create({
      blockHeight: '0',
      time: MAY_31,
    });

    // 2. Create oracle prices at height 0 and 1
    await OraclePriceTable.create({
      ...defaultOraclePrice,
      effectiveAtHeight: '0',
      effectiveAt: MAY_31,
      price: BTC_PRICES.HEIGHT_0,
    });

    await OraclePriceTable.create({
      ...defaultOraclePrice,
      effectiveAtHeight: '1',
      effectiveAt: JUNE_1,
      price: BTC_PRICES.HEIGHT_1,
    });

    // 3. Create a transfer at height 1
    await TransferTable.create({
      ...defaultTransfer,
      createdAtHeight: '1',
      createdAt: defaultBlock.time,
    });
  });

  afterAll(async () => {
    await dbHelpers.teardown();
    jest.resetAllMocks();
  });

  afterEach(async () => {
    await dbHelpers.clearData();
  });

  async function createBlocks(heights: { [height: string]: string }) {
    const promises = Object.entries(heights).map(([height, time]) => BlockTable.create({
      blockHeight: height,
      time,
    }),
    );
    await Promise.all(promises);
  }

  async function createOraclePrices(
    marketId: number,
    prices: { [height: string]: { price: string, time: string } },
  ) {
    const promises = Object.entries(prices).map(
      ([height, { price, time }]) => OraclePriceTable.create({
        marketId,
        price,
        effectiveAtHeight: height,
        effectiveAt: time,
      }),
    );
    await Promise.all(promises);
  }

  async function createFundingPayments(
    payments: { height: string, time: string, amount: string }[],
  ) {
    const promises = payments.map(({ height, time, amount }) => FundingPaymentsTable.create({
      ...defaultFundingPayment,
      createdAtHeight: height,
      createdAt: time,
      payment: amount,
    }),
    );
    await Promise.all(promises);
  }

  function verifyPnlRecord(
    record: any,
    expectedValues: {
      deltaFundingPayments: string,
      deltaPositionEffects: string,
      totalPnl: string,
    },
  ) {
    expect(record.deltaFundingPayments).toBe(expectedValues.deltaFundingPayments);
    expect(record.deltaPositionEffects).toBe(expectedValues.deltaPositionEffects);
    expect(record.totalPnl).toBe(expectedValues.totalPnl);
  }

  async function verifyCache(expectedHeight: string) {
    const cache = await PersistentCacheTable.findById(
      PersistentCacheKeys.PNL_LAST_PROCESSED_HEIGHT,
    );
    expect(cache).toBeDefined();
    expect(cache?.value).toBe(expectedHeight);
  }

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

  it('computes zero PnL for subaccounts with transfers but no positions', async () => {
    await FundingPaymentsTable.create(defaultFundingPayment);

    await updatePnlTask();

    // Check that PNL entries were created for both subaccounts
    const pnlRecords = await PnlTable.findAll({}, []);

    // Expect 2 records (one for each subaccount)
    expect(pnlRecords.results.length).toBe(2);

    const { subaccount1Pnl, subaccount2Pnl } = findPnlRecords(pnlRecords.results, '1');

    expect(subaccount1Pnl).toBeDefined();
    expect(subaccount2Pnl).toBeDefined();

    // Verify PNL values
    verifyPnlRecord(subaccount1Pnl, {
      deltaFundingPayments: '5',
      deltaPositionEffects: '0',
      totalPnl: '5',
    });

    verifyPnlRecord(subaccount2Pnl, {
      deltaFundingPayments: '0',
      deltaPositionEffects: '0',
      totalPnl: '0',
    });

    // Verify createdAt and createdAtHeight values
    for (const record of pnlRecords.results) {
      // createdAtHeight should be the end height (1)
      expect(record.createdAtHeight).toBe('1');
      // createdAt should match the timestamp of the oracle price at height 1
      expect(record.createdAt).toBe(JUNE_1);
    }

    await verifyCache('1');
  });

  it('processes multiple periods of pnl calculations', async () => {
    await OraclePriceTable.create({
      ...defaultOraclePrice,
      effectiveAtHeight: '2',
      effectiveAt: JUNE_5,
    });

    await createFundingPayments([
      { height: '1', time: defaultBlock.time, amount: '5' },
      { height: '2', time: defaultBlock2.time, amount: '10' },
    ]);

    await updatePnlTask();

    // Check that PNL entries were created
    const pnlRecords = await PnlTable.findAll({}, []);

    // Sort records by height
    const sortedRecords = pnlRecords.results.sort((a, b) => parseInt(a.createdAtHeight, 10) -
      parseInt(b.createdAtHeight, 10));

    // We should have records for both heights, 2 heights * 2 subaccounts = 4 records total
    expect(sortedRecords.length).toBe(4);

    const { recordsAtHeight: recordsAtHeight1, subaccount1Pnl: subaccount1PnlAtHeight1, subaccount2Pnl: subaccount2PnlAtHeight1 } = findPnlRecords(sortedRecords, '1');

    const { recordsAtHeight: recordsAtHeight2, subaccount1Pnl: subaccount1PnlAtHeight2, subaccount2Pnl: subaccount2PnlAtHeight2 } = findPnlRecords(sortedRecords, '2');

    expect(recordsAtHeight1.length).toBe(2);
    expect(recordsAtHeight2.length).toBe(2);

    // Verify PNL values at height 1
    verifyPnlRecord(subaccount1PnlAtHeight1, {
      deltaFundingPayments: '5',
      deltaPositionEffects: '0',
      totalPnl: '5',
    });

    verifyPnlRecord(subaccount2PnlAtHeight1, {
      deltaFundingPayments: '0',
      deltaPositionEffects: '0',
      totalPnl: '0',
    });

    // Verify PNL values at height 2
    verifyPnlRecord(subaccount1PnlAtHeight2, {
      deltaFundingPayments: '10',
      deltaPositionEffects: '0',
      totalPnl: '15', // 5 from height 1 + 10 from height 2
    });

    verifyPnlRecord(subaccount2PnlAtHeight2, {
      deltaFundingPayments: '0',
      deltaPositionEffects: '0',
      totalPnl: '0',
    });

    // Verify createdAt matches oracle price timestamps
    for (const record of recordsAtHeight1) {
      expect(record.createdAt).toBe(JUNE_1);
    }

    for (const record of recordsAtHeight2) {
      expect(record.createdAt).toBe(JUNE_5);
    }

    await verifyCache('2');
  });

  it('correctly sums multiple funding payments for the same subaccount at the same height', async () => {
    // Create oracle price at height 2
    await OraclePriceTable.create({
      ...defaultOraclePrice,
      effectiveAtHeight: '2',
      effectiveAt: JUNE_5,
    });

    // Create oracle prices for the second market
    await createOraclePrices(defaultMarket2.id, {
      0: { price: ETH_PRICES.HEIGHT_0, time: MAY_31 },
      1: { price: ETH_PRICES.HEIGHT_1, time: JUNE_1 },
      2: { price: ETH_PRICES.HEIGHT_5, time: JUNE_5 },
    });

    // Create multiple funding payments for the same subaccount at height 2
    await FundingPaymentsTable.create({
      ...defaultFundingPayment,
      subaccountId: defaultSubaccountId,
      createdAtHeight: '2',
      createdAt: JUNE_5,
      payment: '10',
      perpetualId: defaultPerpetualMarket.id,
      ticker: defaultPerpetualMarket.ticker,
    });

    await FundingPaymentsTable.create({
      ...defaultFundingPayment,
      subaccountId: defaultSubaccountId,
      createdAtHeight: '2',
      createdAt: JUNE_5,
      payment: '5',
      perpetualId: defaultPerpetualMarket2.id,
      ticker: defaultPerpetualMarket2.ticker,
    });

    await FundingPaymentsTable.create({
      ...defaultFundingPayment,
      subaccountId: defaultSubaccountId,
      createdAtHeight: '2',
      createdAt: JUNE_5,
      payment: '3',
      perpetualId: '3', // A third market
      ticker: 'ISO-USD',
    });

    await updatePnlTask();

    const pnlRecords = await PnlTable.findAll({}, []);

    const { recordsAtHeight: recordsAtHeight2, subaccount1Pnl } = findPnlRecords(pnlRecords.results, '2');

    expect(recordsAtHeight2.length).toBe(2); // One for each subaccount in the transfer

    // Verify the first subaccount has the sum of all three funding payments (10 + 5 + 3 = 18)
    verifyPnlRecord(subaccount1Pnl, {
      deltaFundingPayments: '18', // Sum of 10 + 5 + 3
      deltaPositionEffects: '0', // No position effects in this test
      totalPnl: '18', // Total equals funding payments
    });

    // Verify cache was updated
    await verifyCache('2');
  });

  it('calculates position effects correctly for open positions in the same subaccount', async () => {
    await createBlocks({
      7: JUNE_7,
      10: JUNE_10,
    });

    await OraclePriceTable.create({
      ...defaultOraclePrice,
      effectiveAtHeight: '7',
      effectiveAt: JUNE_7,
      price: BTC_PRICES.HEIGHT_7,
    });

    await OraclePriceTable.create({
      ...defaultOraclePrice,
      effectiveAtHeight: '10',
      effectiveAt: JUNE_10,
      price: BTC_PRICES.HEIGHT_10,
    });

    // Create ETH oracle prices
    await createOraclePrices(defaultMarket2.id, {
      0: { price: ETH_PRICES.HEIGHT_0, time: MAY_31 },
      1: { price: ETH_PRICES.HEIGHT_1, time: JUNE_1 },
      7: { price: ETH_PRICES.HEIGHT_7, time: JUNE_7 },
      10: { price: ETH_PRICES.HEIGHT_10, time: JUNE_10 },
    });

    // Create funding payments
    await createFundingPayments([
      { height: '1', time: defaultBlock.time, amount: '2' },
      { height: '10', time: defaultBlock10.time, amount: '5' },
    ]);

    // Create positions
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
      createdAt: JUNE_7,
      createdAtHeight: '7',
      openEventId: defaultTendermintEventId2,
      lastEventId: defaultTendermintEventId2,
      settledFunding: '0',
    });

    await updatePnlTask();

    // Check PNL records
    const pnlRecords = await PnlTable.findAll({}, []);
    const { recordsAtHeight: recordsAtHeight10, subaccount1Pnl: subaccountWithPositionsPnl, subaccount2Pnl: subaccountWithoutPositionsPnl } = findPnlRecords(pnlRecords.results, '10');

    expect(recordsAtHeight10.length).toBe(2); // One for each subaccount

    // Verify PNL values for subaccount with positions
    verifyPnlRecord(subaccountWithPositionsPnl, {
      deltaPositionEffects: '4500', // (12000-10000)*2 + (900-800)*5 = 4000 + 500 = 4500
      deltaFundingPayments: '5',
      totalPnl: '4507', // 4500 + 2 + 5 = 4507
    });

    // Verify PNL values for subaccount without positions
    verifyPnlRecord(subaccountWithoutPositionsPnl, {
      deltaPositionEffects: '0',
      deltaFundingPayments: '0',
      totalPnl: '0',
    });

    await verifyCache('10');
  });

  it('calculates position effects correctly for closed positions', async () => {
    await createBlocks({
      5: JUNE_5,
      8: JUNE_8,
      10: JUNE_10,
    });

    await createOraclePrices(defaultOraclePrice.marketId, {
      5: { price: BTC_PRICES.HEIGHT_5, time: JUNE_5 },
      8: { price: BTC_PRICES.HEIGHT_8, time: JUNE_8 },
      10: { price: BTC_PRICES.HEIGHT_10, time: JUNE_10 },
    });

    // Create ETH oracle prices
    await createOraclePrices(defaultMarket2.id, {
      0: { price: ETH_PRICES.HEIGHT_0, time: MAY_31 },
      1: { price: ETH_PRICES.HEIGHT_1, time: JUNE_1 },
      5: { price: ETH_PRICES.HEIGHT_5, time: JUNE_5 },
      8: { price: ETH_PRICES.HEIGHT_8, time: JUNE_8 },
      10: { price: ETH_PRICES.HEIGHT_10, time: JUNE_10 },
    });

    // Create funding payments at heights 1 and 10
    await createFundingPayments([
      { height: '1', time: defaultBlock.time, amount: '2' },
      { height: '10', time: defaultBlock10.time, amount: '5' },
    ]);

    // Position 1: Created at height 1, closed at height 8
    // BTC LONG position opened at height 1 with price $10,000
    // and closed at height 8 with price $9,500
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
      closedAt: JUNE_8,
      closedAtHeight: '8',
      openEventId: defaultTendermintEventId,
      lastEventId: defaultTendermintEventId,
      settledFunding: '0',
    });

    // Position 2: Created at height 5, closed at height 10
    // ETH SHORT position opened at height 5 with price $1200
    // and closed at height 10 with price $850
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
      createdAt: JUNE_5,
      createdAtHeight: '5',
      closedAt: JUNE_10,
      closedAtHeight: '10',
      openEventId: defaultTendermintEventId2,
      lastEventId: defaultTendermintEventId2,
      settledFunding: '0',
    });

    // Run the task
    await updatePnlTask();

    // Check PNL records
    const pnlRecords = await PnlTable.findAll({}, []);
    const { recordsAtHeight: recordsAtHeight10, subaccount1Pnl: subaccountWithPositionsPnl, subaccount2Pnl: subaccountWithoutPositionsPnl } = findPnlRecords(pnlRecords.results, '10');

    expect(recordsAtHeight10.length).toBe(2);

    // Verify PNL values for subaccount with positions
    verifyPnlRecord(subaccountWithPositionsPnl, {
      deltaPositionEffects: '-100', // (9500-10000)*3 + (1200-850)*4 = -1500 + 1400 = -100
      deltaFundingPayments: '5',
      totalPnl: '-93', // -100 + 2 + 5 = -93
    });

    // Verify PNL values for subaccount without positions
    verifyPnlRecord(subaccountWithoutPositionsPnl, {
      deltaPositionEffects: '0',
      deltaFundingPayments: '0',
      totalPnl: '0',
    });

    await verifyCache('10');
  });

  it('calculates comprehensive PNL with funding payments, open and closed positions across multiple periods', async () => {
    await createBlocks({
      5: JUNE_5,
      8: JUNE_8,
      10: JUNE_10,
    });

    await createOraclePrices(defaultOraclePrice.marketId, {
      5: { price: BTC_PRICES.HEIGHT_5, time: JUNE_5 },
      8: { price: BTC_PRICES.HEIGHT_8, time: JUNE_8 },
      10: { price: BTC_PRICES.HEIGHT_10, time: JUNE_10 },
    });

    // Create ETH oracle prices
    await createOraclePrices(defaultMarket2.id, {
      0: { price: ETH_PRICES.HEIGHT_0, time: MAY_31 },
      1: { price: ETH_PRICES.HEIGHT_1, time: JUNE_1 },
      5: { price: ETH_PRICES.HEIGHT_5, time: JUNE_5 },
      8: { price: ETH_PRICES.HEIGHT_8, time: JUNE_8 },
      10: { price: ETH_PRICES.HEIGHT_10, time: JUNE_10 },
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
      closedAt: JUNE_8,
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
      createdAt: JUNE_5,
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
      createdAt: JUNE_5,
      createdAtHeight: '5',
      closedAt: JUNE_10,
      closedAtHeight: '10',
      openEventId: defaultTendermintEventId4,
      lastEventId: defaultTendermintEventId4,
      settledFunding: '0',
    });

    // Create funding payments at heights 5 and 10
    await createFundingPayments([
      { height: '5', time: JUNE_5, amount: '10' },
      { height: '10', time: JUNE_10, amount: '15' },
    ]);

    // Run the task
    await updatePnlTask();

    // Check PNL records
    const pnlRecords = await PnlTable.findAll({}, []);

    const { recordsAtHeight: recordsAtHeight5, subaccount1Pnl: subaccount1PnlAtHeight5, subaccount2Pnl: subaccount2PnlAtHeight5 } = findPnlRecords(pnlRecords.results, '5');

    const { recordsAtHeight: recordsAtHeight10, subaccount1Pnl: subaccount1PnlAtHeight10, subaccount2Pnl: subaccount2PnlAtHeight10 } = findPnlRecords(pnlRecords.results, '10');

    expect(recordsAtHeight5.length).toBe(2);
    expect(recordsAtHeight10.length).toBe(2);

    // ===== PNL CALCULATION AT HEIGHT 5 =====

    // At height 5, the position effects are calculated from height 0 to height 5
    // Open position: BTC LONG, entry at $10,000, current price $11,000, size 2
    // Position effect = (11000 - 10000) * 2 = $2,000
    verifyPnlRecord(subaccount1PnlAtHeight5, {
      deltaPositionEffects: '2000',
      deltaFundingPayments: '10',
      totalPnl: '2010', // 2000 + 10 = 2010
    });

    // Subaccount2 has just opened positions at height 5, so no PnL yet
    verifyPnlRecord(subaccount2PnlAtHeight5, {
      deltaPositionEffects: '0',
      deltaFundingPayments: '0',
      totalPnl: '0',
    });

    // ===== PNL CALCULATION AT HEIGHT 10 =====

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
    verifyPnlRecord(subaccount1PnlAtHeight10, {
      deltaPositionEffects: '2900',
      deltaFundingPayments: '15',
      totalPnl: '4925', // 2010 + 2900 + 15 = 4925
    });

    // At height 10:
    // BTC SHORT: (11000 - 12000) * 4 = -$4,000 (new since height 5)
    // ETH LONG closed: (800 - 1200) * 5 = -$2,000 (new since height 5)
    // So deltaPositionEffects = -$4,000 - $2,000 = -$6,000
    verifyPnlRecord(subaccount2PnlAtHeight10, {
      deltaPositionEffects: '-6000',
      deltaFundingPayments: '0',
      totalPnl: '-6000', // 0 + (-6000) + 0 = -6000
    });

    await verifyCache('10');
  });
});
