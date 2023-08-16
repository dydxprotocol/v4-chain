import { MARKETS_WEBSOCKET_MESSAGE_VERSION, WebsocketTopics } from '@dydxprotocol-indexer/kafka';
import { MarketMessage } from '@dydxprotocol-indexer/v4-protos';
import { ProducerRecord } from 'kafkajs';

export function expectMarketWebsocketMessage(
  producerSendMock: jest.SpyInstance,
  contents: string,
): void {
  expect(producerSendMock).toHaveBeenCalledTimes(1);
  const marketProducerRecord: ProducerRecord = producerSendMock.mock.calls[0][0];
  expect(marketProducerRecord.topic).toEqual(WebsocketTopics.TO_WEBSOCKETS_MARKETS);
  const marketMessageValueBinary: Uint8Array = new Uint8Array(
    marketProducerRecord.messages[0].value as Buffer,
  );

  const marketMessage: MarketMessage = MarketMessage.decode(
    marketMessageValueBinary,
  );
  const expectedMarketMessage: MarketMessage = MarketMessage.fromPartial({
    contents,
    version: MARKETS_WEBSOCKET_MESSAGE_VERSION,
  });
  expect(marketMessage).toEqual(expectedMarketMessage);
}
