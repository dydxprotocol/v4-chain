import { RequestMethod } from '../../../../src/types';
import request from 'supertest';
import { sendRequest } from '../../../helpers/helpers';

describe('test-controller#V4', () => {
  describe('GET /referral_code', () => {
    it('should return referral code for a valid address string', async () => {
      const address = 'some_address';
      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/affiliates/referral_code?address=${address}`,
      });

      expect(response.status).toBe(200);
      expect(response.body).toEqual({
        referralCode: 'TempCode123',
      });
    });
  });
});
