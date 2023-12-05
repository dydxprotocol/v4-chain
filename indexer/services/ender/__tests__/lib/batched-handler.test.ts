import { IndexerTendermintBlock, IndexerTendermintEvent } from '@dydxprotocol-indexer/v4-protos';
import { Handler } from '../../src/handlers/handler';
import { BatchedHandlers } from '../../src/lib/batched-handlers';
import { ConsolidatedKafkaEvent, DydxIndexerSubtypes, EventMessage } from '../../src/lib/types';
import {
  defaultHeight,
  defaultTime,
  defaultTxHash,
} from '../helpers/constants';
import { createIndexerTendermintBlock, createIndexerTendermintEvent } from '../helpers/indexer-proto-helpers';

class FakeHandler extends Handler<EventMessage> {
  eventType: string = 'FakeEvent';
  parallelizationIds: string[] = [];

  public setParallelizationIds(parallelizationIds: string[]): void {
    this.parallelizationIds = parallelizationIds;
  }

  public getParallelizationIds(): string[] {
    return this.parallelizationIds;
  }

  public validate(): void {}

  public async internalHandle(): Promise<ConsolidatedKafkaEvent[]> {
    return Promise.resolve([]);
  }
}

function generateFakeHandler(parallelizationIds: string[]): FakeHandler {
  const block: IndexerTendermintBlock = createIndexerTendermintBlock(
    defaultHeight,
    defaultTime,
    [],
    [defaultTxHash],
  );

  const defaultTransactionIndex: number = 0;
  const defaultEventIndex: number = 0;
  const fakeTxId: number = 0;
  const defaultEvent: IndexerTendermintEvent = createIndexerTendermintEvent(
    DydxIndexerSubtypes.SUBACCOUNT_UPDATE,
    Uint8Array.from([]),
    defaultTransactionIndex,
    defaultEventIndex,
  );

  const handler: FakeHandler = new FakeHandler(
    block,
    0,
    defaultEvent,
    fakeTxId,
    {},
  );
  handler.setParallelizationIds(parallelizationIds);
  return handler;
}

describe('batchedHandler', () => {
  describe('addHandler', () => {
    it('successfully adds handler', () => {
      const batchedHandlers: BatchedHandlers = new BatchedHandlers();

      const pIds: string[] = ['1'];
      const handler: FakeHandler = generateFakeHandler(pIds);
      batchedHandlers.addHandler(handler);

      expect(batchedHandlers.batchedHandlers).toEqual([[handler]]);
      expect(batchedHandlers.pIdBatches).toEqual([new Set(pIds)]);
    });

    it('successfully adds handler to separate batch if parallelization ids overlap', () => {
      const batchedHandlers: BatchedHandlers = new BatchedHandlers();

      const pIds: string[] = ['1'];
      const handler: FakeHandler = generateFakeHandler(pIds);
      batchedHandlers.addHandler(handler);

      const handler2: FakeHandler = generateFakeHandler(pIds);
      batchedHandlers.addHandler(handler2);

      expect(batchedHandlers.batchedHandlers).toEqual([[handler], [handler2]]);
      expect(batchedHandlers.pIdBatches).toEqual([
        new Set(pIds),
        new Set(pIds),
      ]);
    });

    it('successfully adds handler to same batch if parallelization ids have no overlap', () => {
      const batchedHandlers: BatchedHandlers = new BatchedHandlers();

      const pIds: string[] = ['1'];
      const handler: FakeHandler = generateFakeHandler(pIds);
      batchedHandlers.addHandler(handler);

      const pIds2: string[] = ['2'];
      const handler2: FakeHandler = generateFakeHandler(pIds2);
      batchedHandlers.addHandler(handler2);

      expect(batchedHandlers.batchedHandlers).toEqual([[handler, handler2]]);
      expect(batchedHandlers.pIdBatches).toEqual([
        new Set(pIds.concat(pIds2)),
      ]);
    });

    /**
     * For the case where the following orders are processed:
     * 1. OrderFillEvent for Subaccount A and B
     * 2. OrderFillEvent for Subaccount B and C
     * 3. OrderFillEvent for Subaccount C and D
     * The handlers should be processed sequentially in the order they were added instead of having
     * OrderFillEvent 1 and 3 processed in parallel then OrderFillEvent 2.
     */
    it('successfully orders multiple handlers', () => {
      const batchedHandlers: BatchedHandlers = new BatchedHandlers();

      const pIds: string[] = ['A', 'B'];
      const handler: FakeHandler = generateFakeHandler(pIds);
      batchedHandlers.addHandler(handler);

      const pIds2: string[] = ['B', 'C'];
      const handler2: FakeHandler = generateFakeHandler(pIds2);
      batchedHandlers.addHandler(handler2);

      const pIds3: string[] = ['C', 'D'];
      const handler3: FakeHandler = generateFakeHandler(pIds3);
      batchedHandlers.addHandler(handler3);

      expect(batchedHandlers.batchedHandlers).toEqual([[handler], [handler2], [handler3]]);
      expect(batchedHandlers.pIdBatches).toEqual([
        new Set(pIds.concat(pIds2).concat(pIds3)),
        new Set(pIds2.concat(pIds3)),
        new Set(pIds3),
      ]);
    });
  });
});
