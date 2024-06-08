import { ValidatorClient } from '../../../src/clients/validator-client';
import { Network } from '../../../src/clients/constants';
import { DYDX_TEST_ADDRESS } from './constants';

describe('Validator Client', () => {
  let client: ValidatorClient;
  describe('Get', () => {
    beforeEach(async () => {
      client = await ValidatorClient.connect(Network.testnet().validatorConfig);
    });

    it('Account', async () => {
      const account = await client.get.getAccount(DYDX_TEST_ADDRESS);
      expect(account.address).toBe(DYDX_TEST_ADDRESS);
    });

    it('Balance', async () => {
      const balances = await client.get.getAccountBalances(DYDX_TEST_ADDRESS);
      expect(balances).not.toBeUndefined();
    });

    it('All Subaccounts', async () => {
      const subaccounts = await client.get.getSubaccounts();
      expect(subaccounts.subaccount).not.toBeUndefined();
    });

    it('Subaccount', async () => {
      const subaccount = await client.get.getSubaccount(DYDX_TEST_ADDRESS, 0);
      expect(subaccount.subaccount).not.toBeUndefined();
    });

    it('Clob pairs', async () => {
      const clobpairs = await client.get.getAllClobPairs();
      expect(clobpairs.clobPair).not.toBeUndefined();
      expect(clobpairs.clobPair[0].id).toBe(0);
    });

    it('Prices', async () => {
      const prices = await client.get.getAllPrices();
      expect(prices.marketPrices).not.toBeUndefined();
      expect(prices.marketPrices[0].id).toBe(0);
    });

    it('Equity tier limit configuration', async () => {
      const response = await client.get.getEquityTierLimitConfiguration();
      expect(response.equityTierLimitConfig).not.toBeUndefined();
      expect(response.equityTierLimitConfig?.shortTermOrderEquityTiers[0].limit).toBe(0);
    });
  });
});
