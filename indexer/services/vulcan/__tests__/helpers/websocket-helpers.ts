import { KafkaTopics } from '@dydxprotocol-indexer/kafka';
import { OrderbookMessage, SubaccountMessage } from '@dydxprotocol-indexer/v4-protos';
import { ProducerRecord } from 'kafkajs';

export function expectWebsocketSubaccountMessage(
  subaccountProducerRecord: ProducerRecord,
  expectedSubaccountMessage: SubaccountMessage,
): void {
  expect(subaccountProducerRecord.topic).toEqual(KafkaTopics.TO_WEBSOCKETS_SUBACCOUNTS);
  const subaccountMessageValueBinary: Uint8Array = new Uint8Array(
    subaccountProducerRecord.messages[0].value as Buffer,
  );
  const subaccountMessage: SubaccountMessage = SubaccountMessage.decode(
    subaccountMessageValueBinary,
  );
  expect(subaccountMessage).toEqual(expectedSubaccountMessage);
}

export function expectWebsocketOrderbookMessage(
  orderbookProducerRecord: ProducerRecord,
  expectedOrderbookMessage: OrderbookMessage,
): void {
  expect(orderbookProducerRecord.topic).toEqual(KafkaTopics.TO_WEBSOCKETS_ORDERBOOKS);
  const orderbookMessageValueBinary: Uint8Array = new Uint8Array(
    orderbookProducerRecord.messages[0].value as Buffer,
  );
  const orderbookMessage: OrderbookMessage = OrderbookMessage.decode(
    orderbookMessageValueBinary,
  );
  expect(orderbookMessage).toEqual(expectedOrderbookMessage);
}
