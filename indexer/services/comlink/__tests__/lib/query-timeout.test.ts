import express from 'express';

import { createQueryTimeoutResponse, isQueryCanceledError } from '../../src/lib/query-timeout';

describe('query-timeout', () => {
  describe('isQueryCanceledError', () => {
    it('detects the driver error from a raw query', () => {
      const error: Error & { code?: string } = new Error(
        'canceling statement due to statement timeout',
      );
      error.code = '57014';
      expect(isQueryCanceledError(error)).toBe(true);
    });

    it('detects the driver error wrapped by objection', () => {
      const error: Error & { nativeError?: { code: string } } = new Error('DBError');
      error.nativeError = { code: '57014' };
      expect(isQueryCanceledError(error)).toBe(true);
    });

    it.each([
      ['another SQLSTATE', Object.assign(new Error('unique violation'), { code: '23505' })],
      ['a plain error', new Error('boom')],
      ['undefined', undefined],
      ['null', null],
      ['a string', '57014'],
    ])('ignores %s', (_name: string, error: unknown) => {
      expect(isQueryCanceledError(error)).toBe(false);
    });
  });

  describe('createQueryTimeoutResponse', () => {
    it('answers 504 with an error body', () => {
      const json: jest.Mock = jest.fn().mockReturnThis();
      const status: jest.Mock = jest.fn().mockReturnValue({ json });
      const res = { status } as unknown as express.Response;

      createQueryTimeoutResponse(res);

      expect(status).toHaveBeenCalledWith(504);
      expect(json).toHaveBeenCalledWith({ errors: [{ msg: 'Request timed out' }] });
    });
  });
});
