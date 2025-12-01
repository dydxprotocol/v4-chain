import { logger } from '@dydxprotocol-indexer/base';

/**
 * @function sanitizeArray
 * @param input input value from query
 * @description Handles both comma-separated strings and repeated query parameters.
 * If input is a string, converts to uppercase and splits by comma.
 * If input is already an array, uppercases each element.
 * Returns null for empty values.
 */
export function sanitizeArray(
  input: string | string[],
): string[] | null {
  try {
    // Handle array input (repeated query parameters: ?includeTypes=LIMIT&includeTypes=MARKET)
    if (Array.isArray(input)) {
      if (input.length === 0) {
        return null;
      }
      return input.map((item) => item.toUpperCase());
    }

    // Handle string input (comma-separated: ?includeTypes=LIMIT,MARKET)
    if (input === '') {
      return null;
    }
    return input.toUpperCase().split(',');
  } catch (error) {
    logger.error({
      at: 'request-helpers#sanitizeArray',
      message: 'Failed to sanitize array',
      input,
    });
    throw error;
  }
}
