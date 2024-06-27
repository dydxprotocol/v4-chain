import { IndexerClient } from '../../../src';
import { Network } from '../../../src/clients/constants';
import { DYDX_TEST_ADDRESS } from './constants';

describe('IndexerClient', () => {
  const client = new IndexerClient(Network.testnet().indexerConfig);

  describe('Utility Endpoints', () => {
    it('getTime', async () => {
      const response = await client.utility.getTime();
      const iso = response.iso;
      expect(iso).not.toBeUndefined();
    });

    it('getHeight', async () => {
      const response = await client.utility.getHeight();
      const height = response.height;
      const time = response.time;
      expect(height).not.toBeUndefined();
      expect(time).not.toBeUndefined();
    });

    it('Screen Address', async () => {
      const response = await client.utility.screen(DYDX_TEST_ADDRESS);
      const { restricted } = response ?? {};
      expect(restricted).toBeDefined();
    });
  });
});
