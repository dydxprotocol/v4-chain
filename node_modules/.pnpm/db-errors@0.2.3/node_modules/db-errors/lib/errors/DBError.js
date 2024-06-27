'use strict';

class DBError extends Error {

  constructor(args) {
    super(args.nativeError.message);

    this.name = this.constructor.name;
    this.nativeError = args.nativeError;
    this.client = args.client;
  }
}

module.exports = DBError;