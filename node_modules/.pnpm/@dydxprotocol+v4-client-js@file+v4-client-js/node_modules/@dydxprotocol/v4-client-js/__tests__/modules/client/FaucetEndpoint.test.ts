import { FaucetApiHost } from '../../../src/clients/constants';
import { FaucetClient } from '../../../src/clients/faucet-client';
import { DYDX_TEST_ADDRESS } from './constants';

describe('FaucetClient', () => {
  const client = new FaucetClient(FaucetApiHost.TESTNET);

  describe('Faucet Endpoints', () => {
    it('Fill', async () => {
      const response = await client.fill(DYDX_TEST_ADDRESS, 0, 2000);
      expect(response?.status).toBe(202);
    });
  });
});
