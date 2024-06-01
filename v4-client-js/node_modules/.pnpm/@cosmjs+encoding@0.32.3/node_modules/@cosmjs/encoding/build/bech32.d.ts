export declare function toBech32(prefix: string, data: Uint8Array, limit?: number): string;
export declare function fromBech32(address: string, limit?: number): {
    readonly prefix: string;
    readonly data: Uint8Array;
};
/**
 * Takes a bech32 address and returns a normalized (i.e. lower case) representation of it.
 *
 * The input is validated along the way, which makes this significantly safer than
 * using `address.toLowerCase()`.
 */
export declare function normalizeBech32(address: string): string;
