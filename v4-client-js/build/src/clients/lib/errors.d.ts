import { BroadcastTxSyncResponse } from '@cosmjs/tendermint-rpc/build/tendermint37';
/**
 * An edge-case was hit in the client that should never have been reached.
 */
export declare class UnexpectedClientError extends Error {
    constructor();
}
/**
 * An error occurred during the broadcasting process.
 */
export declare class BroadcastErrorObject extends Error {
    result: BroadcastTxSyncResponse;
    code: number;
    codespace?: string;
    constructor(message: string, result: BroadcastTxSyncResponse);
}
/**
 * User error occurred during a client operation.
 */
export declare class UserError extends Error {
    constructor(message: string);
}
/**
 * @description Base class for custom errors.
 */
export declare class CustomError extends Error {
    constructor(message: string);
}
/**
 * @description Base class for a custom error which wraps another error.
 */
export declare class WrappedError extends CustomError {
    readonly originalError: Error;
    constructor(message: string, originalError: Error);
}
