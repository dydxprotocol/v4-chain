import { DEFAULT_API_TIMEOUT } from '../constants';
import { generateQueryPath } from '../helpers/request-helpers';
import { RequestMethod, Response, request } from '../lib/axios';
import { Data } from '../types';

export default class RestClient {
    readonly host: string;
    readonly apiTimeout: Number;

    constructor(host: string, apiTimeout?:Number | null) {
      if (host.endsWith('/')) {
        this.host = host.slice(0, -1);
      } else {
        this.host = host;
      }
      this.apiTimeout = apiTimeout || DEFAULT_API_TIMEOUT;
    }

    async get(
      requestPath: string,
      params: {} = {},
    ): Promise<Data> {
      const url = `${this.host}${generateQueryPath(requestPath, params)}`;
      const response = await request(url);
      return response.data;
    }

    async post(
      requestPath: string,
      params: {} = {},
      body?: unknown | null,
      headers: {} = {},
    ): Promise<Response> {
      const url = `${this.host}${generateQueryPath(requestPath, params)}`;
      return request(url, RequestMethod.POST, body, headers);
    }
}
