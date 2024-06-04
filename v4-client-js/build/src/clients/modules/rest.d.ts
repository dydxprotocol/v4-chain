import { Response } from '../lib/axios';
import { Data } from '../types';
export default class RestClient {
    readonly host: string;
    readonly apiTimeout: Number;
    constructor(host: string, apiTimeout?: Number | null);
    get(requestPath: string, params?: {}): Promise<Data>;
    post(requestPath: string, params?: {}, body?: unknown | null, headers?: {}): Promise<Response>;
}
