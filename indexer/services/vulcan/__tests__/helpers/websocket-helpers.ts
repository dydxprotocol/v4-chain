import { KafkaTopics } from '@dydxprotocol-indexer/kafka';
import { OffChainUpdateV1, OrderbookMessage, SubaccountMessage } from '@dydxprotocol-indexer/v4-protos';
import { IHeaders, ProducerRecord } from 'kafkajs';

export function expectWebsocketSubaccountMessage(
  subaccountProducerRecord: ProducerRecord,
  expectedSubaccountMessage: SubaccountMessage,
  expectedHeaders: IHeaders,
): void {
  expect(subaccountProducerRecord.topic).toEqual(KafkaTopics.TO_WEBSOCKETS_SUBACCOUNTS);
  const subaccountMessageValueBinary: Uint8Array = new Uint8Array(
    subaccountProducerRecord.messages[0].value as Buffer,
  );
  const headers: IHeaders | undefined = subaccountProducerRecord.messages[0].headers;
  const subaccountMessage: SubaccountMessage = SubaccountMessage.decode(
    subaccountMessageValueBinary,
  );
  expect(headers).toEqual(expectedHeaders);
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

export function expectOffchainUpdateMessage(
  offchainUpdateProducerRecord: ProducerRecord,
  expectedKey: Buffer,
  expectedOffchainUpdate: OffChainUpdateV1,
): void {
  expect(offchainUpdateProducerRecord.topic).toEqual(KafkaTopics.TO_VULCAN);
  const offchainUpdateMessageValueBinary: Uint8Array = new Uint8Array(
    offchainUpdateProducerRecord.messages[0].value as Buffer,
  );
  const key: Buffer = offchainUpdateProducerRecord.messages[0].key as Buffer;
  const offchainUpdate: OffChainUpdateV1 = OffChainUpdateV1.decode(
    offchainUpdateMessageValueBinary,
  );
  expect(offchainUpdate).toEqual(expectedOffchainUpdate);
  expect(key).toEqual(expectedKey);
}
