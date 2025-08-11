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
import { sendRequest } from '../../../helpers/helpers';
import { AccountVerificationRequiredAction } from '../../../../src/helpers/compliance/compliance-utils';

import request from 'supertest';
import { ExtendedSecp256k1Signature, Secp256k1 } from '@cosmjs/crypto';
import { verifyADR36Amino } from '@keplr-wallet/cosmos';
import { stats } from '@dydxprotocol-indexer/base';
import { DateTime } from 'luxon';
import { toBech32 } from '@cosmjs/encoding';

jest.mock('@cosmjs/crypto', () => ({
  ...jest.requireActual('@cosmjs/crypto'),
  Secp256k1: {
    verifySignature: jest.fn(),
  },
  ExtendedSecp256k1Signature: {
    fromFixedLength: jest.fn(),
  },
}));

jest.mock('@keplr-wallet/cosmos', () => ({
  ...jest.requireActual('@keplr-wallet/cosmos'),
  verifyADR36Amino: jest.fn(),
}));

jest.mock('@cosmjs/encoding', () => ({
  toBech32: jest.fn(),
}));

describe('affiliates-controller#V4', () => {

  const verifySignatureMock = Secp256k1.verifySignature as jest.Mock;
  const fromFixedLengthMock = ExtendedSecp256k1Signature.fromFixedLength as jest.Mock;
  const verifyADR36AminoMock = verifyADR36Amino as jest.Mock;
  const toBech32Mock = toBech32 as jest.Mock;
  beforeAll(async () => {
    await dbHelpers.migrate();
  });

  afterAll(async () => {
    await dbHelpers.teardown();
  });

  describe('GET /metadata', () => {
    beforeEach(async () => {
      await testMocks.seedData();
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

  describe('POST /referralCode', () => {
    const defaultSubaccountId = SubaccountTable.uuid(
      testConstants.defaultSubaccount.address,
      testConstants.defaultSubaccount.subaccountNumber,
    );

    beforeEach(async () => {
      await testMocks.seedData();
      verifySignatureMock.mockResolvedValue(true);
      fromFixedLengthMock.mockResolvedValue({} as ExtendedSecp256k1Signature);
      verifyADR36AminoMock.mockReturnValue(true);
      toBech32Mock.mockReturnValue(testConstants.defaultAddress);
      jest.spyOn(DateTime, 'now').mockReturnValue(DateTime.fromSeconds(1620000000)); // Mock current time
      jest.spyOn(stats, 'increment');
    });

    afterEach(async () => {
      await dbHelpers.clearData();
    });

    const mockCreateCodeRequest = (code: string, address?: string) => ({
      address: address || testConstants.defaultWallet.address,
      newCode: code,
      action: AccountVerificationRequiredAction.UPDATE_CODE,
      signedMessage: 'signedMessage',
      pubKey: address || testConstants.defaultWallet.address,
      timestamp: DateTime.now().toSeconds(),
    });

    it('should update a referral code for a valid address if it already exists', async () => {
      const newCode = 'NewCode12345';
      const response2: request.Response = await sendRequest({
        type: RequestMethod.POST,
        path: '/v4/affiliates/referralCode',
        body: mockCreateCodeRequest(newCode),
        expectedStatus: 200,
      });

      expect(response2.body).toEqual({
        referralCode: newCode,
      });

      // query the database to check if the referral code was updated
      const usernameRow = await SubaccountUsernamesTable.findByUsername(newCode);
      expect(usernameRow).toEqual({
        username: newCode,
        subaccountId: defaultSubaccountId,
      });

      const usernameRowByAddress = await SubaccountUsernamesTable.findByAddress([
        testConstants.defaultWallet.address,
      ]);
      expect(usernameRowByAddress).toEqual([{
        username: newCode,
        address: testConstants.defaultWallet.address,
      }]);
    });

    it('should update a referral code for a valid address for keplr if it already exists', async () => {
      const newCode = 'NewCode12345';
      const newRequest = mockCreateCodeRequest(newCode);
      const response2: request.Response = await sendRequest({
        type: RequestMethod.POST,
        path: '/v4/affiliates/referralCode-keplr',
        body: {
          ...newRequest,
          // these two are not used in the keplr version of the request
          timestamp: undefined,
          action: undefined,
        },
        expectedStatus: 200,
      });

      expect(response2.body).toEqual({
        referralCode: newCode,
      });

      // query the database to check if the referral code was updated
      const usernameRow = await SubaccountUsernamesTable.findByUsername(newCode);
      expect(usernameRow).toEqual({
        username: newCode,
        subaccountId: defaultSubaccountId,
      });

      const usernameRowByAddress = await SubaccountUsernamesTable.findByAddress([
        testConstants.defaultWallet.address,
      ]);
      expect(usernameRowByAddress).toEqual([{
        username: newCode,
        address: testConstants.defaultWallet.address,
      }]);
    });

    it('should fail to create a referral code for an invalid address', async () => {
      const invalidAddress = 'invalidAddress';
      const newCode = 'NewCode123';
      // mock the signature verification to return true
      verifySignatureMock.mockResolvedValue(true);
      fromFixedLengthMock.mockResolvedValue({} as ExtendedSecp256k1Signature);
      verifyADR36AminoMock.mockReturnValue(true);
      jest.spyOn(DateTime, 'now').mockReturnValue(DateTime.fromSeconds(1620000000)); // Mock current time
      jest.spyOn(stats, 'increment');
      await sendRequest({
        type: RequestMethod.POST,
        path: '/v4/affiliates/referralCode',
        body: mockCreateCodeRequest(newCode, invalidAddress),
        expectedStatus: 400,  // helper performs expect on status
      });
    });

    it('should fail to update an existing referral code if another user has it', async () => {
      const newCode = 'NewCode123';
      await sendRequest({
        type: RequestMethod.POST,
        path: '/v4/affiliates/referralCode',
        body: mockCreateCodeRequest(newCode, testConstants.defaultWallet.address),
        expectedStatus: 200,
      });

      toBech32Mock.mockReturnValue(testConstants.defaultWallet2.address);
      await sendRequest({
        type: RequestMethod.POST,
        path: '/v4/affiliates/referralCode',
        body: mockCreateCodeRequest(newCode, testConstants.defaultWallet2.address),
        expectedStatus: 400,
        errorMsg: 'Referral code already exists',
      });
    });

    it('should fail for invalid codes and succeed for valid codes', async () => {
      const validCodes = [
        '1234567890',
        'foobar123',
        'foobar3319',
      ];
      const invalidCodes = [
        '1',
        '1234567890123456789012345678901234567890',
        'foobar*123',
      ];
      for (const code of validCodes) {
        await sendRequest({
          type: RequestMethod.POST,
          path: '/v4/affiliates/referralCode',
          body: mockCreateCodeRequest(code, testConstants.defaultWallet.address),
          expectedStatus: 200,
        });
      }
      for (const code of invalidCodes) {
        await sendRequest({
          type: RequestMethod.POST,
          path: '/v4/affiliates/referralCode',
          body: mockCreateCodeRequest(code, testConstants.defaultWallet.address),
          expectedStatus: 400,
        });
      }
    });
  });

  describe('GET /snapshot', () => {
    const defaultInfo: AffiliateInfoCreateObject = testConstants.defaultAffiliateInfo;
    const defaultInfo2: AffiliateInfoCreateObject = testConstants.defaultAffiliateInfo2;
    const defaultInfo3: AffiliateInfoCreateObject = testConstants.defaultAffiliateInfo3;

    beforeEach(async () => {
      await testMocks.seedData();
      // Create defaultWallet2, subaccount, and username
      await WalletTable.create(testConstants.defaultWallet2);
      await SubaccountTable.create(testConstants.defaultSubaccountDefaultWalletAddress);

      // Create defaultWallet3, create subaccount, create username
      await WalletTable.create(testConstants.defaultWallet3);
      await SubaccountTable.create(testConstants.defaultSubaccountWithAlternateAddress);

      // Create affiliate infos
      await Promise.all([
        AffiliateInfoTable.create(defaultInfo),
        AffiliateInfoTable.create(defaultInfo2),
        AffiliateInfoTable.create(defaultInfo3),
      ]);
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

    it('should handle no results', async () => {
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
    affiliateTotalReferredFees: Number(info.totalReferredMakerFees) +
    Number(info.totalReferredTakerFees) +
    Number(info.totalReferredMakerRebates),
    affiliateReferredUsers: Number(info.totalReferredUsers),
    affiliateReferredNetProtocolEarnings: Number(info.totalReferredMakerFees) +
    Number(info.totalReferredTakerFees) +
    Number(info.totalReferredMakerRebates) -
    Number(info.affiliateEarnings),
    affiliateReferredTotalVolume: Number(info.referredTotalVolume),
    affiliateReferredMakerFees: Number(info.totalReferredMakerFees),
    affiliateReferredTakerFees: Number(info.totalReferredTakerFees),
    affiliateReferredMakerRebates: Number(info.totalReferredMakerRebates),
  };
}
