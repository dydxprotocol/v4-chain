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
    timeout: config.AXIOS_TIMEOUT_MS,
  });
}
