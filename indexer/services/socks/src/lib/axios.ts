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
    // Comlink serves brotli to clients that accept it. The axios that makes this request is the
    // one @dydxprotocol-indexer/base resolves (1.2.1), which decodes brotli with zlib flush options
    // and intermittently fails with "unexpected end of file" (fixed in axios 1.2.2, axios#5353).
    // Ask for gzip until base's axios is upgraded; bumping socks' own axios does not help.
    headers: { ...axiosConfig.headers, 'Accept-Encoding': 'gzip' },
    timeout: config.AXIOS_TIMEOUT_MS,
  });
}
