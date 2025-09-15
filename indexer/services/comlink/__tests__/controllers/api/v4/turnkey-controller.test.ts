import { dbHelpers, TurnkeyUserCreateObject, TurnkeyUsersTable } from '@dydxprotocol-indexer/postgres';
import { TurnkeyApiClient } from '@turnkey/sdk-server';
import { TurnkeyController } from '../../../../src/controllers/api/v4/turnkey-controller';
import { generatePrivateKey, privateKeyToAccount } from 'viem/accounts';
import { SigninMethod } from '../../../../src/types';
import { PolicyEngine } from '../../../../src/helpers/policy-engine';

jest.mock('../../../../src/config', () => ({
  ...jest.requireActual('../../../../src/config'),
  MASTER_SIGNER_PUBLIC: '0x1234567890123456789012345678901234567890',
  INDEXER_INTERNAL_IPS: '127.0.0.1',
}));

describe('TurnkeyController', () => {
  let mockParentApiClient: TurnkeyApiClient;
  let mockBridgeSenderApiClient: TurnkeyApiClient;
  let mockPolicyEngine: jest.Mocked<PolicyEngine>;
  let controller: TurnkeyController;
  const mockUser: TurnkeyUserCreateObject = {
    suborg_id: 'test-suborg-id',
    email: 'test@example.com',
    salt: 'test-salt',
    created_at: new Date().toISOString(),
    evm_address: '0x1234567890123456789012345678901234567890',
    svm_address: 'svm1234567890123456789012345678901234567890',
    dydx_address: 'dydx1234567890123456789012345678901234567890',
  };

  const testPrivateKey = generatePrivateKey();
  const generatedEvmWallet = privateKeyToAccount(testPrivateKey);
  const mockUser2: TurnkeyUserCreateObject = {
    suborg_id: 'test-suborg-id-3',
    email: 'test3@example.com',
    salt: 'test-salt',
    created_at: new Date().toISOString(),
    evm_address: generatedEvmWallet.address,
    svm_address: 'svm1234567890123456789012345678901234567891',
  };

  beforeAll(async () => {
    await dbHelpers.clearData();
    await dbHelpers.migrate();
    mockParentApiClient = {
      createSubOrganization: jest.fn(),
      emailAuth: jest.fn(),
      oauthLogin: jest.fn(),
      getSubOrgIds: jest.fn(),
      getUser: jest.fn(),
      getUsers: jest.fn(),
      updateRootQuorum: jest.fn(),
    } as unknown as TurnkeyApiClient;
    mockBridgeSenderApiClient = {
      createSubOrganization: jest.fn(),
      emailAuth: jest.fn(),
      oauthLogin: jest.fn(),
      getSubOrgIds: jest.fn(),
      getUser: jest.fn(),
      getUsers: jest.fn(),
      updateRootQuorum: jest.fn().mockResolvedValue({}),
      createPolicy: jest.fn().mockResolvedValue({}),
    } as unknown as TurnkeyApiClient;

    // Create mock PolicyEngine
    mockPolicyEngine = {
      configurePolicy: jest.fn().mockResolvedValue(undefined),
      configureSolanaPolicy: jest.fn().mockResolvedValue(undefined),
      removeSelfFromRootQuorum: jest.fn().mockResolvedValue(undefined),
      getAPIUserId: jest.fn().mockResolvedValue('mock-api-user-id'),
    } as unknown as jest.Mocked<PolicyEngine>;

    controller = new TurnkeyController(mockParentApiClient, mockBridgeSenderApiClient);

    // Replace the private policyEngine property with our mock
    (controller as any).policyEngine = mockPolicyEngine;

    await TurnkeyUsersTable.create(mockUser);
    await TurnkeyUsersTable.create(mockUser2);
  });

  afterAll(async () => {
    // No cleanup needed for these tests
    await dbHelpers.clearData();
    await dbHelpers.teardown();
  });

  describe('POST /signin', () => {
    describe('Email signin', () => {
      it('should successfully sign in existing user with email', async () => {

        jest.mocked(mockParentApiClient.emailAuth).mockResolvedValue({
          activity: {
            result: {
              emailAuthResult: {
                apiKeyId: 'api-key-id',
                userId: 'user-id',
              },
            },
          },
        } as any);

        const response = await controller.signIn({
          signinMethod: SigninMethod.EMAIL,
          userEmail: 'test@example.com',
          targetPublicKey: 'target-public-key',
        });

        expect(response.apiKeyId).toEqual('api-key-id');
        expect(response.userId).toEqual('user-id');
        expect(response.salt).toEqual('test-salt');
      });

      it('should create new user and sign in with email', async () => {
        // No existing user in database
        jest.mocked(mockParentApiClient.createSubOrganization).mockResolvedValue({
          subOrganizationId: 'test-suborg-id-2',
          rootUserIds: ['user-id', 'user-id-2'],
          wallet: {
            addresses: ['0x123', 'svm-address'],
          },
        } as any);
        jest.mocked(mockParentApiClient.emailAuth).mockResolvedValue({
          activity: {
            result: {
              emailAuthResult: {
                apiKeyId: 'api-key-id',
                userId: 'user-id',
              },
            },
          },
        } as any);
        jest.mocked(mockBridgeSenderApiClient.getUser).mockResolvedValue({
          user: { userId: 'user-id', userEmail: 'test2@example.com', authenticators: [] },
        } as any);

        const response = await controller.signIn({
          signinMethod: SigninMethod.EMAIL,
          userEmail: 'test2@example.com',
          targetPublicKey: 'target-public-key',
        });

        expect(jest.mocked(mockParentApiClient.createSubOrganization)).toHaveBeenCalled();

        // Verify user was created in database
        const createdUser = await TurnkeyUsersTable.findByEmail('test2@example.com');
        expect(createdUser).toBeDefined();
        expect(createdUser?.email).toEqual('test2@example.com');
        expect(createdUser?.suborg_id).toEqual('test-suborg-id-2');
        expect(createdUser?.evm_address).toEqual('0x123');
        expect(createdUser?.svm_address).toEqual('svm-address');
        expect(createdUser?.salt).toBeDefined();
        expect(response?.apiKeyId).toEqual('api-key-id');
        expect(response?.userId).toEqual('user-id');
        expect(response?.organizationId).toEqual('test-suborg-id-2');
      });

      it('should throw error for invalid email format', async () => {
        await expect(controller.signIn({
          signinMethod: SigninMethod.EMAIL,
          userEmail: 'invalid-email',
          targetPublicKey: 'target-public-key',
        })).rejects.toThrow();
      });

      it('should return the dydx address if it exists', async () => {
        jest.mocked(mockParentApiClient.emailAuth).mockResolvedValue({
          activity: {
            result: {
              emailAuthResult: {
                apiKeyId: 'api-key-id',
                userId: 'user-id',
              },
            },
          },
        } as any);
        const response = await controller.signIn({
          signinMethod: SigninMethod.EMAIL,
          userEmail: 'test@example.com',
          targetPublicKey: 'target-public-key',
        });

        expect(response?.dydxAddress).toEqual('dydx1234567890123456789012345678901234567890');
      });

      it('should throw error when required fields are missing', async () => {
        await expect(controller.signIn({
          signinMethod: SigninMethod.EMAIL,
          userEmail: 'test@example.com',
        })).rejects.toThrow();
      });
    });

    describe('Social signin', () => {
      it('should successfully sign in existing user with social', async () => {
        jest.mocked(mockParentApiClient.getSubOrgIds).mockResolvedValueOnce({
          organizationIds: ['test-suborg-id'],
        } as any);
        jest.mocked(mockParentApiClient.oauthLogin).mockResolvedValue({
          activity: {
            result: {
              oauthLoginResult: {
                session: 'session-token',
              },
            },
          },
        } as any);
        const response = await controller.signIn({
          signinMethod: SigninMethod.SOCIAL,
          targetPublicKey: 'target-public-key',
          provider: 'google',
          oidcToken: 'oidc-token',
        });

        expect(response?.session).toEqual('session-token');
      });

      it('should return the dydx address if it exists', async () => {
        jest.mocked(mockParentApiClient.getSubOrgIds).mockResolvedValueOnce({
          organizationIds: ['test-suborg-id'],
        } as any);
        jest.mocked(mockParentApiClient.oauthLogin).mockResolvedValue({
          activity: {
            result: {
              oauthLoginResult: {
                session: 'session-token',
              },
            },
          },
        } as any);
        const response = await controller.signIn({
          signinMethod: SigninMethod.SOCIAL,
          targetPublicKey: 'target-public-key',
          provider: 'google',
          oidcToken: 'oidc-token',
        });

        expect(response?.dydxAddress).toEqual('dydx1234567890123456789012345678901234567890');
      });

      it('should throw error when required fields are missing', async () => {
        await expect(controller.signIn({
          signinMethod: SigninMethod.SOCIAL,
          targetPublicKey: 'target-public-key',
          provider: 'google',
        })).rejects.toThrow();
      });
    });

    describe('Passkey signin', () => {
      const mockAttestation = {
        credentialId: 'credential-id',
        clientDataJson: 'client-data',
        attestationObject: 'attestation-object',
        transports: ['AUTHENTICATOR_TRANSPORT_USB'],
      };

      it('should create new user', async () => {
        jest.mocked(mockParentApiClient.createSubOrganization).mockResolvedValueOnce({
          subOrganizationId: 'test-suborg-id-2',
          rootUserIds: ['user-id', 'user-id-2'],
          wallet: {
            addresses: ['0x1232', 'svm-address2'],
          },
        } as any);
        jest.mocked(mockParentApiClient.getSubOrgIds).mockResolvedValueOnce({
          organizationIds: ['test-suborg-id-2'],
        } as any);
        jest.mocked(mockBridgeSenderApiClient.getUser).mockResolvedValue({
          user: { userId: 'user-id-2', authenticators: [{ credentialId: 'credential-id' }] },
        } as any);

        const response = await controller.signIn({
          signinMethod: SigninMethod.PASSKEY,
          challenge: 'challenge',
          attestation: mockAttestation as any,
        });

        expect(response?.organizationId).toEqual('test-suborg-id-2');
        expect(response?.salt).toBeDefined();
      });

      it('should successfully sign in existing user with passkey', async () => {
        jest.mocked(mockParentApiClient.getSubOrgIds).mockResolvedValueOnce({
          organizationIds: ['test-suborg-id'],
        } as any);

        const response = await controller.signIn({
          signinMethod: SigninMethod.PASSKEY,
          challenge: 'challenge',
          attestation: mockAttestation as any,
        });

        expect(response?.organizationId).toEqual('test-suborg-id');
      });

      it('should throw error when required fields are missing', async () => {
        await expect(controller.signIn({
          signinMethod: SigninMethod.PASSKEY,
          challenge: 'challenge',
          attestation: undefined,
        })).rejects.toThrow();
      });

      it('should return the dydx address if it exists', async () => {
        jest.mocked(mockParentApiClient.getSubOrgIds).mockResolvedValueOnce({
          organizationIds: ['test-suborg-id'],
        } as any);
        // await TurnkeyUsersTable.create(mockUser);

        const response = await controller.signIn({
          signinMethod: SigninMethod.PASSKEY,
          challenge: 'challenge',
          attestation: mockAttestation as any,
        });

        expect(response?.dydxAddress).toEqual('dydx1234567890123456789012345678901234567890');
      });
    });
  });

  describe('POST /uploadDydxAddress', () => {
    it('should upload the dydx address', async () => {
      const newDydxAddress = 'dydx1234567890123456789012345678901234567891';
      const signature = await generatedEvmWallet.signMessage({ message: newDydxAddress });
      jest.mocked(mockBridgeSenderApiClient.getUsers).mockResolvedValue({
        users: [
          { userId: 'user-id-2', userName: 'API User' },
        ],
      } as any);
      const response = await controller.uploadAddress({
        dydxAddress: newDydxAddress,
        signature,
      });

      expect(response).toEqual({ success: true });

      // verify the dydx address is updated in the database
      const user = await TurnkeyUsersTable.findByEmail('test3@example.com');
      expect(user?.dydx_address).toEqual('dydx1234567890123456789012345678901234567891');

      // verify that PolicyEngine methods were called
      expect(mockPolicyEngine.configureSolanaPolicy).toHaveBeenCalledWith(
        newDydxAddress,
        user?.suborg_id,
      );
      expect(mockPolicyEngine.configurePolicy).toHaveBeenCalledWith(
        user?.suborg_id,
        user?.evm_address,
      );
      // expect(mockPolicyEngine.removeSelfFromRootQuorum).toHaveBeenCalledWith(
      //   user?.suborg_id,
      // );
    });
  });
});
