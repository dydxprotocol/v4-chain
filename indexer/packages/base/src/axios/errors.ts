import { WrappedError } from '../errors';

export interface AxiosOriginalError extends Error {
  isAxiosError: true,
  toJSON(): Error,
}

export interface AxiosErrorResponse {
  status: number,
  statusText: string,
  data: {},
}

/**
 * An error thrown by axios.
 *
 * Depending on your use case, if logging errors, you may want to catch axios errors and sanitize
 * them to remove the request and response objects, or sensitive fields. For example:
 *
 *   this.originalError = _.omit(originalError.toJSON(), 'config')
 */
export class AxiosError extends WrappedError<AxiosOriginalError> {}

/**
 * Axios error with response error fields.
 */
export class AxiosServerError extends AxiosError {
  public readonly status: number;
  public readonly statusText: string;
  public readonly data: {};

  constructor(
    response: AxiosErrorResponse,
    originalError: AxiosOriginalError,
  ) {
    super(
      `${response.status}: ${response.statusText} - ${JSON.stringify(response.data, null, 2)}`,
      originalError,
    );
    this.status = response.status;
    this.statusText = response.statusText;
    this.data = response.data;
  }
}

// https://github.com/axios/axios/blob/cd7ff042b0b80f6f02e5564d184019131c90cacd/lib/core/enhanceError.js#L23
export type AxiosSafeErrorObject = {};

/**
 * An error thrown by axios made safer with large fields removed and the original error converted
 * into an object with its config removed.
 */
export class AxiosSafeError extends WrappedError<AxiosSafeErrorObject> {}

/**
  * Axios error with only status, statusText and data response error fields and a smaller original
  * error.
  */
export class AxiosSafeServerError extends AxiosSafeError {
  public readonly status: number;
  public readonly statusText: string;
  public readonly data: {};

  constructor(
    response: { data: {}, status: number, statusText: string },
    originalError: AxiosSafeErrorObject,
  ) {
    super(
      `${response.status}: ${response.statusText} - ${JSON.stringify(response.data, null, 2)}`,
      originalError,
    );
    this.status = response.status;
    this.statusText = response.statusText;
    this.data = response.data;
  }
}

/**
 * Helper functions to make an Error look like an error returned by Axios in tests.
 */
export function makeAxiosOriginalError(
  message: string,
  additionalFields: {} = {},
): AxiosOriginalError {
  const error = new Error(message) as AxiosOriginalError;
  Object.assign(error, additionalFields);
  error.isAxiosError = true;
  error.toJSON = () => error;
  return error;
}

export function makeAxiosError(
  message: string,
  additionalFields: {} = {},
): AxiosError {
  return new AxiosError(message, makeAxiosOriginalError(message, additionalFields));
}

export function makeAxiosServerError(
  status: number,
  statusText: string,
  message: string,
  additionalFields: {} = {},
): AxiosServerError {
  return new AxiosServerError(
    {
      status,
      statusText,
      data: {},
    },
    makeAxiosOriginalError(message, additionalFields),
  );
}

export function makeAxiosSafeServerError(
  status: number,
  statusText: string,
  message: string,
  additionalFields: {} = {},
): AxiosSafeServerError {
  return new AxiosSafeServerError(
    {
      status,
      statusText,
      data: {},
    },
    makeAxiosOriginalError(message, additionalFields),
  );
}
