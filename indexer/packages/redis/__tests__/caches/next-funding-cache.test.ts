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
      expect(await getNextFunding(client, [
        ['BTC', '0.0001'],
        ['ETH', '0'],
      ])).toEqual(
        {
          BTC: new Big('0.00025'), // 0.00015 + 0.0001
          ETH: new Big('0.0005'), // 0.0005 + 0
        },
      );
    });

    it('clear funding samples', async () => {
      await addFundingSample('BTC', new Big('0.0001'), client);
      await addFundingSample('BTC', new Big('0.0002'), client);  // avg = 0.00015
      await clearFundingSamples('BTC', client);
      await addFundingSample('ETH', new Big('0.0005'), client);  // avg = 0.0005
      expect(await getNextFunding(client, [
        ['BTC', '0.0001'],
        ['ETH', '0.00015'],
      ])).toEqual(
        {
          BTC: undefined, // no samples
          ETH: new Big('0.00065'), // 0.0005 + 0.00015
        },
      );
    });

    it('get next funding with no values', async () => {
      expect(await getNextFunding(client, [
        ['BTC', '0.001'],
      ])).toEqual(
        // Even though default funding rate is 0.001,
        // return undefined since there are no samples
        { BTC: undefined },
      );
    });
  });
});
