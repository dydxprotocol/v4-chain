import { dbHelpers, TurnkeyUserCreateObject, TurnkeyUsersTable } from '@dydxprotocol-indexer/postgres';
import { TurnkeyApiClient } from '@turnkey/sdk-server';
import { TurnkeyController } from '../../../../src/controllers/api/v4/turnkey-controller';
import { generatePrivateKey, privateKeyToAccount } from 'viem/accounts';
import { SigninMethod } from '../../../../src/types';
import { PolicyEngine } from '../../../../src/helpers/policy-engine';
import { extractEmailFromOidcToken } from '../../../../src/lib/turnkey-helpers';

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
      expect(mockPolicyEngine.configurePolicy).toHaveBeenCalledWith(
        user?.suborg_id,
        user?.evm_address,
        user?.dydx_address,
      );
      expect(mockPolicyEngine.removeSelfFromRootQuorum).toHaveBeenCalledWith(
        user?.suborg_id,
      );
    });
  });
});

describe('extractEmailFromOidcToken', () => {
  // Helper to create a mock JWT token
  const createMockJwtToken = (payload: any): string => {
    const header = { alg: 'RS256', kid: '1670273806824' };
    const signature = 'mock-signature';

    const encodedHeader = Buffer.from(JSON.stringify(header)).toString('base64url');
    const encodedPayload = Buffer.from(JSON.stringify(payload)).toString('base64url');
    const encodedSignature = Buffer.from(signature).toString('base64url');

    return `${encodedHeader}.${encodedPayload}.${encodedSignature}`;
  };

  const mockGooglePayload = {
    iss: 'accounts.google.com',
    aud: 'your-client-id.googleusercontent.com',
    sub: '1234567890',
    email: 'user@example.com',
    email_verified: true,
    name: 'John Doe',
    iat: 1670273806,
    exp: 1670277406,
  };

  const mockApplePayload = {
    iss: 'https://appleid.apple.com',
    aud: 'your-client-id',
    sub: '1234567890',
    email: 'user@privaterelay.appleid.com',
    email_verified: true,
    iat: 1670273806,
    exp: 1670277406,
  };

  it('should extract email from valid Google OIDC token', () => {
    const token = createMockJwtToken(mockGooglePayload);
    const extractedEmail = extractEmailFromOidcToken(token, 'google');

    expect(extractedEmail).toBe('user@example.com');
  });

  it('should extract email from valid Apple OIDC token', () => {
    const token = createMockJwtToken(mockApplePayload);
    const extractedEmail = extractEmailFromOidcToken(token, 'apple');

    expect(extractedEmail).toBe('user@privaterelay.appleid.com');
  });

  it('should handle case insensitive provider names', () => {
    const token = createMockJwtToken(mockGooglePayload);

    const upperCaseResult = extractEmailFromOidcToken(token, 'GOOGLE');
    const mixedCaseResult = extractEmailFromOidcToken(token, 'Google');

    expect(upperCaseResult).toBe('user@example.com');
    expect(mixedCaseResult).toBe('user@example.com');
  });

  it('should return undefined for non-Google/Apple providers', () => {
    const token = createMockJwtToken(mockGooglePayload);
    const extractedEmail = extractEmailFromOidcToken(token, 'facebook');

    expect(extractedEmail).toBeUndefined();
  });

  it('should return undefined when email is missing from payload', () => {
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    const { email, ...payloadWithoutEmail } = mockGooglePayload;

    const token = createMockJwtToken(payloadWithoutEmail);
    const extractedEmail = extractEmailFromOidcToken(token, 'google');

    expect(extractedEmail).toBeUndefined();
  });

  it('should return undefined when email is not a string', () => {
    const payloadWithInvalidEmail = { ...mockGooglePayload, email: 123 };

    const token = createMockJwtToken(payloadWithInvalidEmail);
    const extractedEmail = extractEmailFromOidcToken(token, 'google');

    expect(extractedEmail).toBeUndefined();
  });

  it('should handle invalid JWT format (not 3 parts)', () => {
    const invalidToken = 'invalid.token';
    const extractedEmail = extractEmailFromOidcToken(invalidToken, 'google');

    expect(extractedEmail).toBeUndefined();
  });

  it('should handle malformed base64 payload', () => {
    const invalidToken = 'header.invalid-base64-payload.signature';
    const extractedEmail = extractEmailFromOidcToken(invalidToken, 'google');

    expect(extractedEmail).toBeUndefined();
  });

  it('should handle invalid JSON in payload', () => {
    const invalidJsonPayload = Buffer.from('invalid-json').toString('base64url');
    const invalidToken = `header.${invalidJsonPayload}.signature`;

    const extractedEmail = extractEmailFromOidcToken(invalidToken, 'google');

    expect(extractedEmail).toBeUndefined();
  });

  it('should handle empty email string', () => {
    const payloadWithEmptyEmail = { ...mockGooglePayload, email: '' };

    const token = createMockJwtToken(payloadWithEmptyEmail);
    const extractedEmail = extractEmailFromOidcToken(token, 'google');

    expect(extractedEmail).toBeUndefined();
  });

  it('should handle null email', () => {
    const payloadWithNullEmail = { ...mockGooglePayload, email: null };

    const token = createMockJwtToken(payloadWithNullEmail);
    const extractedEmail = extractEmailFromOidcToken(token, 'google');

    expect(extractedEmail).toBeUndefined();
  });

  it('should handle real-world Google token structure', () => {
    const realWorldPayload = {
      iss: 'accounts.google.com',
      azp: 'your-client-id.googleusercontent.com',
      aud: 'your-client-id.googleusercontent.com',
      sub: '123456789012345678901',
      email: 'john.doe@gmail.com',
      email_verified: true,
      at_hash: 'abc123def456',
      name: 'John Doe',
      picture: 'https://lh3.googleusercontent.com/a/default-user=s96-c',
      given_name: 'John',
      family_name: 'Doe',
      locale: 'en',
      iat: 1670273806,
      exp: 1670277406,
      jti: 'abc123def456ghi789',
    };

    const token = createMockJwtToken(realWorldPayload);
    const extractedEmail = extractEmailFromOidcToken(token, 'google');

    expect(extractedEmail).toBe('john.doe@gmail.com');
  });

  it('should handle tokens with special characters in email', () => {
    const payloadWithSpecialEmail = {
      ...mockGooglePayload,
      email: 'user+test@example-domain.co.uk',
    };

    const token = createMockJwtToken(payloadWithSpecialEmail);
    const extractedEmail = extractEmailFromOidcToken(token, 'google');

    expect(extractedEmail).toBe('user+test@example-domain.co.uk');
  });

  it('should handle very long email addresses', () => {
    const longEmail = 'very.long.email.address.with.many.dots@very-long-domain-name.example.com';
    const payloadWithLongEmail = {
      ...mockGooglePayload,
      email: longEmail,
    };

    const token = createMockJwtToken(payloadWithLongEmail);
    const extractedEmail = extractEmailFromOidcToken(token, 'google');

    expect(extractedEmail).toBe(longEmail);
  });
});
