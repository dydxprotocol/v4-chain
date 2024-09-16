import {
  dbHelpers,
  testConstants,
  testMocks,
  SubaccountTable,
  SubaccountUsernamesTable,
  WalletTable,
  AffiliateReferredUsersTable,
  AffiliateInfoTable,
  AffiliateInfoCreateObject,
} from '@dydxprotocol-indexer/postgres';
import {
  AffiliateSnapshotRequest,
  AffiliateSnapshotResponse,
  RequestMethod,
  AffiliateSnapshotResponseObject,
} from '../../../../src/types';
import request from 'supertest';
import { sendRequest } from '../../../helpers/helpers';

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
      const nonExistentAddress = 'adgsakhasgt';
      await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/affiliates/metadata?address=${nonExistentAddress}`,
        expectedStatus: 404, // helper performs expect on status
      });
    });

    it('should classify not volume eligible', async () => {
      await WalletTable.update(
        {
          address: testConstants.defaultWallet.address,
          totalVolume: '0',
          totalTradingRewards: '0',
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
          totalVolume: '100000',
          totalTradingRewards: '0',
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
        affiliateAddress: testConstants.defaultWallet.address,
        refereeAddress: testConstants.defaultWallet2.address,
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
      await WalletTable.create(testConstants.defaultWallet2);
      await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/affiliates/metadata?address=${testConstants.defaultWallet2.address}`,
        expectedStatus: 500,  // helper performs expect on status
      });
    });
  });

  describe('GET /address', () => {
    beforeEach(async () => {
      await testMocks.seedData();
      await SubaccountUsernamesTable.create(testConstants.defaultSubaccountUsername);
    });

    afterEach(async () => {
      await dbHelpers.clearData();
    });

    it('should return address for a valid referral code', async () => {
      const referralCode: string = testConstants.defaultSubaccountUsername.username;
      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/affiliates/address?referralCode=${referralCode}`,
        expectedStatus: 200,  // helper performs expect on status
      });

      expect(response.body).toEqual({
        address: testConstants.defaultWallet.address,
      });
    });

    it('should fail when referral code not found', async () => {
      const nonExistentReferralCode = 'BadCode123';
      await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/affiliates/address?referralCode=${nonExistentReferralCode}`,
        expectedStatus: 404,  // helper performs expect on status
      });
    });
  });

  describe('GET /snapshot', () => {
    const defaultInfo: AffiliateInfoCreateObject = testConstants.defaultAffiliateInfo;
    const defaultInfo2: AffiliateInfoCreateObject = testConstants.defaultAffiliateInfo2;
    const defaultInfo3: AffiliateInfoCreateObject = testConstants.defaultAffiliateInfo3;

    beforeEach(async () => {
      await testMocks.seedData();
      // Create username for defaultWallet
      await SubaccountUsernamesTable.create(testConstants.defaultSubaccountUsername);

      // Create defaultWallet2, subaccount, and username
      await WalletTable.create(testConstants.defaultWallet2);
      await SubaccountTable.create(testConstants.defaultSubaccountDefaultWalletAddress);
      await SubaccountUsernamesTable.create(
        testConstants.subaccountUsernameWithDefaultWalletAddress,
      );

      // Create defaultWallet3, create subaccount, create username
      await WalletTable.create(testConstants.defaultWallet3);
      await SubaccountTable.create(testConstants.defaultSubaccountWithAlternateAddress);
      await SubaccountUsernamesTable.create(testConstants.subaccountUsernameWithAlternativeAddress);

      // Create affiliate infos
      await AffiliateInfoTable.create(defaultInfo);
      await AffiliateInfoTable.create(defaultInfo2);
      await AffiliateInfoTable.create(defaultInfo3);
    });

    afterEach(async () => {
      await dbHelpers.clearData();
    });

    it('should return snapshots when optional params not specified', async () => {
      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: '/v4/affiliates/snapshot',
      });

      expect(response.status).toBe(200);
      expect(response.body.affiliateList).toHaveLength(3);
      expect(response.body.currentOffset).toEqual(0);
      expect(response.body.total).toEqual(3);
    });

    it('should filter by address', async () => {
      const req: AffiliateSnapshotRequest = {
        addressFilter: [testConstants.defaultWallet.address],
      };
      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/affiliates/snapshot?addressFilter=${req.addressFilter!.join(',')}`,
        expectedStatus: 200,  // helper performs expect on status,
      });

      const expectedResponse: AffiliateSnapshotResponse = {
        affiliateList: [
          affiliateInfoCreateToResponseObject(
            defaultInfo, testConstants.defaultSubaccountUsername.username,
          ),
        ],
        total: 1,
        currentOffset: 0,
      };
      expect(response.body.affiliateList).toHaveLength(1);
      expect(response.body.affiliateList[0]).toEqual(expectedResponse.affiliateList[0]);
      expect(response.body.currentOffset).toEqual(expectedResponse.currentOffset);
      expect(response.body.total).toEqual(expectedResponse.total);
    });

    it('should handle no results when filter by address', async () => {
      const req: AffiliateSnapshotRequest = {
        addressFilter: ['nonexistentaddress'],
      };
      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/affiliates/snapshot?addressFilter=${req.addressFilter!.join(',')}`,
        expectedStatus: 200,  // helper performs expect on status,
      });

      const expectedResponse: AffiliateSnapshotResponse = {
        affiliateList: [],
        total: 0,
        currentOffset: 0,
      };
      expect(response.body.affiliateList).toHaveLength(0);
      expect(response.body.affiliateList[0]).toEqual(expectedResponse.affiliateList[0]);
      expect(response.body.currentOffset).toEqual(expectedResponse.currentOffset);
      expect(response.body.total).toEqual(expectedResponse.total);
    });

    it('should handle offset out of bounds', async () => {
      const offset = 5;
      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/affiliates/snapshot?offset=${offset}`,
        expectedStatus: 200,  // helper performs expect on status,
      });

      const expectedResponse: AffiliateSnapshotResponse = {
        affiliateList: [],
        total: 0,
        currentOffset: offset,
      };
      expect(response.body.affiliateList).toHaveLength(0);
      expect(response.body.affiliateList[0]).toEqual(expectedResponse.affiliateList[0]);
      expect(response.body.currentOffset).toEqual(expectedResponse.currentOffset);
      expect(response.body.total).toEqual(expectedResponse.total);
    });

    it('should return snapshots when all params specified', async () => {
      const req: AffiliateSnapshotRequest = {
        addressFilter: [testConstants.defaultWallet.address, testConstants.defaultWallet2.address],
        sortByAffiliateEarning: true,
      };
      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/affiliates/snapshot?${req.addressFilter!.map((address) => `addressFilter[]=${address}`).join('&')}&offset=1&limit=1&sortByAffiliateEarning=${req.sortByAffiliateEarning}`,
        expectedStatus: 200,  // helper performs expect on status
      });

      // addressFilter removes defaultInfo3
      // sortorder -> [defaultInfo2, defaultInfo]
      // offset=1 -> defaultInfo
      const expectedResponse: AffiliateSnapshotResponse = {
        affiliateList: [
          affiliateInfoCreateToResponseObject(
            defaultInfo, testConstants.defaultSubaccountUsername.username,
          ),
        ],
        total: 1,
        currentOffset: 1,
      };

      expect(response.body.affiliateList).toHaveLength(1);
      expect(response.body.currentOffset).toEqual(expectedResponse.currentOffset);
      expect(response.body.total).toEqual(expectedResponse.total);
      expect(response.body.affiliateList[0]).toEqual(expectedResponse.affiliateList[0]);
    });

  });

  describe('GET /total_volume', () => {
    beforeEach(async () => {
      await testMocks.seedData();
      await WalletTable.update(
        {
          address: testConstants.defaultWallet.address,
          totalVolume: '100000',
          totalTradingRewards: '0',
        },
      );
    });

    afterEach(async () => {
      await dbHelpers.clearData();
    });

    it('should return total volume for a valid address', async () => {
      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/affiliates/total_volume?address=${testConstants.defaultWallet.address}`,
        expectedStatus: 200, // helper performs expect on status
      });

      expect(response.body).toEqual({
        totalVolume: 100000,
      });
    });

    it('should fail if address does not exist', async () => {
      const nonExistentAddress = 'adgsakhasgt';
      await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/affiliates/metadata?address=${nonExistentAddress}`,
        expectedStatus: 404, // helper performs expect on status
      });
    });
  });
});

function affiliateInfoCreateToResponseObject(
  info: AffiliateInfoCreateObject,
  username: string,
): AffiliateSnapshotResponseObject {
  return {
    affiliateAddress: info.address,
    affiliateReferralCode: username,
    affiliateEarnings: Number(info.affiliateEarnings),
    affiliateReferredTrades:
      Number(info.referredTakerTrades) + Number(info.referredMakerTrades),
    affiliateTotalReferredFees: Number(info.totalReferredFees),
    affiliateReferredUsers: Number(info.totalReferredUsers),
    affiliateReferredNetProtocolEarnings: Number(info.referredNetProtocolEarnings),
    affiliateReferredTotalVolume: Number(info.referredTotalVolume),
  };
}
