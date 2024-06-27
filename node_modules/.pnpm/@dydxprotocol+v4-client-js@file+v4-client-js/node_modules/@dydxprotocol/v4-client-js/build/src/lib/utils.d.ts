/**
 * Returns a random integer value between 0 and (n-1).
 */
export declare function randomInt(n: number): number;
/**
 * Generate a random clientId.
 */
export declare function generateRandomClientId(): number;
/**
 * Deterministically generate a valid clientId from an arbitrary string by performing a
 * quick hashing function on the string.
 */
export declare function clientIdFromString(input: string): number;
/**
 * Pauses the execution of the program for a specified time.
 * @param ms - The number of milliseconds to pause the program.
 * @returns A promise that resolves after the specified number of milliseconds.
 */
export declare function sleep(ms: number): Promise<void>;
