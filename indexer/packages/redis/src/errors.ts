export class InvalidRedisOrderError extends Error {
  constructor(message: string) {
    super(`Invalid redis order: ${message}`);
    this.name = this.constructor.name;
    Error.captureStackTrace(this, this.constructor);
  }
}

export class InvalidTotalFilledQuantumsError extends Error {
  constructor(message: string) {
    super(`Invalid total filled quantums: ${message}`);
    this.name = this.constructor.name;
    Error.captureStackTrace(this, this.constructor);
  }
}

export class InvalidOptionsError extends Error {
  constructor(message: string) {
    super(`Invalid options passed in: ${message}`);
    this.name = this.constructor.name;
    Error.captureStackTrace(this, this.constructor);
  }
}
