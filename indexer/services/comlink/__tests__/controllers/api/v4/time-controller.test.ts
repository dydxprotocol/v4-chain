import { RequestMethod } from '../../../../src/types';
import { sendRequest } from '../../../helpers/helpers';

describe('time-controller#V4', () => {
  describe('GET', () => {
    it('Get /time', async () => {
      const time: any = await sendRequest({
        type: RequestMethod.GET,
        path: '/v4/time',
      });

      expect(time.body.iso).not.toBeNull();
      expect(time.body.iso).not.toBeUndefined();
      expect(time.body.epoch).not.toBeNull();
      expect(time.body.epoch).not.toBeUndefined();
    });
  });
});
