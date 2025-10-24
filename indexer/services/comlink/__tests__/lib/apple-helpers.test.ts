import { AppleHelpers } from '../../src/lib/apple-helpers';
import { TurnkeyError } from '../../src/lib/errors';
import { AppleTokenResponse } from '../../src/types';
import { SignJWT, importPKCS8 } from 'jose';
import fetch from 'node-fetch';

// Mock dependencies
jest.mock('jose');
jest.mock('node-fetch');
jest.mock('@dydxprotocol-indexer/base', () => ({
  logger: {
    error: jest.fn(),
    warning: jest.fn(),
  },
}));

const mockSignJWT = SignJWT as jest.MockedClass<typeof SignJWT>;
const mockImportPKCS8 = importPKCS8 as jest.MockedFunction<typeof importPKCS8>;
const mockFetch = fetch as jest.MockedFunction<typeof fetch>;

describe('AppleHelpers', () => {
  const mockTeamId = 'TEAM123';
  const mockServiceId = 'com.example.app';
  const mockKeyId = 'KEY123';
  const mockPrivateKey = `-----BEGIN PRIVATE KEY-----
MIGHAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBG0wawIBAQQg...
-----END PRIVATE KEY-----`;
  const mockCode = 'auth_code_123';
  const mockKeyLike = { kty: 'EC' } as any;

  beforeEach(() => {
    jest.clearAllMocks();
  });

  describe('generateClientSecret', () => {
    it('should generate a valid JWT client secret', async () => {
      const mockJwt = 'mock.jwt.token';
      const mockSignJWTInstance = {
        setProtectedHeader: jest.fn().mockReturnThis(),
        sign: jest.fn().mockResolvedValue(mockJwt),
      };

      mockSignJWT.mockImplementation(() => mockSignJWTInstance as any);
      mockImportPKCS8.mockResolvedValue(mockKeyLike);

      const result = await AppleHelpers.generateClientSecret(
        mockTeamId,
        mockServiceId,
        mockKeyId,
        mockPrivateKey,
      );

      expect(result).toBe(mockJwt);
      expect(mockImportPKCS8).toHaveBeenCalledWith(mockPrivateKey, 'ES256');
      expect(mockSignJWT).toHaveBeenCalledWith({
        iss: mockTeamId,
        iat: expect.any(Number),
        exp: expect.any(Number),
        aud: 'https://appleid.apple.com',
        sub: mockServiceId,
      });
      expect(mockSignJWTInstance.setProtectedHeader).toHaveBeenCalledWith({
        alg: 'ES256',
        kid: mockKeyId,
      });
      expect(mockSignJWTInstance.sign).toHaveBeenCalledWith(mockKeyLike);
    });

    it('should set correct expiration time (6 months)', async () => {
      const mockJwt = 'mock.jwt.token';
      const mockSignJWTInstance = {
        setProtectedHeader: jest.fn().mockReturnThis(),
        sign: jest.fn().mockResolvedValue(mockJwt),
      };

      mockSignJWT.mockImplementation(() => mockSignJWTInstance as any);
      mockImportPKCS8.mockResolvedValue(mockKeyLike);

      const now = Math.floor(Date.now() / 1000);

      await AppleHelpers.generateClientSecret(
        mockTeamId,
        mockServiceId,
        mockKeyId,
        mockPrivateKey,
      );

      const expectedExp = now + (60 * 60 * 24 * 180); // 6 months
      expect(mockSignJWT).toHaveBeenCalledWith({
        iss: mockTeamId,
        iat: now,
        exp: expectedExp,
        aud: 'https://appleid.apple.com',
        sub: mockServiceId,
      });
    });

    it('should throw TurnkeyError when private key parsing fails', async () => {
      const errorMessage = 'Invalid private key format';
      mockImportPKCS8.mockRejectedValue(new Error(errorMessage));

      await expect(
        AppleHelpers.generateClientSecret(
          mockTeamId,
          mockServiceId,
          mockKeyId,
          mockPrivateKey,
        ),
      ).rejects.toThrow(TurnkeyError);

      await expect(
        AppleHelpers.generateClientSecret(
          mockTeamId,
          mockServiceId,
          mockKeyId,
          mockPrivateKey,
        ),
      ).rejects.toThrow(`Failed to generate Apple client secret: Failed to parse Apple private key: ${errorMessage}`);
    });

    it('should throw TurnkeyError when JWT signing fails', async () => {
      const errorMessage = 'JWT signing failed';
      const mockSignJWTInstance = {
        setProtectedHeader: jest.fn().mockReturnThis(),
        sign: jest.fn().mockRejectedValue(new Error(errorMessage)),
      };

      mockSignJWT.mockImplementation(() => mockSignJWTInstance as any);
      mockImportPKCS8.mockResolvedValue(mockKeyLike);

      await expect(
        AppleHelpers.generateClientSecret(
          mockTeamId,
          mockServiceId,
          mockKeyId,
          mockPrivateKey,
        ),
      ).rejects.toThrow(TurnkeyError);

      await expect(
        AppleHelpers.generateClientSecret(
          mockTeamId,
          mockServiceId,
          mockKeyId,
          mockPrivateKey,
        ),
      ).rejects.toThrow(`Failed to generate Apple client secret: ${errorMessage}`);
    });
  });

  describe('fetchTokenFromCode', () => {
    const mockTokenResponse: AppleTokenResponse = {
      access_token: 'access_token_123',
      token_type: 'Bearer',
      expires_in: 3600,
      refresh_token: 'refresh_token_123',
      id_token: 'id_token_123',
    };

    it('should successfully exchange code for token', async () => {
      const mockJwt = 'mock.jwt.token';
      const mockSignJWTInstance = {
        setProtectedHeader: jest.fn().mockReturnThis(),
        sign: jest.fn().mockResolvedValue(mockJwt),
      };

      mockSignJWT.mockImplementation(() => mockSignJWTInstance as any);
      mockImportPKCS8.mockResolvedValue(mockKeyLike);

      const mockResponse = {
        ok: true,
        json: jest.fn().mockResolvedValue(mockTokenResponse),
      };
      mockFetch.mockResolvedValue(mockResponse as any);

      const result = await AppleHelpers.fetchTokenFromCode(
        mockCode,
        mockTeamId,
        mockServiceId,
        mockKeyId,
        mockPrivateKey,
      );

      expect(result).toEqual(mockTokenResponse);
      expect(mockFetch).toHaveBeenCalledWith('https://appleid.apple.com/auth/token', {
        method: 'POST',
        headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
        body: expect.stringContaining('client_id=com.example.app'),
      });
    });

    it('should include correct form parameters in request', async () => {
      const mockJwt = 'mock.jwt.token';
      const mockSignJWTInstance = {
        setProtectedHeader: jest.fn().mockReturnThis(),
        sign: jest.fn().mockResolvedValue(mockJwt),
      };

      mockSignJWT.mockImplementation(() => mockSignJWTInstance as any);
      mockImportPKCS8.mockResolvedValue(mockKeyLike);

      const mockResponse = {
        ok: true,
        json: jest.fn().mockResolvedValue(mockTokenResponse),
      };
      mockFetch.mockResolvedValue(mockResponse as any);

      await AppleHelpers.fetchTokenFromCode(
        mockCode,
        mockTeamId,
        mockServiceId,
        mockKeyId,
        mockPrivateKey,
      );

      expect(mockFetch).toHaveBeenCalledWith(
        'https://appleid.apple.com/auth/token',
        expect.objectContaining({
          method: 'POST',
          headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
          body: expect.stringMatching(/client_id=com\.example\.app&client_secret=mock\.jwt\.token&code=auth_code_123&grant_type=authorization_code/),
        }),
      );
    });

    it('should throw TurnkeyError when Apple API returns error', async () => {
      const mockJwt = 'mock.jwt.token';
      const mockSignJWTInstance = {
        setProtectedHeader: jest.fn().mockReturnThis(),
        sign: jest.fn().mockResolvedValue(mockJwt),
      };

      mockSignJWT.mockImplementation(() => mockSignJWTInstance as any);
      mockImportPKCS8.mockResolvedValue(mockKeyLike);

      const errorText = 'invalid_grant';
      const mockResponse = {
        ok: false,
        status: 400,
        text: jest.fn().mockResolvedValue(errorText),
      };
      mockFetch.mockResolvedValue(mockResponse as any);

      await expect(
        AppleHelpers.fetchTokenFromCode(
          mockCode,
          mockTeamId,
          mockServiceId,
          mockKeyId,
          mockPrivateKey,
        ),
      ).rejects.toThrow(TurnkeyError);

      await expect(
        AppleHelpers.fetchTokenFromCode(
          mockCode,
          mockTeamId,
          mockServiceId,
          mockKeyId,
          mockPrivateKey,
        ),
      ).rejects.toThrow(`Apple token exchange failed: 400 ${errorText}`);
    });

    it('should throw TurnkeyError when response has no id_token', async () => {
      const mockJwt = 'mock.jwt.token';
      const mockSignJWTInstance = {
        setProtectedHeader: jest.fn().mockReturnThis(),
        sign: jest.fn().mockResolvedValue(mockJwt),
      };

      mockSignJWT.mockImplementation(() => mockSignJWTInstance as any);
      mockImportPKCS8.mockResolvedValue(mockKeyLike);

      const responseWithoutIdToken = {
        access_token: 'access_token_123',
        token_type: 'Bearer',
        expires_in: 3600,
      };

      const mockResponse = {
        ok: true,
        json: jest.fn().mockResolvedValue(responseWithoutIdToken),
      };
      mockFetch.mockResolvedValue(mockResponse as any);

      await expect(
        AppleHelpers.fetchTokenFromCode(
          mockCode,
          mockTeamId,
          mockServiceId,
          mockKeyId,
          mockPrivateKey,
        ),
      ).rejects.toThrow(TurnkeyError);

      await expect(
        AppleHelpers.fetchTokenFromCode(
          mockCode,
          mockTeamId,
          mockServiceId,
          mockKeyId,
          mockPrivateKey,
        ),
      ).rejects.toThrow('No ID token received from Apple');
    });

    it('should throw TurnkeyError when fetch fails', async () => {
      const errorMessage = 'Network error';
      mockImportPKCS8.mockRejectedValue(new Error(errorMessage));

      await expect(
        AppleHelpers.fetchTokenFromCode(
          mockCode,
          mockTeamId,
          mockServiceId,
          mockKeyId,
          mockPrivateKey,
        ),
      ).rejects.toThrow(TurnkeyError);

      await expect(
        AppleHelpers.fetchTokenFromCode(
          mockCode,
          mockTeamId,
          mockServiceId,
          mockKeyId,
          mockPrivateKey,
        ),
      ).rejects.toThrow(`Failed to fetch Apple token: Failed to generate Apple client secret: Failed to parse Apple private key: ${errorMessage}`);
    });
  });

  describe('parsePrivateKey', () => {
    it('should successfully parse valid private key', async () => {
      mockImportPKCS8.mockResolvedValue(mockKeyLike);

      const result = await AppleHelpers.parsePrivateKey(mockPrivateKey);

      expect(result).toBe(mockKeyLike);
      expect(mockImportPKCS8).toHaveBeenCalledWith(mockPrivateKey, 'ES256');
    });

    it('should throw TurnkeyError when private key is invalid', async () => {
      const errorMessage = 'Invalid private key format';
      mockImportPKCS8.mockRejectedValue(new Error(errorMessage));

      await expect(
        AppleHelpers.parsePrivateKey('invalid_key'),
      ).rejects.toThrow(TurnkeyError);

      await expect(
        AppleHelpers.parsePrivateKey('invalid_key'),
      ).rejects.toThrow(`Failed to parse Apple private key: ${errorMessage}`);
    });

    it('should throw TurnkeyError when private key parsing throws non-Error', async () => {
      const errorMessage = 'Unknown error';
      mockImportPKCS8.mockRejectedValue(errorMessage);

      await expect(
        AppleHelpers.parsePrivateKey('invalid_key'),
      ).rejects.toThrow(TurnkeyError);

      await expect(
        AppleHelpers.parsePrivateKey('invalid_key'),
      ).rejects.toThrow(`Failed to parse Apple private key: ${errorMessage}`);
    });
  });

});
