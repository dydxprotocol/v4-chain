/* ------- UTILITY TYPES ------- */
import { RawBinding } from 'knex';

export type IsoString = string;

export type RegexPattern = string;

export enum IsolationLevel {
  SERIALIZABLE = 'SERIALIZABLE',
  REPEATABLE_READ = 'REPEATABLE READ',
  READ_COMMITTED = 'READ COMMITTED',
  READ_UNCOMMITTED = 'READ UNCOMMITTED',
}

export interface Options {
  txId?: number,
  forUpdate?: boolean,
  noWait?: boolean,
  orderBy?: [string, Ordering][],
  readReplica?: boolean,
  random?: boolean,
  bindings?: readonly RawBinding[],
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  sqlOptions?: Readonly<{ [key: string]: any }>,
}

export enum Ordering {
  ASC = 'ASC',
  DESC = 'DESC',
}
