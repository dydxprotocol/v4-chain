import { ProcessingQueue } from '../../src/helpers/processing-helper';

// Mock logger to avoid noise in tests
jest.mock('@dydxprotocol-indexer/base/build/src/logger', () => ({
  info: jest.fn(),
  warning: jest.fn(),
  error: jest.fn(),
}));

describe('ProcessingQueue', () => {
  let processingQueue: ProcessingQueue;

  beforeEach(() => {
    // Get a fresh instance for each test
    processingQueue = ProcessingQueue.getInstance();
    // Clear the queue before each test
    processingQueue.clearQueue();
    // Clear all mocks
    jest.clearAllMocks();
  });

  afterEach(() => {
    // Clean up after each test
    processingQueue.clearQueue();
  });

  describe('Singleton Pattern', () => {
    it('should return the same instance when getInstance is called multiple times', () => {
      const instance1 = ProcessingQueue.getInstance();
      const instance2 = ProcessingQueue.getInstance();
      const instance3 = ProcessingQueue.getInstance();

      expect(instance1).toBe(instance2);
      expect(instance2).toBe(instance3);
      expect(instance1).toBe(processingQueue);
    });
  });

  describe('Basic Queue Operations', () => {
    it('should start with an empty queue', () => {
      expect(processingQueue.getQueueSize()).toBe(0);
      expect(processingQueue.getProcessingAddresses()).toEqual([]);
    });

    it('should add an address to the queue successfully', () => {
      const address = 'test-address-1';
      const result = processingQueue.addToQueue(address);

      expect(result).toBe(true);
      expect(processingQueue.getQueueSize()).toBe(1);
      expect(processingQueue.isProcessing(address)).toBe(true);
      expect(processingQueue.getProcessingAddresses()).toContain(address);
    });

    it('should not add the same address twice', () => {
      const address = 'test-address-1';

      const firstAdd = processingQueue.addToQueue(address);
      const secondAdd = processingQueue.addToQueue(address);

      expect(firstAdd).toBe(true);
      expect(secondAdd).toBe(false);
      expect(processingQueue.getQueueSize()).toBe(1);
    });

    it('should remove an address from the queue', () => {
      const address = 'test-address-1';

      processingQueue.addToQueue(address);
      expect(processingQueue.isProcessing(address)).toBe(true);

      processingQueue.removeFromQueue(address);
      expect(processingQueue.isProcessing(address)).toBe(false);
      expect(processingQueue.getQueueSize()).toBe(0);
    });

    it('should handle removing an address that is not in the queue', () => {
      const address = 'non-existent-address';

      expect(() => processingQueue.removeFromQueue(address)).not.toThrow();
      expect(processingQueue.getQueueSize()).toBe(0);
    });

    it('should handle multiple addresses in the queue', () => {
      const addresses = ['address-1', 'address-2', 'address-3'];

      addresses.forEach((address) => {
        const result = processingQueue.addToQueue(address);
        expect(result).toBe(true);
      });

      expect(processingQueue.getQueueSize()).toBe(3);
      addresses.forEach((address) => {
        expect(processingQueue.isProcessing(address)).toBe(true);
      });

      // Remove one address
      processingQueue.removeFromQueue(addresses[1]);
      expect(processingQueue.getQueueSize()).toBe(2);
      expect(processingQueue.isProcessing(addresses[1])).toBe(false);
      expect(processingQueue.isProcessing(addresses[0])).toBe(true);
      expect(processingQueue.isProcessing(addresses[2])).toBe(true);
    });

    it('should clear all addresses from the queue', () => {
      const addresses = ['address-1', 'address-2', 'address-3'];

      addresses.forEach((address) => processingQueue.addToQueue(address));
      expect(processingQueue.getQueueSize()).toBe(3);

      processingQueue.clearQueue();
      expect(processingQueue.getQueueSize()).toBe(0);
      addresses.forEach((address) => {
        expect(processingQueue.isProcessing(address)).toBe(false);
      });
    });
  });

  describe('Timeout Functionality', () => {
    beforeEach(() => {
      jest.useFakeTimers();
    });

    afterEach(() => {
      jest.useRealTimers();
    });

    it('should automatically remove address after timeout', () => {
      const address = 'timeout-test-address';

      processingQueue.addToQueue(address);
      expect(processingQueue.isProcessing(address)).toBe(true);

      // Fast-forward time to trigger timeout (5 minutes)
      jest.advanceTimersByTime(5 * 60 * 1000);

      expect(processingQueue.isProcessing(address)).toBe(false);
      expect(processingQueue.getQueueSize()).toBe(0);
    });

    it('should not timeout if address is removed before timeout', () => {
      const address = 'no-timeout-test-address';

      processingQueue.addToQueue(address);
      expect(processingQueue.isProcessing(address)).toBe(true);

      // Remove before timeout
      processingQueue.removeFromQueue(address);
      expect(processingQueue.isProcessing(address)).toBe(false);

      // Fast-forward time past timeout
      jest.advanceTimersByTime(5 * 60 * 1000);

      // Should still be false (not cause any issues)
      expect(processingQueue.isProcessing(address)).toBe(false);
      expect(processingQueue.getQueueSize()).toBe(0);
    });

    it('should handle multiple addresses with different timeout times', () => {
      const address1 = 'timeout-address-1';
      const address2 = 'timeout-address-2';

      processingQueue.addToQueue(address1);

      // Add second address 1 minute later
      jest.advanceTimersByTime(1 * 60 * 1000);
      processingQueue.addToQueue(address2);

      expect(processingQueue.getQueueSize()).toBe(2);

      // Fast-forward 4 more minutes (total 5 minutes for address1, 4 for address2)
      jest.advanceTimersByTime(4 * 60 * 1000);

      // First address should timeout
      expect(processingQueue.isProcessing(address1)).toBe(false);
      expect(processingQueue.isProcessing(address2)).toBe(true);
      expect(processingQueue.getQueueSize()).toBe(1);

      // Fast-forward 1 more minute (total 5 minutes for address2)
      jest.advanceTimersByTime(1 * 60 * 1000);

      // Both addresses should be timed out
      expect(processingQueue.isProcessing(address1)).toBe(false);
      expect(processingQueue.isProcessing(address2)).toBe(false);
      expect(processingQueue.getQueueSize()).toBe(0);
    });

    it('should clear timeouts when queue is cleared', () => {
      const addresses = ['timeout-address-1', 'timeout-address-2'];

      addresses.forEach((address) => processingQueue.addToQueue(address));
      expect(processingQueue.getQueueSize()).toBe(2);

      processingQueue.clearQueue();
      expect(processingQueue.getQueueSize()).toBe(0);

      // Fast-forward past timeout
      jest.advanceTimersByTime(5 * 60 * 1000);

      // Should still be empty (timeouts were cleared)
      expect(processingQueue.getQueueSize()).toBe(0);
      addresses.forEach((address) => {
        expect(processingQueue.isProcessing(address)).toBe(false);
      });
    });
  });

  describe('Edge Cases', () => {
    it('should handle empty string addresses', () => {
      const result = processingQueue.addToQueue('');
      expect(result).toBe(true);
      expect(processingQueue.isProcessing('')).toBe(true);
      expect(processingQueue.getQueueSize()).toBe(1);
    });

    it('should handle very long addresses', () => {
      const longAddress = 'a'.repeat(1000);
      const result = processingQueue.addToQueue(longAddress);

      expect(result).toBe(true);
      expect(processingQueue.isProcessing(longAddress)).toBe(true);
    });

    it('should handle special characters in addresses', () => {
      const specialAddress = '0x123-456_789.test@example.com';
      const result = processingQueue.addToQueue(specialAddress);

      expect(result).toBe(true);
      expect(processingQueue.isProcessing(specialAddress)).toBe(true);
    });

    it('should handle many addresses without performance issues', () => {
      const numberOfAddresses = 1000;
      const addresses = Array.from({ length: numberOfAddresses }, (_, i) => `address-${i}`);

      const startTime = Date.now();
      addresses.forEach((address) => {
        processingQueue.addToQueue(address);
      });
      const addTime = Date.now() - startTime;

      expect(processingQueue.getQueueSize()).toBe(numberOfAddresses);
      expect(addTime).toBeLessThan(1000); // Should complete within 1 second

      const checkStartTime = Date.now();
      addresses.forEach((address) => {
        expect(processingQueue.isProcessing(address)).toBe(true);
      });
      const checkTime = Date.now() - checkStartTime;

      expect(checkTime).toBeLessThan(1000); // Should complete within 1 second
    });
  });
});
