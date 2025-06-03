import {
  FundingPaymentsCreateObject,
  FundingPaymentsFromDatabase,
} from '../../src/types';
import * as FundingPaymentsTable from '../../src/stores/funding-payments-table';
import { clearData, migrate, teardown } from '../../src/helpers/db-helpers';
import { seedData } from '../helpers/mock-generators';
import { defaultFundingPayment, defaultFundingPayment2 } from '../helpers/constants';

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
    expect(fundingPayments[0]).toEqual(expect.objectContaining(defaultFundingPayment2));
    expect(fundingPayments[1]).toEqual(expect.objectContaining(defaultFundingPayment));
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
});
