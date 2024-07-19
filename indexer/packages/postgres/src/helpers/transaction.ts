import { Model, Transaction as ObjectionTx } from 'objection';

import { IsolationLevel } from '../types';

export default class Transaction {
  static transactions: { [id: number]: ObjectionTx } = {};
  static nextTxId: number = 1; // Start at non-falsey value.

  public static async start(): Promise<number> {
    const id = Transaction.nextTxId;
    Transaction.nextTxId += 1;
    Transaction.transactions[id] = await Model.startTransaction();
    return id;
  }

  public static async setIsolationLevel(id: number, level: IsolationLevel): Promise<void> {
    await Transaction.transactions[id].raw(`SET TRANSACTION ISOLATION LEVEL ${level}`);
  }

  public static async commit(id: number): Promise<void> {
    if (!Transaction.transactions[id]) {
      throw new Error('No transaction to commit');
    }
    await Transaction.transactions[id].commit();
    delete Transaction.transactions[id];
  }

  public static async rollback(id: number): Promise<void> {
    if (!Transaction.transactions[id]) {
      throw new Error('No transaction to rollback');
    }
    await Transaction.transactions[id].rollback();
    delete Transaction.transactions[id];
  }

  public static get(txId: number | undefined): ObjectionTx | undefined {
    if (typeof txId === 'number') {
      if (!Transaction.transactions[txId]) {
        throw new Error('No transaction found');
      }
      return Transaction.transactions[txId];
    }
    return undefined;
  }
}
