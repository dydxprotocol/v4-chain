import { stats } from '@dydxprotocol-indexer/base';
import {
  dbHelpers,
  testMocks,
  PersistentCacheTable,
  PersistentCacheKeys,
  BlockTable,
  Transaction,
  FillTable,
  WalletTable,
  OrderTable,
  FundingPaymentsTable,
  Liquidity,
  OrderSide,
  FillType,
  OrderType,
  FundingIndexUpdatesTable,
  PositionSide,
} from '@dydxprotocol-indexer/postgres';
import updateFundingPaymentsTask from '../../src/tasks/update-funding-payments';
import { createdDateTime, defaultFill, defaultFundingIndexUpdate, defaultOrder, defaultOrderId, defaultPerpetualMarket, defaultSubaccountId, defaultTendermintEventId, defaultWallet } from '@dydxprotocol-indexer/postgres/build/__tests__/helpers/constants';
import { ORDER_FLAG_SHORT_TERM } from '@dydxprotocol-indexer/v4-proto-parser';

describe('update-funding-payments', () => {
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

  it('Creates funding payment from fills only (no snapshot)', async () => {
    // seed some data.
    await OrderTable.create(defaultOrder);
    await FillTable.create(defaultFill);
    await FundingIndexUpdatesTable.create(defaultFundingIndexUpdate);
    // Run task
    await updateFundingPaymentsTask();

    // Check that persistent cache was updated
    const persistentCache = await PersistentCacheTable.findById(
      PersistentCacheKeys.FUNDING_PAYMENTS_LAST_PROCESSED_HEIGHT,
    );
    expect(persistentCache).toBeDefined();
    expect(persistentCache?.value).toEqual('2');

    // check that funding payments were created, one should have been created.
    const fundingPayments = await FundingPaymentsTable.findAll(
        {},
        [],
    );
    expect(fundingPayments).toBeDefined();
    expect(fundingPayments?.length).toEqual(1);
  });

  it('Creates funding payment from snapshot only (no fills)', async () => {
    // Create initial funding payment as snapshot
    const snapshotFundingPayment = {
      subaccountId: defaultSubaccountId,
      createdAt: new Date().toISOString(),
      createdAtHeight: '2',
      perpetualId: defaultFundingIndexUpdate.perpetualId,
      ticker: 'BTC-USD',
      oraclePrice: '50000',
      size: '1',
      side: PositionSide.LONG,
      rate: '0.0001',
      payment: '5',
    };
    await BlockTable.create({
      blockHeight: '3',
      time: new Date().toISOString(),
    });
    await FundingPaymentsTable.create(snapshotFundingPayment);
    await FundingIndexUpdatesTable.create({
      perpetualId: defaultPerpetualMarket.id,
      eventId: defaultTendermintEventId,
      rate: '0.0004',
      oraclePrice: '10000',
      fundingIndex: '10050',
      effectiveAt: createdDateTime.toISO(),
      effectiveAtHeight: '3',
    });

    // Set initial persistent cache value to last funding payment height
    await PersistentCacheTable.upsert({
      key: PersistentCacheKeys.FUNDING_PAYMENTS_LAST_PROCESSED_HEIGHT,
      value: '2',
    });

    // Verify initial cache state
    const initialCache = await PersistentCacheTable.findById(
      PersistentCacheKeys.FUNDING_PAYMENTS_LAST_PROCESSED_HEIGHT,
    );
    expect(initialCache).toBeDefined();
    expect(initialCache?.value).toEqual('2');
    // Run task
    await updateFundingPaymentsTask();

    // Verify funding payments
    const fundingPayments = await FundingPaymentsTable.findAll(
      {},
      [],
    );
    expect(fundingPayments.length).toEqual(2); // Original snapshot + new payment
    expect(fundingPayments[1]).toMatchObject({
      subaccountId: defaultSubaccountId,
      perpetualId: defaultFundingIndexUpdate.perpetualId,
      size: '1', // Should maintain size from snapshot
      side: PositionSide.LONG,
    });

    // Verify persistent cache
    const persistentCache = await PersistentCacheTable.findById(
      PersistentCacheKeys.FUNDING_PAYMENTS_LAST_PROCESSED_HEIGHT,
    );
    expect(persistentCache?.value).toEqual('3');
  });

  it('Creates funding payment from fills only (no snapshot)', async () => {
    await OrderTable.create(defaultOrder);
    // Create fills without any existing funding payment
    const fill = {
      ...defaultFill,
      side: OrderSide.BUY,
      size: '2',
      createdAtHeight: '2',
    };
    await FillTable.create(fill);
    await FundingIndexUpdatesTable.create(defaultFundingIndexUpdate);

    // Run task
    await updateFundingPaymentsTask();

    // Verify funding payments
    const fundingPayments = await FundingPaymentsTable.findAll(
      {},
      [],
    );
    expect(fundingPayments.length).toEqual(1);
    expect(fundingPayments[0]).toMatchObject({
      subaccountId: defaultSubaccountId,
      perpetualId: defaultFundingIndexUpdate.perpetualId,
      size: '2', // Should match fill size
      side: PositionSide.LONG,
    });

    // Verify persistent cache
    const persistentCache = await PersistentCacheTable.findById(
      PersistentCacheKeys.FUNDING_PAYMENTS_LAST_PROCESSED_HEIGHT,
    );
    expect(persistentCache?.value).toEqual('2');
  });

  it('Maintains persistent cache value across multiple runs and does not create funding payments if no unprocessedfunding index update', async () => {
    // Initial setup
    await OrderTable.create(defaultOrder);
    await FillTable.create(defaultFill);
    await FundingIndexUpdatesTable.create(defaultFundingIndexUpdate);

    // First run
    await updateFundingPaymentsTask();
    const firstRunCache = await PersistentCacheTable.findById(
      PersistentCacheKeys.FUNDING_PAYMENTS_LAST_PROCESSED_HEIGHT,
    );
    expect(firstRunCache?.value).toEqual('2');

    // Create new block
    await BlockTable.create({
      blockHeight: '3',
      time: new Date().toISOString(),
    });

    // Second run
    await updateFundingPaymentsTask();
    const secondRunCache = await PersistentCacheTable.findById(
      PersistentCacheKeys.FUNDING_PAYMENTS_LAST_PROCESSED_HEIGHT,
    );
    expect(secondRunCache?.value).toEqual('3');

    // Verify funding payments were created for both runs
    const fundingPayments = await FundingPaymentsTable.findAll(
      {},
      [],
    );
    expect(fundingPayments.length).toEqual(1);
  });
});
