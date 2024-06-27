import { ValidationChain } from '../chain';
import { Middleware, Request } from '../base';
import { Result } from '../validation-result';
export declare type OneOfCustomMessageBuilder = (options: {
    req: Request;
}) => any;
/**
 * Creates a middleware that will ensure that at least one of the given validation chains
 * or validation chain groups are valid.
 *
 * If none are, a single error for a pseudo-field `_error` is added to the request,
 * with the errors of each chain made available under the `nestedErrors` property.
 *
 * @param chains an array of validation chains to check if are valid.
 *               If any of the items of `chains` is an array of validation chains, then all of them
 *               must be valid together for the request to be considered valid.
 * @param message a function for creating a custom error message in case none of the chains are valid
 */
export declare function oneOf(chains: (ValidationChain | ValidationChain[])[], message?: OneOfCustomMessageBuilder): Middleware & {
    run: (req: Request) => Promise<Result>;
};
/**
 * Creates a middleware that will ensure that at least one of the given validation chains
 * or validation chain groups are valid.
 *
 * If none are, a single error for a pseudo-field `_error` is added to the request,
 * with the errors of each chain made available under the `nestedErrors` property.
 *
 * @param chains an array of validation chains to check if are valid.
 *               If any of the items of `chains` is an array of validation chains, then all of them
 *               must be valid together for the request to be considered valid.
 * @param message an error message to use in case none of the chains are valid
 */
export declare function oneOf(chains: (ValidationChain | ValidationChain[])[], message?: any): Middleware & {
    run: (req: Request) => Promise<Result>;
};
