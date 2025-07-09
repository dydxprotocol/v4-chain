import { INDEXER_COMPLIANCE_BLOCKED_PAYLOAD } from '@dydxprotocol-indexer/compliance';

import { WebsocketTopic } from '../types';

export class InvalidForwardMessageError extends Error {
  constructor(message: string) {
    super(`Invalid forwarded message. Error: ${message}.`);
    this.name = 'InvalidForwardMessageError';
  }
}

export class InvalidChannelError extends Error {
  constructor(channel: string) {
    super(`Invalid channel: ${channel}`);
    this.name = 'InvalidChannelError';
  }
}

export class InvalidTopicError extends Error {
  constructor(topic: WebsocketTopic) {
    super(`Invalid topic: ${topic}`);
    this.name = 'InvalidTopicError';
  }
}

export class BlockedError extends Error {
  constructor() {
    super(INDEXER_COMPLIANCE_BLOCKED_PAYLOAD);
    this.name = 'BlockedError';
  }
}
