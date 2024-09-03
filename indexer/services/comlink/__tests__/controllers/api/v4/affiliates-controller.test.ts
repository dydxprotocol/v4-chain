import { AffiliateSnapshotRequest, RequestMethod } from '../../../../src/types';
import request from 'supertest';
import { sendRequest } from '../../../helpers/helpers';

describe('affiliates-controller#V4', () => {
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

  describe('GET /address', () => {
    it('should return address for a valid referral code string', async () => {
      const referralCode = 'TempCode123';
      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/affiliates/address?referralCode=${referralCode}`,
      });

      expect(response.status).toBe(200);
      expect(response.body).toEqual({
        address: 'some_address',
      });
    });
  });

  describe('GET /snapshot', () => {
    it('should return snapshots when all params specified', async () => {
      const req: AffiliateSnapshotRequest = {
        limit: 10,
        offset: 10,
        sortByReferredFees: true,
      };
      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/affiliates/snapshot?limit=${req.limit}&offset=${req.offset}&sortByReferredFees=${req.sortByReferredFees}`,
      });

      expect(response.status).toBe(200);
      expect(response.body.affiliateList).toHaveLength(10);
      expect(response.body.currentOffset).toBe(10);
      expect(response.body.total).toBe(10);
    });

    it('should return snapshots when optional params not specified', async () => {
      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: '/v4/affiliates/snapshot',
      });

      expect(response.status).toBe(200);
      expect(response.body.affiliateList).toHaveLength(1000);
      expect(response.body.currentOffset).toBe(0);
      expect(response.body.total).toBe(1000);
    });
  });

  describe('GET /total_volume', () => {
    it('should return total_volume for a valid address', async () => {
      const address = 'some_address';
      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/affiliates/total_volume?address=${address}`,
      });

      expect(response.status).toBe(200);
      expect(response.body).toEqual({
        totalVolume: 111.1,
      });
    });
  });
});
