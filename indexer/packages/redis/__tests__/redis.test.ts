import { logger } from '@dydxprotocol-indexer/base';
import { RetryStrategyOptions } from 'redis';

import config from '../src/config';
import {
  createRetryConnectionStrategy,
  deleteAllAsync,
  deleteAsync,
  getAsync,
  hDelAsync,
  hGetAllAsync,
  hGetAsync,
  hincrbyAsync,
  hincrbyFloatAsync,
  hMGetAsync,
  hSetAsync,
  hSetnxAsync,
  incrbyAsync,
  lRangeAsync,
  pttl,
  rPushAsync,
  setAsync,
  setexAsync,
  setExpiry,
  ttl,
  zAddAsync,
  zRangeByScoreAsync,
  zScoreAsync,
  zRemRangeByScoreAsync,
} from '../src/helpers/redis';
import {
  callRetryStrategy,
  redis,
} from './helpers/utils';

const HASH = 'hash';
const KEY = 'foo';
const VAL = 'val';

describe('Redis', () => {
  beforeEach(async () => {
    await deleteAllAsync(redis);
  });

  it('get when null', async () => {
    const get = await getAsync(KEY, redis);
    expect(get).toEqual(null);

    const hGet = await hGetAsync({ hash: HASH, key: KEY }, redis);
    expect(hGet).toEqual(null);

    const hGetAll = await hGetAllAsync(HASH, redis);
    expect(hGetAll).toEqual({});
  });

  it('get when not null/deleted', async () => {
    await setAsync({ key: KEY, value: VAL }, redis);
    await hSetAsync({ hash: HASH, pairs: { [KEY]: VAL } }, redis);

    const get = await getAsync(KEY, redis);
    expect(get).toEqual(VAL);

    const hGet = await hGetAsync({ hash: HASH, key: KEY }, redis);
    expect(hGet).toEqual(VAL);

    const hGetAll = await hGetAllAsync(HASH, redis);
    expect(hGetAll).toEqual({ [KEY]: VAL });

    await deleteAsync(KEY, redis);
    await setexAsync({ key: KEY, value: VAL, timeToLiveSeconds: 10000 }, redis);
    const getex = await getAsync(KEY, redis);
    expect(getex).toEqual(VAL);

    await deleteAllAsync(redis);
    const getFinal = await getAsync(KEY, redis);
    expect(getFinal).toEqual(null);
  });

  it('hsetnx', async () => {
    const result1 = await hSetnxAsync({ hash: HASH, key: KEY, value: VAL }, redis);
    expect(result1).toEqual(1);

    // set when already set
    const result2 = await hSetnxAsync({ hash: HASH, key: KEY, value: VAL }, redis);
    expect(result2).toEqual(0);

    // set when key has already been set
    const result3 = await hSetnxAsync({ hash: HASH, key: KEY, value: 'v2' }, redis);
    expect(result3).toEqual(0);

    const result4 = await hGetAsync({ hash: HASH, key: KEY }, redis);
    expect(result4).toEqual(VAL);
  });

  it('hset, hdel and both multi', async () => {
    // single insert/delete
    await hSetAsync({ hash: HASH, pairs: { [KEY]: VAL } }, redis);
    const result1 = await hGetAsync({ hash: HASH, key: KEY }, redis);
    expect(result1).toEqual(VAL);

    await hDelAsync({ hash: HASH, keys: [KEY] }, redis);
    const result2 = await hGetAsync({ hash: HASH, key: KEY }, redis);
    expect(result2).toBeNull();

    // multi insert/delete
    await hSetAsync({ hash: HASH, pairs: { [KEY]: VAL, key2: 'foo2' } }, redis);
    const result3 = await hGetAsync({ hash: HASH, key: 'key2' }, redis);
    expect(result3).toEqual('foo2');

    await hDelAsync({ hash: HASH, keys: [KEY, 'key2'] }, redis);
    const result4 = await hGetAsync({ hash: HASH, key: KEY }, redis);
    expect(result4).toBeNull();
  });

  it('hincrby', async () => {
    const result = await hincrbyAsync({
      hash: HASH,
      key: KEY,
      changeBy: '10',
    }, redis);
    expect(result).toEqual(10);

    const result2 = await hincrbyAsync({
      hash: HASH,
      key: KEY,
      changeBy: '5',
    }, redis);
    expect(result2).toEqual(15);

    const result3 = await hincrbyAsync({
      hash: HASH,
      key: KEY,
      changeBy: '-30',
    }, redis);
    expect(result3).toEqual(-15);
  });

  it('hincrbyfloat', async () => {
    const result = await hincrbyFloatAsync({
      hash: HASH,
      key: KEY,
      changeBy: '10.00',
    }, redis);
    expect(result).toEqual('10');

    const result2 = await hincrbyFloatAsync({
      hash: HASH,
      key: KEY,
      changeBy: '5.73',
    }, redis);
    expect(result2).toEqual('15.73');

    const result3 = await hincrbyFloatAsync({
      hash: HASH,
      key: KEY,
      changeBy: '-30',
    }, redis);
    expect(result3).toEqual('-14.27');
  });

  it('incrby', async () => {
    const result = await incrbyAsync({
      key: KEY,
      changeBy: '10',
    }, redis);
    expect(result).toEqual(10);

    const result2 = await incrbyAsync({
      key: KEY,
      changeBy: '5',
    }, redis);
    expect(result2).toEqual(15);

    const result3 = await incrbyAsync({
      key: KEY,
      changeBy: '-30',
    }, redis);
    expect(result3).toEqual(-15);
  });

  it('rPushAsync/lRangeAsync', async () => {
    const result = await rPushAsync({
      key: KEY,
      value: 'a',
    }, redis);
    expect(result).toEqual(1);

    const result2 = await rPushAsync({
      key: KEY,
      value: 'b',
    }, redis);
    expect(result2).toEqual(2);

    const result3 = await lRangeAsync(KEY, redis);
    expect(result3).toEqual(['a', 'b']);

    const result4 = await deleteAsync(KEY, redis);
    expect(result4).toEqual(1);

    const result5 = await lRangeAsync(KEY, redis);
    expect(result5).toEqual([]);
  });

  it('hmget', async () => {
    await hSetAsync({
      hash: KEY,
      pairs: {
        a: 'A',
        b: 'B',
      },
    }, redis);

    const result = await hMGetAsync({
      hash: KEY,
      fields: ['a', 'b', 'c'],
    }, redis);
    expect(result).toEqual(['A', 'B', null]);
  });

  it('ttl, pttl, setExpiry', async () => {
    await setAsync({
      key: 'a',
      value: 'A',
    }, redis);

    await setExpiry(redis, 'a', 10);
    const secLeft: number = await ttl(redis, 'a');
    expect(secLeft).toBeGreaterThan(9);

    const msLeft: number = await pttl(redis, 'a');
    expect(msLeft).toBeGreaterThan(9000);
  });

  describe('reconnection strategy', () => {
    const timeout: number = 10;
    const retryStrategy: (options: RetryStrategyOptions) => number = createRetryConnectionStrategy(
      'url', timeout,
    );
    const errorObj: NodeJS.ErrnoException = new Error();
    let loggerErrorSpy: jest.SpyInstance;
    let loggerInfoSpy: jest.SpyInstance;

    beforeEach(() => {
      loggerErrorSpy = jest.spyOn(logger, 'error');
      loggerInfoSpy = jest.spyOn(logger, 'info');
    });

    afterAll(() => {
      jest.resetAllMocks();
    });

    it('returns timeout', () => {
      expect(retryStrategy({
        total_retry_time: 1, times_connected: 1, error: errorObj, attempt: 1,
      }))
        .toEqual(timeout);
    });

    it('logs info if error is null and not more than attempts threshold', () => {
      callRetryStrategy(retryStrategy, null);
      expect(loggerInfoSpy).toHaveBeenCalledTimes(1);
      expect(loggerErrorSpy).not.toHaveBeenCalled();
    });

    it('logs info if error is undefined and not more than attempts threshold', () => {
      callRetryStrategy(retryStrategy, undefined);
      expect(loggerInfoSpy).toHaveBeenCalledTimes(1);
      expect(loggerErrorSpy).not.toHaveBeenCalled();
    });

    it('logs error if error is not null/undefined and not more than attempts threshold', () => {
      callRetryStrategy(retryStrategy, errorObj);
      expect(loggerInfoSpy).not.toHaveBeenCalled();
      expect(loggerErrorSpy).toHaveBeenCalledTimes(1);
    });

    it('logs error if error is null and more than attempts threshold', () => {
      callRetryStrategy(retryStrategy, null, config.REDIS_RECONNECT_ATTEMPT_ERROR_THRESHOLD + 1);
      expect(loggerInfoSpy).not.toHaveBeenCalled();
      expect(loggerErrorSpy).toHaveBeenCalledTimes(1);
    });

    it('logs error if error is undefined and more than attempts threshold', () => {
      callRetryStrategy(
        retryStrategy, undefined, config.REDIS_RECONNECT_ATTEMPT_ERROR_THRESHOLD + 1,
      );
      expect(loggerInfoSpy).not.toHaveBeenCalled();
      expect(loggerErrorSpy).toHaveBeenCalledTimes(1);
    });
  });

  describe('zAddAsync ->', () => {
    beforeEach(async () => {
      await zAddAsync({ key: HASH, value: 'val-10', score: -10 }, redis);
      await zAddAsync({ key: HASH, value: 'val-1', score: -1 }, redis);
      await zAddAsync({ key: HASH, value: 'val0', score: 0 }, redis);
      await zAddAsync({ key: HASH, value: 'val1', score: 1 }, redis);
      await zAddAsync({ key: HASH, value: 'val10', score: 10 }, redis);
    });

    afterEach(async () => {
      await deleteAsync(HASH, redis);
    });

    it('zRemRangeByScoreAsync', async () => {
      const expectedWithScores: string[] = [
        'val-10', '-10',
        'val-1', '-1',
        'val0', '0',
        'val1', '1',
        'val10', '10',
      ];
      const actual: string[] = await zRangeByScoreAsync({
        key: HASH,
        start: -Infinity,
        startIsInclusive: false,
        end: Infinity,
        endIsInclusive: false,
        withScores: true,
      }, redis);
      expect(actual).toEqual(expectedWithScores);
      const deleted: number = await zRemRangeByScoreAsync({
        key: HASH,
        start: -Infinity,
        startIsInclusive: false,
        end: 0,
        endIsInclusive: true,
      }, redis);
      expect(deleted).toEqual(3);
      const afterRemoval: string[] = await zRangeByScoreAsync({
        key: HASH,
        start: -Infinity,
        startIsInclusive: false,
        end: Infinity,
        endIsInclusive: false,
        withScores: true,
      }, redis);
      expect(afterRemoval).toEqual(['val1', '1', 'val10', '10']);
    });

    describe('zScore', () => {
      it.each([
        ['finds score', 'val10', '10'],
        ['finds 0 score', 'val0', '0'],
        ['not found returns null', 'val', null],
      ])('%s', async (_name: string, member: string, expected: string|null) => {
        const actual: string|null = await zScoreAsync({ hash: HASH, key: member }, redis);
        expect(actual).toStrictEqual(expected);
      });
    });

    describe('zRangeByScoreAsync', () => {
      const expectedWithoutScores: string[] = ['val-10', 'val-1', 'val0', 'val1', 'val10'];
      const expectedWithScores: string[] = [
        'val-10', '-10',
        'val-1', '-1',
        'val0', '0',
        'val1', '1',
        'val10', '10',
      ];

      it.each([
        // range: -infinity, infinity
        ['all values w/ scores', -Infinity, false, Infinity, false, true, expectedWithScores],
        ['all values w/o scores', -Infinity, false, Infinity, false, false, expectedWithoutScores],
        // range: bounded, infinity
        ['with inclusive start set', -1, true, Infinity, false, true, expectedWithScores.slice(2)],
        [
          'with inclusive start set w/o scores',
          -1,
          true,
          Infinity,
          false,
          false,
          expectedWithoutScores.slice(1),
        ],
        ['with exclusive start set', -1, false, Infinity, false, true, expectedWithScores.slice(4)],
        [
          'with exclusive start set w/o scores',
          -1,
          false,
          Infinity,
          false,
          false,
          expectedWithoutScores.slice(2),
        ],
        // range: -infinity, bounded
        ['with inclusive end set', -Infinity, false, 1, true, true, expectedWithScores.slice(0, -2)],
        [
          'with inclusive end set w/o scores',
          -Infinity,
          false,
          1,
          true,
          false,
          expectedWithoutScores.slice(0, -1),
        ],
        ['with exclusive end set', -Infinity, false, 1, false, true, expectedWithScores.slice(0, -4)],
        [
          'with exclusive end set w/o scores',
          -Infinity,
          false,
          1,
          false,
          false,
          expectedWithoutScores.slice(0, -2),
        ],
        // range: bounded, bounded
        [
          'with inclusive start, inclusive end',
          -1,
          true,
          1,
          true,
          true,
          expectedWithScores.slice(2, -2),
        ],
        [
          'with inclusive start, exclusive end',
          -1,
          true,
          1,
          false,
          true,
          expectedWithScores.slice(2, -4),
        ],
        [
          'with exclusive start, inclusive end',
          -1,
          false,
          1,
          true,
          true,
          expectedWithScores.slice(4, -2),
        ],
        [
          'with exclusive start, exclusive end',
          -1,
          false,
          1,
          false,
          true,
          expectedWithScores.slice(4, -4),
        ],
        [
          'with inclusive start, inclusive end w/o scores',
          -1,
          true,
          1,
          true,
          false,
          expectedWithoutScores.slice(1, -1),
        ],
        [
          'with inclusive start, exclusive end w/o scores',
          -1,
          true,
          1,
          false,
          false,
          expectedWithoutScores.slice(1, -2),
        ],
        [
          'with exclusive start, inclusive end w/o scores',
          -1,
          false,
          1,
          true,
          false,
          expectedWithoutScores.slice(2, -1),
        ],
        [
          'with exclusive start, exclusive end w/o scores',
          -1,
          false,
          1,
          false,
          false,
          expectedWithoutScores.slice(2, -2),
        ],
      ])('can retrieve %s', async (
        _name: string,
        start: number,
        startIsInclusive: boolean,
        end: number,
        endIsInclusive: boolean,
        withScores: boolean,
        expected: string[],
      ) => {
        const actual: string[] = await zRangeByScoreAsync({
          key: HASH,
          start,
          startIsInclusive,
          end,
          endIsInclusive,
          withScores,
        }, redis);
        expect(actual).toEqual(expected);
      });
    });
  });
});
