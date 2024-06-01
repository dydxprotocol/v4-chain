import { Data } from '../../types';
import { RequestMethod } from './types';
export interface Response {
    status: number;
    data: Data;
    headers: {};
}
export declare function request(url: string, method?: RequestMethod, body?: unknown | null, headers?: {}): Promise<Response>;
