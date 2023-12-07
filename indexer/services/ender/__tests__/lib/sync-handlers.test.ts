import { SyncHandlers } from '../../src/lib/sync-handlers';
import { DydxIndexerSubtypes } from '../../src/lib/types';
import { KafkaPublisher } from '../../src/lib/kafka-publisher';
import * as pg from 'pg';
import { mock, MockProxy } from 'jest-mock-extended';
import { Handler } from '../../src/handlers/handler';

describe('syncHandler', () => {
  describe('addHandler/process', () => {
    it('successfully adds handler', async () => {
      const firstHandler: MockProxy<Handler<string>> = mock<Handler<string>>();
      firstHandler.blockEventIndex = 0;
      const secondHandler: MockProxy<Handler<string>> = mock<Handler<string>>();
      secondHandler.blockEventIndex = 1;
      const handlerNotInvoked: MockProxy<Handler<string>> = mock<Handler<string>>();
      handlerNotInvoked.blockEventIndex = 2;
      const syncHandlers: SyncHandlers = new SyncHandlers();

      // handlers are processed in the order in which they are received.
      syncHandlers.addHandler(DydxIndexerSubtypes.MARKET, firstHandler);
      syncHandlers.addHandler(DydxIndexerSubtypes.ASSET, secondHandler);
      // should be ignored, because transfers are not handled by syncHandlers
      syncHandlers.addHandler(DydxIndexerSubtypes.TRANSFER, handlerNotInvoked);

      const resultRow: pg.QueryResultRow = [
        'forFirstHandler',
        'forSecondHandler',
        'forNotInvokedHandler',
      ];
      const kafkaPublisher: KafkaPublisher = new KafkaPublisher();
      await syncHandlers.process(kafkaPublisher, resultRow);

      expect(firstHandler.handle).toHaveBeenCalledWith('forFirstHandler');
      expect(firstHandler.handle).toHaveBeenCalledTimes(1);
      expect(secondHandler.handle).toHaveBeenCalledWith('forSecondHandler');
      expect(secondHandler.handle).toHaveBeenCalledTimes(1);
      expect(firstHandler.handle.mock.invocationCallOrder[0]).toBeLessThan(
        secondHandler.handle.mock.invocationCallOrder[0]);
      expect(handlerNotInvoked.handle).toHaveBeenCalledTimes(0);
    });
  });
});
