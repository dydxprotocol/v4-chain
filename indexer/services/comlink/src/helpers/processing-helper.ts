import logger from '@dydxprotocol-indexer/base/build/src/logger';

/**
 * ProcessingQueue manages a queue of addresses currently being processed
 * to prevent duplicate processing of the same address
 */
export class ProcessingQueue {
  private static instance: ProcessingQueue;
  private processingAddresses: Set<string> = new Set();
  private readonly processingTimeouts: Map<string, NodeJS.Timeout> = new Map();
  private readonly TIMEOUT_MS = 5 * 60 * 1000; // 5 minutes timeout

  private constructor() {
    // Private constructor for singleton pattern
  }

  /**
   * Get the singleton instance of ProcessingQueue
   */
  public static getInstance(): ProcessingQueue {
    if (!ProcessingQueue.instance) {
      ProcessingQueue.instance = new ProcessingQueue();
    }
    return ProcessingQueue.instance;
  }

  /**
   * Check if an address is currently being processed
   * @param fromAddress The address to check
   * @returns true if the address is being processed, false otherwise
   */
  public isProcessing(fromAddress: string): boolean {
    return this.processingAddresses.has(fromAddress);
  }

  /**
   * Add an address to the processing queue
   * @param fromAddress The address to add to the queue
   * @returns true if successfully added, false if already in queue
   */
  public addToQueue(fromAddress: string): boolean {
    if (this.processingAddresses.has(fromAddress)) {
      logger.info({
        at: 'ProcessingQueue#addToQueue',
        message: 'Address already in processing queue, skipping',
        fromAddress,
      });
      return false;
    }

    this.processingAddresses.add(fromAddress);

    // Set a timeout to automatically remove the address if processing takes too long
    const timeout = setTimeout(() => {
      this.removeFromQueue(fromAddress);
      logger.warning({
        at: 'ProcessingQueue#addToQueue',
        message: 'Address removed from queue due to timeout',
        fromAddress,
        timeoutMs: this.TIMEOUT_MS,
      });
    }, this.TIMEOUT_MS);

    this.processingTimeouts.set(fromAddress, timeout);

    logger.info({
      at: 'ProcessingQueue#addToQueue',
      message: 'Address added to processing queue',
      fromAddress,
      queueSize: this.processingAddresses.size,
    });

    return true;
  }

  /**
   * Remove an address from the processing queue
   * @param fromAddress The address to remove from the queue
   */
  public removeFromQueue(fromAddress: string): void {
    const wasProcessing = this.processingAddresses.has(fromAddress);
    this.processingAddresses.delete(fromAddress);

    // Clear the timeout if it exists
    const timeout = this.processingTimeouts.get(fromAddress);
    if (timeout) {
      clearTimeout(timeout);
      this.processingTimeouts.delete(fromAddress);
    }

    if (wasProcessing) {
      logger.info({
        at: 'ProcessingQueue#removeFromQueue',
        message: 'Address removed from processing queue',
        fromAddress,
        queueSize: this.processingAddresses.size,
      });
    }
  }

  /**
   * Get the current size of the processing queue
   * @returns The number of addresses currently being processed
   */
  public getQueueSize(): number {
    return this.processingAddresses.size;
  }

  /**
   * Get all addresses currently in the queue (for debugging)
   * @returns Array of addresses currently being processed
   */
  public getProcessingAddresses(): string[] {
    return Array.from(this.processingAddresses);
  }

  /**
   * Clear all addresses from the queue (for testing or emergency cleanup)
   */
  public clearQueue(): void {
    // Clear all timeouts
    for (const timeout of this.processingTimeouts.values()) {
      clearTimeout(timeout);
    }
    this.processingTimeouts.clear();
    this.processingAddresses.clear();

    logger.info({
      at: 'ProcessingQueue#clearQueue',
      message: 'Processing queue cleared',
    });
  }
}
