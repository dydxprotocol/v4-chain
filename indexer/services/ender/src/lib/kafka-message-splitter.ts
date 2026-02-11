/**
 * Splits Kafka producer messages that exceed the broker's max message size
 * into multiple smaller messages. Supports SubaccountMessage (splits fills/orders)
 * and TradeMessage (splits trades array).
 */

import {
  KafkaTopics,
  ProducerMessage,
} from '@dydxprotocol-indexer/kafka';
import {
  SubaccountMessageContents,
  TradeMessageContents,
} from '@dydxprotocol-indexer/postgres';
import {
  SubaccountMessage as SubaccountMessageCodec,
  TradeMessage as TradeMessageCodec,
} from '@dydxprotocol-indexer/v4-protos';
import type { SubaccountMessage, TradeMessage } from '@dydxprotocol-indexer/v4-protos';

/**
 * Returns the total size of a producer message in bytes (key + value).
 */
function messageSize(msg: ProducerMessage): number {
  return (msg.key?.byteLength ?? 0) + msg.value.byteLength;
}

/**
 * Splits a single SubaccountMessage that exceeds maxBytes by splitting
 * the fills array, then orders array, into chunks. Returns multiple
 * ProducerMessages each under maxBytes. If the message cannot be split
 * (e.g. no splittable arrays or single item too large), returns [original].
 */
function splitSubaccountMessage(
  message: ProducerMessage,
  maxBytes: number,
): ProducerMessage[] {
  const decoded: SubaccountMessage = SubaccountMessageCodec.decode(
    new Uint8Array(message.value),
  );
  let contents: SubaccountMessageContents;
  try {
    contents = JSON.parse(decoded.contents) as SubaccountMessageContents;
  } catch {
    return [message];
  }

  const baseMeta: Pick<SubaccountMessage, 'blockHeight' | 'transactionIndex' | 'eventIndex' | 'subaccountId' | 'version'> = {
    blockHeight: decoded.blockHeight,
    transactionIndex: decoded.transactionIndex,
    eventIndex: decoded.eventIndex,
    subaccountId: decoded.subaccountId,
    version: decoded.version,
  };

  // Try splitting fills first (most common large payload)
  if (contents.fills && contents.fills.length > 1) {
    const chunks = chunkArrayToFitSize(
      contents.fills,
      maxBytes,
      (chunk) => {
        const newContents: SubaccountMessageContents = {
          ...contents,
          fills: chunk,
        };
        const proto: SubaccountMessage = {
          ...baseMeta,
          contents: JSON.stringify(newContents),
        };
        return Buffer.from(
          SubaccountMessageCodec.encode(proto).finish(),
        ).byteLength;
      },
    );
    if (chunks.length > 1) {
      return chunks.map((fillChunk) => {
        const newContents: SubaccountMessageContents = {
          ...contents,
          fills: fillChunk,
        };
        const proto: SubaccountMessage = {
          ...baseMeta,
          contents: JSON.stringify(newContents),
        };
        return {
          key: message.key,
          value: Buffer.from(SubaccountMessageCodec.encode(proto).finish()),
        };
      });
    }
  }

  // Try splitting orders
  if (contents.orders && contents.orders.length > 1) {
    const chunks = chunkArrayToFitSize(
      contents.orders,
      maxBytes,
      (chunk) => {
        const newContents: SubaccountMessageContents = {
          ...contents,
          orders: chunk,
        };
        const proto: SubaccountMessage = {
          ...baseMeta,
          contents: JSON.stringify(newContents),
        };
        return Buffer.from(
          SubaccountMessageCodec.encode(proto).finish(),
        ).byteLength;
      },
    );
    if (chunks.length > 1) {
      return chunks.map((orderChunk) => {
        const newContents: SubaccountMessageContents = {
          ...contents,
          orders: orderChunk,
        };
        const proto: SubaccountMessage = {
          ...baseMeta,
          contents: JSON.stringify(newContents),
        };
        return {
          key: message.key,
          value: Buffer.from(SubaccountMessageCodec.encode(proto).finish()),
        };
      });
    }
  }

  return [message];
}

/**
 * Splits a single TradeMessage that exceeds maxBytes by splitting
 * the trades array into chunks. Returns multiple ProducerMessages.
 */
function splitTradeMessage(
  message: ProducerMessage,
  maxBytes: number,
): ProducerMessage[] {
  const decoded: TradeMessage = TradeMessageCodec.decode(
    new Uint8Array(message.value),
  );
  let contents: TradeMessageContents;
  try {
    contents = JSON.parse(decoded.contents) as TradeMessageContents;
  } catch {
    return [message];
  }

  if (!contents.trades || contents.trades.length <= 1) {
    return [message];
  }

  const chunks = chunkArrayToFitSize(
    contents.trades,
    maxBytes,
    (chunk) => {
      const newContents: TradeMessageContents = { trades: chunk };
      const proto: TradeMessage = {
        ...decoded,
        contents: JSON.stringify(newContents),
      };
      return Buffer.from(TradeMessageCodec.encode(proto).finish()).byteLength;
    },
  );

  if (chunks.length <= 1) {
    return [message];
  }

  return chunks.map((tradeChunk) => {
    const newContents: TradeMessageContents = { trades: tradeChunk };
    const proto: TradeMessage = {
      ...decoded,
      contents: JSON.stringify(newContents),
    };
    return {
      value: Buffer.from(TradeMessageCodec.encode(proto).finish()),
    };
  });
}

/**
 * Splits an array into chunks such that each chunk's encoded size (as measured by
 * sizeFn) fits within maxBytes. sizeFn(chunk) returns the size in bytes of that chunk.
 * Uses a greedy approach: add items until the next would exceed maxBytes, then start a new chunk.
 */
function chunkArrayToFitSize<T>(
  arr: T[],
  maxBytes: number,
  sizeFn: (chunk: T[]) => number,
): T[][] {
  const result: T[][] = [];
  let current: T[] = [];

  for (const item of arr) {
    const candidate = current.concat([item]);
    const candidateSize = sizeFn(candidate);
    if (candidateSize > maxBytes && current.length > 0) {
      result.push(current);
      current = [item];
    } else {
      current = candidate;
    }
  }
  if (current.length > 0) {
    result.push(current);
  }
  return result;
}

/**
 * For a given topic and list of producer messages, splits any message that exceeds
 * maxBytes into multiple messages. Only topics with splittable message types
 * (TO_WEBSOCKETS_SUBACCOUNTS, TO_WEBSOCKETS_TRADES) are split; others are returned unchanged.
 */
export function splitOversizedMessages(
  topic: KafkaTopics,
  messages: ProducerMessage[],
  maxBytes: number,
): ProducerMessage[] {
  if (messages.length === 0) return messages;

  const result: ProducerMessage[] = [];
  for (const msg of messages) {
    const size = messageSize(msg);
    if (size <= maxBytes) {
      result.push(msg);
      continue;
    }
    const keyBytes = msg.key?.byteLength ?? 0;
    const maxValueBytes = Math.max(1, maxBytes - keyBytes);
    if (topic === KafkaTopics.TO_WEBSOCKETS_SUBACCOUNTS) {
      result.push(...splitSubaccountMessage(msg, maxValueBytes));
    } else if (topic === KafkaTopics.TO_WEBSOCKETS_TRADES) {
      result.push(...splitTradeMessage(msg, maxValueBytes));
    } else {
      result.push(msg);
    }
  }
  return result;
}
