// Mock the amplitude SDK
jest.mock('@amplitude/analytics-node', () => ({
  init: jest.fn(),
  track: jest.fn(),
}));

// Mock the logger
jest.mock('@dydxprotocol-indexer/base', () => ({
  logger: {
    info: jest.fn(),
    debug: jest.fn(),
    error: jest.fn(),
  },
}));

// Mock the config
jest.mock('../../src/config', () => ({
  AMPLITUDE_API_KEY: 'test-api-key',
  AMPLITUDE_SERVER_URL: 'https://api.eu.amplitude.com/2/httpapi',
}));

describe('amplitude-helpers', () => {
  let mockInit: jest.MockedFunction<any>;
  let mockTrack: jest.MockedFunction<any>;
  let mockLogger: jest.Mocked<any>;

  beforeEach(() => {
    jest.clearAllMocks();
    // Reset module state by clearing require cache
    jest.resetModules();

    // Re-define mocks after resetModules
    mockInit = jest.mocked(require('@amplitude/analytics-node').init);
    mockTrack = jest.mocked(require('@amplitude/analytics-node').track);
    mockLogger = jest.mocked(require('@dydxprotocol-indexer/base').logger);
  });

  describe('initializeAmplitude', () => {
    it('should initialize amplitude with correct configuration when API key is provided', () => {
      const { initializeAmplitude } = require('../../src/lib/amplitude-helpers');

      initializeAmplitude();

      expect(mockInit).toHaveBeenCalledWith('test-api-key', {
        serverUrl: 'https://api.eu.amplitude.com/2/httpapi',
        flushQueueSize: 1,
        flushIntervalMillis: 1000,
        serverZone: 'EU',
      });

      expect(mockLogger.info).toHaveBeenCalledWith({
        at: 'amplitude-helpers#initializeAmplitude',
        message: 'Amplitude client initialized successfully',
      });
    });

    it('should not reinitialize if already initialized', () => {
      const { initializeAmplitude } = require('../../src/lib/amplitude-helpers');

      // Call twice
      initializeAmplitude();
      initializeAmplitude();

      // Should only be called once
      expect(mockInit).toHaveBeenCalledTimes(1);
    });

    it('should handle initialization errors gracefully', () => {
      const error = new Error('Initialization failed');
      mockInit.mockImplementation(() => {
        throw error;
      });

      const { initializeAmplitude } = require('../../src/lib/amplitude-helpers');

      initializeAmplitude();

      expect(mockLogger.error).toHaveBeenCalledWith({
        at: 'amplitude-helpers#initializeAmplitude',
        message: 'Failed to initialize Amplitude client',
        error,
      });
    });
  });

  describe('trackAmplitudeEvent', () => {
    it('should track event successfully when amplitude is initialized', async () => {
      const { trackAmplitudeEvent, initializeAmplitude } = require('../../src/lib/amplitude-helpers');

      // Initialize first
      initializeAmplitude();

      const eventType = 'TestEvent';
      const userId = 'test-user-123';
      const eventProperties = { testProp: 'testValue' };

      await trackAmplitudeEvent(eventType, userId, eventProperties);

      expect(mockTrack).toHaveBeenCalledWith(eventType, eventProperties, {
        user_id: userId,
      });

      expect(mockLogger.debug).toHaveBeenCalledWith({
        at: 'amplitude-helpers#trackAmplitudeEvent',
        message: 'Amplitude event tracked successfully',
        eventType,
        userId,
        eventProperties,
      });
    });

    it('should track event without userId when not provided', async () => {
      const { trackAmplitudeEvent, initializeAmplitude } = require('../../src/lib/amplitude-helpers');

      // Initialize first
      initializeAmplitude();

      const eventType = 'TestEvent';
      const eventProperties = { testProp: 'testValue' };

      await trackAmplitudeEvent(eventType, undefined, eventProperties);

      expect(mockTrack).toHaveBeenCalledWith(eventType, eventProperties, {
        user_id: undefined,
      });
    });

    it('should track event without eventProperties when not provided', async () => {
      const { trackAmplitudeEvent, initializeAmplitude } = require('../../src/lib/amplitude-helpers');

      // Initialize first
      initializeAmplitude();

      const eventType = 'TestEvent';
      const userId = 'test-user-123';

      await trackAmplitudeEvent(eventType, userId);

      expect(mockTrack).toHaveBeenCalledWith(eventType, undefined, {
        user_id: userId,
      });
    });

    it('should handle tracking errors gracefully', async () => {
      const error = new Error('Tracking failed');
      (mockTrack as jest.Mock).mockRejectedValueOnce(error);

      const { trackAmplitudeEvent, initializeAmplitude } = require('../../src/lib/amplitude-helpers');

      // Initialize first
      initializeAmplitude();

      const eventType = 'TestEvent';
      const userId = 'test-user-123';
      const eventProperties = { testProp: 'testValue' };

      await trackAmplitudeEvent(eventType, userId, eventProperties);

      expect(mockLogger.error).toHaveBeenCalledWith({
        at: 'amplitude-helpers#trackAmplitudeEvent',
        message: 'Failed to track Amplitude event',
        error,
        eventType,
        userId,
        eventProperties,
      });
    });
  });

  describe('trackTurnkeyDepositSubmitted', () => {
    it('should track TurnKey deposit event with all parameters', async () => {
      const { trackTurnkeyDepositSubmitted, initializeAmplitude } = require('../../src/lib/amplitude-helpers');

      // Initialize first
      initializeAmplitude();

      const dydxAddress = 'dydx1abc123';
      const chainId = 'dydxprotocol';
      const amount = '1000';
      const transactionHash = '0x1234567890abcdef';
      const sourceAssetDenom = 'uusdc';

      await trackTurnkeyDepositSubmitted(
        dydxAddress,
        chainId,
        amount,
        transactionHash,
        sourceAssetDenom,
      );

      expect(mockTrack).toHaveBeenCalledWith(
        'TurnkeyDepositSubmitted',
        {
          chain_id: chainId,
          amount,
          transaction_hash: transactionHash,
          source_asset_denom: sourceAssetDenom,
          timestamp: expect.any(String),
        },
        {
          user_id: dydxAddress,
        },
      );
    });

    it('should track TurnKey deposit event without sourceAssetDenom', async () => {
      const { trackTurnkeyDepositSubmitted, initializeAmplitude } = require('../../src/lib/amplitude-helpers');

      // Initialize first
      initializeAmplitude();

      const dydxAddress = 'dydx1abc123';
      const chainId = 'dydxprotocol';
      const amount = '1000';
      const transactionHash = '0x1234567890abcdef';

      await trackTurnkeyDepositSubmitted(
        dydxAddress,
        chainId,
        amount,
        transactionHash,
      );

      expect(mockTrack).toHaveBeenCalledWith(
        'TurnkeyDepositSubmitted',
        {
          chain_id: chainId,
          amount,
          transaction_hash: transactionHash,
          source_asset_denom: undefined,
          timestamp: expect.any(String),
        },
        {
          user_id: dydxAddress,
        },
      );
    });

    it('should include current timestamp in event properties', async () => {
      const { trackTurnkeyDepositSubmitted, initializeAmplitude } = require('../../src/lib/amplitude-helpers');

      // Initialize first
      initializeAmplitude();

      const beforeTime = new Date();

      await trackTurnkeyDepositSubmitted(
        'dydx1abc123',
        'dydxprotocol',
        '1000',
        '0x1234567890abcdef',
      );

      const afterTime = new Date();

      expect(mockTrack).toHaveBeenCalledWith(
        'TurnkeyDepositSubmitted',
        expect.objectContaining({
          timestamp: expect.stringMatching(/^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{3}Z$/),
        }),
        expect.any(Object),
      );

      // Verify timestamp is within expected range
      const callArgs = mockTrack.mock.calls[0];
      const eventProperties = callArgs[1] as Record<string, unknown>;
      const timestamp = new Date(eventProperties.timestamp as string);
      expect(timestamp.getTime()).toBeGreaterThanOrEqual(beforeTime.getTime());
      expect(timestamp.getTime()).toBeLessThanOrEqual(afterTime.getTime());
    });

    it('should handle tracking errors gracefully', async () => {
      const error = new Error('Tracking failed');
      (mockTrack as jest.Mock).mockRejectedValueOnce(error);

      const { trackTurnkeyDepositSubmitted, initializeAmplitude } = require('../../src/lib/amplitude-helpers');

      // Initialize first
      initializeAmplitude();

      await trackTurnkeyDepositSubmitted(
        'dydx1abc123',
        'dydxprotocol',
        '1000',
        '0x1234567890abcdef',
      );

      expect(mockLogger.error).toHaveBeenCalledWith({
        at: 'amplitude-helpers#trackAmplitudeEvent',
        message: 'Failed to track Amplitude event',
        error,
        eventType: 'TurnkeyDepositSubmitted',
        userId: 'dydx1abc123',
        eventProperties: expect.objectContaining({
          chain_id: 'dydxprotocol',
          amount: '1000',
          transaction_hash: '0x1234567890abcdef',
          source_asset_denom: undefined,
          timestamp: expect.any(String),
        }),
      });
    });
  });

  describe('integration tests', () => {
    it('should work end-to-end with proper initialization and tracking', async () => {
      const {
        initializeAmplitude,
        trackAmplitudeEvent,
        trackTurnkeyDepositSubmitted,
      } = require('../../src/lib/amplitude-helpers');

      // Initialize amplitude
      initializeAmplitude();
      expect(mockInit).toHaveBeenCalledTimes(1);

      // Track a generic event
      await trackAmplitudeEvent('TestEvent', 'user123', { test: 'value' });
      expect(mockTrack).toHaveBeenCalledTimes(1);

      // Track a TurnKey deposit event
      await trackTurnkeyDepositSubmitted(
        'dydx1abc123',
        'dydxprotocol',
        '1000',
        '0x1234567890abcdef',
        'uusdc',
      );
      expect(mockTrack).toHaveBeenCalledTimes(2);

      // Verify all calls were made correctly
      expect(mockInit).toHaveBeenCalledWith('test-api-key', {
        serverUrl: 'https://api.eu.amplitude.com/2/httpapi',
        flushQueueSize: 1,
        flushIntervalMillis: 1000,
        serverZone: 'EU',
      });

      expect(mockTrack).toHaveBeenNthCalledWith(1, 'TestEvent', { test: 'value' }, {
        user_id: 'user123',
      });

      expect(mockTrack).toHaveBeenNthCalledWith(2, 'TurnkeyDepositSubmitted', {
        chain_id: 'dydxprotocol',
        amount: '1000',
        transaction_hash: '0x1234567890abcdef',
        source_asset_denom: 'uusdc',
        timestamp: expect.any(String),
      }, {
        user_id: 'dydx1abc123',
      });
    });
  });
});
