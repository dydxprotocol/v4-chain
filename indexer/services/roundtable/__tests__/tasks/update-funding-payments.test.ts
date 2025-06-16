import {
  dbHelpers,
  testMocks,
  PersistentCacheTable,
  PersistentCacheKeys,
  BlockTable,
  FillTable,
  OrderTable,
  FundingPaymentsTable,
  OrderSide,
  FundingIndexUpdatesTable,
  PositionSide,
  FillType,
  Liquidity,
  Ordering,
} from '@dydxprotocol-indexer/postgres';
import updateFundingPaymentsTask from '../../src/tasks/update-funding-payments';
import {
  createdDateTime,
  defaultFill,
  defaultFundingIndexUpdate,
  defaultOrder,
  defaultOrderId,
  defaultPerpetualMarket,
  defaultSubaccountId,
  defaultTendermintEventId,
} from '@dydxprotocol-indexer/postgres/build/__tests__/helpers/constants';

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
      fundingIndex: '5',
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
      fundingIndex: '4',
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
    const { results: fundingPayments } = await FundingPaymentsTable.findAll({}, []);
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
    // Create fills without any existing funding payments
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
    const { results: fundingPayments } = await FundingPaymentsTable.findAll({}, []);
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

  it('Maintains persistent cache value across multiple runs and does not create funding payments if no unprocessed funding index update', async () => {
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

    // Second run doesn't create funding payments because no new funding index updates.
    await updateFundingPaymentsTask();
    const secondRunCache = await PersistentCacheTable.findById(
      PersistentCacheKeys.FUNDING_PAYMENTS_LAST_PROCESSED_HEIGHT,
    );
    expect(secondRunCache?.value).toEqual('2');

    // Verify funding payments were created for both runs
    const { results: fundingPayments } = await FundingPaymentsTable.findAll({}, []);
    expect(fundingPayments.length).toEqual(1);
  });

  it('Creates funding payment from fills and snapshot', async () => {
    await OrderTable.create(defaultOrder);
    await FillTable.create(defaultFill);
    await FundingIndexUpdatesTable.create(defaultFundingIndexUpdate);

    // Run task
    await updateFundingPaymentsTask();

    // Verify funding payments
    const { results: fundingPayments } = await FundingPaymentsTable.findAll({}, []);
    expect(fundingPayments.length).toEqual(1);
    expect(fundingPayments[0]).toMatchObject({
      subaccountId: defaultSubaccountId,
      perpetualId: defaultFundingIndexUpdate.perpetualId,
    });

    // Verify persistent cache
    const persistentCache = await PersistentCacheTable.findById(
      PersistentCacheKeys.FUNDING_PAYMENTS_LAST_PROCESSED_HEIGHT,
    );
    expect(persistentCache?.value).toEqual('2');
  });

  it('Backfills all funding payments from genesis', async () => {
    // First funding payment
    await OrderTable.create(defaultOrder);
    await FillTable.create(defaultFill);
    await FundingIndexUpdatesTable.create(defaultFundingIndexUpdate);

    // second funding payment directly as a snapshot.
    await BlockTable.create({
      blockHeight: '3',
      time: new Date().toISOString(),
    });
    await FundingIndexUpdatesTable.create({
      perpetualId: defaultPerpetualMarket.id,
      eventId: defaultTendermintEventId,
      rate: '0.0004',
      oraclePrice: '10000',
      fundingIndex: '10054',
      effectiveAt: new Date().toISOString(),
      effectiveAtHeight: '3',
    });

    // third funding payment created from fills + snapshot with fills as same height as
    // third index update which means it will be included.
    await BlockTable.create({
      blockHeight: '4',
      time: new Date().toISOString(),
    });
    await FillTable.create({
      subaccountId: defaultSubaccountId,
      side: OrderSide.BUY,
      liquidity: Liquidity.MAKER,
      type: FillType.LIMIT,
      clobPairId: '1',
      orderId: defaultOrderId,
      size: '10',
      price: '20000',
      quoteAmount: '200000',
      eventId: defaultTendermintEventId,
      transactionHash: '', // TODO: Add a real transaction Hash
      createdAt: new Date().toISOString(),
      createdAtHeight: '4',
      clientMetadata: '0',
      fee: '1.1',
      affiliateRevShare: '1.10',
    });
    await FundingIndexUpdatesTable.create({
      perpetualId: defaultPerpetualMarket.id,
      eventId: defaultTendermintEventId,
      rate: '0.0004',
      oraclePrice: '10000',
      fundingIndex: '10058',
      effectiveAt: new Date().toISOString(),
      effectiveAtHeight: '4',
    });

    // Run task
    await updateFundingPaymentsTask();

    // assert funding payments
    const { results: fundingPayments } = await FundingPaymentsTable.findAll(
      {},
      [],
      {
        orderBy: [['createdAtHeight', Ordering.ASC]],
      },
    );
    expect(fundingPayments.length).toEqual(3);
    expect(fundingPayments).toEqual([
      expect.objectContaining({
        subaccountId: defaultSubaccountId,
        perpetualId: defaultPerpetualMarket.id,
        size: '10',
        side: PositionSide.LONG,
        payment: '-100500',
        createdAtHeight: '2',
        oraclePrice: '10000',
        rate: '0.0004',
        ticker: 'BTC-USD',
      }),
      expect.objectContaining({
        subaccountId: defaultSubaccountId,
        perpetualId: defaultPerpetualMarket.id,
        size: '10',
        side: PositionSide.LONG,
        payment: '-40',
        createdAtHeight: '3',
        oraclePrice: '10000',
        rate: '0.0004',
        ticker: 'BTC-USD',
      }),
      expect.objectContaining({
        subaccountId: defaultSubaccountId,
        perpetualId: defaultPerpetualMarket.id,
        size: '20',
        side: PositionSide.LONG,
        payment: '-80',
        createdAtHeight: '4',
        oraclePrice: '10000',
        rate: '0.0004',
        ticker: 'BTC-USD',
      }),
    ]);
  });
});
