import { ClobPairStatus } from '@dydxprotocol-indexer/v4-protos';

/**
 * Custom error types.
 */

/**
 * Base class for custom errors.
 */
class CustomError extends Error {
  constructor(message: string) {
    super(message);
    // Set a more specific name. This will show up in e.g. console.log.
    this.name = this.constructor.name;
  }
}

export class RequiredFieldMissing extends CustomError {
  constructor(field: string) {
    super(`Required field '${field}' is missing`);
  }
}

export class ValidationError extends CustomError {
}

/**
 * Custom errors for converting protocol events into SQL.
 */

export class InvalidClobPairStatusError extends Error {
  constructor(clobPairStatus: ClobPairStatus) {
    super(`Invalid clob pair status: ${clobPairStatus}`);
    this.name = 'InvalidClobPairStatusError';
  }
}

export class PerpetualDoesNotExistError extends Error {
  constructor(perpetualId: number, clobPairId: number) {
    super(
      `Perpetual with id ${perpetualId} does not exist. Referenced in clob pair with id ` +
      `${clobPairId}.`,
    );
    this.name = 'PerpetualDoesNotExistError';
  }
}

export class MarketDoesNotExistError extends Error {
  constructor(marketId: number, perpetualId: number) {
    super(
      `Market with id ${marketId} does not exist. Referenced in perpetual with id ` +
      `${perpetualId}.`,
    );
    this.name = 'MarketDoesNotExistError';
  }
}

export class LiquidityTierDoesNotExistError extends Error {
  constructor(liquidityTier: number, perpetualId: number) {
    super(
      `Liquidity tier with id ${liquidityTier} does not exist. Referenced in perpetual with id ` +
      `${perpetualId}.`,
    );
    this.name = 'LiquidityTierDoesNotExistError';
  }
}
