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
  constructor(topic: string) {
    super(`Invalid topic: ${topic}`);
    this.name = 'InvalidTopicError';
  }
}
