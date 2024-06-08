import { MAX_UINT_32 } from '../clients/constants';

/**
 * Returns a random integer value between 0 and (n-1).
 */
export function randomInt(
  n: number,
): number {
  return Math.floor(Math.random() * n);
}

/**
 * Generate a random clientId.
 */
export function generateRandomClientId(): number {
  return randomInt(MAX_UINT_32 + 1);
}

/**
 * Deterministically generate a valid clientId from an arbitrary string by performing a
 * quick hashing function on the string.
 */
export function clientIdFromString(
  input: string,
): number {
  let hash: number = 0;
  if (input.length === 0) return hash;
  for (let i = 0; i < input.length; i++) {
    hash = ((hash << 5) - hash) + input.charCodeAt(i); // eslint-disable-line no-bitwise
    hash |= 0; // eslint-disable-line no-bitwise
  }

  // Bitwise operators covert the value to a 32-bit integer.
  // We must coerce this into a 32-bit unsigned integer.
  return hash + (2 ** 31);
}

/**
 * Pauses the execution of the program for a specified time.
 * @param ms - The number of milliseconds to pause the program.
 * @returns A promise that resolves after the specified number of milliseconds.
 */
export async function sleep(ms: number): Promise<void> {
  return new Promise((resolve) => setTimeout(resolve, ms));
}
