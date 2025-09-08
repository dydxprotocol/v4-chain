import fetch, { Headers, Request, Response } from 'node-fetch';

// Polyfill Headers for Node.js 16
if (!globalThis.fetch) {
  // @ts-ignore
  globalThis.fetch = fetch;
  // @ts-ignore
  globalThis.Headers = Headers;
  // @ts-ignore
  globalThis.Request = Request;
  // @ts-ignore
  globalThis.Response = Response;
}

// Add global declaration for types
declare global {
  interface Array<T> {
    findLast(predicate: (value: T, index: number, obj: T[]) => unknown): T | undefined,
  }
}

// Polyfill
if (!Array.prototype.findLast) {
  // eslint-disable-next-line no-extend-native
  Array.prototype.findLast = function findLast<T>(
    this: T[],
    callback: (element: T, index: number, array: T[]) => unknown,
  ): T | undefined {
    if (this == null) {
      throw new TypeError('this is null or not defined');
    }
    const len = this.length;
    for (let i = len - 1; i >= 0; i--) {
      if (callback(this[i], i, this)) {
        return this[i];
      }
    }
    return undefined;
  };
}
