import Knex from 'knex';

import config from '../../src/config';
import { getConfigForHost, getSessionSettingsSql } from '../../src/db/knexfile';

function runAfterCreate(knexConfig: Knex.Config, queryError: Error | null = null): {
  sql: string,
  done: jest.Mock,
  conn: object,
} {
  const query: jest.Mock = jest.fn((_sql: string, cb: Function) => cb(queryError));
  const conn = { query };
  const done: jest.Mock = jest.fn();
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  (knexConfig.pool as any).afterCreate(conn, done);
  return { sql: query.mock.calls[0][0], done, conn };
}

describe('knexfile', () => {
  const idleInTx: string = `SET idle_in_transaction_session_timeout = ${config.PG_IDLE_IN_TX_TIMEOUT_MS}`;

  describe('getSessionSettingsSql', () => {
    it('sets only the idle-in-transaction timeout when the statement timeout is disabled', () => {
      expect(getSessionSettingsSql(0)).toEqual(idleInTx);
    });

    it('adds statement_timeout when it is positive', () => {
      expect(getSessionSettingsSql(55_000)).toEqual(`${idleInTx}; SET statement_timeout = 55000`);
    });
  });

  describe('getConfigForHost', () => {
    it.each([
      [0, idleInTx],
      [55_000, `${idleInTx}; SET statement_timeout = 55000`],
    ])('runs the session settings on each new connection (statement timeout %d)', (
      statementTimeoutMs: number,
      expectedSql: string,
    ) => {
      const { sql, done, conn } = runAfterCreate(getConfigForHost('localhost', statementTimeoutMs));

      expect(sql).toEqual(expectedSql);
      expect(done).toHaveBeenCalledWith(null, conn);
    });

    it('defaults to no statement timeout', () => {
      expect(runAfterCreate(getConfigForHost('localhost')).sql).toEqual(idleInTx);
    });

    it('passes a failed SET to the pool so the connection is not used', () => {
      const error: Error = new Error('invalid value for parameter "statement_timeout"');
      const { done, conn } = runAfterCreate(getConfigForHost('localhost', 55_000), error);

      expect(done).toHaveBeenCalledWith(error, conn);
    });
  });

  describe('env wiring', () => {
    const envKey: string = 'PG_READ_REPLICA_STATEMENT_TIMEOUT_MS';
    const original: string | undefined = process.env[envKey];

    afterEach(() => {
      if (original === undefined) {
        delete process.env[envKey];
      } else {
        process.env[envKey] = original;
      }
    });

    it('applies the statement timeout to the read replica only', () => {
      process.env[envKey] = '1234';
      jest.isolateModules(() => {
        // eslint-disable-next-line global-require, @typescript-eslint/no-var-requires
        const knexfile = require('../../src/db/knexfile');

        expect(runAfterCreate(knexfile.knexReadReplicaConfigForEnv).sql)
          .toEqual(`${idleInTx}; SET statement_timeout = 1234`);
        expect(runAfterCreate(knexfile.knexPrimaryConfigForEnv).sql).toEqual(idleInTx);
      });
    });
  });
});
