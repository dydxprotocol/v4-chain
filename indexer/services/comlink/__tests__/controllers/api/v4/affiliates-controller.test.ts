import {
  dbHelpers,
  testConstants,
  testMocks,
  SubaccountUsernamesTable,
  WalletTable,
  AffiliateReferredUsersTable,
} from '@dydxprotocol-indexer/postgres';
import { AffiliateSnapshotRequest, RequestMethod } from '../../../../src/types';
import request from 'supertest';
import { sendRequest } from '../../../helpers/helpers';
import { defaultWallet, defaultWallet2 } from '@dydxprotocol-indexer/postgres/build/__tests__/helpers/constants';

describe('affiliates-controller#V4', () => {
  beforeAll(async () => {
    await dbHelpers.migrate();
  });

  afterAll(async () => {
    await dbHelpers.teardown();
  });

  describe('GET /metadata', () => {
    beforeEach(async () => {
      await testMocks.seedData();
      await SubaccountUsernamesTable.create(testConstants.defaultSubaccountUsername);
    });
  
    afterEach(async () => {
      await dbHelpers.clearData();
    });

    it('should return referral code for address with username', async () => {
      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/affiliates/metadata?address=${testConstants.defaultWallet.address}`,
        expectedStatus: 200,  // helper performs expect on status
      });

      expect(response.body).toEqual({
        // username is the referral code
        referralCode: testConstants.defaultSubaccountUsername.username,
        isVolumeEligible: false,
        isAffiliate: false,
      });
    });

    it('should fail if address does not exist', async () => {
      const nonExistentAddress = 'adgsakhasgt'
      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/affiliates/metadata?address=${nonExistentAddress}`,
        expectedStatus: 404, // helper performs expect on status
      });
    });

    it('should classify not volume eligible', async () => {      
      await WalletTable.update(
        { 
          address: testConstants.defaultWallet.address,
          totalVolume: "0",
          totalTradingRewards: "0",
        },
      );
      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/affiliates/metadata?address=${testConstants.defaultWallet.address}`,
        expectedStatus: 200,  // helper performs expect on status
      });
      expect(response.body).toEqual({
        referralCode: testConstants.defaultSubaccountUsername.username,
        isVolumeEligible: false,
        isAffiliate: false,
      });
    });

    it('should classify volume eligible', async () => {      
      await WalletTable.update(
        { 
          address: testConstants.defaultWallet.address,
          totalVolume: "100000",
          totalTradingRewards: "0",
        },
      );
      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/affiliates/metadata?address=${testConstants.defaultWallet.address}`,
        expectedStatus: 200,  // helper performs expect on status
      });
      expect(response.body).toEqual({
        referralCode: testConstants.defaultSubaccountUsername.username,
        isVolumeEligible: true,
        isAffiliate: false,
      });
    });  

    it('should classify is not affiliate', async () => {
      // AffiliateReferredUsersTable is empty
      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/affiliates/metadata?address=${testConstants.defaultWallet.address}`,
        expectedStatus: 200,  // helper performs expect on status
      });
      expect(response.body).toEqual({
        referralCode: testConstants.defaultSubaccountUsername.username,
        isVolumeEligible: false,
        isAffiliate: false,
      });
    });

    it('should classify is affiliate', async () => {
      await AffiliateReferredUsersTable.create({
        affiliateAddress: defaultWallet.address,
        refereeAddress: defaultWallet2.address,
        referredAtBlock: '1',
      });
      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/affiliates/metadata?address=${testConstants.defaultWallet.address}`,
        expectedStatus: 200,  // helper performs expect on status
      });
      expect(response.body).toEqual({
        referralCode: testConstants.defaultSubaccountUsername.username,
        isVolumeEligible: false,
        isAffiliate: true,
      });
    });

    it('should fail if subaccount username not found', async () => {
      // create defaultWallet2 without subaccount username
      WalletTable.create(testConstants.defaultWallet2);
      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/affiliates/metadata?address=${testConstants.defaultWallet2.address}`,
        expectedStatus: 500,  // helper performs expect on status
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
