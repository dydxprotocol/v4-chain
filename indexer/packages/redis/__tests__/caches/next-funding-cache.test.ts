import { addFundingSample, clearFundingSamples, getNextFunding } from '../../src/caches/next-funding-cache';
import { deleteAllAsync } from '../../src/helpers/redis';
import { redis as client } from '../helpers/utils';
import Big from 'big.js';

describe('nextFundingCache', () => {
  beforeEach(async () => {
    await deleteAllAsync(client);
  });

  afterEach(async () => {
    await deleteAllAsync(client);
  });

  describe('getNextFunding', () => {
    it('get next funding', async () => {
      await addFundingSample('BTC', new Big('0.0001'), client);
      await addFundingSample('BTC', new Big('0.0002'), client);  // avg = 0.00015
      await addFundingSample('ETH', new Big('0.0005'), client);  // avg = 0.0005
      expect(await getNextFunding(client, ['BTC', 'ETH'])).toEqual(
        { BTC: new Big('0.00015'), ETH: new Big('0.0005') },
      );
    });

    it('clear funding samples', async () => {
      await addFundingSample('BTC', new Big('0.0001'), client);
      await addFundingSample('BTC', new Big('0.0002'), client);  // avg = 0.00015
      await clearFundingSamples('BTC', client);
      await addFundingSample('ETH', new Big('0.0005'), client);  // avg = 0.0005
      expect(await getNextFunding(client, ['BTC', 'ETH'])).toEqual(
        { BTC: undefined, ETH: new Big('0.0005') },
      );
    });

    it('get next funding with no values', async () => {
      expect(await getNextFunding(client, ['BTC'])).toEqual(
        { BTC: undefined },
      );
    });
  });
});
