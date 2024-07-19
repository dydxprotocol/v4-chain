import { RequestMethod } from '../../../../src/types';
import { sendRequest } from '../../../helpers/helpers';
import {
  BlockFromDatabase,
  BlockTable,
  dbHelpers,
  testMocks,
} from '@dydxprotocol-indexer/postgres';
import { stats } from '@dydxprotocol-indexer/base';

describe('height-controller#V4', () => {
  beforeAll(async () => {
    await dbHelpers.migrate();
    jest.spyOn(stats, 'increment');
    jest.spyOn(stats, 'timing');
  });

  afterAll(async () => {
    await dbHelpers.teardown();
  });

  describe('GET', () => {
    afterEach(async () => {
      await dbHelpers.clearData();
    });

    it('Get /height gets latest block', async () => {
      await testMocks.seedData();
      const latestBlock: BlockFromDatabase = await BlockTable.getLatest();
      const block: any = await sendRequest({
        type: RequestMethod.GET,
        path: '/v4/height',
      });

      expect(block.body.height).toEqual(latestBlock!.blockHeight);
      expect(block.body.time).toEqual(latestBlock!.time);
      expect(stats.timing).toBeCalledTimes(1);
      expect(stats.increment).toHaveBeenCalledWith('comlink.height-controller.response_status_code.200', 1,
        {
          path: '/',
          method: 'GET',
        });
    });

    it('Get /height returns 404 if no blocks', async () => {
      await sendRequest({
        type: RequestMethod.GET,
        path: '/v4/height',
        expectedStatus: 404,
      });
      expect(stats.timing).toBeCalledTimes(1);
      expect(stats.increment).toHaveBeenCalledWith('comlink.height-controller.response_status_code.404', 1,
        {
          path: '/',
          method: 'GET',
        });
    });
  });
});
