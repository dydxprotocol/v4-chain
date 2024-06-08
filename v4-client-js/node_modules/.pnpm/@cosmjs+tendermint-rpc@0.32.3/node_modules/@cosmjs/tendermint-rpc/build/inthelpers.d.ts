/**
 * Takes an integer value from the Tendermint RPC API and
 * returns it as number.
 *
 * Only works within the safe integer range.
 */
export declare function apiToSmallInt(input: string | number): number;
/**
 * Takes an integer value from the Tendermint RPC API and
 * returns it as BigInt.
 *
 * This supports the full uint64 and int64 ranges.
 */
export declare function apiToBigInt(input: string): bigint;
/**
 * Takes an integer in the safe integer range and returns
 * a string representation to be used in the Tendermint RPC API.
 */
export declare function smallIntToApi(num: number): string;
