import { WrappedError } from '../errors';
interface AxiosOriginalError extends Error {
    isAxiosError: true;
    toJSON(): {};
}
interface AxiosErrorResponse {
    status: number;
    statusText: string;
    data: {};
}
/**
 * @description An error thrown by axios.
 *
 * Depending on your use case, if logging errors, you may want to catch axios errors and sanitize
 * them to remove the request and response objects, or sensitive fields. For example:
 *
 *   this.originalError = _.omit(originalError.toJSON(), 'config')
 */
export declare class AxiosError extends WrappedError {
}
/**
 * @description Axios error with response error fields.
 */
export declare class AxiosServerError extends AxiosError {
    readonly status: number;
    readonly statusText: string;
    readonly data: {};
    constructor(response: AxiosErrorResponse, originalError: AxiosOriginalError);
}
export {};
