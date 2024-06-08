/**
 * Utilities for writing unit tests with jest.
 */

/* eslint-disable @typescript-eslint/no-explicit-any */
type Fn = (...args: any[]) => any;
type FnMock<F extends Fn> = jest.Mock<ReturnType<F>, Parameters<F>>;

type ObjMock<T extends {}> = {
  [K in keyof T]: T[K] extends Fn ? FnMock<T[K]> : ObjMock<T[K]>;
};

type Module = { [param: string]: Fn };
type ModuleMock<M extends Module> = {
  [K in keyof M]: M[K] extends Fn ? FnMock<M[K]> : ObjMock<M[K]>;
};

/**
 * Allow pending jobs in Node's PromiseJobs queue to run.
 *
 * May be needed when testing tasks triggered by timers.
 *
 * See https://stackoverflow.com/questions/52177631/jest-timer-and-promise-dont-work-well-settimeout-and-async-function
 */
export async function advancePromiseQueue(n: number = 1): Promise<void> {
  for (let i = 0; i < n; i += 1) {
    await Promise.resolve();
  }
}

/**
 * Wrap a mocked function or module with the appropriate Jest mock typings.
 */
export function asMock<F extends Fn>(mock: F): FnMock<F>;
export function asMock<M extends Module>(mock: M): ModuleMock<M>;
export function asMock(mock: Fn | Module): FnMock<Fn> | ModuleMock<Module> {
  if (typeof mock === 'function') {
    return mock as FnMock<typeof mock>;
  }
  return mock as ModuleMock<typeof mock>;
}

/**
 * Mock the value returned by Date.now().
 */
export function mockNow(now: Date | number): FnMock<typeof Date.now> {
  const nowMs: number = new Date(now).getTime();
  jest.spyOn(Date, 'now').mockImplementation(() => nowMs);
  return asMock(Date.now);
}
