import { safeAxiosRequest } from '@dydxprotocol-indexer/base';
import { AxiosRequestConfig } from 'axios';

import config from '../config';

// keep as a wrapper for tests
// eslint-disable-next-line  @typescript-eslint/require-await
export async function axiosRequest(
  axiosConfig: AxiosRequestConfig,
// eslint-disable-next-line  @typescript-eslint/no-explicit-any
): Promise<any> {
  return safeAxiosRequest({
    ...axiosConfig,
    // Comlink compresses responses, and serves brotli to clients that accept it. axios 1.2.x
    // advertises brotli but intermittently fails to decode it ("unexpected end of file"), so ask
    // for gzip, which it decodes reliably.
    headers: { ...axiosConfig.headers, 'Accept-Encoding': 'gzip' },
    timeout: config.AXIOS_TIMEOUT_MS,
  });
}
