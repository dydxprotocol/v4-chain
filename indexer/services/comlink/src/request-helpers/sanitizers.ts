import { logger } from '@dydxprotocol-indexer/base';

/**
 * @function sanitizeArray
 * @param input input value from query
 * @description Checks if the input is empty and if it isn't, set the string to upper case and split
 * it using `,` as the delimiter.
 */
export function sanitizeArray(
  input: string,
): string[] | null {
  try {
    return (
      // eslint-disable-next-line no-mixed-operators
      (input !== '') && input.toUpperCase().split(',') || null);
  } catch (error) {
    logger.error({
      at: 'request-helpers#sanitizeArray',
      message: 'Failed to sanitize array',
      input,
    });
    throw error;
  }
}
