import config from '../../src/config';
import { getConfigForHost, getSessionSettingsSql } from '../../src/db/knexfile';

describe('knexfile', () => {
  describe('getSessionSettingsSql', () => {
    it('sets only the idle-in-transaction timeout when the statement timeout is disabled', () => {
      expect(getSessionSettingsSql(0)).toEqual(
        `SET idle_in_transaction_session_timeout = ${config.PG_IDLE_IN_TX_TIMEOUT_MS}`,
      );
    });

    it('adds statement_timeout when it is positive', () => {
      expect(getSessionSettingsSql(55_000)).toEqual(
        `SET idle_in_transaction_session_timeout = ${config.PG_IDLE_IN_TX_TIMEOUT_MS}; ` +
        'SET statement_timeout = 55000',
      );
    });
  });

  describe('getConfigForHost', () => {
    it.each([
      [0, getSessionSettingsSql(0)],
      [55_000, getSessionSettingsSql(55_000)],
    ])('runs the session settings on each new connection (statement timeout %d)', (
      statementTimeoutMs: number,
      expectedSql: string,
    ) => {
      const query: jest.Mock = jest.fn((_sql: string, cb: Function) => cb(null));
      const conn = { query };
      const done: jest.Mock = jest.fn();

      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      const pool: any = getConfigForHost('localhost', statementTimeoutMs).pool;
      pool.afterCreate(conn, done);

      expect(query).toHaveBeenCalledWith(expectedSql, expect.any(Function));
      expect(done).toHaveBeenCalledWith(null, conn);
    });

    it('defaults to no statement timeout', () => {
      const query: jest.Mock = jest.fn((_sql: string, cb: Function) => cb(null));
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      const pool: any = getConfigForHost('localhost').pool;
      pool.afterCreate({ query }, jest.fn());

      expect(query.mock.calls[0][0]).not.toContain('statement_timeout =');
    });
  });
});
