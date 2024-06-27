export declare function toUtf8(str: string): Uint8Array;
/**
 * Takes UTF-8 data and decodes it to a string.
 *
 * In lossy mode, the [REPLACEMENT CHARACTER](https://en.wikipedia.org/wiki/Specials_(Unicode_block))
 * is used to substitude invalid encodings.
 * By default lossy mode is off and invalid data will lead to exceptions.
 */
export declare function fromUtf8(data: Uint8Array, lossy?: boolean): string;
