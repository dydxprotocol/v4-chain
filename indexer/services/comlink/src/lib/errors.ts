export class V4ClientError extends Error {
  constructor() {
    super('An error occurred while querying the V4 application.');
    this.name = 'V4ClientError';
  }
}

export class UnexpectedServerError extends Error {
  constructor(message: string) {
    super(message);
    this.name = 'UnexpectedServerError';
  }
}

export class NotFoundError extends Error {
  constructor(message: string) {
    super(message);
    this.name = 'NotFoundError';
  }
}
