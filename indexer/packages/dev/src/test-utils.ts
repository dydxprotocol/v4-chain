/**
 * Utilities for writing unit tests with jest.
 */

/* eslint-disable @typescript-eslint/no-explicit-any */
type Fn = (...args: any[]) => any;
type FnMock<F extends Fn> = jest.Mock<ReturnType<F>, Parameters<F>>;

type ObjMock<T extends {}> = {
  [K in keyof T]: T[K] extends Fn
    ? FnMock<T[K]>
    : T[K] extends object
      ? ObjMock<T[K]>
      : T[K];
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

/**
 * Function to make `wrapBackgroundTask` a synchronous function for testing.
 * Requires code like this at the top of the file:
 *
 * jest.mock('@dydxprotocol-indexer/base', () => ({
 *   ...jest.requireActual('@dydxprotocol-indexer/base'),
 *   wrapBackgroundTask: jest.fn(),
 * }));
 *
 * @param wrapBackgroundTask
 */
export function synchronizeWrapBackgroundTask(
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  wrapBackgroundTask: (t: Promise<any>, sre: boolean, tn: string | null) => void,
): void {
  asMock(wrapBackgroundTask).mockImplementation(async (
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    task: Promise<any>,
    _shouldRethrowErrors: boolean,
    _taskName?: string | null,
  ) => {
    await task;
  });
}
