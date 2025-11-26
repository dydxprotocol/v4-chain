import {
  FundingPaymentsCreateObject,
  FundingPaymentsFromDatabase,
} from '../../src/types';
import * as FundingPaymentsTable from '../../src/stores/funding-payments-table';
import { clearData, migrate, teardown } from '../../src/helpers/db-helpers';
import { seedData } from '../helpers/mock-generators';
import {
  defaultSubaccountId,
  defaultSubaccountId2,
  defaultSubaccountId3,
  defaultFundingPayment,
  defaultFundingPayment2,
} from '../helpers/constants';

describe('funding payments store', () => {
  const updatedHeight: string = '5';

  beforeEach(async () => {
    await clearData();
    await seedData();
  });

  beforeAll(async () => {
    await migrate();
  });

  afterEach(async () => {
    await clearData();
  });

  afterAll(async () => {
    await teardown();
  });

  it('Successfully creates a funding payment', async () => {
    await FundingPaymentsTable.create(defaultFundingPayment);
  });

  it('Successfully finds all FundingPayments', async () => {
    await Promise.all([
      FundingPaymentsTable.create(defaultFundingPayment),
      FundingPaymentsTable.create(defaultFundingPayment2),
    ]);

    const { results: fundingPayments } = await FundingPaymentsTable.findAll({}, [], {});

    expect(fundingPayments.length).toEqual(2);
    expect(fundingPayments[0]).toEqual(expect.objectContaining(defaultFundingPayment));
    expect(fundingPayments[1]).toEqual(expect.objectContaining(defaultFundingPayment2));
  });

  it('Successfully finds FundingPayments with createdAtHeight', async () => {
    await FundingPaymentsTable.create(defaultFundingPayment);

    const { results: fundingPayments } = await FundingPaymentsTable.findAll(
      {
        createdAtHeight: defaultFundingPayment.createdAtHeight,
      },
      [],
      { readReplica: true },
    );

    expect(fundingPayments.length).toEqual(1);
    expect(fundingPayments[0]).toEqual(
      expect.objectContaining({
        ...defaultFundingPayment,
      }),
    );
  });

  it('Successfully finds all FundingPayments created before or at height', async () => {
    await Promise.all([
      FundingPaymentsTable.create(defaultFundingPayment),
      FundingPaymentsTable.create({
        ...defaultFundingPayment2,
        createdAtHeight: updatedHeight,
      }),
    ]);

    const { results: fundingPayments } = await FundingPaymentsTable.findAll(
      {
        createdBeforeOrAtHeight: defaultFundingPayment.createdAtHeight,
      },
      [],
      {},
    );

    expect(fundingPayments.length).toEqual(1);
    expect(fundingPayments[0]).toEqual(expect.objectContaining(defaultFundingPayment));
  });

  it('Successfully finds all FundingPayments created before or at time', async () => {
    const fundingPayment2: FundingPaymentsCreateObject = {
      ...defaultFundingPayment,
      createdAtHeight: updatedHeight,
      createdAt: '1982-05-25T00:00:00.000Z',
    };
    await Promise.all([
      FundingPaymentsTable.create(defaultFundingPayment),
      FundingPaymentsTable.create(fundingPayment2),
    ]);

    const { results: fundingPayments } = await FundingPaymentsTable.findAll(
      {
        createdBeforeOrAt: '2000-05-25T00:00:00.000Z',
      },
      [],
      {},
    );

    expect(fundingPayments.length).toEqual(1);
    expect(fundingPayments[0]).toEqual(expect.objectContaining(fundingPayment2));
  });

  it('Successfully finds a FundingPayment by id', async () => {
    await FundingPaymentsTable.create(defaultFundingPayment);

    const fPayment: FundingPaymentsFromDatabase | undefined = await FundingPaymentsTable.findById(
      defaultFundingPayment.subaccountId,
      defaultFundingPayment.createdAt,
      defaultFundingPayment.ticker,
    );
    expect(fPayment).toEqual(expect.objectContaining(defaultFundingPayment));
  });

  it('supports pagination', async () => {
    await Promise.all([
      FundingPaymentsTable.create(defaultFundingPayment),
      FundingPaymentsTable.create(defaultFundingPayment2),
    ]);

    const { results: fundingPayments } = await FundingPaymentsTable.findAll(
      {
        page: 1,
        limit: 1,
      },
      [],
    );
    expect(fundingPayments.length).toEqual(1);
    expect(fundingPayments[0]).toEqual(expect.objectContaining(defaultFundingPayment));
    const { results: fundingPayments2 } = await FundingPaymentsTable.findAll(
      {
        page: 2,
        limit: 1,
      },
      [],
    );
    expect(fundingPayments2.length).toEqual(1);
    expect(fundingPayments2[0]).toEqual(expect.objectContaining(defaultFundingPayment2));
  });

  it('returns correct net funding payments for a subaccount between block heights', async () => {
    const fundingPayment1: FundingPaymentsCreateObject = {
      ...defaultFundingPayment,
      createdAt: '2023-01-01T00:00:00.000Z',
      createdAtHeight: '10',
      payment: '8.5',
    };

    const fundingPayment2: FundingPaymentsCreateObject = {
      ...defaultFundingPayment,
      createdAt: '2023-01-02T00:00:00.000Z',
      createdAtHeight: '20',
      payment: '-3.2',
    };

    await Promise.all([
      FundingPaymentsTable.create(fundingPayment1),
      FundingPaymentsTable.create(fundingPayment2),
    ]);

    const netPayments = await
    FundingPaymentsTable.getNetFundingPaymentsBetweenBlockHeightsForSubaccount(
      defaultSubaccountId,
      '10',
      '30',
    );

    // fundingPayment1 is at height 10 which is excluded
    // fundingPayment2 has payment: '-3.2'
    expect(netPayments.toString()).toEqual('-3.2');
  });

  it('correctly sums positive and negative payments in the specified range', async () => {
    const fundingPayment1: FundingPaymentsCreateObject = {
      ...defaultFundingPayment,
      createdAt: '2023-02-01T00:00:00.000Z',
      createdAtHeight: '15',
      payment: '-2.5', // Negative payment
    };

    const fundingPayment2: FundingPaymentsCreateObject = {
      ...defaultFundingPayment,
      createdAt: '2023-02-02T00:00:00.000Z',
      createdAtHeight: '25',
      payment: '7.8', // Positive payment
    };

    await Promise.all([
      FundingPaymentsTable.create(fundingPayment1),
      FundingPaymentsTable.create(fundingPayment2),
    ]);

    const netPayments = await
    FundingPaymentsTable.getNetFundingPaymentsBetweenBlockHeightsForSubaccount(
      defaultSubaccountId,
      '10',
      '30',
    );

    // Sum of payments: -2.5 + 7.8 = 5.3
    expect(netPayments.toString()).toEqual('5.3');
  });

  it('only includes funding payments for the specified subaccount with mixed payment signs', async () => {
    // First subaccount payments
    const fundingPayment1: FundingPaymentsCreateObject = {
      ...defaultFundingPayment,
      createdAt: '2023-03-01T00:00:00.000Z',
      createdAtHeight: '15',
      payment: '-4.2',
    };

    const fundingPayment1a: FundingPaymentsCreateObject = {
      ...defaultFundingPayment,
      createdAt: '2023-03-01T12:00:00.000Z',
      createdAtHeight: '20',
      payment: '2.5',
    };

    const fundingPayment1b: FundingPaymentsCreateObject = {
      ...defaultFundingPayment,
      createdAt: '2023-03-02T00:00:00.000Z',
      createdAtHeight: '25',
      payment: '-1.8',
    };

    // Second subaccount payments
    const fundingPayment2: FundingPaymentsCreateObject = {
      ...defaultFundingPayment2, // Uses defaultSubaccountId2
      createdAt: '2023-03-02T00:00:00.000Z',
      createdAtHeight: '25',
      payment: '6.7',
    };

    const fundingPayment2a: FundingPaymentsCreateObject = {
      ...defaultFundingPayment2,
      createdAt: '2023-03-02T06:00:00.000Z',
      createdAtHeight: '27',
      payment: '-2.3',
    };

    const fundingPayment2b: FundingPaymentsCreateObject = {
      ...defaultFundingPayment2,
      createdAt: '2023-03-02T12:00:00.000Z',
      createdAtHeight: '29',
      payment: '1.5',
    };

    // Payment outside the block height range to ensure it's excluded
    const fundingPaymentOutsideRange: FundingPaymentsCreateObject = {
      ...defaultFundingPayment,
      createdAt: '2023-03-03T00:00:00.000Z',
      createdAtHeight: '35',
      payment: '10.0', // Should be excluded based on block height
    };

    await Promise.all([
      FundingPaymentsTable.create(fundingPayment1),
      FundingPaymentsTable.create(fundingPayment1a),
      FundingPaymentsTable.create(fundingPayment1b),
      FundingPaymentsTable.create(fundingPayment2),
      FundingPaymentsTable.create(fundingPayment2a),
      FundingPaymentsTable.create(fundingPayment2b),
      FundingPaymentsTable.create(fundingPaymentOutsideRange),
    ]);

    // Test with defaultSubaccountId
    const netPayments1 = await
    FundingPaymentsTable.getNetFundingPaymentsBetweenBlockHeightsForSubaccount(
      defaultSubaccountId,
      '10',
      '30',
    );

    // Sum of all payments for defaultSubaccountId in the range: -4.2 + 2.5 + (-1.8) = -3.5
    expect(netPayments1.toString()).toEqual('-3.5');

    // Test with defaultSubaccountId2
    const netPayments2 = await
    FundingPaymentsTable.getNetFundingPaymentsBetweenBlockHeightsForSubaccount(
      defaultSubaccountId2,
      '10',
      '30',
    );

    // Sum of all payments for defaultSubaccountId2 in the range: 6.7 + (-2.3) + 1.5 = 5.9
    expect(netPayments2.toString()).toEqual('5.9');
  });

  it('returns 0 for subaccounts with no funding payments in the specified range', async () => {
    // Create a payment outside the block height range we'll query
    const fundingPaymentOutsideRange: FundingPaymentsCreateObject = {
      ...defaultFundingPayment,
      createdAt: '2023-03-01T00:00:00.000Z',
      createdAtHeight: '50',
      payment: '5.0',
    };

    // Create a payment for a different subaccount
    const fundingPaymentDifferentSubaccount: FundingPaymentsCreateObject = {
      ...defaultFundingPayment,
      subaccountId: defaultSubaccountId2,
      createdAt: '2023-03-01T00:00:00.000Z',
      createdAtHeight: '25',
      payment: '3.0',
    };

    await Promise.all([
      FundingPaymentsTable.create(fundingPaymentOutsideRange),
      FundingPaymentsTable.create(fundingPaymentDifferentSubaccount),
    ]);

    // Test case 1: No payments in the specified block height range
    const netPaymentsOutsideRange = await
    FundingPaymentsTable.getNetFundingPaymentsBetweenBlockHeightsForSubaccount(
      defaultSubaccountId,
      '10',
      '40',
    );

    // Should return 0 when no payments are found in the range
    expect(netPaymentsOutsideRange.toString()).toEqual('0');

    // Test case 2: No payments for this subaccount
    const netPaymentsNonExistentSubaccount = await
    FundingPaymentsTable.getNetFundingPaymentsBetweenBlockHeightsForSubaccount(
      defaultSubaccountId3,
      '10',
      '60',
    );

    // Should return 0 for a subaccount with no payments
    expect(netPaymentsNonExistentSubaccount.toString()).toEqual('0');
  });
});
