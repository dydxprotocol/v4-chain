/**
 * Returns the Node.js crypto module when available and `undefined`
 * otherwise.
 *
 * Detects an unimplemented fallback module from Webpack 5 and returns
 * `undefined` in that case.
 */
export declare function getNodeCrypto(): Promise<any | undefined>;
export declare function getSubtle(): Promise<any | undefined>;
export declare function pbkdf2Sha512Subtle(subtle: any, secret: Uint8Array, salt: Uint8Array, iterations: number, keylen: number): Promise<Uint8Array>;
/**
 * Implements pbkdf2-sha512 using the Node.js crypro module (`import "crypto"`).
 * This does not use subtle from [Crypto](https://developer.mozilla.org/en-US/docs/Web/API/Crypto).
 */
export declare function pbkdf2Sha512NodeCrypto(nodeCrypto: any, secret: Uint8Array, salt: Uint8Array, iterations: number, keylen: number): Promise<Uint8Array>;
export declare function pbkdf2Sha512Noble(secret: Uint8Array, salt: Uint8Array, iterations: number, keylen: number): Promise<Uint8Array>;
/**
 * A pbkdf2 implementation for BIP39. This is not exported at package level and thus a private API.
 */
export declare function pbkdf2Sha512(secret: Uint8Array, salt: Uint8Array, iterations: number, keylen: number): Promise<Uint8Array>;
