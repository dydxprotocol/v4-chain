import { safeAxiosRequest } from '@dydxprotocol-indexer/base';

import config from '../../src/config';
import { axiosRequest } from '../../src/lib/axios';

jest.mock('@dydxprotocol-indexer/base', () => ({
  ...jest.requireActual('@dydxprotocol-indexer/base'),
  safeAxiosRequest: jest.fn(),
}));

describe('axiosRequest', () => {
  it('asks comlink for gzip and keeps the caller\'s headers', async () => {
    await axiosRequest({
      url: 'http://comlink/v4/trades/perpetualMarket/BTC-USD',
      headers: { 'cf-ipcountry': 'GB' },
    });

    expect(safeAxiosRequest).toHaveBeenCalledWith({
      url: 'http://comlink/v4/trades/perpetualMarket/BTC-USD',
      headers: { 'cf-ipcountry': 'GB', 'Accept-Encoding': 'gzip' },
      timeout: config.AXIOS_TIMEOUT_MS,
    });
  });
});
