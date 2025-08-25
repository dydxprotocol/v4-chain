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

export class BadRequestError extends Error {
  constructor(message: string) {
    super(message);
    this.name = 'BadRequestError';
  }
}

export class DatabaseError extends Error {
  constructor(message: string) {
    super(message);
    this.name = 'DatabaseError';
  }
}

export class InvalidParamError extends Error {
  constructor(message: string) {
    super(message);
    this.name = 'InvalidParamError';
  }
}

export class TurnkeyError extends Error {
  constructor(message: string) {
    super(message);
    this.name = 'TurnkeyError';
  }
}
