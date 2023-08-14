import { TransactionFromDatabase } from '../../src/types';
import * as TransactionTable from '../../src/stores/transaction-table';
import {
  clearData,
  migrate,
  teardown,
} from '../../src/helpers/db-helpers';
import { defaultTransaction, defaultTransactionId } from '../helpers/constants';

describe('Transaction store', () => {
  beforeAll(async () => {
    await migrate();
  });

  afterEach(async () => {
    await clearData();
  });

  afterAll(async () => {
    await teardown();
  });

  it('Successfully creates a Transaction', async () => {
    await TransactionTable.create(defaultTransaction);
  });

  it('Successfully finds all Transactions', async () => {
    await Promise.all([
      TransactionTable.create(defaultTransaction),
      TransactionTable.create({
        ...defaultTransaction,
        transactionIndex: 1,
      }),
    ]);

    const transactions: TransactionFromDatabase[] = await TransactionTable.findAll(
      {},
      [],
      { readReplica: true },
    );

    expect(transactions.length).toEqual(2);
    expect(transactions[0]).toEqual(expect.objectContaining(defaultTransaction));
    expect(transactions[1]).toEqual(expect.objectContaining({
      ...defaultTransaction,
      transactionIndex: 1,
    }));
  });

  it.each([
    { blockHeight: ['1'] },
    { transactionIndex: [0] },
    { transactionHash: ['3ac776f8-1900-43de-ac38-7fb516f7d6d0'] },
  ])('Successfully finds Transaction', async (searchQuery) => {
    await Promise.all([
      TransactionTable.create(defaultTransaction),
      TransactionTable.create({
        blockHeight: '2',
        transactionIndex: 1,
        transactionHash: '0a63ca75-526c-42bb-8b79-1e71a39d1fa9',
      }),
    ]);

    const transactions: TransactionFromDatabase[] = await TransactionTable.findAll(
      searchQuery,
      [],
      { readReplica: true },
    );

    expect(transactions.length).toEqual(1);
    expect(transactions[0]).toEqual(expect.objectContaining(defaultTransaction));
  });

  it('Successfully finds a Transaction', async () => {
    await TransactionTable.create(defaultTransaction);

    const transaction: TransactionFromDatabase | undefined = await TransactionTable.findById(
      defaultTransactionId,
    );

    expect(transaction).toEqual(expect.objectContaining(defaultTransaction));
  });

  it('Unable finds a Transaction', async () => {
    const transaction: TransactionFromDatabase | undefined = await TransactionTable.findById(
      defaultTransactionId,
    );
    expect(transaction).toEqual(undefined);
  });
});
